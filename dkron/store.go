package dkron

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
	"github.com/victorcoder/dkron/cron"
)

const MaxExecutions = 100

type Store struct {
	Client   store.Store
	agent    *AgentCommand
	keyspace string
	backend  string
}

func init() {
	etcd.Register()
	consul.Register()
	zookeeper.Register()
}

func NewStore(backend string, machines []string, a *AgentCommand, keyspace string) *Store {
	s, err := libkv.NewStore(store.Backend(backend), machines, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(logrus.Fields{
		"backend":  backend,
		"machines": machines,
		"keyspace": keyspace,
	}).Debug("store: Backend config")

	_, err = s.List(keyspace)
	if err != store.ErrKeyNotFound && err != nil {
		log.WithError(err).Fatal("store: Store backend not reachable")
	}

	return &Store{Client: s, agent: a, keyspace: keyspace, backend: backend}
}

// Store a job
func (s *Store) SetJob(job *Job) error {
	// Sanitize the job name
	job.Name = generateSlug(job.Name)
	jobKey := fmt.Sprintf("%s/jobs/%s", s.keyspace, job.Name)

	// Init the job agent
	job.Agent = s.agent

	if err := s.validateJob(job); err != nil {
		return err
	}

	// Get if the requested job already exist
	ej, err := s.GetJob(job.Name)
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
	}

	jobJSON, _ := json.Marshal(job)

	log.WithFields(logrus.Fields{
		"job":  job.Name,
		"json": string(jobJSON),
	}).Debug("store: Setting job")

	if err := s.Client.Put(jobKey, jobJSON, nil); err != nil {
		return err
	}

	return nil
}

// Set the depencency tree for a job given the job and the previous version
// of the Job or nil if it's new.
func (s *Store) SetJobDependencyTree(job *Job, previousJob *Job) error {
	// Existing job that doesn't have parent job set and it's being set
	if previousJob != nil && previousJob.ParentJob == "" && job.ParentJob != "" {
		pj, err := job.GetParent()
		if err != nil {
			return err
		}
		pj.Lock()
		defer pj.Unlock()

		pj.DependentJobs = append(pj.DependentJobs, job.Name)
		if err := s.SetJob(pj); err != nil {
			return err
		}
	}

	// Existing job that has parent job set and it's being removed
	if previousJob != nil && previousJob.ParentJob != "" && job.ParentJob == "" {
		pj, err := previousJob.GetParent()
		if err != nil {
			return err
		}
		pj.Lock()
		defer pj.Unlock()

		ndx := 0
		for i, djn := range pj.DependentJobs {
			if djn == job.Name {
				ndx = i
				break
			}
		}
		pj.DependentJobs = append(pj.DependentJobs[:ndx], pj.DependentJobs[ndx+1:]...)
		if err := s.SetJob(pj); err != nil {
			return err
		}
	}

	// New job that has parent job set
	if previousJob == nil && job.ParentJob != "" {
		pj, err := job.GetParent()
		if err != nil {
			return err
		}
		pj.Lock()
		defer pj.Unlock()

		pj.DependentJobs = append(pj.DependentJobs, job.Name)
		if err := s.SetJob(pj); err != nil {
			return err
		}
	}

	return nil
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

	if job.Command == "" {
		return ErrNoCommand
	}

	if job.Concurrency != ConcurrencyAllow && job.Concurrency != ConcurrencyForbid && job.Concurrency != "" {
		return ErrWrongConcurrency
	}

	return nil
}

// GetJobs returns all jobs
func (s *Store) GetJobs() ([]*Job, error) {
	res, err := s.Client.List(s.keyspace + "/jobs/")
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
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

// Get a job
func (s *Store) GetJob(name string) (*Job, error) {
	res, err := s.Client.Get(s.keyspace + "/jobs/" + name)
	if err != nil {
		return nil, err
	}

	var job Job
	if err = json.Unmarshal([]byte(res.Value), &job); err != nil {
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"job": job.Name,
	}).Debug("store: Retrieved job from datastore")

	job.Agent = s.agent
	return &job, nil
}

func (s *Store) DeleteJob(name string) (*Job, error) {
	job, err := s.GetJob(name)
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
	res, err := s.Client.List(prefix)
	if err != nil {
		return nil, err
	}

	var executions []*Execution

	for _, node := range res {
		if store.Backend(s.backend) != store.ZK {
			path := store.SplitKey(node.Key)
			dir := path[len(path)-2]
			if dir != jobName {
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

func (s *Store) GetLastExecutionGroup(jobName string) ([]*Execution, error) {
	res, err := s.Client.List(fmt.Sprintf("%s/executions/%s", s.keyspace, jobName))
	if err != nil {
		return nil, err
	}

	var ex Execution
	err = json.Unmarshal([]byte(res[len(res)-1].Value), &ex)
	if err != nil {
		return nil, err
	}
	return s.GetExecutionGroup(&ex)
}

func (s *Store) GetExecutionGroup(execution *Execution) ([]*Execution, error) {
	res, err := s.Client.List(fmt.Sprintf("%s/executions/%s", s.keyspace, execution.JobName))
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
	sort.Sort(byGroup)

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
		return "", err
	}

	execs, err := s.GetExecutions(execution.JobName)
	if err != nil {
		log.Errorf("store: No executions found for job %s", execution.JobName)
	}

	// Get and ordered array of all execution groups
	var byGroup int64arr
	for _, ex := range execs {
		byGroup = append(byGroup, ex.Group)
	}
	sort.Sort(byGroup)

	// Delete all execution results over the limit, starting from olders
	if len(byGroup) > MaxExecutions {
		for i := range byGroup[MaxExecutions:] {
			err := s.Client.Delete(fmt.Sprintf("%s/executions/%s/%s", s.keyspace, execs[i].JobName, execs[i].Key()))
			if err != nil {
				log.Errorf("store: Trying to delete overflowed execution %s", execs[i].Key())
			}
		}
	}

	return key, nil
}

// Removes all executions of a job
func (s *Store) DeleteExecutions(jobName string) error {
	return s.Client.DeleteTree(fmt.Sprintf("%s/executions/%s", s.keyspace, jobName))
}

// Retrieve the leader from the store
func (s *Store) GetLeader() []byte {
	res, err := s.Client.Get(s.LeaderKey())
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
