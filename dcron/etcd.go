package dcron

import (
	"encoding/json"
	etcdc "github.com/coreos/go-etcd/etcd"
)

var machines = []string{"http://127.0.0.1:2379"}
var etcd = NewClient(machines)
var keyspace = "/dcron"

type etcdClient struct {
	Client *etcdc.Client
}

func NewClient(machines []string) *etcdClient {
	return &etcdClient{Client: etcdc.NewClient(machines)}
}

func (e *etcdClient) SetJob(job *Job) error {
	jobJson, _ := json.Marshal(job)
	if _, err := e.Client.Set(keyspace+"/jobs/"+job.Name+"/job", string(jobJson), 0); err != nil {
		return err
	}

	return nil
}
