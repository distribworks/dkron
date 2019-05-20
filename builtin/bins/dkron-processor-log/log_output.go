package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/victorcoder/dkron/plugintypes"
)

type LogOutput struct {
	forward bool
}

func (l *LogOutput) Process(args *plugintypes.ExecutionProcessorArgs) plugintypes.Execution {
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

func (l *LogOutput) parseConfig(config plugintypes.PluginConfig) {
	forward, ok := config["forward"].(bool)
	if ok {
		l.forward = forward
		log.Infof("Forwarding set to: %t", forward)
	} else {
		log.Error("Incorrect format in forward param")
	}
}
