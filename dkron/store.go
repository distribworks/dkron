package dkron

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	dkronpb "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

const (
	// MaxExecutions to maintain in the storage
	MaxExecutions = 100

	jobsPrefix       = "jobs"
	executionsPrefix = "executions"
)

var (
	// ErrDependentJobs is returned when deleting a job that has dependent jobs
	ErrDependentJobs = errors.New("store: could not delete job with dependent jobs, delete children first")
)

// Store is the local implementation of the Storage interface.
// It gives dkron the ability to manipulate its embedded storage
// BuntDB.
type Store struct {
	db   *buntdb.DB
	lock *sync.Mutex

	tracer trace.Tracer
	logger *logrus.Entry
}

// JobOptions additional options to apply when loading a Job.
type JobOptions struct {
	Metadata map[string]string `json:"tags"`
	Sort     string
	Order    string
	Query    string
	Status   string
	Disabled string
}

// ExecutionOptions additional options like "Sort" will be ready for JSON marshall
type ExecutionOptions struct {
	Sort     string
	Order    string
	Timezone *time.Location
}

type kv struct {
	Key   string
	Value []byte
}

// NewStore creates a new Storage instance.
func NewStore(logger *logrus.Entry, tracer trace.Tracer) (*Store, error) {
	db, err := buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}
	_ = db.CreateIndex("name", jobsPrefix+":*", buntdb.IndexJSON("name"))
	_ = db.CreateIndex("started_at", executionsPrefix+":*", buntdb.IndexJSON("started_at"))
	_ = db.CreateIndex("finished_at", executionsPrefix+":*", buntdb.IndexJSON("finished_at"))
	_ = db.CreateIndex("attempt", executionsPrefix+":*", buntdb.IndexJSON("attempt"))
	_ = db.CreateIndex("displayname", jobsPrefix+":*", buntdb.IndexJSON("displayname"))
	_ = db.CreateIndex("schedule", jobsPrefix+":*", buntdb.IndexJSON("schedule"))
	_ = db.CreateIndex("success_count", jobsPrefix+":*", buntdb.IndexJSON("success_count"))
	_ = db.CreateIndex("error_count", jobsPrefix+":*", buntdb.IndexJSON("error_count"))
	_ = db.CreateIndex("last_success", jobsPrefix+":*", buntdb.IndexJSON("last_success"))
	_ = db.CreateIndex("last_error", jobsPrefix+":*", buntdb.IndexJSON("last_error"))
	_ = db.CreateIndex("next", jobsPrefix+":*", buntdb.IndexJSON("next"))

	store := &Store{
		db:     db,
		lock:   &sync.Mutex{},
		logger: logger,
		tracer: tracer,
	}

	return store, nil
}

func (s *Store) setJobTxFunc(pbj *dkronpb.Job) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		jobKey := fmt.Sprintf("%s:%s", jobsPrefix, pbj.Name)

		jb, err := json.Marshal(pbj)
		if err != nil {
			return err
		}
		s.logger.WithField("job", pbj.Name).Debug("store: Setting job")

		if _, _, err := tx.Set(jobKey, string(jb), nil); err != nil {
			return err
		}

		return nil
	}
}

// DB is the getter for the BuntDB instance
func (s *Store) DB() *buntdb.DB {
	return s.db
}

