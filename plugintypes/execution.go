package plugintypes

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/victorcoder/dkron/proto"
)

// Execution type holds all of the details of a specific Execution.
type Execution struct {
	// Name of the job this executions refers to.
	JobName string `json:"job_name,omitempty"`

	// Start time of the execution.
	StartedAt time.Time `json:"started_at,omitempty"`

	// When the execution finished running.
	FinishedAt time.Time `json:"finished_at,omitempty"`

	// If this execution executed succesfully.
	Success bool `json:"success,omitempty"`

	// Partial output of the execution.
	Output []byte `json:"output,omitempty"`

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
func NewExecutionFromProto(edr *proto.ExecutionDoneRequest) *Execution {
	startedAt, _ := ptypes.Timestamp(edr.GetStartedAt())
	finishedAt, _ := ptypes.Timestamp(edr.GetFinishedAt())
	return &Execution{
		JobName:    edr.JobName,
		Success:    edr.Success,
		Output:     edr.Output,
		NodeName:   edr.NodeName,
		Group:      edr.Group,
		Attempt:    uint(edr.Attempt),
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
	}
}

func (e *Execution) ToProto() *proto.ExecutionDoneRequest {
	startedAt, _ := ptypes.TimestampProto(e.StartedAt)
	finishedAt, _ := ptypes.TimestampProto(e.FinishedAt)
	return &proto.ExecutionDoneRequest{
		JobName:    e.JobName,
		Success:    e.Success,
		Output:     e.Output,
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

func (e *Execution) GetGroup() string {
	return strconv.FormatInt(e.Group, 10)
}

// ExecList stores a slice of Executions.
// This slice can be sorted to provide a time ordered slice of Executions.
type ExecList []*Execution

func (el ExecList) Len() int {
	return len(el)
}

func (el ExecList) Swap(i, j int) {
	el[i], el[j] = el[j], el[i]
}

func (el ExecList) Less(i, j int) bool {
	return el[i].StartedAt.Before(el[j].StartedAt)
}
