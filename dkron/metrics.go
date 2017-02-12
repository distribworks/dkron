package dkron

import (
	"fmt"
	"time"

	"github.com/armon/go-metrics"
	"github.com/armon/go-metrics/datadog"
)

func initMetrics(a *AgentCommand) int {
	// Setup the inmem sink and signal handler
	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	metrics.DefaultInmemSignal(inm)

	var fanout metrics.FanoutSink

	// Configure the statsd sink
	if a.config.StatsdAddr != "" {
		sink, err := metrics.NewStatsdSink(a.config.StatsdAddr)
		if err != nil {
			a.Ui.Error(fmt.Sprintf("Failed to start statsd sink. Got: %s", err))
			return 1
		}
		fanout = append(fanout, sink)
	}

	// Configure the DogStatsd sink
	if a.config.DogStatsdAddr != "" {
		var tags []string

		if a.config.DogStatsdTags != nil {
			tags = a.config.DogStatsdTags
		}

		sink, err := datadog.NewDogStatsdSink(a.config.DogStatsdAddr, a.config.NodeName)
		if err != nil {
			a.Ui.Error(fmt.Sprintf("Failed to start DogStatsd sink. Got: %s", err))
			return 1
		}
		sink.SetTags(tags)
		fanout = append(fanout, sink)
	}

	// Initialize the global sink
	if len(fanout) > 0 {
		fanout = append(fanout, inm)
		metrics.NewGlobal(metrics.DefaultConfig("dkron"), fanout)
	} else {
		metrics.NewGlobal(metrics.DefaultConfig("dkron"), inm)
	}

	return 0
}
