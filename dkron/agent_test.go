package dkron

import (
	"os"
	"testing"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/hashicorp/serf/testutil"
	"github.com/mitchellh/cli"
)

var logLevel = "debug"
var etcdAddr = getEnvWithDefault()

func getEnvWithDefault() string {
	ea := os.Getenv("DKRON_BACKEND_MACHINE")
	if ea == "" {
		return "127.0.0.1:2379"
	}
	return ea
}

func TestAgentCommand_implements(t *testing.T) {
	var _ cli.Command = new(AgentCommand)
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
		"-bind", testutil.GetBindAddr().String(),
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

	ui := new(cli.MockUi)
	a1 := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	client, err := libkv.NewStore("etcd", []string{etcdAddr}, &store.Config{})
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
		"-bind", a1Addr,
		"-join", a2Addr,
		"-node", a1Name,
		"-server",
		"-log-level", logLevel,
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- a1.Run(args)
	}()

	// Wait for the first agent to start and set itself as leader
	time.Sleep(2 * time.Second)

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
		"-bind", a2Addr,
		"-join", a1Addr + ":8946",
		"-join", a1Addr + ":8946",
		"-node", a2Name,
		"-server",
		"-log-level", logLevel,
	}

	resultCh2 := make(chan int)
	go func() {
		resultCh2 <- a2.Run(args2)
	}()

	kv, _ := client.Get("dkron/leader")
	leader := string(kv.Value)
	log.Printf("%s is the current leader", leader)
	if leader != a1Name {
		t.Errorf("Expected %s to be the leader, got %s", a1Name, leader)
	}

	// Send a shutdown request
	shutdownCh <- struct{}{}

	// Wait until test2 steps as leader
	time.Sleep(30 * time.Second)

	kv, _ = client.Get("dkron/leader")
	leader = string(kv.Value)
	log.Printf("%s is the current leader", leader)
	if leader != a2Name {
		t.Errorf("Expected %s to be the leader, got %s", a2Name, leader)
	}
}

func Test_processFilteredNodes(t *testing.T) {
	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	client, err := libkv.NewStore("etcd", []string{etcdAddr}, &store.Config{})
	err = client.DeleteTree("dkron")
	if err != nil {
		if err == store.ErrNotReachable {
			t.Fatal("etcd server needed to run tests")
		}
	}

	a1Addr := testutil.GetBindAddr().String()
	a2Addr := testutil.GetBindAddr().String()

	args := []string{
		"-bind", a1Addr,
		"-join", a2Addr,
		"-node", "test1",
		"-server",
		"-tag", "role=test",
		"-log-level", logLevel,
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- a.Run(args)
	}()

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
		"-bind", a2Addr,
		"-join", a1Addr,
		"-node", "test2",
		"-server",
		"-tag", "role=test",
		"-log-level", logLevel,
	}

	resultCh2 := make(chan int)
	go func() {
		resultCh2 <- a2.Run(args2)
	}()

	job := &Job{
		Name: "test_job_1",
		Tags: map[string]string{
			"role": "test:2",
		},
	}

	time.Sleep(2 * time.Second)
	nodes, tags, err := a.processFilteredNodes(job)

	if nodes[0] != "test1" || nodes[1] != "test2" {
		t.Fatal("Not expected returned nodes")
	}

	if tags["role"] != "test" {
		t.Fatalf("Tags error, expected: test, got %s", tags["role"])
	}

	if job.Tags["role"] != "test:2" {
		t.Fatalf("Job tags error, expected: 'role=test:2', got 'role=%s'", job.Tags["role"])
	}

	// Send a shutdown request
	shutdownCh <- struct{}{}
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
		"-bind", testutil.GetBindAddr().String(),
		"-node", "test1",
		"-server",
		"-tag", "role=test",
		"-encrypt", "kPpdjphiipNSsjd4QHWbkA==",
		"-log-level", logLevel,
	}

	go a.Run(args)
	time.Sleep(2 * time.Second)

	if !a.serf.EncryptionEnabled() {
		t.Fatal("Encryption not enabled for serf")
	}
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
		"-bind", a1Addr.String() + ":5000",
		"-node", "test1",
		"-server",
		"-tag", "role=test",
		"-log-level", logLevel,
	}

	go a.Run(args)
	time.Sleep(2 * time.Second)

	getRPCAddr := a.getRPCAddr()
	exRPCAddr := a1Addr.String() + ":6868"

	if exRPCAddr != getRPCAddr {
		t.Fatalf("Expected address was: %s got %s", exRPCAddr, getRPCAddr)
	}

	shutdownCh <- struct{}{}
}
