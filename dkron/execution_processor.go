package dkron

import (
	"github.com/distribworks/dkron/v2/plugin"
)

// ExecutionProcessor is an interface that wraps the Process method.
// Plugins must implement this interface.
type ExecutionProcessor interface {
	// Main plugin method, will be called when an execution is done.
	Process(args *plugin.ExecutionProcessorArgs) Execution
}
