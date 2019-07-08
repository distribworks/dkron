package dkron

import (
	"crypto/tls"
)

// WithPlugins option to set plugins to the agent
func WithPlugins(plugins *Plugins) AgentOption {
	return func(agent *Agent) {
		if plugins != nil {
			agent.ProcessorPlugins = plugins.Processors
			agent.ExecutorPlugins = plugins.Executors
		}
	}
}

// WithEstablishLeadershipFunc set an extra function to run after leadership acquisition
func WithEstablishLeadershipFunc(establishLeadershipFunc func() error) AgentOption {
	return func(agent *Agent) {
		agent.establishLeadershipFuncs = append(agent.establishLeadershipFuncs, establishLeadershipFunc)
	}
}

// WithRevokeLeadershipFunc set an extra function to run after leadership rovokation
func WithRevokeLeadershipFunc(revokeLeadershipFunc func() error) AgentOption {
	return func(agent *Agent) {
		agent.revokeLeadershipFuncs = append(agent.revokeLeadershipFuncs, revokeLeadershipFunc)
	}
}

// WithTransportCredentials set tls config in the agent
func WithTransportCredentials(tls *tls.Config) AgentOption {
	return func(agent *Agent) {
		agent.TLSConfig = tls
	}
}
