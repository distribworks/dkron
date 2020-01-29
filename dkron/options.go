package dkron

import (
	"crypto/tls"
)

// WithPlugins option to set plugins to the agent
func WithPlugins(plugins Plugins) AgentOption {
	return func(agent *Agent) {
		agent.ProcessorPlugins = plugins.Processors
		agent.ExecutorPlugins = plugins.Executors
	}
}

// WithTransportCredentials set tls config in the agent
func WithTransportCredentials(tls *tls.Config) AgentOption {
	return func(agent *Agent) {
		agent.TLSConfig = tls
	}
}

// WithStore set store in the agent
func WithStore(store Storage) AgentOption {
	return func(agent *Agent) {
		agent.Store = store
	}
}
