package dkron

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	dkronpb "github.com/distribworks/dkron/v3/plugin/types"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
)

const (
	// MaxExecutions to maintain in the storage
	MaxExecutions = 100

	jobsPrefix       = "jobs"
	executionsPrefix = "executions"
)

var (
	// ErrDependentJobs is returned when deleting a job that has dependent jobs
	ErrDependentJobs = errors.New("store: could not delete job with dependent jobs, delete childs first")
)

// Store is the local implementation of the Storage interface.
// It gives dkron the ability to manipulate its embedded storage
// BuntDB.
type Store struct {
	db   *buntdb.DB
	lock *sync.Mutex // for
}

// JobOptions additional options to apply when loading a Job.
type JobOptions struct {
	Metadata map[string]string `json:"tags"`
}

type kv struct {
	Key   string
	Value []byte
}

// NewStore creates a new Storage instance.
func NewStore() (*Store, error) {
	db, err := buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}

	store := &Store{
		db:   db,
		lock: &sync.Mutex{},
	}

	return store, nil
}

func (s *Store) setJobTxFunc(pbj *dkronpb.Job) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		jobKey := fmt.Sprintf("%s:%s", jobsPrefix, pbj.Name)

		jb, err := proto.Marshal(pbj)
		if err != nil {
			return err
		}
		log.WithField("job", pbj.Name).Debug("store: Setting job")

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
func (s *Store) SetJob(job *Job, copyDependentJobs bool) error {
	var pbej dkronpb.Job
	var ej *Job

	if err := job.Validate(); err != nil {
		return err
	}

	// Abort if parent not found before committing job to the store
	if job.ParentJob != "" {
		if j, _ := s.GetJob(job.ParentJob, nil); j == nil {
			return ErrParentJobNotFound
		}
	}

	err := s.db.Update(func(tx *buntdb.Tx) error {
		// Get if the requested job already exist
		err := s.getJobTxFunc(job.Name, &pbej)(tx)
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}

		ej = NewJobFromProto(&pbej)

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
			// If comming from a backup us the previous value, don't allow overwriting this
			if job.Next.Before(ej.Next) {
				job.Next = ej.Next
			}
		}

		pbj := job.ToProto()
		s.setJobTxFunc(pbj)(tx)
		return nil
	})
	if err != nil {
		return err
	}

	// If the parent job changed update the parents of the old (if any) and new jobs
	if job.ParentJob != ej.ParentJob {
		if err := s.removeFromParent(ej); err != nil {
			return err
		}
		if err := s.addToParent(job); err != nil {
			return err
		}
	}

	return nil
}

// Removes the given job from its parent.
// Does nothing if nil is passed as child.
func (s *Store) removeFromParent(child *Job) error {
	// Do nothing if no job was given or job has no parent
	if child == nil || child.ParentJob == "" {
		return nil
	}

	parent, err := child.GetParent(s)
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
	if err := s.SetJob(parent, false); err != nil {
		return err
	}

	return nil
}

// Adds the given job to its parent.
func (s *Store) addToParent(child *Job) error {
	// Do nothing if job has no parent
	if child.ParentJob == "" {
		return nil
	}

	parent, err := child.GetParent(s)
	if err != nil {
		return err
	}

	parent.DependentJobs = append(parent.DependentJobs, child.Name)
	if err := s.SetJob(parent, false); err != nil {
		return err
	}

	return nil
}

// SetExecutionDone saves the execution and updates the job with the corresponding
// results
func (s *Store) SetExecutionDone(execution *Execution) (bool, error) {
	err := s.db.Update(func(tx *buntdb.Tx) error {
		// Load the job from the store
		var pbj dkronpb.Job
		if err := s.getJobTxFunc(execution.JobName, &pbj)(tx); err != nil {
			if err == buntdb.ErrNotFound {
				log.Warning(ErrExecutionDoneForDeletedJob)
				return ErrExecutionDoneForDeletedJob
			}
			log.WithError(err).Fatal(err)
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
		log.WithError(err).Error("store: Error in SetExecutionDone")
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
func (s *Store) GetJobs(options *JobOptions) ([]*Job, error) {
	jobs := make([]*Job, 0)

	err := s.db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			if strings.HasPrefix(key, jobsPrefix+":") {
				var pbj dkronpb.Job
				_ = proto.Unmarshal([]byte(value), &pbj)
				job := NewJobFromProto(&pbj)

				if options == nil || (options.Metadata == nil || len(options.Metadata) == 0 || s.jobHasMetadata(job, options.Metadata)) {
					jobs = append(jobs, job)
				}
			}
			return true
		})
		return err
	})

	return jobs, err
}

// GetJob finds and return a Job from the store
func (s *Store) GetJob(name string, options *JobOptions) (*Job, error) {
	var pbj dkronpb.Job

	err := s.db.View(s.getJobTxFunc(name, &pbj))
	if err != nil {
		return nil, err
	}

	job := NewJobFromProto(&pbj)

	return job, nil
}

// This will allow reuse this code to avoid nesting transactions
func (s *Store) getJobTxFunc(name string, pbj *dkronpb.Job) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		item, err := tx.Get(fmt.Sprintf("%s:%s", jobsPrefix, name))
		if err != nil {
			return err
		}

		if err := proto.Unmarshal([]byte(item), pbj); err != nil {
			return err
		}

		log.WithFields(logrus.Fields{
			"job": pbj.Name,
		}).Debug("store: Retrieved job from datastore")

		return nil
	}
}

