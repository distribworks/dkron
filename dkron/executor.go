package dkron

// Executor is the interface that we're exposing as a plugin.
type Executor interface {
	Execute(args *ExecuteRequest) ([]byte, error)
}

// ExecutorPluginConfig is the plugin config
type ExecutorPluginConfig map[string]string
