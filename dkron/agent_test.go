package dkron

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
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

func lastSelector(nodes []Node) int {
	return len(nodes) - 1
}

func Test_getTargetNodes(t *testing.T) {
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

	t.Run("Test cardinality of 2 returns correct nodes", func(t *testing.T) {
		tags := map[string]string{"tag": "test:2"}

		nodes := a1.getTargetNodes(tags, lastSelector)

		sort.Slice(nodes, func(i, j int) bool { return nodes[i].Name < nodes[j].Name })
		assert.Exactly(t, "test1", nodes[0].Name)
		assert.Exactly(t, "test2", nodes[1].Name)
		assert.Len(t, nodes, 2)
	})

	t.Run("Test cardinality of 1 with two qualified nodes returns 1 node", func(t *testing.T) {
		tags2 := map[string]string{"tag": "test:1"}

		nodes := a1.getTargetNodes(tags2, defaultSelector)

		assert.Len(t, nodes, 1)
	})

	t.Run("Test no cardinality specified, all nodes returned", func(t *testing.T) {
		var tags3 map[string]string

		nodes := a1.getTargetNodes(tags3, lastSelector)

		sort.Slice(nodes, func(i, j int) bool { return nodes[i].Name < nodes[j].Name })
		assert.Len(t, nodes, 3)
		assert.Exactly(t, "test1", nodes[0].Name)
		assert.Exactly(t, "test2", nodes[1].Name)
		assert.Exactly(t, "test3", nodes[2].Name)
	})

	t.Run("Test exclusive tag returns correct node", func(t *testing.T) {
		tags4 := map[string]string{"tag": "test_client:1"}

		nodes := a1.getTargetNodes(tags4, defaultSelector)

		assert.Len(t, nodes, 1)
		assert.Exactly(t, "test3", nodes[0].Name)
	})

	t.Run("Test existing tag but no matching value returns no nodes", func(t *testing.T) {
		tags5 := map[string]string{"tag": "no_tag"}

		nodes := a1.getTargetNodes(tags5, defaultSelector)

		assert.Len(t, nodes, 0)
	})

	t.Run("Test 1 matching and 1 not matching tag returns no nodes", func(t *testing.T) {
		tags6 := map[string]string{
			"foo": "bar:1",
			"tag": "test:2",
		}

		nodes := a1.getTargetNodes(tags6, defaultSelector)

		assert.Len(t, nodes, 0)
	})

	t.Run("Test matching tags with cardinality of 2 but only 1 matching node returns correct node", func(t *testing.T) {
		tags7 := map[string]string{
			"tag":   "test:2",
			"extra": "tag:2",
		}

		nodes := a1.getTargetNodes(tags7, defaultSelector)

		assert.Len(t, nodes, 1)
		assert.Exactly(t, "test2", nodes[0].Name)
	})

	t.Run("Test invalid cardinality yields 0 nodes", func(t *testing.T) {
		tags9 := map[string]string{
			"tag": "test:invalid",
		}

		nodes := a1.getTargetNodes(tags9, defaultSelector)

		assert.Len(t, nodes, 0)
	})

	t.Run("Test two tags matching same 3 servers and cardinality of 1 should always return 1 server", func(t *testing.T) {
		// Do this multiple times: an old bug caused this to sometimes succeed and
		// sometimes fail (=return no nodes at all) due to the use of math.rand
		// Statistically, about 33% should succeed and the rest should fail if
		// the code is buggy.
		// Another bug caused one node to be favored over the others. With a large
		// enough number of attempts, each node should be chosen 1/3 of the time.
		tags8 := map[string]string{
			"additional":  "value:1",
			"additional2": "value2:1",
		}
		distrib := make(map[string]int)

		// Modified version of getTargetNodes
		faked_getTargetNodes := func(tags map[string]string, selectFunc func(nodes []Node) int) []Node {
			bareTags, card := cleanTags(tags, a1.logger)
			allNodes := a1.serf.Members()

			// Sort the nodes: serf.Members() doesn't always return the nodes in the same order, which skews the results.
			sort.Slice(allNodes, func(i, j int) bool { return allNodes[i].Name < allNodes[j].Name })

			nodes := a1.getQualifyingNodes(allNodes, bareTags)
			return selectNodes(nodes, card, selectFunc)
		}

		var sampleSize = 999
		for i := 0; i < sampleSize; i++ {
			roundRobinSelector := func(nodes []Node) int { return i % len(nodes) }

			nodes := faked_getTargetNodes(tags8, roundRobinSelector)

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
	})

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

func Test_getQualifyingNodes(t *testing.T) {
	n1 := Node{
		Status: serf.StatusAlive,
		Tags: map[string]string{
			"region":  "global",
			"tag":     "test",
			"just1":   "value",
			"tagfor2": "2",
		},
	}
	n2 := Node{
		Status: serf.StatusAlive,
		Tags: map[string]string{
			"region":  "global",
			"tag":     "test",
			"just2":   "value",
			"tagfor2": "2",
		},
	}
	n3 := Node{
		Status: serf.StatusAlive,
		Tags: map[string]string{
			"region": "global",
			"tag":    "test",
			"just3":  "value",
		},
	}
	n4 := Node{
		Status: serf.StatusNone,
		Tags: map[string]string{
			"region": "global",
			"dead":   "true",
			"just1":  "value",
		},
	}
	n5 := Node{
		Status: serf.StatusAlive,
		Tags: map[string]string{
			"region": "atlantis",
			"just1":  "value",
		},
	}
	tests := []struct {
		name    string
		inNodes []Node
		inTags  map[string]string
		want    []Node
	}{
		{
			name:    "All nodes match tag",
			inNodes: []Node{n1, n2, n3},
			inTags:  map[string]string{"tag": "test"},
			want:    []Node{n1, n2, n3},
		},
		{
			name:    "Only node1 matches tag",
			inNodes: []Node{n1, n2, n3},
			inTags:  map[string]string{"just1": "value"},
			want:    []Node{n1},
		},
		{
			name:    "Only node2 matches tag",
			inNodes: []Node{n1, n2, n3},
			inTags:  map[string]string{"just2": "value"},
			want:    []Node{n2},
		},
		{
			name:    "Tag matches two nodes",
			inNodes: []Node{n1, n2, n3},
			inTags:  map[string]string{"tagfor2": "2"},
			want:    []Node{n1, n2},
		},
		{
			name:    "No nodes match tag",
			inNodes: []Node{n1, n2, n3},
			inTags:  map[string]string{"unknown": "value"},
			want:    []Node{},
		},
		{
			name:    "Dead nodes don't match",
			inNodes: []Node{n1, n4},
			inTags:  map[string]string{},
			want:    []Node{n1},
		},
		{
			name:    "No nodes returns no nodes",
			inNodes: []Node{},
			inTags:  map[string]string{"just1": "value"},
			want:    []Node{},
		},
		{
			name:    "No tags matches all nodes",
			inNodes: []Node{n1, n2, n3},
			inTags:  map[string]string{},
			want:    []Node{n1, n2, n3},
		},
		{
			name:    "Nodes out of region don't match",
			inNodes: []Node{n1, n5},
			inTags:  map[string]string{},
			want:    []Node{n1},
		},
		{
			name:    "Multiple tags match correct nodes",
			inNodes: []Node{n1, n2, n3},
			inTags:  map[string]string{"tag": "test", "tagfor2": "2"},
			want:    []Node{n1, n2},
		},
	}
	agentStub := NewAgent(DefaultConfig())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := agentStub.getQualifyingNodes(tt.inNodes, tt.inTags)
			assert.Len(t, actual, len(tt.want))
			for _, expectedItem := range tt.want {
				assert.Contains(t, actual, expectedItem)
			}
		})
	}
}

