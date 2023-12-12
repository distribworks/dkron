package dkron

import (
	"fmt"
	"strconv"
	"time"

	proto "github.com/distribworks/dkron/v3/plugin/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Execution type holds all of the details of a specific Execution.
type Execution struct {
	// Id is the Key for this execution
	Id string `json:"id,omitempty"`

	// Name of the job this executions refers to.
	JobName string `json:"job_name,omitempty"`

	// Start time of the execution.
	StartedAt time.Time `json:"started_at,omitempty"`

	// When the execution finished running.
	FinishedAt time.Time `json:"finished_at,omitempty"`

	// If this execution executed successfully.
	Success bool `json:"success"`

	// Partial output of the execution.
	Output string `json:"output,omitempty"`

	// Node name of the node that run this execution.
	NodeName string `json:"node_name,omitempty"`

	// Execution group to what this execution belongs to.
	Group int64 `json:"group,omitempty"`

	// Retry attempt of this execution.
	Attempt uint `json:"attempt,omitempty"`
}

// NewExecution creates a new execution.
func NewExecution(jobName string) *Execution {
	return &Execution{
		JobName: jobName,
		Group:   time.Now().UnixNano(),
		Attempt: 1,
	}
}

// NewExecutionFromProto maps a proto.ExecutionDoneRequest to an Execution object
func NewExecutionFromProto(e *proto.Execution) *Execution {
	startedAt := e.GetStartedAt().AsTime()
	finishedAt := e.GetFinishedAt().AsTime()
	return &Execution{
		Id:         e.Key(),
		JobName:    e.JobName,
		Success:    e.Success,
		Output:     string(e.Output),
		NodeName:   e.NodeName,
		Group:      e.Group,
		Attempt:    uint(e.Attempt),
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
	}
}

// ToProto returns the protobuf struct corresponding to
// the representation of the current execution.
func (e *Execution) ToProto() *proto.Execution {
	startedAt := timestamppb.New(e.StartedAt)
	finishedAt := timestamppb.New(e.FinishedAt)
	return &proto.Execution{
		JobName:    e.JobName,
		Success:    e.Success,
		Output:     []byte(e.Output),
		NodeName:   e.NodeName,
		Group:      e.Group,
		Attempt:    uint32(e.Attempt),
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
	}
}

// Key wil generate the execution Id for an execution.
func (e *Execution) Key() string {
	return fmt.Sprintf("%d-%s", e.StartedAt.UnixNano(), e.NodeName)
}

// GetGroup is the getter for the execution group.
func (e *Execution) GetGroup() string {
	return strconv.FormatInt(e.Group, 10)
}
