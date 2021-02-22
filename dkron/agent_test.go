package dkron

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/serf/serf"
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

	time.Sleep(2 * time.Second)

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
	if err := a2.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

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
	if err := a3.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)

	// Send a shutdown request
	a1.Stop()

	// Wait until a follower steps as leader
	time.Sleep(2 * time.Second)
	assert.True(t, (a2.IsLeader() || a3.IsLeader()))
	log.Println(a3.IsLeader())

	a2.Stop()
	a3.Stop()
}

func lastSelector(nodes []serf.Member) int {
	return len(nodes) - 1
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

	nodes, err := a1.processFilteredNodes(job, lastSelector)
	require.NoError(t, err)

	assert.Exactly(t, "test1", nodes[0].Name)
	assert.Exactly(t, "test2", nodes[1].Name)
	assert.Len(t, nodes, 2)

	// Test cardinality of 1 with two qualified nodes returns 1 node
	job2 := &Job{
		Name: "test_job_2",
		Tags: map[string]string{
			"tag": "test:1",
		},
	}

	nodes, err = a1.processFilteredNodes(job2, defaultSelector)
	require.NoError(t, err)

	assert.Len(t, nodes, 1)

	// Test no cardinality specified, all nodes returned
	job3 := &Job{
		Name: "test_job_3",
	}

	nodes, err = a1.processFilteredNodes(job3, lastSelector)
	require.NoError(t, err)

	assert.Len(t, nodes, 3)
	assert.Exactly(t, "test1", nodes[0].Name)
	assert.Exactly(t, "test2", nodes[1].Name)
	assert.Exactly(t, "test3", nodes[2].Name)

	// Test exclusive tag returns correct node
	job4 := &Job{
		Name: "test_job_4",
		Tags: map[string]string{
			"tag": "test_client:1",
		},
	}

	nodes, err = a1.processFilteredNodes(job4, defaultSelector)
	require.NoError(t, err)

	assert.Len(t, nodes, 1)
	assert.Exactly(t, "test3", nodes[0].Name)

	// Test existing tag but no matching value returns no nodes
	job5 := &Job{
		Name: "test_job_5",
		Tags: map[string]string{
			"tag": "no_tag",
		},
	}

	nodes, err = a1.processFilteredNodes(job5, defaultSelector)
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

	nodes, err = a1.processFilteredNodes(job6, defaultSelector)
	require.NoError(t, err)

	assert.Len(t, nodes, 0)

	// Test matching tags with cardinality of 2 but only 1 matching node returns correct node
	job7 := &Job{
		Name: "test_job_7",
		Tags: map[string]string{
			"tag":   "test:2",
			"extra": "tag:2",
		},
	}

	nodes, err = a1.processFilteredNodes(job7, defaultSelector)
	require.NoError(t, err)

	assert.Len(t, nodes, 1)
	assert.Exactly(t, "test2", nodes[0].Name)

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
	var sampleSize = 999
	for i := 0; i < sampleSize; i++ {
		// round-robin on the selected nodes to come out at an exactly equal distribution
		roundRobinSelector := func(nodes []serf.Member) int { return i % len(nodes) }
		nodes, err = a1.processFilteredNodes(job8, roundRobinSelector)
		require.NoError(t, err)

		assert.Len(t, nodes, 1)
		distrib[nodes[0].Name]++
	}

	// Each node must have been chosen 1/3 of the time.
	for name, count := range distrib {
		fmt.Println(name, float64(count)/float64(sampleSize)*100.0, "%", count)
	}
	assert.Exactly(t, sampleSize/3, distrib["test1"])
	assert.Exactly(t, sampleSize/3, distrib["test2"])
	assert.Exactly(t, sampleSize/3, distrib["test3"])

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

func Test_filterNodes(t *testing.T) {
	nodes := []serf.Member{
		{
			Tags: map[string]string{
				"region":  "global",
				"tag":     "test",
				"just1":   "value",
				"tagfor2": "2",
			},
		},
		{
			Tags: map[string]string{
				"region":  "global",
				"tag":     "test",
				"just2":   "value",
				"tagfor2": "2",
			},
		},
		{
			Tags: map[string]string{
				"region": "global",
				"tag":    "test",
				"just3":  "value",
			},
		},
	}
	type args struct {
		execNodes []serf.Member
		tags      map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []serf.Member
		want2   int
		wantErr bool
	}{
		{
			name: "All nodes tag",
			args: args{
				execNodes: nodes,
				tags:      map[string]string{"tag": "test"},
			},
			want:    nodes,
			want2:   3,
			wantErr: false,
		},
		{
			name: "Just node1 tag",
			args: args{
				execNodes: nodes,
				tags:      map[string]string{"just1": "value"},
			},
			want:    []serf.Member{nodes[0]},
			want2:   1,
			wantErr: false,
		},
		{
			name: "Just node2 tag",
			args: args{
				execNodes: nodes,
				tags:      map[string]string{"just2": "value"},
			},
			want:    []serf.Member{nodes[1]},
			want2:   1,
			wantErr: false,
		},
		{
			name: "Matching 2 nodes",
			args: args{
				execNodes: nodes,
				tags:      map[string]string{"tagfor2": "2"},
			},
			want:    []serf.Member{nodes[0], nodes[1]},
			want2:   2,
			wantErr: false,
		},
		{
			name: "No matching nodes",
			args: args{
				execNodes: nodes,
				tags:      map[string]string{"unknown": "value"},
			},
			want:    []serf.Member{},
			want2:   0,
			wantErr: false,
		},
		{
			name: "All nodes low cardinality",
			args: args{
				execNodes: nodes,
				tags:      map[string]string{"tag": "test:1"},
			},
			want:    nodes,
			want2:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2, err := filterNodes(tt.args.execNodes, tt.args.tags)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterNodes() got = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("filterNodes() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
