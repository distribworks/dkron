package dcron

import (
	"github.com/coreos/go-etcd/etcd"
)

var machines = []string{"http://127.0.0.1:2379"}
var etcdClient = etcd.NewClient(machines)
