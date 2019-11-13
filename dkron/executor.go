package dkron

type StatusHelper interface {
	Update(float32, []byte, bool) (int64, error)
}

// Executor is the interface that we're exposing as a plugin.
type Executor interface {
	Execute(args *ExecuteRequest, cb StatusHelper) (*ExecuteResponse, error)
}

// ExecutorPluginConfig is the plugin config
type ExecutorPluginConfig map[string]string
