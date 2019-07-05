package dkron

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/distribworks/dkron/cron"
	dkronpb "github.com/distribworks/dkron/proto"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

const (
	MaxExecutions            = 100
	defaultUpdateMaxAttempts = 5
	defaultGCInterval        = 5 * time.Minute
	defaultGCDiscardRatio    = 0.7
)

var (
	// ErrTooManyUpdateConflicts is returned when all update attempts fails
	ErrTooManyUpdateConflicts = errors.New("badger: too many transaction conflicts")
)

type Store struct {
	agent  *Agent
	db     *badger.DB
	lock   *sync.Mutex // for
	closed bool
}

type JobOptions struct {
	ComputeStatus bool
	Metadata      map[string]string `json:"tags"`
}

func NewStore(a *Agent, dir string) (*Store, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	store := &Store{
		db:    db,
		agent: a,
		lock:  &sync.Mutex{},
	}

	go store.runGcLoop()

	return store, nil
}

func (s *Store) runGcLoop() {
	ticker := time.NewTicker(defaultGCInterval)
	defer ticker.Stop()
	for range ticker.C {
		s.lock.Lock()
		closed := s.closed
		s.lock.Unlock()
		if closed {
			break
		}

		// One call would only result in removal of at max one log file.
		// As an optimization, you could also immediately re-run it whenever it returns nil error
		//(indicating a successful value log GC), as shown below.
	again:
		err := s.db.RunValueLogGC(defaultGCDiscardRatio)
		if err == nil {
			goto again
		}
	}
}

