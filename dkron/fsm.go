package dkron

import (
	"io"
	"sync"

	dkronpb "github.com/distribworks/dkron/v2/proto"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"
)

// MessageType is the type to encode FSM commands.
type MessageType uint8

const (
	// SetJobType is the command used to store a job in the store.
	SetJobType MessageType = iota
	// DeleteJobType is the command used to delete a Job from the store.
	DeleteJobType
	// SetExecutionType is the command used to store an Execution to the store.
	SetExecutionType
	// DeleteExecutionsType is the command used to delete executions from the store.
	DeleteExecutionsType
	// ExecutionDoneType is the command to perform the logic needed once an exeuction
	// is done.
	ExecutionDoneType
)

type dkronFSM struct {
	store Storage
	mu    sync.Mutex
}

// NewFSM is used to construct a new FSM with a blank state
func NewFSM(store Storage) *dkronFSM {
	return &dkronFSM{
		store: store,
	}
}

// Apply applies a Raft log entry to the key-value store.
func (d *dkronFSM) Apply(l *raft.Log) interface{} {
	buf := l.Data
	msgType := MessageType(buf[0])

	log.WithField("command", msgType).Debug("fsm: received command")

	switch msgType {
	case SetJobType:
		return d.applySetJob(buf[1:])
	case DeleteJobType:
		return d.applyDeleteJob(buf[1:])
	case ExecutionDoneType:
		return d.applyExecutionDone(buf[1:])
	case SetExecutionType:
		return d.applySetExecution(buf[1:])
	}

	return nil
}

func (d *dkronFSM) applySetJob(buf []byte) interface{} {
	var pj dkronpb.Job
	if err := proto.Unmarshal(buf, &pj); err != nil {
		return err
	}
	job := NewJobFromProto(&pj)
	if err := d.store.SetJob(job, false); err != nil {
		return err
	}
	return nil
}

func (d *dkronFSM) applyDeleteJob(buf []byte) interface{} {
	var djr dkronpb.DeleteJobRequest
	if err := proto.Unmarshal(buf, &djr); err != nil {
		return err
	}
	job, err := d.store.DeleteJob(djr.GetJobName())
	if err != nil {
		return err
	}
	return job
}

func (d *dkronFSM) applyExecutionDone(buf []byte) interface{} {
	var execDoneReq dkronpb.ExecutionDoneRequest
	if err := proto.Unmarshal(buf, &execDoneReq); err != nil {
		return err
	}
	execution := NewExecutionFromProto(execDoneReq.Execution)

	log.WithField("execution", execution.Key()).
		WithField("output", string(execution.Output)).
		Debug("fsm: Setting execution")
	_, err := d.store.SetExecutionDone(execution)

	return err
}

func (d *dkronFSM) applySetExecution(buf []byte) interface{} {
	var pbex dkronpb.Execution
	if err := proto.Unmarshal(buf, &pbex); err != nil {
		return err
	}
	execution := NewExecutionFromProto(&pbex)
	key, err := d.store.SetExecution(execution)
	if err != nil {
		return err
	}
	return key
}

// Snapshot returns a snapshot of the key-value store. We wrap
// the things we need in dkronSnapshot and then send that over to Persist.
// Persist encodes the needed data from dkronSnapshot and transport it to
// Restore where the necessary data is replicated into the finite state machine.
// This allows the consensus algorithm to truncate the replicated log.
func (d *dkronFSM) Snapshot() (raft.FSMSnapshot, error) {
	return &dkronSnapshot{store: d.store}, nil
}

// Restore stores the key-value store to a previous state.
func (d *dkronFSM) Restore(r io.ReadCloser) error {
	defer r.Close()
	return d.store.Restore(r)
}

type dkronSnapshot struct {
	store Storage
}

func (d *dkronSnapshot) Persist(sink raft.SnapshotSink) error {
	if err := d.store.Snapshot(sink); err != nil {
		sink.Cancel()
		return err
	}

	// Close the sink.
	if err := sink.Close(); err != nil {
		return err
	}

	return nil
}

func (d *dkronSnapshot) Release() {}