// SetJob stores a job in the storage
func (s *Store) SetJob(ctx context.Context, job *Job, copyDependentJobs bool) error {
	ctx, span := s.tracer.Start(ctx, "buntdb.set.job")
	defer span.End()

	var pbej dkronpb.Job
	var ej *Job

	if err := job.Validate(); err != nil {
		return err
	}

	// Abort if parent not found before committing job to the store
	if job.ParentJob != "" {
		if j, _ := s.GetJob(ctx, job.ParentJob, nil); j == nil {
			return ErrParentJobNotFound
		}
	}

	err := s.db.Update(func(tx *buntdb.Tx) error {
		// Get if the requested job already exist
		err := s.getJobTxFunc(job.Name, &pbej)(tx)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}

		ej = NewJobFromProto(&pbej, s.logger)

		if ej.Name != "" {
			// When the job runs, these status vars are updated
			// otherwise use the ones that are stored
			if ej.LastError.After(job.LastError) {
				job.LastError = ej.LastError
			}
			if ej.LastSuccess.After(job.LastSuccess) {
				job.LastSuccess = ej.LastSuccess
			}
			if ej.SuccessCount > job.SuccessCount {
				job.SuccessCount = ej.SuccessCount
			}
			if ej.ErrorCount > job.ErrorCount {
				job.ErrorCount = ej.ErrorCount
			}
			if len(ej.DependentJobs) != 0 && copyDependentJobs {
				job.DependentJobs = ej.DependentJobs
			}
			if ej.Status != "" {
				job.Status = ej.Status
			}
		}

		if job.Schedule != ej.Schedule {
			job.Next, err = job.GetNext()
			if err != nil {
				return err
			}
		} else {
			// If coming from a backup us the previous value, don't allow overwriting this
			if job.Next.Before(ej.Next) {
				job.Next = ej.Next
			}
		}

		pbj := job.ToProto()
		if err := s.setJobTxFunc(pbj)(tx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// If the parent job changed update the parents of the old (if any) and new jobs
	if job.ParentJob != ej.ParentJob {
		if err := s.removeFromParent(ctx, ej); err != nil {
			return err
		}
		if err := s.addToParent(ctx, job); err != nil {
			return err
		}
	}

	return nil
}

// Removes the given job from its parent.
// Does nothing if nil is passed as child.
func (s *Store) removeFromParent(ctx context.Context, child *Job) error {
	ctx, span := s.tracer.Start(ctx, "buntdb.remove_from_parent")
	defer span.End()

	// Do nothing if no job was given or job has no parent
	if child == nil || child.ParentJob == "" {
		return nil
	}

	parent, err := child.GetParent(ctx, s)
	if err != nil {
		return err
	}

	// Remove all occurrences from the parent, not just one.
	// Due to an old bug (in v1), a parent can have the same child more than once.
	djs := []string{}
	for _, djn := range parent.DependentJobs {
		if djn != child.Name {
			djs = append(djs, djn)
		}
	}
	parent.DependentJobs = djs
	if err := s.SetJob(ctx, parent, false); err != nil {
		return err
	}

	return nil
}

// Adds the given job to its parent.
func (s *Store) addToParent(ctx context.Context, child *Job) error {
	ctx, span := s.tracer.Start(ctx, "buntdb.add_to_parent")
	defer span.End()

	// Do nothing if job has no parent
	if child.ParentJob == "" {
		return nil
	}

	parent, err := child.GetParent(ctx, s)
	if err != nil {
		return err
	}

	parent.DependentJobs = append(parent.DependentJobs, child.Name)
	if err := s.SetJob(ctx, parent, false); err != nil {
		return err
	}

	return nil
}

// SetExecutionDone saves the execution and updates the job with the corresponding
// results
func (s *Store) SetExecutionDone(ctx context.Context, execution *Execution) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "buntdb.set.execution_done")
	defer span.End()

	err := s.db.Update(func(tx *buntdb.Tx) error {
		// Load the job from the store
		var pbj dkronpb.Job
		if err := s.getJobTxFunc(execution.JobName, &pbj)(tx); err != nil {
			if err == buntdb.ErrNotFound {
				s.logger.Warn(ErrExecutionDoneForDeletedJob)
				return ErrExecutionDoneForDeletedJob
			}
			s.logger.WithError(err).Fatal(err)
			return err
		}

		key := fmt.Sprintf("%s:%s:%s", executionsPrefix, execution.JobName, execution.Key())

		// Save the execution to store
		pbe := execution.ToProto()
		if err := s.setExecutionTxFunc(key, pbe)(tx); err != nil {
			return err
		}

		if pbe.Success {
			pbj.LastSuccess.HasValue = true
			pbj.LastSuccess.Time = pbe.FinishedAt
			pbj.SuccessCount++
		} else {
			pbj.LastError.HasValue = true
			pbj.LastError.Time = pbe.FinishedAt
			pbj.ErrorCount++
		}

		status, err := s.computeStatus(pbj.Name, pbe.Group, tx)
		if err != nil {
			return err
		}
		pbj.Status = status

		if err := s.setJobTxFunc(&pbj)(tx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.logger.WithError(err).Error("store: Error in SetExecutionDone")
		return false, err
	}

	return true, nil
}

func (s *Store) jobHasMetadata(job *Job, metadata map[string]string) bool {
	if job == nil || job.Metadata == nil || len(job.Metadata) == 0 {
		return false
	}

	for k, v := range metadata {
		if val, ok := job.Metadata[k]; !ok || v != val {
			return false
		}
	}

	return true
}

// GetJobs returns all jobs
func (s *Store) GetJobs(ctx context.Context, options *JobOptions) ([]*Job, error) {
	ctx, span := s.tracer.Start(ctx, "buntdb.get.jobs")
	defer span.End()

	if options == nil {
		options = &JobOptions{
			Sort: "name",
		}
	}

	jobs := make([]*Job, 0)
	jobsFn := func(key, item string) bool {
		var pbj dkronpb.Job
		// [TODO] This condition is temporary while we migrate to JSON marshalling for jobs
		// so we can use BuntDb indexes. To be removed in future versions.
		if err := proto.Unmarshal([]byte(item), &pbj); err != nil {
			if err := json.Unmarshal([]byte(item), &pbj); err != nil {
				return false
			}
		}
		job := NewJobFromProto(&pbj, s.logger)
		job.logger = s.logger

		if options == nil ||
			(options.Metadata == nil || len(options.Metadata) == 0 || s.jobHasMetadata(job, options.Metadata)) &&
				(options.Query == "" || strings.Contains(job.Name, options.Query) || strings.Contains(job.DisplayName, options.Query)) &&
				(options.Disabled == "" || strconv.FormatBool(job.Disabled) == options.Disabled) &&
				((options.Status == "untriggered" && job.Status == "") || (options.Status == "" || job.Status == options.Status)) {

			jobs = append(jobs, job)
		}
		return true
	}

	err := s.db.View(func(tx *buntdb.Tx) error {
		var err error
		if options.Order == "DESC" {
			err = tx.Descend(options.Sort, jobsFn)
		} else {
			err = tx.Ascend(options.Sort, jobsFn)
		}
		return err
	})

	return jobs, err
}

// GetJob finds and return a Job from the store
func (s *Store) GetJob(ctx context.Context, name string, options *JobOptions) (*Job, error) {
	ctx, span := s.tracer.Start(ctx, "buntdb.get.job", trace.WithAttributes(attribute.String("job_name", name)))
	defer span.End()

	var pbj dkronpb.Job

	err := s.db.View(s.getJobTxFunc(name, &pbj))
	if err != nil {
		return nil, err
	}

	job := NewJobFromProto(&pbj, s.logger)
	job.logger = s.logger

	return job, nil
}

// This will allow reuse this code to avoid nesting transactions
func (s *Store) getJobTxFunc(name string, pbj *dkronpb.Job) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		item, err := tx.Get(fmt.Sprintf("%s:%s", jobsPrefix, name))
		if err != nil {
			return err
		}

		// [TODO] This condition is temporary while we migrate to JSON marshalling for jobs
		// so we can use BuntDb indexes. To be removed in future versions.
		if err := proto.Unmarshal([]byte(item), pbj); err != nil {
			if err := json.Unmarshal([]byte(item), pbj); err != nil {
				return err
			}
		}

		s.logger.WithFields(logrus.Fields{
			"job": pbj.Name,
		}).Debug("store: Retrieved job from datastore")

		return nil
	}
}

