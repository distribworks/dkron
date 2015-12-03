package dkron

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	"github.com/docker/libkv/store/etcd"
	"github.com/docker/libkv/store/zookeeper"
)

type Leader struct {
	Key       []byte
	LastIndex uint64
}

type Store struct {
	Client   store.Store
	agent    *AgentCommand
	keyspace string
}

func init() {
	etcd.Register()
	consul.Register()
	zookeeper.Register()
}

func NewStore(backend string, machines []string, a *AgentCommand, keyspace string) *Store {
	store, err := libkv.NewStore(store.Backend(backend), machines, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &Store{Client: store, agent: a, keyspace: keyspace}
}

// Store a job
func (s *Store) SetJob(job *Job) error {
	jobJson, _ := json.Marshal(job)

	log.WithFields(logrus.Fields{
		"job":  job.Name,
		"json": string(jobJson),
	}).Debug("store: Setting job")

	if err := s.Client.Put(s.keyspace+"/jobs/"+job.Name, jobJson, nil); err != nil {
		return err
	}

	return nil
}

// Get all jobs
func (s *Store) GetJobs() ([]*Job, error) {
	res, err := s.Client.List(s.keyspace + "/jobs/")
	if err != nil {
		if err == store.ErrKeyNotFound {
			log.Debug("store: No jobs found")
			return nil, nil
		}
		return nil, err
	}

	var jobs []*Job
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

	if err := s.Client.Delete(s.keyspace + "/jobs/" + name); err != nil {
		return nil, err
	}

	return job, nil
}

func (s *Store) GetExecutions(jobName string) ([]*Execution, error) {
	res, err := s.Client.List(fmt.Sprintf("%s/executions/%s", s.keyspace, jobName))
	if err != nil {
		return nil, err
	}

	var executions []*Execution

	for _, node := range res {
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

// Save a new execution and returns the key of the new saved item or an error.
func (s *Store) SetExecution(execution *Execution) (string, error) {
	exJson, _ := json.Marshal(execution)
	key := fmt.Sprintf("%d-%s", execution.StartedAt.UnixNano(), execution.NodeName)

	log.WithFields(logrus.Fields{
		"job":       execution.JobName,
		"execution": key,
	}).Debug("store: Setting key")

	err := s.Client.Put(fmt.Sprintf("%s/executions/%s/%s", s.keyspace, execution.JobName, key), exJson, nil)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *Store) GetLeader() *Leader {
	res, err := s.Client.Get(s.keyspace + "/leader")
	if err != nil {
		if err == store.ErrNotReachable {
			log.Fatal("Store not reachable, be sure you have an existing key-value store running is running and is reachable.")
		} else if err != store.ErrKeyNotFound {
			log.Error(err)
		}
		return nil
	}

	log.WithFields(logrus.Fields{
		"key": string(res.Value),
	}).Debug("store: Retrieved leader from datastore")

	return &Leader{Key: res.Value, LastIndex: res.LastIndex}
}

func (s *Store) TryLeaderSwap(newKey string, old *Leader) (bool, error) {
	oldKV := &store.KVPair{
		LastIndex: old.LastIndex,
	}
	success, _, err := s.Client.AtomicPut(s.keyspace+"/leader", []byte(newKey), oldKV, nil)

	log.WithFields(logrus.Fields{
		"old_leader": string(old.Key),
		"new_leader": newKey,
	}).Debug("store: Leader Swap")

	return success, err
}

func (s *Store) SetLeader(leader string) error {
	err := s.Client.Put(s.keyspace+"/leader", []byte(leader), nil)
	if err != nil {
		return err
	}

	return nil
}
