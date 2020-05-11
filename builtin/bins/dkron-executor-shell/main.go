package main

import (
	"net/http"
	"os"

	dkplugin "github.com/distribworks/dkron/v3/plugin"
	"github.com/hashicorp/go-plugin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	finish := make(chan bool)
	promServer := http.NewServeMux()
	promServer.Handle("/metrics", promhttp.Handler())

	go func() {
		http.ListenAndServe(":"+getEnv("PROMETHEUS_PORT"), promServer)
	}()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: dkplugin.Handshake,
		Plugins: map[string]plugin.Plugin{
			"executor": &dkplugin.ExecutorPlugin{Executor: &Shell{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
	<-finish
}

func getEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if v == "" {
		log.Warningf("empty value for environment variable %s", key)
		return "set_my_env_var"
	}
	if !ok {
		log.Warningf("environment variable %s is not set", key)
		return "var_is_empty"
	}
	return v
}