// DeleteJob deletes the given job from the store, along with
// all its executions and references to it.
func (s *Store) DeleteJob(ctx context.Context, name string) (*Job, error) {
	ctx, span := s.tracer.Start(ctx, "buntdb.delete.job", trace.WithAttributes(attribute.String("job_name", name)))
	defer span.End()

	var job *Job
	err := s.db.Update(func(tx *buntdb.Tx) error {
		// Get the job
		var pbj dkronpb.Job
		if err := s.getJobTxFunc(name, &pbj)(tx); err != nil {
			return err
		}
		// Check if the job has dependent jobs
		// and return an error indicating to remove childs
		// first.
		if len(pbj.DependentJobs) > 0 {
			return ErrDependentJobs
		}
		job = NewJobFromProto(&pbj, s.logger)

		if err := s.deleteExecutionsTxFunc(name)(tx); err != nil {
			return err
		}

		_, err := tx.Delete(fmt.Sprintf("%s:%s", jobsPrefix, name))
		return err
	})
	if err != nil {
		return nil, err
	}

	// If the transaction succeeded, remove from parent
	if job.ParentJob != "" {
		if err := s.removeFromParent(ctx, job); err != nil {
			return nil, err
		}
	}

	return job, nil
}

// GetExecutions returns the executions given a Job name.
func (s *Store) GetExecutions(ctx context.Context, jobName string, opts *ExecutionOptions) ([]*Execution, error) {
	ctx, span := s.tracer.Start(ctx, "buntdb.get.executions", trace.WithAttributes(attribute.String("job_name", jobName)))
	defer span.End()

	prefix := fmt.Sprintf("%s:%s:", executionsPrefix, jobName)

	kvs, err := s.list(prefix, true, opts)
	if err != nil {
		return nil, err
	}

	return s.unmarshalExecutions(kvs, opts.Timezone)
}

