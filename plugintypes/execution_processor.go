package plugintypes

// ExecutionProcessor is an interface that wraps the Process method.
// Plugins must implement this interface.
type ExecutionProcessor interface {
	// Main plugin method, will be called when an execution is done.
	Process(args *ExecutionProcessorArgs) Execution
}

// ExecutionProcessorArgs holds the Execution and PluginConfig for an ExecutionProcessor.
type ExecutionProcessorArgs struct {
	// The execution to pass to the processor
	Execution Execution
	// The configuration for this plugin call
	Config PluginConfig
}

// PluginConfig holds a map of the plugin configuration data structure.
type PluginConfig map[string]interface{}
