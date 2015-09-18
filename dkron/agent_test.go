package dkron

import (
	"bytes"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libkv/store"
	"github.com/hashicorp/serf/testutil"
	"github.com/mitchellh/cli"
)

func TestAgentCommand_implements(t *testing.T) {
	var _ cli.Command = new(AgentCommand)
}

func TestAgentCommandRun(t *testing.T) {
	log.Level = logrus.FatalLevel

	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	args := []string{
		"-bind", testutil.GetBindAddr().String(),
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

func TestAgentCommandElectLeader(t *testing.T) {
	log.Level = logrus.ErrorLevel

	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a1 := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	s := NewStore("etcd", []string{}, nil, "dkron")
	err := s.Client.DeleteTree("dkron")
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
	}

	resultCh := make(chan int)
	go func() {
		resultCh <- a1.Run(args)
	}()

	// Wait for the first agent to start and set itself as leader
	time.Sleep(2 * time.Second)
	test1Key := a1.config.Tags["key"]
	t.Logf("test1 key %s", test1Key)

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
		"-join", a1Addr,
		"-node", "test2",
		"-server",
	}

	resultCh2 := make(chan int)
	go func() {
		resultCh2 <- a2.Run(args2)
	}()

	time.Sleep(2 * time.Second)
	test2Key := a2.config.Tags["key"]
	t.Logf("test2 key %s", test2Key)

	// Send a shutdown request
	shutdownCh <- struct{}{}

	time.Sleep(2 * time.Second)

	// Listen for leader key changes or timeout
	stopCh := make(chan struct{})
	receiver, err := s.Client.Watch("/dkron/leader", stopCh)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for the new leader election
	for {
		select {
		case res := <-receiver:
			if res != nil {
				if bytes.Equal(res.Value, []byte(test2Key)) {
					t.Logf("Leader changed: %s", res.Value)
					stopCh <- struct{}{}
					return
				}
				if bytes.Equal(res.Value, []byte(test1Key)) {
					t.Logf("Leader set to agent1: %s", res.Value)
				}
			}
		case <-time.After(10 * time.Second):
			t.Fatal("No leader swap occurred")
			stopCh <- struct{}{}
			return
		}
	}
}

func Test_processFilteredNodes(t *testing.T) {
	log.Level = logrus.ErrorLevel

	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	s := NewStore("etcd", []string{}, nil, "dkron")
	err := s.Client.DeleteTree("dkron")
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
	nodes, err := a.processFilteredNodes(job)

	if nodes[0] != "test1" || nodes[1] != "test2" {
		t.Fatal("Not expected returned nodes")
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
	log.Level = logrus.ErrorLevel

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
	}

	go func() {
		a.Run(args)
	}()
	time.Sleep(2 * time.Second)

	if !a.serf.EncryptionEnabled() {
		t.Fatal("Encryption not enabled for serf")
	}
	shutdownCh <- struct{}{}
}
