package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/distribworks/dkron/v3/plugin"
	"github.com/distribworks/dkron/v3/plugin/types"
	log "github.com/sirupsen/logrus"
)

// LogOutput represent a LogOutputter
type LogOutput struct {
	forward bool
}

// Process method prints the execution output to the stdout
func (l *LogOutput) Process(args *plugin.ProcessorArgs) types.Execution {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	l.parseConfig(args.Config)

	// Output to stdout in case of success or to stderr on failure
	if args.Execution.Success {
		fmt.Printf("----- BEGIN OUTPUT job=%s execution=%s -----\n", args.Execution.JobName, args.Execution.Key())
		fmt.Print(string(args.Execution.Output))
		fmt.Printf("\n----- END OUTPUT -----\n")
	} else {
		fmt.Fprintf(os.Stderr, "----- BEGIN OUTPUT job=%s execution=%s -----\n", args.Execution.JobName, args.Execution.Key())
		fmt.Fprint(os.Stderr, string(args.Execution.Output))
		fmt.Fprintf(os.Stderr, "\n----- END OUTPUT -----\n")
	}

	// Override output if not forwarding
	if !l.forward {
		args.Execution.Output = []byte("Output in dkron log")
	}

	return args.Execution
}

func (l *LogOutput) parseConfig(config plugin.Config) {
	forward, err := strconv.ParseBool(config["forward"])
	if err != nil {
		l.forward = false
		log.WithField("param", "forward").Warning("Incorrect format or param not found.")
	} else {
		l.forward = forward
		log.Infof("Forwarding set to: %t", forward)
	}
}
