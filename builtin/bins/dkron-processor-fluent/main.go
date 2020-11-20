package main

import (
	"os"
	"strconv"

	"github.com/distribworks/dkron/v3/plugin"
	"github.com/fluent/fluent-logger-golang/fluent"
	log "github.com/sirupsen/logrus"
)

// FluentOutput is good
type FluentOutput struct {
	forward bool
	tag     string
	fluent  *fluent.Fluent
}

func main() {

	fh := getEnv("FLUENTBIT_HOST")
	fp, _ := strconv.Atoi(getEnv("FLUENTBIT_PORT"))
	ft := getEnv("FLUENTBIT_TAG")

	logger, _ := fluent.New(fluent.Config{
		FluentHost: fh,
		FluentPort: fp,
		RequestAck: false,
	})

	fo := FluentOutput{fluent: logger, tag: ft}

	plugin.Serve(&plugin.ServeOpts{
		Processor: &fo,
	})
}

func getEnv(key string) string {
	v, ok := os.LookupEnv(key)

	if key == "FLUENT_PORT" {
		r, e := strconv.Atoi(v)
		if e != nil {
			return strconv.Itoa(r)
		}
		log.Warningf("non integer value '%s' for FLUENT_PORT variable. Using port 24224 as replacement", v)
		return "24224"
	}

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