func (s *Store) list(prefix string, checkRoot bool, opts *ExecutionOptions) ([]kv, error) {
	var found bool
	kvs := []kv{}

	err := s.db.View(s.listTxFunc(prefix, &kvs, &found, opts))
	if err == nil && !found && checkRoot {
		return nil, buntdb.ErrNotFound
	}

	return kvs, err
}

func (*Store) listTxFunc(prefix string, kvs *[]kv, found *bool, opts *ExecutionOptions) func(tx *buntdb.Tx) error {
	fnc := func(key, value string) bool {
		if strings.HasPrefix(key, prefix) {
			*found = true
			// ignore self in listing
			if !bytes.Equal(trimDirectoryKey([]byte(key)), []byte(prefix)) {
				kv := kv{Key: key, Value: []byte(value)}
				*kvs = append(*kvs, kv)
			}
		}
		return true
	}

	return func(tx *buntdb.Tx) (err error) {
		if opts.Order == "DESC" {
			err = tx.Descend(opts.Sort, fnc)
		} else {
			err = tx.Ascend(opts.Sort, fnc)
		}
		return err
	}
}

// GetExecutionGroup returns all executions in the same group of a given execution
func (s *Store) GetExecutionGroup(ctx context.Context, execution *Execution, opts *ExecutionOptions) ([]*Execution, error) {
	ctx, span := s.tracer.Start(ctx, "buntdb.get.executions_group")
	defer span.End()

	res, err := s.GetExecutions(ctx, execution.JobName, opts)
	if err != nil {
		return nil, err
	}

	var executions []*Execution
	for _, ex := range res {
		if ex.Group == execution.Group {
			executions = append(executions, ex)
		}
	}
	return executions, nil
}

// GetGroupedExecutions returns executions for a job grouped and with an ordered index
// to facilitate access.
func (s *Store) GetGroupedExecutions(ctx context.Context, jobName string, opts *ExecutionOptions) (map[int64][]*Execution, []int64, error) {
	ctx, span := s.tracer.Start(ctx, "buntdb.get.grouped_executions", trace.WithAttributes(attribute.String("job_name", jobName)))
	defer span.End()

	execs, err := s.GetExecutions(ctx, jobName, opts)
	if err != nil {
		return nil, nil, err
	}
	groups := make(map[int64][]*Execution)
	for _, exec := range execs {
		groups[exec.Group] = append(groups[exec.Group], exec)
	}

	// Build a separate data structure to show in order
	var byGroup int64arr
	for key := range groups {
		byGroup = append(byGroup, key)
	}
	sort.Sort(sort.Reverse(byGroup))

	return groups, byGroup, nil
}

func (*Store) setExecutionTxFunc(key string, pbe *dkronpb.Execution) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		// Get previous execution
		i, err := tx.Get(key)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		// Do nothing if a previous execution exists and is
		// more recent, avoiding non ordered execution set
		if i != "" {
			var p dkronpb.Execution
			// [TODO] This condition is temporary while we migrate to JSON marshalling for executions
			// so we can use BuntDb indexes. To be removed in future versions.
			if err := proto.Unmarshal([]byte(i), &p); err != nil {
				if err := json.Unmarshal([]byte(i), &p); err != nil {
					return err
				}
			}
			// Compare existing execution
			if p.GetFinishedAt().Seconds > pbe.GetFinishedAt().Seconds {
				return nil
			}
		}

		eb, err := json.Marshal(pbe)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(key, string(eb), nil)
		return err
	}
}