// DeleteJob deletes the given job from the store, along with
// all its executions and references to it.
func (s *Store) DeleteJob(name string) (*Job, error) {
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
		job = NewJobFromProto(&pbj)

		if err := s.deleteExecutionsTxFunc(name)(tx); err != nil {
			return err
		}

		_, err := tx.Delete(fmt.Sprintf("%s:%s", jobsPrefix, name))
		return err
	})
	if err != nil {
		return nil, err
	}

	// If the transaction succeded, remove from parent
	if job.ParentJob != "" {
		if err := s.removeFromParent(job); err != nil {
			return nil, err
		}
	}

	return job, nil
}

// GetExecutions returns the exections given a Job name.
func (s *Store) GetExecutions(jobName string) ([]*Execution, error) {
	prefix := fmt.Sprintf("%s:%s:", executionsPrefix, jobName)

	kvs, err := s.list(prefix, true)
	if err != nil {
		return nil, err
	}

	return s.unmarshalExecutions(kvs)
}

func (s *Store) list(prefix string, checkRoot bool) ([]kv, error) {
	var found bool
	kvs := []kv{}

	err := s.db.View(s.listTxFunc(prefix, &kvs, &found))
	if err == nil && !found && checkRoot {
		return nil, buntdb.ErrNotFound
	}

	return kvs, err
}

func (*Store) listTxFunc(prefix string, kvs *[]kv, found *bool) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			if strings.HasPrefix(key, prefix) {
				*found = true
				// ignore self in listing
				if !bytes.Equal(trimDirectoryKey([]byte(key)), []byte(prefix)) {
					kv := kv{Key: key, Value: []byte(value)}
					*kvs = append(*kvs, kv)
				}
			}
			return true
		})
		return err
	}
}

// GetLastExecutionGroup get last execution group given the Job name.
func (s *Store) GetLastExecutionGroup(jobName string) ([]*Execution, error) {
	executions, byGroup, err := s.GetGroupedExecutions(jobName)
	if err != nil {
		return nil, err
	}

	if len(executions) > 0 && len(byGroup) > 0 {
		return executions[byGroup[0]], nil
	}

	return nil, nil
}

// GetExecutionGroup returns all executions in the same group of a given execution
func (s *Store) GetExecutionGroup(execution *Execution) ([]*Execution, error) {
	res, err := s.GetExecutions(execution.JobName)
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
func (s *Store) GetGroupedExecutions(jobName string) (map[int64][]*Execution, []int64, error) {
	execs, err := s.GetExecutions(jobName)
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
			if err := proto.Unmarshal([]byte(i), &p); err != nil {
				return err
			}
			// Compare existing execution
			if p.GetFinishedAt().Seconds > pbe.GetFinishedAt().Seconds {
				return nil
			}
		}

		eb, err := proto.Marshal(pbe)
		if err != nil {
			return err
		}

		_, _, err = tx.Set(key, string(eb), nil)
		return err
	}
}

// SetExecution Save a new execution and returns the key of the new saved item or an error.
func (s *Store) SetExecution(execution *Execution) (string, error) {
	pbe := execution.ToProto()
	key := fmt.Sprintf("%s:%s:%s", executionsPrefix, execution.JobName, execution.Key())

	log.WithFields(logrus.Fields{
		"job":       execution.JobName,
		"execution": key,
		"finished":  execution.FinishedAt.String(),
	}).Debug("store: Setting key")

	err := s.db.Update(s.setExecutionTxFunc(key, pbe))

	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"job":       execution.JobName,
			"execution": key,
		}).Debug("store: Failed to set key")
		return "", err
	}

	execs, err := s.GetExecutions(execution.JobName)
	if err != nil && err != buntdb.ErrNotFound {
		log.WithError(err).
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
			log.WithFields(logrus.Fields{
				"job":       execs[i].JobName,
				"execution": execs[i].Key(),
			}).Debug("store: to detele key")
			err = s.db.Update(func(tx *buntdb.Tx) error {
				k := fmt.Sprintf("%s:%s:%s", executionsPrefix, execs[i].JobName, execs[i].Key())
				_, err := tx.Delete(k)
				return err
			})
			if err != nil {
				log.WithError(err).
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
		tx.Ascend("", func(key, value string) bool {
			if strings.HasPrefix(key, prefix) {
				delkeys = append(delkeys, key)
			}
			return true
		})

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
	return s.db.Save(w)
}

// Restore load data created with backup in to Bunt
func (s *Store) Restore(r io.ReadCloser) error {
	return s.db.Load(r)
}

func (s *Store) unmarshalExecutions(items []kv) ([]*Execution, error) {
	var executions []*Execution
	for _, item := range items {
		var pbe dkronpb.Execution

		if err := proto.Unmarshal(item.Value, &pbe); err != nil {
			log.WithError(err).WithField("key", item.Key).Debug("error unmarshaling")
			return nil, err
		}
		execution := NewExecutionFromProto(&pbe)
		executions = append(executions, execution)
	}
	return executions, nil
}

func (s *Store) computeStatus(jobName string, exGroup int64, tx *buntdb.Tx) (string, error) {
	// compute job status based on execution group
	kvs := []kv{}
	found := false
	prefix := fmt.Sprintf("%s:%s:", executionsPrefix, jobName)

	if err := s.listTxFunc(prefix, &kvs, &found)(tx); err != nil {
		return "", err
	}

	execs, err := s.unmarshalExecutions(kvs)
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
		status = StatusPartialyFailed
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
