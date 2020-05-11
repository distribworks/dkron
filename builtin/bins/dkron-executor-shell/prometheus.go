package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "dkron_job"
)

var (
	cpuUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "cpu_usage",
		Help:      "CPU usage by job",
	},
		[]string{"job_name"})

	memUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "mem_usage_kb",
		Help:      "Current memory consumed by job",
	},
		[]string{"job_name"})

	executionTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "execution_time_seconds",
		Help:      "Duration of job execution",
	})
		
	exitCode = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "exit_code",
		Help:      "Exit code of a job",
	})
		

	lastExecutionTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "last_execution_unixtimestamp",
	})
		
)

func updateMetric(jobName string, metricName *prometheus.GaugeVec, value float64) {
	metricName.WithLabelValues(jobName).Set(value)
}

