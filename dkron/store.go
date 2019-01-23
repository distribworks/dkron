package dkron

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"
	"github.com/abronan/valkeyrie/store/boltdb"
	"github.com/abronan/valkeyrie/store/consul"
	"github.com/abronan/valkeyrie/store/dynamodb"
	"github.com/abronan/valkeyrie/store/etcd/v2"
	"github.com/abronan/valkeyrie/store/etcd/v3"
	"github.com/abronan/valkeyrie/store/redis"
	"github.com/abronan/valkeyrie/store/zookeeper"
	"github.com/sirupsen/logrus"
	"github.com/victorcoder/dkron/cron"
)

const MaxExecutions = 100

type Storage interface {
	SetJob(job *Job) error
	AtomicJobPut(job *Job, prevJobKVPair *store.KVPair) (bool, error)
	SetJobDependencyTree(job *Job, previousJob *Job) error
	GetJobs() ([]*Job, error)
	GetJob(name string, options *JobOptions) (*Job, error)
	GetJobWithKVPair(name string, options *JobOptions) (*Job, *store.KVPair, error)
	DeleteJob(name string) (*Job, error)
	GetExecutions(jobName string) ([]*Execution, error)
	GetLastExecutionGroup(jobName string) ([]*Execution, error)
	GetExecutionGroup(execution *Execution) ([]*Execution, error)
	GetGroupedExecutions(jobName string) (map[int64][]*Execution, []int64, error)
	SetExecution(execution *Execution) (string, error)
	DeleteExecutions(jobName string) error
	GetLeader() []byte
	LeaderKey() string
}

type Store struct {
	Client   store.Store
	agent    *Agent
	keyspace string
	backend  store.Backend
}

type JobOptions struct {
	ComputeStatus bool
	Tags          map[string]string `json:"tags"`
}

func init() {
	etcd.Register()
	etcdv3.Register()
	consul.Register()
	zookeeper.Register()
	redis.Register()
	boltdb.Register()
	dynamodb.Register()
}

func NewStore(backend store.Backend, machines []string, a *Agent, keyspace string, config *store.Config) *Store {
	s, err := valkeyrie.NewStore(store.Backend(backend), machines, config)
	if err != nil {
		log.Error(err)
	}

	log.WithFields(logrus.Fields{
		"backend":  backend,
		"machines": machines,
		"keyspace": keyspace,
	}).Debug("store: Backend config")

	return &Store{Client: s, agent: a, keyspace: keyspace, backend: backend}
}

func (s *Store) Healthy() error {
	_, err := s.Client.List(s.keyspace, nil)
	if err != store.ErrKeyNotFound && err != nil {
		return err
	}
	return nil
}

