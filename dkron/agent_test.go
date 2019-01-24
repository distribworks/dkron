package dkron

import (
	"os"
	"testing"
	"time"

	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"
	"github.com/hashicorp/serf/testutil"
	"github.com/stretchr/testify/assert"
)

var (
	logLevel = "error"
	etcdAddr = getEnvWithDefault()
)

func getEnvWithDefault() string {
	ea := os.Getenv("DKRON_BACKEND_MACHINE")
	if ea == "" {
		return "127.0.0.1:2379"
	}
	return ea
}

func TestAgentCommand_runForElection(t *testing.T) {
	a1Name := "test1"
	a2Name := "test2"
	a1Addr := testutil.GetBindAddr().String()
	a2Addr := testutil.GetBindAddr().String()
	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	// Override leader TTL
	defaultLeaderTTL = 2 * time.Second

	client, err := valkeyrie.NewStore("etcdv3", []string{etcdAddr}, &store.Config{})
	if err != nil {
		panic(err)
	}
	err = client.DeleteTree("dkron")
	if err != nil {
		if err != store.ErrKeyNotFound {
			panic(err)
		}
	}

	c := DefaultConfig()
	c.BindAddr = a1Addr
	c.StartJoin = []string{a2Addr}
	c.NodeName = a1Name
	c.Server = true
	c.LogLevel = logLevel
	c.Backend = "etcdv3"
	c.BackendMachines = []string{os.Getenv("DKRON_BACKEND_MACHINE")}

	a1 := NewAgent(c, nil)
	if err := a1.Start(); err != nil {
		t.Fatal(err)
	}

	// Wait for the first agent to start and set itself as leader
	kv1, err := watchOrDie(client, "dkron/leader")
	if err != nil {
		t.Fatal(err)
	}
	leaderA1 := string(kv1.Value)
	t.Logf("%s is the current leader", leaderA1)
	assert.Equal(t, a1Name, leaderA1)

	// Start another agent
	c = DefaultConfig()
	c.BindAddr = a2Addr
	c.StartJoin = []string{a1Addr + ":8946"}
	c.NodeName = a2Name
	c.Server = true
	c.LogLevel = logLevel
	c.Backend = "etcdv3"
	c.BackendMachines = []string{os.Getenv("DKRON_BACKEND_MACHINE")}

	a2 := NewAgent(c, nil)
	a2.Start()

	// Send a shutdown request
	a1.Stop()

	// Wait until test2 steps as leader
rewatch:
	kv2, err := watchOrDie(client, "dkron/leader")
	if err != nil {
		t.Fatal(err)
	}
	if len(kv2.Value) == 0 || string(kv2.Value) == a1Name {
		goto rewatch
	}
	t.Logf("%s is the current leader", kv2.Value)
	assert.Equal(t, a2Name, string(kv2.Value))
	a2.Stop()
}

func watchOrDie(client store.Store, key string) (*store.KVPair, error) {
	for {
		resultCh, err := client.Watch(key, nil, nil)
		if err != nil {
			if err == store.ErrKeyNotFound {
				continue
			}
			return nil, err
		}

		// If the channel worked, read the actual value
		kv := <-resultCh
		return kv, nil
	}
}

func Test_processFilteredNodes(t *testing.T) {
	client, err := valkeyrie.NewStore("etcdv3", []string{etcdAddr}, &store.Config{})
	err = client.DeleteTree("dkron")
	if err != nil {
		if err == store.ErrNotReachable {
			t.Fatal("etcd server needed to run tests")
		}
	}

	a1Addr := testutil.GetBindAddr().String()
	a2Addr := testutil.GetBindAddr().String()

	c := DefaultConfig()
	c.BindAddr = a1Addr
	c.StartJoin = []string{a2Addr}
	c.NodeName = "test1"
	c.Server = true
	c.LogLevel = logLevel
	c.Tags = map[string]string{"role": "test"}
	c.Backend = "etcdv3"
	c.BackendMachines = []string{os.Getenv("DKRON_BACKEND_MACHINE")}

	a1 := NewAgent(c, nil)
	a1.Start()

	time.Sleep(2 * time.Second)

	// Start another agent
	c = DefaultConfig()
	c.BindAddr = a2Addr
	c.StartJoin = []string{a1Addr + ":8946"}
	c.NodeName = "test2"
	c.Server = true
	c.LogLevel = logLevel
	c.Tags = map[string]string{"role": "test"}
	c.Backend = "etcdv3"
	c.BackendMachines = []string{os.Getenv("DKRON_BACKEND_MACHINE")}

	a2 := NewAgent(c, nil)
	a2.Start()

	time.Sleep(2 * time.Second)

	job := &Job{
		Name: "test_job_1",
		Tags: map[string]string{
			"foo":  "bar:1",
			"role": "test:2",
		},
	}

	nodes, tags, err := a1.processFilteredNodes(job)

	assert.Contains(t, nodes, "test1")
	assert.Contains(t, nodes, "test2")
	assert.Equal(t, tags["role"], "test")
	assert.Equal(t, job.Tags["role"], "test:2")
	assert.Equal(t, job.Tags["foo"], "bar:1")

	a1.Stop()
	a2.Stop()
}

func TestEncrypt(t *testing.T) {
	c := DefaultConfig()
	c.BindAddr = testutil.GetBindAddr().String()
	c.NodeName = "test1"
	c.Server = true
	c.Tags = map[string]string{"role": "test"}
	c.EncryptKey = "kPpdjphiipNSsjd4QHWbkA=="
	c.LogLevel = logLevel
	c.Backend = "etcdv3"
	c.BackendMachines = []string{os.Getenv("DKRON_BACKEND_MACHINE")}

	a := NewAgent(c, nil)
	a.Start()

	time.Sleep(2 * time.Second)

	assert.True(t, a.serf.EncryptionEnabled())
	a.Stop()
}

func Test_getRPCAddr(t *testing.T) {
	a1Addr := testutil.GetBindAddr()

	c := DefaultConfig()
	c.BindAddr = a1Addr.String() + ":5000"
	c.NodeName = "test1"
	c.Server = true
	c.Tags = map[string]string{"role": "test"}
	c.LogLevel = logLevel
	c.Backend = "etcdv3"
	c.BackendMachines = []string{os.Getenv("DKRON_BACKEND_MACHINE")}

	a := NewAgent(c, nil)
	a.Start()

	time.Sleep(2 * time.Second)

	getRPCAddr := a.getRPCAddr()
	exRPCAddr := a1Addr.String() + ":6868"

	assert.Equal(t, exRPCAddr, getRPCAddr)
	a.Stop()
}

func TestAgentConfig(t *testing.T) {
	advAddr := testutil.GetBindAddr().String()

	c := DefaultConfig()
	c.BindAddr = testutil.GetBindAddr().String()
	c.AdvertiseAddr = advAddr
	c.LogLevel = logLevel

	a := NewAgent(c, nil)
	a.Start()

	time.Sleep(2 * time.Second)

	assert.NotEqual(t, a.config.AdvertiseAddr, a.config.BindAddr)
	assert.NotEmpty(t, a.config.AdvertiseAddr)
	assert.Equal(t, advAddr, a.config.AdvertiseAddr)

	a.Stop()
}