// Store a job
func (s *Store) SetJob(job *Job, copyDependentJobs bool) error {
	//Existing job that has children, let's keep it's children

	jobKey := fmt.Sprintf("jobs/%s", job.Name)

	// Init the job agent
	job.Agent = s.agent

	if err := s.validateJob(job); err != nil {
		return err
	}

	err := s.db.Update(func(txn *badger.Txn) error {
		// Get if the requested job already exist
		ej, err := s.GetJob(job.Name, nil)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if ej != nil {
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
		}

		pbj := job.ToProto()
		jb, err := proto.Marshal(pbj)
		if err != nil {
			return err
		}
		log.WithField("job", job.Name).Debug("store: Setting job")

		if err := txn.Set([]byte(jobKey), jb); err != nil {
			return err
		}

		if ej != nil {
			// Existing job that doesn't have parent job set and it's being set
			if ej.ParentJob == "" && job.ParentJob != "" {
				pj, err := job.GetParent()
				if err != nil {
					return err
				}

				pj.DependentJobs = append(pj.DependentJobs, job.Name)
				if err := s.SetJob(pj, false); err != nil {
					return err
				}
			}

			// Existing job that has parent job set and it's being removed
			if ej.ParentJob != "" && job.ParentJob == "" {
				pj, err := ej.GetParent()
				if err != nil {
					return err
				}

				ndx := 0
				for i, djn := range pj.DependentJobs {
					if djn == job.Name {
						ndx = i
						break
					}
				}
				pj.DependentJobs = append(pj.DependentJobs[:ndx], pj.DependentJobs[ndx+1:]...)
				if err := s.SetJob(pj, false); err != nil {
					return err
				}
			}
		}

		// New job that has parent job set
		if ej == nil && job.ParentJob != "" {
			pj, err := job.GetParent()
			if err != nil {
				return err
			}

			pj.DependentJobs = append(pj.DependentJobs, job.Name)
			if err := s.SetJob(pj, false); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (s *Store) validateTimeZone(timezone string) error {
	if timezone == "" {
		return nil
	}
	_, err := time.LoadLocation(timezone)
	return err
}

// SetExecutionDone saves the execution and updates the job with the corresponding
// results
func (s *Store) SetExecutionDone(execution *Execution) (bool, error) {
	err := s.db.Update(func(txn *badger.Txn) error {
		// Load the job from the store
		job, err := s.GetJob(execution.JobName, &JobOptions{
			ComputeStatus: true,
		})

		if err != nil {
			if err == badger.ErrKeyNotFound {
				log.Warning(ErrExecutionDoneForDeletedJob)
				return ErrExecutionDoneForDeletedJob
			}
			log.WithError(err).Fatal(err)
			return err
		}

		// Save the execution to store
		if _, err := s.SetExecution(execution); err != nil {
			return err
		}

		if execution.Success {
			job.LastSuccess = execution.FinishedAt
			job.SuccessCount++
		} else {
			job.LastError = execution.FinishedAt
			job.ErrorCount++
		}

		if err := s.SetJob(job, false); err != nil {
			log.WithError(err).Fatal("store: Error in SetExecutionDone")
			return err
		}

		return nil
	})

	return true, err
}

func (s *Store) validateJob(job *Job) error {
	if job.ParentJob == job.Name {
		return ErrSameParent
	}

	// Only validate the schedule if it doesn't have a parent
	if job.ParentJob == "" {
		if _, err := cron.Parse(job.Schedule); err != nil {
			return fmt.Errorf("%s: %s", ErrScheduleParse.Error(), err)
		}
	}

	if job.Concurrency != ConcurrencyAllow && job.Concurrency != ConcurrencyForbid && job.Concurrency != "" {
		return ErrWrongConcurrency
	}
	if err := s.validateTimeZone(job.Timezone); err != nil {
		return err
	}

	return nil
}

func (s *Store) jobHasMetadata(job *Job, metadata map[string]string) bool {
	if job == nil || job.Metadata == nil || len(job.Metadata) == 0 {
		return false
	}

	res := true
	for k, v := range metadata {
		var found bool

		if val, ok := job.Metadata[k]; ok && v == val {
			found = true
		}

		res = res && found

		if !res {
			break
		}
	}

	return res
}

// GetJobs returns all jobs
func (s *Store) GetJobs(options *JobOptions) ([]*Job, error) {
	jobs := make([]*Job, 0)

	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("jobs")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			var pbj dkronpb.Job
			if err := proto.Unmarshal(v, &pbj); err != nil {
				return err
			}
			job := NewJobFromProto(&pbj)

			job.Agent = s.agent
			if options != nil {
				if options.Metadata != nil && len(options.Metadata) > 0 && !s.jobHasMetadata(job, options.Metadata) {
					continue
				}
				if options.ComputeStatus {
					job.Status = job.GetStatus()
				}
			}

			jobs = append(jobs, job)
		}
		return nil
	})

	return jobs, err
}

// GetJob finds and return a Job from the store
func (s *Store) GetJob(name string, options *JobOptions) (*Job, error) {
	var job *Job

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("jobs/" + name))
		if err != nil {
			return err
		}

		res, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		var pbj dkronpb.Job
		if err := proto.Unmarshal(res, &pbj); err != nil {
			return err
		}
		job = NewJobFromProto(&pbj)

		log.WithFields(logrus.Fields{
			"job": job.Name,
		}).Debug("store: Retrieved job from datastore")

		job.Agent = s.agent
		if options != nil && options.ComputeStatus {
			job.Status = job.GetStatus()
		}

		return nil
	})

	return job, err
}

func (s *Store) DeleteJob(name string) (*Job, error) {
	var job *Job
	err := s.db.Update(func(txn *badger.Txn) error {
		j, err := s.GetJob(name, nil)
		if err != nil {
			return err
		}
		job = j

		if err := s.DeleteExecutions(name); err != nil {
			if err != nil {
				return err
			}
		}

		return txn.Delete([]byte("jobs/" + name))
	})

	return job, err
}

func (s *Store) GetExecutions(jobName string) ([]*Execution, error) {
	prefix := fmt.Sprintf("executions/%s", jobName)

	kvs, err := s.list(prefix, true)
	if err != nil {
		return nil, err
	}

	return s.unmarshalExecutions(kvs, jobName)
}

