package plugintypes

// Executor is the interface that we're exposing as a plugin.
type Executor interface {
	Execute(args *ExecuteRequest) (*ExecuteResponse, error)
}

// ExecutorPluginConfig is the plugin config
type ExecutorPluginConfig map[string]string