// SetExecution Save a new execution and returns the key of the new saved item or an error.
func (s *Store) SetExecution(ctx context.Context, execution *Execution) (string, error) {
	ctx, span := s.tracer.Start(ctx, "buntdb.set.executions")
	defer span.End()

	pbe := execution.ToProto()
	key := fmt.Sprintf("%s:%s:%s", executionsPrefix, execution.JobName, execution.Key())

	s.logger.WithFields(logrus.Fields{
		"job":       execution.JobName,
		"execution": key,
		"finished":  execution.FinishedAt.String(),
	}).Debug("store: Setting key")

	err := s.db.Update(s.setExecutionTxFunc(key, pbe))

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"job":       execution.JobName,
			"execution": key,
		}).Debug("store: Failed to set key")
		return "", err
	}

	execs, err := s.GetExecutions(ctx, execution.JobName, &ExecutionOptions{})
	if err != nil && err != buntdb.ErrNotFound {
		s.logger.WithError(err).
			WithField("job", execution.JobName).
			Error("store: Error getting executions for job")
	}

	// Delete all execution results over the limit, starting from olders
	if len(execs) > MaxExecutions {
		//sort the array of all execution groups by StartedAt time
		sort.Slice(execs, func(i, j int) bool {
			return execs[i].StartedAt.Before(execs[j].StartedAt)
		})

		for i := 0; i < len(execs)-MaxExecutions; i++ {
			s.logger.WithFields(logrus.Fields{
				"job":       execs[i].JobName,
				"execution": execs[i].Key(),
			}).Debug("store: to delete key")
			err = s.db.Update(func(tx *buntdb.Tx) error {
				k := fmt.Sprintf("%s:%s:%s", executionsPrefix, execs[i].JobName, execs[i].Key())
				_, err := tx.Delete(k)
				return err
			})
			if err != nil {
				s.logger.WithError(err).
					WithField("execution", execs[i].Key()).
					Error("store: Error trying to delete overflowed execution")
			}
		}
	}

	return key, nil
}

// DeleteExecutions removes all executions of a job
func (s *Store) deleteExecutionsTxFunc(jobName string) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		var delkeys []string
		prefix := fmt.Sprintf("%s:%s", executionsPrefix, jobName)
		if err := tx.Ascend("", func(key, value string) bool {
			if strings.HasPrefix(key, prefix) {
				delkeys = append(delkeys, key)
			}
			return true
		}); err != nil {
			return err
		}

		for _, k := range delkeys {
			_, _ = tx.Delete(k)
		}

		return nil
	}
}

// Shutdown close the KV store
func (s *Store) Shutdown() error {
	return s.db.Close()
}

// Snapshot creates a backup of the data stored in BuntDB
func (s *Store) Snapshot(w io.WriteCloser) error {
	_, span := s.tracer.Start(context.Background(), "buntdb.create.snapshot")
	defer span.End()

	return s.db.Save(w)
}

// Restore load data created with backup in to Bunt
func (s *Store) Restore(r io.ReadCloser) error {
	_, span := s.tracer.Start(context.Background(), "buntdb.restore.snapshot")
	defer span.End()

	return s.db.Load(r)
}

func (s *Store) unmarshalExecutions(items []kv, timezone *time.Location) ([]*Execution, error) {
	var executions []*Execution
	for _, item := range items {
		var pbe dkronpb.Execution

		// [TODO] This condition is temporary while we migrate to JSON marshalling for jobs
		// so we can use BuntDb indexes. To be removed in future versions.
		if err := proto.Unmarshal([]byte(item.Value), &pbe); err != nil {
			if err := json.Unmarshal(item.Value, &pbe); err != nil {
				s.logger.WithError(err).WithField("key", item.Key).Debug("error unmarshaling JSON")
				return nil, err
			}
		}
		execution := NewExecutionFromProto(&pbe)
		if timezone != nil {
			execution.FinishedAt = execution.FinishedAt.In(timezone)
			execution.StartedAt = execution.StartedAt.In(timezone)
		}
		executions = append(executions, execution)
	}
	return executions, nil
}

func (s *Store) computeStatus(jobName string, exGroup int64, tx *buntdb.Tx) (string, error) {
	// compute job status based on execution group
	kvs := []kv{}
	found := false
	prefix := fmt.Sprintf("%s:%s:", executionsPrefix, jobName)

	if err := s.listTxFunc(prefix, &kvs, &found, &ExecutionOptions{})(tx); err != nil {
		return "", err
	}

	execs, err := s.unmarshalExecutions(kvs, nil)
	if err != nil {
		return "", err
	}

	var executions []*Execution
	for _, ex := range execs {
		if ex.Group == exGroup {
			executions = append(executions, ex)
		}
	}

	success := 0
	failed := 0

	var status string
	for _, ex := range executions {
		if ex.Success {
			success = success + 1
		} else {
			failed = failed + 1
		}
	}

	if failed == 0 {
		status = StatusSuccess
	} else if failed > 0 && success == 0 {
		status = StatusFailed
	} else if failed > 0 && success > 0 {
		status = StatusPartiallyFailed
	}

	return status, nil
}

func trimDirectoryKey(key []byte) []byte {
	if isDirectoryKey(key) {
		return key[:len(key)-1]
	}

	return key
}

func isDirectoryKey(key []byte) bool {
	return len(key) > 0 && key[len(key)-1] == ':'
}
