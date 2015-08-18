package dkron

import (
	"encoding/json"
	"fmt"
	stdlog "log"
	"time"

	etcdc "github.com/coreos/go-etcd/etcd"
)

const keyspace = "/dkron"

type etcdClient struct {
	Client *etcdc.Client
	agent  *AgentCommand
}

// ServerStats encapsulates various statistics about an EtcdServer and its
// communication with other members of the cluster
type EtcdServerStats struct {
	Name string `json:"name"`
	// TODO(jonboulle): use ID instead of name?
	ID        string    `json:"id"`
	StartTime time.Time `json:"startTime"`

	LeaderInfo struct {
		Name      string    `json:"leader"`
		Uptime    string    `json:"uptime"`
		StartTime time.Time `json:"startTime"`
	} `json:"leaderInfo"`

	RecvAppendRequestCnt uint64  `json:"recvAppendRequestCnt,"`
	RecvingPkgRate       float64 `json:"recvPkgRate,omitempty"`
	RecvingBandwidthRate float64 `json:"recvBandwidthRate,omitempty"`

	SendAppendRequestCnt uint64  `json:"sendAppendRequestCnt"`
	SendingPkgRate       float64 `json:"sendPkgRate,omitempty"`
	SendingBandwidthRate float64 `json:"sendBandwidthRate,omitempty"`
}

func init() {
	etcdc.SetLogger(stdlog.New(log.Writer(), "go-etcd", stdlog.LstdFlags))
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
		eerr := err.(*etcdc.EtcdError)
		if eerr.ErrorCode == 100 {
			log.Info("No jobs found")
			return nil, nil
		}
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

func (e *etcdClient) GetExecutions(jobName string) ([]*Execution, error) {
	res, err := e.Client.Get(fmt.Sprintf("%s/executions/%s", keyspace, jobName), true, false)
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

// Save a new execution and returns the key of the new saved item or an error.
func (e *etcdClient) SetExecution(execution *Execution) (string, error) {
	eJson, _ := json.Marshal(execution)
	key := fmt.Sprintf("%d-%s", execution.StartedAt.UnixNano(), execution.NodeName)

	log.Debugf("Setting etcd key %s: %s", execution.JobName, string(eJson))
	res, err := e.Client.Set(fmt.Sprintf("%s/executions/%s/%s", keyspace, execution.JobName, key), string(eJson), 0)
	if err != nil {
		return "", err
	}

	return res.Node.Key, nil
}

func (e *etcdClient) GetLeader() string {
	res, err := e.Client.Get(keyspace+"/leader", false, false)
	if err != nil {
		if eerr, ok := err.(*etcdc.EtcdError); ok {
			if eerr.ErrorCode == etcdc.ErrCodeEtcdNotReachable {
				log.Fatal("etcd not reachable, be sure etcd is running.\nYou can download etc from https://github.com/coreos/etcd/releases")
			}
		}
		log.Error(err.Error())
		return ""
	}

	log.Debugf("Retrieved leader from datastore: %v", res.Node.Value)
	return res.Node.Value
}
