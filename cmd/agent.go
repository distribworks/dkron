package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/distribworks/dkron/v3/dkron"
	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ShutdownCh chan (struct{})
var agent *dkron.Agent

const (
	// gracefulTimeout controls how long we wait before forcefully terminating
	gracefulTimeout = 3 * time.Hour
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start a dkron agent",
	Long: `Start a dkron agent that schedules jobs, listens for executions and runs executors.
It also runs a web UI.`,
	// Run will execute the main functions of the agent command.
	// This includes the main eventloop and starting the server if enabled.
	//
	// The returned value is the exit code.
	// protoc -I proto/ proto/executor.proto --go_out=plugins=grpc:dkron/
	RunE: func(cmd *cobra.Command, args []string) error {
		return agentRun(args...)
	},
}

func init() {
	dkronCmd.AddCommand(agentCmd)

	agentCmd.Flags().AddFlagSet(dkron.ConfigFlagSet())
	viper.BindPFlags(agentCmd.Flags())
}

func agentRun(args ...string) error {
	// Make sure we clean up any managed plugins at the end of this
	p := &Plugins{
		LogLevel: config.LogLevel,
		NodeName: config.NodeName,
	}
	if err := p.DiscoverPlugins(); err != nil {
		log.Fatal(err)
	}
	plugins := dkron.Plugins{
		Processors: p.Processors,
		Executors:  p.Executors,
	}

	agent = dkron.NewAgent(config, dkron.WithPlugins(plugins))
	if err := agent.Start(); err != nil {
		return err
	}

	exit := handleSignals()
	if exit != 0 {
		return fmt.Errorf("exit status: %d", exit)
	}

	return nil
}

// handleSignals blocks until we get an exit-causing signal
func handleSignals() int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

WAIT:
	// Wait for a signal
	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	case err := <-agent.RetryJoinCh():
		fmt.Println("[ERR] agent: Retry join failed: ", err)
		return 1
	case <-ShutdownCh:
		sig = os.Interrupt
	}
	fmt.Printf("Caught signal: %v", sig)

	// Check if this is a SIGHUP
	if sig == syscall.SIGHUP {
		handleReload()
		goto WAIT
	}

	// Fail fast if not doing a graceful leave
	if sig != syscall.SIGTERM && sig != os.Interrupt {
		return 1
	}

	// Attempt a graceful leave
	log.Info("agent: Gracefully shutting down agent...")
	go func() {
		if err := agent.Stop(); err != nil {
			fmt.Printf("Error: %s", err)
			log.Error(fmt.Sprintf("Error: %s", err))
			return
		}
	}()

	gracefulCh := make(chan struct{})

	for {
		log.Info("Waiting for jobs to finish...")
		if agent.GetRunningJobs() < 1 {
			log.Info("No jobs left. Exiting.")
			break
		}
		time.Sleep(1 * time.Second)
	}

	plugin.CleanupClients()

	close(gracefulCh)

	// Wait for leave or another signal
	select {
	case <-signalCh:
		return 1
	case <-time.After(gracefulTimeout):
		return 1
	case <-gracefulCh:
		return 0
	}
}

// handleReload is invoked when we should reload our configs, e.g. SIGHUP
func handleReload() {
	fmt.Println("Reloading configuration...")
	initConfig()
	//Config reloading will also reload Notification settings
	agent.UpdateTags(config.Tags)
}

// UnmarshalTags is a utility function which takes a slice of strings in
// key=value format and returns them as a tag mapping.
func UnmarshalTags(tags []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, tag := range tags {
		parts := strings.SplitN(tag, "=", 2)
		if len(parts) != 2 || len(parts[0]) == 0 {
			return nil, fmt.Errorf("invalid tag: '%s'", tag)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}
