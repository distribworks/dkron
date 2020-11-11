package dkron

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/serf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	logLevel = "error"
)

func TestAgentCommand_runForElection(t *testing.T) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	a1Name := "test1"
	a2Name := "test2"
	ip1, returnFn1 := testutil.TakeIP()
	a1Addr := ip1.String()
	defer returnFn1()
	ip2, returnFn2 := testutil.TakeIP()
	a2Addr := ip2.String()
	defer returnFn2()

	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	c := DefaultConfig()
	c.BindAddr = a1Addr
	c.StartJoin = []string{a2Addr}
	c.NodeName = a1Name
	c.Server = true
	c.LogLevel = logLevel
	c.BootstrapExpect = 3
	c.DevMode = true
	c.DataDir = dir

	a1 := NewAgent(c)
	if err := a1.Start(); err != nil {
		t.Fatal(err)
	}

	// Wait for the first agent to start and elect itself as leader
	if a1.IsLeader() {
		m, err := a1.leaderMember()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s is the current leader", m.Name)
		assert.Equal(t, a1Name, m.Name)
	}

	// Start another agent
	c = DefaultConfig()
	c.BindAddr = a2Addr
	c.StartJoin = []string{a1Addr + ":8946"}
	c.NodeName = a2Name
	c.Server = true
	c.LogLevel = logLevel
	c.BootstrapExpect = 3
	c.DevMode = true
	c.DataDir = dir

	a2 := NewAgent(c)
	a2.Start()

	// Start another agent
	c = DefaultConfig()
	ip3, returnFn3 := testutil.TakeIP()
	defer returnFn3()
	c.BindAddr = ip3.String()
	c.StartJoin = []string{a1Addr + ":8946"}
	c.NodeName = "test3"
	c.Server = true
	c.LogLevel = logLevel
	c.BootstrapExpect = 3
	c.DevMode = true
	c.DataDir = dir

	a3 := NewAgent(c)
	a3.Start()

	time.Sleep(2 * time.Second)

	// Send a shutdown request
	a1.Stop()

	// Wait until a follower steps as leader
	time.Sleep(2 * time.Second)
	assert.True(t, (a2.IsLeader() || a3.IsLeader()))
	log.Info(a3.IsLeader())

	a2.Stop()
	a3.Stop()
}