// Store a job
func (s *Store) SetJob(job *Job, copyDependentJobs bool) error {
	//Existing job that has children, let's keep it's children

	// Sanitize the job name
	job.Name = generateSlug(job.Name)
	jobKey := fmt.Sprintf("%s/jobs/%s", s.keyspace, job.Name)

	// Init the job agent
	job.Agent = s.agent

	if err := s.validateJob(job); err != nil {
		return err
	}

	// Get if the requested job already exist
	ej, err := s.GetJob(job.Name, nil)
	if err != nil && err != store.ErrKeyNotFound {
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

	jobJSON, _ := json.Marshal(job)

	log.WithFields(logrus.Fields{
		"job":  job.Name,
		"json": string(jobJSON),
	}).Debug("store: Setting job")

	if err := s.Client.Put(jobKey, jobJSON, nil); err != nil {
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
}

func (s *Store) validateTimeZone(timezone string) error {
	if timezone == "" {
		return nil
	}
	_, err := time.LoadLocation(timezone)
	return err
}

func (s *Store) AtomicJobPut(job *Job, prevJobKVPair *store.KVPair) (bool, error) {
	jobKey := fmt.Sprintf("%s/jobs/%s", s.keyspace, job.Name)
	jobJSON, _ := json.Marshal(job)

	ok, _, err := s.Client.AtomicPut(jobKey, jobJSON, prevJobKVPair, nil)

	return ok, err
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

func (s *Store) jobHasTags(job *Job, tags map[string]string) bool {
	if job == nil || job.Tags == nil || len(job.Tags) == 0 {
		return false
	}

	res := true
	for k, v := range tags {
		var found bool

		if val, ok := job.Tags[k]; ok && v == val {
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
	res, err := s.Client.List(s.keyspace+"/jobs/", nil)
	if err != nil {
		if err == store.ErrKeyNotFound {
			log.Debug("store: No jobs found")
			return []*Job{}, nil
		}
		return nil, err
	}

	jobs := make([]*Job, 0)
	for _, node := range res {
		var job Job
		err := json.Unmarshal([]byte(node.Value), &job)
		if err != nil {
			return nil, err
		}
		job.Agent = s.agent
		if options != nil {
			if options.Tags != nil && len(options.Tags) > 0 && !s.jobHasTags(&job, options.Tags) {
				continue
			}
			if options.ComputeStatus {
				job.Status = job.GetStatus()
			}
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

// Get a job
func (s *Store) GetJob(name string, options *JobOptions) (*Job, error) {
	job, _, err := s.GetJobWithKVPair(name, options)
	return job, err
}

func (s *Store) GetJobWithKVPair(name string, options *JobOptions) (*Job, *store.KVPair, error) {
	res, err := s.Client.Get(s.keyspace+"/jobs/"+name, nil)
	if err != nil {
		return nil, nil, err
	}

	var job Job
	if err = json.Unmarshal([]byte(res.Value), &job); err != nil {
		return nil, nil, err
	}

	log.WithFields(logrus.Fields{
		"job": job.Name,
	}).Debug("store: Retrieved job from datastore")

	job.Agent = s.agent
	if options != nil && options.ComputeStatus {
		job.Status = job.GetStatus()
	}

	return &job, res, nil
}

func (s *Store) DeleteJob(name string) (*Job, error) {
	job, err := s.GetJob(name, nil)
	if err != nil {
		return nil, err
	}

	if err := s.DeleteExecutions(name); err != nil {
		if err != store.ErrKeyNotFound {
			return nil, err
		}
	}

	if err := s.Client.Delete(s.keyspace + "/jobs/" + name); err != nil {
		return nil, err
	}

	return job, nil
}

func (s *Store) GetExecutions(jobName string) ([]*Execution, error) {
	prefix := fmt.Sprintf("%s/executions/%s", s.keyspace, jobName)
	res, err := s.Client.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	return s.unmarshalExecutions(res, jobName)
}

func (s *Store) GetLastExecutionGroup(jobName string) ([]*Execution, error) {
	res, err := s.Client.List(fmt.Sprintf("%s/executions/%s", s.keyspace, jobName), nil)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return []*Execution{}, nil
	}

	var lastEx Execution
	var executions []*Execution
	// res does not guarantee any order,
	// so compare them by `StartedAt` time and get the last one
	for _, node := range res {
		var ex Execution
		err := json.Unmarshal([]byte(node.Value), &ex)
		if err != nil {
			return nil, err
		}
		if ex.StartedAt.After(lastEx.StartedAt) {
			lastEx = ex
			executions = []*Execution{&ex}
		} else if ex.Group == lastEx.Group {
			executions = append(executions, &ex)
		}
	}
	return executions, nil
}

func (s *Store) GetExecutionGroup(execution *Execution) ([]*Execution, error) {
	res, err := s.Client.List(fmt.Sprintf("%s/executions/%s", s.keyspace, execution.JobName), nil)
	if err != nil {
		return nil, err
	}

	var executions []*Execution
	for _, node := range res {
		var ex Execution
		err := json.Unmarshal([]byte(node.Value), &ex)
		if err != nil {
			return nil, err
		}

		if ex.Group == execution.Group {
			executions = append(executions, &ex)
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

// Save a new execution and returns the key of the new saved item or an error.
func (s *Store) SetExecution(execution *Execution) (string, error) {
	exJson, _ := json.Marshal(execution)
	key := execution.Key()

	log.WithFields(logrus.Fields{
		"job":       execution.JobName,
		"execution": key,
	}).Debug("store: Setting key")

	err := s.Client.Put(fmt.Sprintf("%s/executions/%s/%s", s.keyspace, execution.JobName, key), exJson, nil)
	if err != nil {
		log.WithFields(logrus.Fields{
			"job":       execution.JobName,
			"execution": key,
			"error":     err,
		}).Debug("store: Failed to set key")
		return "", err
	}

	execs, err := s.GetExecutions(execution.JobName)
	if err != nil {
		log.WithError(err).
			WithField("job", execution.JobName).
			Error("store: Error no executions found for job")
	}

	// Delete all execution results over the limit, starting from olders
	if len(execs) > MaxExecutions {
		//sort the array of all execution groups by StartedAt time
		sort.Sort(ExecList(execs))
		for i := 0; i < len(execs)-MaxExecutions; i++ {
			log.WithFields(logrus.Fields{
				"job":       execs[i].JobName,
				"execution": execs[i].Key(),
			}).Debug("store: to detele key")
			err := s.Client.Delete(fmt.Sprintf("%s/executions/%s/%s", s.keyspace, execs[i].JobName, execs[i].Key()))
			if err != nil {
				log.WithError(err).
					WithField("execution", execs[i].Key()).
					Error("store: Error trying to delete overflowed execution")
			}
		}
	}

	return key, nil
}

func (s *Store) unmarshalExecutions(res []*store.KVPair, stopWord string) ([]*Execution, error) {
	var executions []*Execution
	for _, node := range res {
		if store.Backend(s.backend) != store.ZK {
			path := store.SplitKey(node.Key)
			dir := path[len(path)-2]
			if dir != stopWord {
				continue
			}
		}
		var execution Execution
		err := json.Unmarshal([]byte(node.Value), &execution)
		if err != nil {
			return nil, err
		}
		executions = append(executions, &execution)
	}
	return executions, nil
}

// Removes all executions of a job
func (s *Store) DeleteExecutions(jobName string) error {
	return s.Client.DeleteTree(fmt.Sprintf("%s/executions/%s", s.keyspace, jobName))
}

// Retrieve the leader from the store
func (s *Store) GetLeader() []byte {
	res, err := s.Client.Get(s.LeaderKey(), nil)
	if err != nil {
		if err == store.ErrNotReachable {
			log.Fatal("store: Store not reachable, be sure you have an existing key-value store running is running and is reachable.")
		} else if err != store.ErrKeyNotFound {
			log.Error(err)
		}
		return nil
	}

	log.WithField("node", string(res.Value)).Debug("store: Retrieved leader from datastore")

	return res.Value
}

// Retrieve the leader key used in the KV store to store the leader node
func (s *Store) LeaderKey() string {
	return s.keyspace + "/leader"
}
