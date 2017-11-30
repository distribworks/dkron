package dkron

import (
	"os"
	"testing"
	"time"

	"github.com/abronan/valkeyrie"
	"github.com/abronan/valkeyrie/store"
	"github.com/hashicorp/serf/testutil"
	"github.com/mitchellh/cli"
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

func TestAgentCommandRun(t *testing.T) {
	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	args := []string{
		"-bind-addr", testutil.GetBindAddr().String(),
		"-log-level", logLevel,
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- a.Run(args)
	}()

	time.Sleep(2 * time.Second)

	// Verify it runs "forever"
	select {
	case <-resultCh:
		t.Fatalf("ended too soon, err: %s", ui.ErrorWriter.String())
	case <-time.After(50 * time.Millisecond):
	}

	// Send a shutdown request
	shutdownCh <- struct{}{}

	select {
	case code := <-resultCh:
		if code != 0 {
			t.Fatalf("bad code: %d", code)
		}
	case <-time.After(50 * time.Millisecond):
		t.Fatalf("timeout")
	}
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

	ui := new(cli.MockUi)
	a1 := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	client, err := valkeyrie.NewStore("etcd", []string{etcdAddr}, &store.Config{})
	if err != nil {
		panic(err)
	}
	err = client.DeleteTree("dkron")
	if err != nil {
		if err != store.ErrKeyNotFound {
			panic(err)
		}
	}

	args := []string{
		"-bind-addr", a1Addr,
		"-join", a2Addr,
		"-node-name", a1Name,
		"-server",
		"-log-level", logLevel,
	}

	go a1.Run(args)

	// Wait for the first agent to start and set itself as leader
	kv1, err := watchOrDie(t, client, "dkron/leader")
	if err != nil {
		t.Fatal(err)
	}
	leaderA1 := string(kv1.Value)
	t.Logf("%s is the current leader", leaderA1)
	assert.Equal(t, a1Name, leaderA1)

	// Start another agent
	shutdownCh2 := make(chan struct{})
	defer close(shutdownCh2)

	ui2 := new(cli.MockUi)
	a2 := &AgentCommand{
		Ui:         ui2,
		ShutdownCh: shutdownCh2,
	}
	defer func() { shutdownCh2 <- struct{}{} }()

	args2 := []string{
		"-bind-addr", a2Addr,
		"-join", a1Addr + ":8946",
		"-node-name", a2Name,
		"-server",
		"-log-level", logLevel,
	}

	go a2.Run(args2)

	// Send a shutdown request
	shutdownCh <- struct{}{}
	a1.candidate.Stop()

	// Wait until test2 steps as leader
rewatch:
	kv2, err := watchOrDie(t, client, "dkron/leader")
	if err != nil {
		t.Fatal(err)
	}
	if len(kv2.Value) == 0 || string(kv2.Value) == a1Name {
		goto rewatch
	}
	assert.Equal(t, a2Name, string(kv2.Value))
}

func watchOrDie(t *testing.T, client store.Store, key string) (*store.KVPair, error) {
	for {
		resultCh, err := client.Watch(key, nil, nil)
		if err != nil {
			return nil, err
		}

		_, more := <-resultCh
		if more == false {
			// The channel is closed, recreate the watch
			continue
		}
		// If the channel worked, read the actual value
		kv := <-resultCh
		t.Logf("Value for key %s: %s", key, string(kv.Value))
		return kv, nil
	}
}

func Test_processFilteredNodes(t *testing.T) {
	client, err := valkeyrie.NewStore("etcd", []string{etcdAddr}, &store.Config{})
	err = client.DeleteTree("dkron")
	if err != nil {
		if err == store.ErrNotReachable {
			t.Fatal("etcd server needed to run tests")
		}
	}

	shutdownCh1 := make(chan struct{})
	defer close(shutdownCh1)

	ui := new(cli.MockUi)
	a1 := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh1,
	}

	a1Addr := testutil.GetBindAddr().String()
	a2Addr := testutil.GetBindAddr().String()

	args := []string{
		"-bind-addr", a1Addr,
		"-join", a2Addr,
		"-node-name", "test1",
		"-server",
		"-tag", "role=test",
		"-log-level", logLevel,
	}

	go a1.Run(args)
	time.Sleep(2 * time.Second)

	// Start another agent
	shutdownCh2 := make(chan struct{})
	defer close(shutdownCh2)

	ui2 := new(cli.MockUi)
	a2 := &AgentCommand{
		Ui:         ui2,
		ShutdownCh: shutdownCh2,
	}

	args2 := []string{
		"-bind-addr", a2Addr,
		"-join", a1Addr,
		"-node-name", "test2",
		"-server",
		"-tag", "role=test",
		"-log-level", logLevel,
	}

	go a2.Run(args2)
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

	// Send a shutdown request
	shutdownCh1 <- struct{}{}
	shutdownCh2 <- struct{}{}
}

func Test_UnmarshalTags(t *testing.T) {
	tagPairs := []string{
		"tag1=val1",
		"tag2=val2",
	}

	tags, err := UnmarshalTags(tagPairs)

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if v, ok := tags["tag1"]; !ok || v != "val1" {
		t.Fatalf("bad: %v", tags)
	}
	if v, ok := tags["tag2"]; !ok || v != "val2" {
		t.Fatalf("bad: %v", tags)
	}
}

func TestEncrypt(t *testing.T) {
	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	args := []string{
		"-bind-addr", testutil.GetBindAddr().String(),
		"-node-name", "test1",
		"-server",
		"-tag", "role=test",
		"-encrypt", "kPpdjphiipNSsjd4QHWbkA==",
		"-log-level", logLevel,
	}

	go a.Run(args)
	time.Sleep(2 * time.Second)

	assert.True(t, a.serf.EncryptionEnabled())
	shutdownCh <- struct{}{}
}

func Test_getRPCAddr(t *testing.T) {
	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	a1Addr := testutil.GetBindAddr()

	args := []string{
		"-bind-addr", a1Addr.String() + ":5000",
		"-node-name", "test1",
		"-server",
		"-tag", "role=test",
		"-log-level", logLevel,
	}

	go a.Run(args)
	time.Sleep(2 * time.Second)

	getRPCAddr := a.getRPCAddr()
	exRPCAddr := a1Addr.String() + ":6868"

	assert.Equal(t, exRPCAddr, getRPCAddr)

	shutdownCh <- struct{}{}
}

func TestAgentConfig(t *testing.T) {
	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	advAddr := testutil.GetBindAddr().String()
	args := []string{
		"-bind-addr", testutil.GetBindAddr().String(),
		"-advertise-addr", advAddr,
		"-log-level", logLevel,
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- a.Run(args)
	}()

	time.Sleep(2 * time.Second)

	assert.NotEqual(t, a.config.AdvertiseAddr, a.config.BindAddr)
	assert.NotEmpty(t, a.config.AdvertiseAddr)
	assert.Equal(t, advAddr, a.config.AdvertiseAddr)

	// Send a shutdown request
	shutdownCh <- struct{}{}

	select {
	case code := <-resultCh:
		if code != 0 {
			t.Fatalf("bad code: %d", code)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
}