func Test_processFilteredNodes(t *testing.T) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()
	a1Addr := ip1.String()

	ip2, returnFn2 := testutil.TakeIP()
	defer returnFn2()
	a2Addr := ip2.String()

	c := DefaultConfig()
	c.BindAddr = a1Addr
	c.StartJoin = []string{a2Addr}
	c.NodeName = "test1"
	c.Server = true
	c.LogLevel = logLevel
	c.Tags = map[string]string{
		"tag":         "test",
		"region":      "global",
		"additional":  "value",
		"additional2": "value2",
	}
	c.DevMode = true
	c.DataDir = dir

	a1 := NewAgent(c)
	a1.Start()

	time.Sleep(2 * time.Second)

	// Start another agent
	c = DefaultConfig()
	c.BindAddr = a2Addr
	c.StartJoin = []string{a1Addr + ":8946"}
	c.NodeName = "test2"
	c.Server = true
	c.LogLevel = logLevel
	c.Tags = map[string]string{
		"tag":         "test",
		"extra":       "tag",
		"region":      "global",
		"additional":  "value",
		"additional2": "value2",
	}
	c.DevMode = true
	c.DataDir = dir

	a2 := NewAgent(c)
	a2.Start()

	// Start another agent
	ip3, returnFn3 := testutil.TakeIP()
	defer returnFn3()
	a3Addr := ip3.String()

	c = DefaultConfig()
	c.BindAddr = a3Addr
	c.StartJoin = []string{a1Addr + ":8946"}
	c.NodeName = "test3"
	c.Server = false
	c.LogLevel = logLevel
	c.Tags = map[string]string{
		"tag":         "test_client",
		"extra":       "tag",
		"region":      "global",
		"additional":  "value",
		"additional2": "value2",
	}
	c.DevMode = true
	c.DataDir = dir

	a3 := NewAgent(c)
	a3.Start()

	time.Sleep(2 * time.Second)

	// Test cardinality of 2 returns correct nodes
	job := &Job{
		Name: "test_job_1",
		Tags: map[string]string{
			"tag": "test:2",
		},
	}

	nodes, tags, err := a1.processFilteredNodes(job)
	require.NoError(t, err)

	assert.Contains(t, nodes, "test1")
	assert.Contains(t, nodes, "test2")
	assert.Len(t, nodes, 2)
	assert.Equal(t, tags["tag"], "test")

	// Test cardinality of 1 with two qualified nodes returns 1 node
	job2 := &Job{
		Name: "test_job_2",
		Tags: map[string]string{
			"tag": "test:1",
		},
	}

	nodes, _, err = a1.processFilteredNodes(job2)
	require.NoError(t, err)

	assert.Len(t, nodes, 1)

	// Test no cardinality specified, all nodes returned
	job3 := &Job{
		Name: "test_job_3",
	}

	nodes, _, err = a1.processFilteredNodes(job3)
	require.NoError(t, err)

	assert.Len(t, nodes, 3)
	assert.Contains(t, nodes, "test1")
	assert.Contains(t, nodes, "test2")
	assert.Contains(t, nodes, "test3")

	// Test exclusive tag returns correct node
	job4 := &Job{
		Name: "test_job_4",
		Tags: map[string]string{
			"tag": "test_client:1",
		},
	}

	nodes, _, err = a1.processFilteredNodes(job4)
	require.NoError(t, err)

	assert.Len(t, nodes, 1)
	assert.Contains(t, nodes, "test3")

	// Test existing tag but no matching value returns no nodes
	job5 := &Job{
		Name: "test_job_5",
		Tags: map[string]string{
			"tag": "no_tag",
		},
	}

	nodes, _, err = a1.processFilteredNodes(job5)
	require.NoError(t, err)

	assert.Len(t, nodes, 0)

	// Test 1 matching and 1 not matching tag returns no nodes
	job6 := &Job{
		Name: "test_job_6",
		Tags: map[string]string{
			"foo": "bar:1",
			"tag": "test:2",
		},
	}

	nodes, tags, err = a1.processFilteredNodes(job6)
	require.NoError(t, err)

	assert.Len(t, nodes, 0)
	assert.Equal(t, tags["tag"], "test")

	// Test matching tags with cardinality of 2 but only 1 matching node returns correct node
	job7 := &Job{
		Name: "test_job_7",
		Tags: map[string]string{
			"tag":   "test:2",
			"extra": "tag:2",
		},
	}

	nodes, tags, err = a1.processFilteredNodes(job7)
	require.NoError(t, err)

	assert.Contains(t, nodes, "test2")
	assert.Len(t, nodes, 1)
	assert.Equal(t, tags["tag"], "test")
	assert.Equal(t, tags["extra"], "tag")

	// Test two tags matching same 3 servers and cardinality of 1 should always return 1 server

	// Do this multiple times: an old bug caused this to sometimes succeed and
	// sometimes fail (=return no nodes at all) due to the use of math.rand
	// Statistically, about 33% should succeed and the rest should fail if
	// the code is buggy.
	// Another bug caused one node to be favored over the others. With a
	// large enough number of attempts, each node should be chosen about 1/3
	// of the time.
	job8 := &Job{
		Name: "test_job_8",
		Tags: map[string]string{
			"additional":  "value:1",
			"additional2": "value2:1",
		},
	}
	distrib := make(map[string]int)
	var sampleSize = 1000
	for i := 0; i < sampleSize; i++ {
		nodes, tags, err = a1.processFilteredNodes(job8)
		require.NoError(t, err)

		assert.Len(t, nodes, 1)
		assert.Equal(t, tags["additional"], "value")
		assert.Equal(t, tags["additional2"], "value2")
		for name := range nodes {
			distrib[name] = distrib[name] + 1
		}
	}

	// Each node must have been chosen between 30% and 36% of the time,
	// for the distribution to be considered equal.
	// Note: This test should almost never, but still can, fail even if the
	// code is fine. To fix this, the randomizer ought to be mocked.
	for name, count := range distrib {
		fmt.Println(name, float64(count)/float64(sampleSize)*100.0, "%")
	}
	assert.Greater(t, float64(distrib["test1"])/float64(sampleSize), 0.3)
	assert.Less(t, float64(distrib["test1"])/float64(sampleSize), 0.36)
	assert.Greater(t, float64(distrib["test2"])/float64(sampleSize), 0.3)
	assert.Less(t, float64(distrib["test2"])/float64(sampleSize), 0.36)
	assert.Greater(t, float64(distrib["test3"])/float64(sampleSize), 0.3)
	assert.Less(t, float64(distrib["test3"])/float64(sampleSize), 0.36)

	// Clean up
	a1.Stop()
	a2.Stop()
	a3.Stop()
}

