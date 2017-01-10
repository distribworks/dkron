package dkron

import (
	"fmt"
	"time"

	"github.com/armon/go-metrics"
)

func initMetrics(config *Config) {
	// Setup the inmem sink and signal handler
	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	metrics.DefaultInmemSignal(inm)

	var fanout metrics.FanoutSink
	// Configure the DogStatsd sink
	if config.DogStatsdAddr != "" {
		var tags []string

		if config.DogStatsdTags != nil {
			tags = config.DogStatsdTags
		}

		sink, err := datadog.NewDogStatsdSink(config.DogStatsdAddr, metricsConf.HostName)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Failed to start DogStatsd sink. Got: %s", err))
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
}