func Test_filterArray(t *testing.T) {
	n1 := Node{Name: "node1"}
	n2 := Node{Name: "node2"}
	n3 := Node{Name: "node3"}
	matchAll := func(m Node) bool { return true }
	matchNone := func(m Node) bool { return false }
	filtertests := []struct {
		name   string
		in     []Node
		filter func(Node) bool
		expect []Node
	}{
		{"No items match", []Node{n1, n2, n3}, matchNone, []Node{}},
		{"All items match", []Node{n1, n2, n3}, matchAll, []Node{n1, n2, n3}},
		{"Empty input returns empty output", []Node{}, matchAll, []Node{}},
		{"All but first match", []Node{n1, n2, n3}, func(m Node) bool { return m.Name != "node1" }, []Node{n2, n3}},
		{"All but last match", []Node{n1, n2, n3}, func(m Node) bool { return m.Name != "node3" }, []Node{n1, n2}},
		{"Middle does not match", []Node{n1, n2, n3}, func(m Node) bool { return m.Name != "node2" }, []Node{n1, n3}},
		{"Only middle matches", []Node{n1, n2, n3}, func(m Node) bool { return m.Name == "node2" }, []Node{n2}},
	}

	for _, tt := range filtertests {
		t.Run(tt.name, func(t *testing.T) {
			actual := filterArray(tt.in, tt.filter)
			assert.Len(t, actual, len(tt.expect))
			for _, expectedItem := range tt.expect {
				assert.Contains(t, actual, expectedItem)
			}
		})
	}
}

func Test_selectNodes(t *testing.T) {
	n1 := Node{Name: "node1"}
	n2 := Node{Name: "node2"}
	n3 := Node{Name: "node3"}
	node2Selector := func(nodes []Node) int {
		for i, node := range nodes {
			if node.Name == "node2" {
				return i
			}
		}
		panic("This shouldn't happen")
	}
	selectertests := []struct {
		name        string
		in          []Node
		cardinality int
		selector    func([]Node) int
		expect      []Node
	}{
		{"Cardinality 0 returns none", []Node{n1, n2, n3}, 0, defaultSelector, []Node{}},
		{"Cardinality < 0 returns none", []Node{n1, n2, n3}, -1, defaultSelector, []Node{}},
		{"Cardinality > #nodes returns all", []Node{n1, n2, n3}, 1000, defaultSelector, []Node{n1, n2, n3}},
		{"Cardinality = #nodes returns all", []Node{n1, n2, n3}, 3, defaultSelector, []Node{n1, n2, n3}},
		{"Cardinality = 1 returns one", []Node{n1, n2, n3}, 1, lastSelector, []Node{n3}},
		{"Cardinality = 2 returns two", []Node{n1, n2, n3}, 2, lastSelector, []Node{n2, n3}},
		{"Pick node2", []Node{n1, n2, n3}, 1, node2Selector, []Node{n2}},
		{"No nodes, card>0 returns none", []Node{}, 2, defaultSelector, []Node{}},
		{"No nodes, card=0 returns none", []Node{}, 0, defaultSelector, []Node{}},
	}
	for _, tt := range selectertests {
		t.Run(tt.name, func(t *testing.T) {
			actual := selectNodes(tt.in, tt.cardinality, tt.selector)
			assert.Len(t, actual, len(tt.expect))
			for _, expectedItem := range tt.expect {
				assert.Contains(t, actual, expectedItem)
			}
		})
	}
}