func TestEncrypt(t *testing.T) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()

	c := DefaultConfig()
	c.BindAddr = ip1.String()
	c.NodeName = "test1"
	c.Server = true
	c.Tags = map[string]string{"role": "test"}
	c.EncryptKey = "kPpdjphiipNSsjd4QHWbkA=="
	c.LogLevel = logLevel
	c.DevMode = true
	c.DataDir = dir

	a := NewAgent(c)
	a.Start()

	time.Sleep(2 * time.Second)

	assert.True(t, a.serf.EncryptionEnabled())
	a.Stop()
}

func Test_advertiseRPCAddr(t *testing.T) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()
	a1Addr := ip1.String()

	c := DefaultConfig()
	c.BindAddr = a1Addr + ":5000"
	c.AdvertiseAddr = "8.8.8.8"
	c.NodeName = "test1"
	c.Server = true
	c.Tags = map[string]string{"role": "test"}
	c.LogLevel = logLevel
	c.DevMode = true
	c.DataDir = dir

	a := NewAgent(c)
	a.Start()

	time.Sleep(2 * time.Second)

	advertiseRPCAddr := a.advertiseRPCAddr()
	exRPCAddr := "8.8.8.8:6868"

	assert.Equal(t, exRPCAddr, advertiseRPCAddr)

	a.Stop()
}

func Test_bindRPCAddr(t *testing.T) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()
	a1Addr := ip1.String()

	c := DefaultConfig()
	c.BindAddr = a1Addr + ":5000"
	c.NodeName = "test1"
	c.Server = true
	c.Tags = map[string]string{"role": "test"}
	c.LogLevel = logLevel
	c.DevMode = true
	c.DataDir = dir

	a := NewAgent(c)
	a.Start()

	time.Sleep(2 * time.Second)

	bindRPCAddr := a.bindRPCAddr()
	exRPCAddr := a1Addr + ":6868"

	assert.Equal(t, exRPCAddr, bindRPCAddr)
	a.Stop()
}

func TestAgentConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "dkron-test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()
	advAddr := ip1.String()

	ip2, returnFn2 := testutil.TakeIP()
	defer returnFn2()

	c := DefaultConfig()
	c.BindAddr = ip2.String()
	c.AdvertiseAddr = advAddr
	c.LogLevel = logLevel
	c.DataDir = dir
	c.DevMode = true

	a := NewAgent(c)
	a.Start()

	time.Sleep(2 * time.Second)

	assert.NotEqual(t, a.config.AdvertiseAddr, a.config.BindAddr)
	assert.NotEmpty(t, a.config.AdvertiseAddr)
	assert.Equal(t, advAddr+":8946", a.config.AdvertiseAddr)

	a.Stop()
}
