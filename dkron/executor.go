package dkron

// KV is the interface that we're exposing as a plugin.
type Executor interface {
	Execute(args *ExecutorArgs) error
}

// Arguments for calling an execution processor
type ExecutorArgs struct {
	// The execution to pass to the processor
	Execution Execution
	// The configuration for this plugin call
	Config PluginConfig
}
