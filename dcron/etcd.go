package dcron

import (
	"encoding/json"
	etcdc "github.com/coreos/go-etcd/etcd"
)

const keyspace = "/dcron"

type etcdClient struct {
	Client *etcdc.Client
	agent  *AgentCommand
}

func NewEtcdClient(machines []string, a *AgentCommand) *etcdClient {
	return &etcdClient{Client: etcdc.NewClient(machines), agent: a}
}

func (e *etcdClient) SetJob(job *Job) error {
	jobJson, _ := json.Marshal(job)
	log.Debugf("Setting etcd key %s: %s", job.Name, string(jobJson))
	if _, err := e.Client.Set(keyspace+"/jobs/"+job.Name, string(jobJson), 0); err != nil {
		return err
	}

	return nil
}

func (e *etcdClient) GetJobs() ([]*Job, error) {
	res, err := e.Client.Get(keyspace+"/jobs/", true, false)
	if err != nil {
		return nil, err
	}

	var jobs []*Job
	for _, node := range res.Node.Nodes {
		log.Debug(*node)
		var job Job
		err := json.Unmarshal([]byte(node.Value), &job)
		if err != nil {
			return nil, err
		}
		job.Agent = e.agent
		jobs = append(jobs, &job)
		log.Debug(job)
	}
	return jobs, nil
}

func (e *etcdClient) GetJob(name string) (*Job, error) {
	res, err := e.Client.Get(keyspace+"/jobs/"+name, false, false)
	if err != nil {
		return nil, err
	}

	var job Job
	if err = json.Unmarshal([]byte(res.Node.Value), &job); err != nil {
		return nil, err
	}
	log.Debugf("Retrieved job from datastore: %v", job)
	job.Agent = e.agent
	return &job, nil
}

func (e *etcdClient) GetExecutions() ([]*Execution, error) {
	res, err := e.Client.Get(keyspace+"/executions/", true, false)
	if err != nil {
		return nil, err
	}

	var executions []*Execution
	for _, node := range res.Node.Nodes {
		var execution Execution
		err := json.Unmarshal([]byte(node.Value), &execution)
		if err != nil {
			return nil, err
		}
		executions = append(executions, &execution)
	}
	return executions, nil
}

func (e *etcdClient) GetLeader() string {
	res, err := e.Client.Get(keyspace+"/leader", false, false)
	if err != nil {
		if eerr, ok := err.(*etcdc.EtcdError); ok {
			if eerr.ErrorCode == etcdc.ErrCodeEtcdNotReachable {
				log.Panic(err)
			}
		}
		log.Error(err.Error())
		return ""
	}

	log.Debugf("Retrieved leader from datastore: %v", res.Node.Value)
	return res.Node.Value
}
