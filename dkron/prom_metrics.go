package dkron

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// dkron_agent_event_received_total{event="..."}
	agentEventReceivedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "dkron",
		Subsystem: "agent",
		Name:      "event_received_total",
		Help:      "Count of events received by the agent",
	}, []string{"event"})

	// dkron_job_executions_failed_total{job_name="..."}
	jobExecutionsFailedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "dkron",
		Subsystem: "job",
		Name:      "executions_failed_total",
		Help:      "Total number of failed job executions",
	}, []string{"job_name"})

	// dkron_job_executions_succeeded_total{job_name="..."}
	jobExecutionsSucceededTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "dkron",
		Subsystem: "job",
		Name:      "executions_succeeded_total",
		Help:      "Total number of successful job executions",
	}, []string{"job_name"})
)

func init() {
	// Pre-create common agent event label values so metrics are visible at 0
	agentEventReceivedTotal.WithLabelValues("query_execution_done").Add(0)
	agentEventReceivedTotal.WithLabelValues("query_run_job").Add(0)
}
