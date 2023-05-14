package main

import (
	"net/http"
	"os"

	dkplugin "github.com/distribworks/dkron/v3/plugin"
	"github.com/hashicorp/go-plugin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	prometheusPort := os.Getenv("SHELL_EXECUTOR_PROMETHEUS_PORT")

	if prometheusPort == "" {
		prometheusPort = "9422" // Default shell executor prometheus metrics port
	}

	promServer := http.NewServeMux()
	promServer.Handle("/metrics", promhttp.Handler())

	go func() {
		http.ListenAndServe(":"+prometheusPort, promServer)
	}()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: dkplugin.Handshake,
		Plugins: map[string]plugin.Plugin{
			"executor": &dkplugin.ExecutorPlugin{Executor: &Shell{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
