package dkron

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/hashicorp/go-plugin"
)

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

// Init a new execution
func NewExecution(jobName string) *Execution {
	return &Execution{
		JobName: jobName,
		Group:   time.Now().UnixNano(),
		Attempt: 1,
	}
}

// Used to enerate the execution Id
func (e *Execution) Key() string {
	return fmt.Sprintf("%d-%s", e.StartedAt.UnixNano(), e.NodeName)
}

type Outputter interface {
	Output(execution Execution) string
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"outputer": new(OutputterPlugin),
}

type OutputterPlugin struct{}

func (OutputterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return nil, nil
}

func (OutputterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &OutputterRPC{client: c}, nil
}

// Here is an implementation that talks over RPC
type OutputterRPC struct{ client *rpc.Client }

func (g *OutputterRPC) Output(execution Execution) string {
	var resp string
	err := g.client.Call("Plugin.Output", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}