type kv struct {
	Key   string
	Value []byte
}

func (s *Store) list(prefix string, checkRoot bool) ([]*kv, error) {
	prefix = strings.TrimSuffix(prefix, "/")

	kvs := []*kv{}
	found := false

	err := s.db.View(func(tx *badger.Txn) error {
		it := tx.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(prefix)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			found = true
			item := it.Item()
			k := item.Key()

			// ignore self in listing
			if bytes.Equal(trimDirectoryKey(k), prefix) {
				continue
			}

			body, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			kv := &kv{Key: string(k), Value: body}
			kvs = append(kvs, kv)
		}

		return nil
	})

	if err == nil && !found && checkRoot {
		return nil, badger.ErrKeyNotFound
	}

	return kvs, err
}

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

// Returns executions for a job grouped and with an ordered index
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

// SetExecution Save a new execution and returns the key of the new saved item or an error.
func (s *Store) SetExecution(execution *Execution) (string, error) {
	pbe := execution.ToProto()
	eb, err := proto.Marshal(pbe)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("executions/%s/%s", execution.JobName, execution.Key())

	log.WithFields(logrus.Fields{
		"job":       execution.JobName,
		"execution": key,
		"finished":  execution.FinishedAt.String(),
	}).Debug("store: Setting key")

	err = s.db.Update(func(txn *badger.Txn) error {
		// Get previous execution
		i, err := txn.Get([]byte(key))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		// Do nothing if a previous execution exists and is
		// more recent, avoiding non ordered execution set
		if i != nil {
			v, err := i.ValueCopy(nil)
			if err != nil {
				return err
			}
			var p dkronpb.Execution
			if err := proto.Unmarshal(v, &p); err != nil {
				return err
			}
			// Compare existing execution
			if p.GetFinishedAt().Seconds > pbe.GetFinishedAt().Seconds {
				return nil
			}
		}
		return txn.Set([]byte(key), eb)
	})

	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"job":       execution.JobName,
			"execution": key,
		}).Debug("store: Failed to set key")
		return "", err
	}

	execs, err := s.GetExecutions(execution.JobName)
	if err != nil && err != badger.ErrKeyNotFound {
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
			err = s.db.Update(func(txn *badger.Txn) error {
				k := fmt.Sprintf("executions/%s/%s", execs[i].JobName, execs[i].Key())
				return txn.Delete([]byte(k))
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
func (s *Store) DeleteExecutions(jobName string) error {
	prefix := []byte(jobName)

	// transaction may conflict
ConflictRetry:
	for i := 0; i < defaultUpdateMaxAttempts; i++ {

		// always retry when TxnTooBig is signalled
	TxnTooBigRetry:
		for {
			txn := s.db.NewTransaction(true)
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = false

			it := txn.NewIterator(opts)

			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				k := it.Item().KeyCopy(nil)

				err := txn.Delete(k)
				it.Close()
				if err != badger.ErrTxnTooBig {
					return err
				}

				err = txn.Commit(nil)

				// commit failed with conflict
				if err == badger.ErrConflict {
					continue ConflictRetry
				}

				if err != nil {
					return err
				}

				// open new transaction and continue
				continue TxnTooBigRetry
			}

			it.Close()
			err := txn.Commit(nil)

			// commit failed with conflict
			if err == badger.ErrConflict {
				continue ConflictRetry
			}

			return err
		}
	}

	return ErrTooManyUpdateConflicts
}

func (s *Store) Shutdown() error {
	return s.db.Close()
}

func (s *Store) unmarshalExecutions(items []*kv, stopWord string) ([]*Execution, error) {
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

func trimDirectoryKey(key []byte) []byte {
	if isDirectoryKey(key) {
		return key[:len(key)-1]
	}

	return key
}

func isDirectoryKey(key []byte) bool {
	return len(key) > 0 && key[len(key)-1] == '/'
}
