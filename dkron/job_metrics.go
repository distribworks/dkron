package dkron

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	jobExecutionSuccessCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "dkron",
		Subsystem: "job",
		Name:      "executions_succeeded_total",
		Help:      "Total number of successful job executions",
	}, []string{"job_name"})

	jobExecutionFailureCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "dkron",
		Subsystem: "job",
		Name:      "executions_failed_total",
		Help:      "Total number of failed job executions",
	}, []string{"job_name"})
)
