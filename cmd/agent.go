package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/cli"
	"github.com/victorcoder/dkron/dkron"
)

const (
	// gracefulTimeout controls how long we wait before forcefully terminating
	gracefulTimeout = 3 * time.Second
)

// AgentCommand run dkron agent
type AgentCommand struct {
	Ui         cli.Ui
	ShutdownCh <-chan struct{}

	config *dkron.Config
	agent  *dkron.Agent
}

// Help returns agent command usage to the CLI.
func (a *AgentCommand) Help() string {
	helpText := `
Usage: dkron agent [options]
	Run dkron agent

Options:

  -bind-addr=0.0.0.0:8946         Address to bind network listeners to.
  -advertise-addr=bind_addr       Address used to advertise to other nodes in the cluster. By default, the bind address is advertised.
  -http-addr=0.0.0.0:8080         Address to bind the UI web server to. Only used when server.
  -discover=cluster               A cluster name used to discovery peers. On
                                  networks that support multicast, this can be used to have
                                  peers join each other without an explicit join.
  -join=addr                      An initial agent to join with. This flag can be
                                  specified multiple times.
  -node=hostname                  Name of this node. Must be unique in the cluster
  -profile=[lan|wan|local]        Profile is used to control the timing profiles used.
                                  The default if not provided is lan.
  -server=false                   This node is running in server mode.
  -tag key=value                  Tag can be specified multiple times to attach multiple
                                  key/value tag pairs to the given node.
  -keyspace=dkron                 The keyspace to use. A prefix under all data is stored
                                  for this instance.
  -backend=[etcd|consul|zk]       Backend storage to use, etcd, consul or zookeeper. The default
                                  is etcd.
  -backend-machine=127.0.0.1:2379 Backend storage servers addresses to connect to. This flag can be
                                  specified multiple times.
  -encrypt                        Key for encrypting network traffic.
                                  Must be a base64-encoded 16-byte key.
  -rpc-port=6868                  RPC Port used to communicate with clients. Only used when server.
                                  The RPC IP Address will be the same as the bind address.
  -advertise-rpc-port             Use the value of -rpc-port by default

  -mail-host                      Mail server host address to use for notifications.
  -mail-port                      Mail server port.
  -mail-username                  Mail server username used for authentication.
  -mail-password                  Mail server password to use.
  -mail-from                      From email address to use.

  -webhook-url                    Webhook url to call for notifications.
  -webhook-payload                Body of the POST request to send on webhook call.
  -webhook-header                 Headers to use when calling the webhook URL. Can be specified multiple times.

  -log-level=info                 Log level (debug, info, warn, error, fatal, panic). Default to info.

  -dog-statsd-addr                DataDog Agent address
  -dog-statsd-tags                Datadog tags, specified as key:value
  -statsd-addr                    Statsd Address
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns the purpose of the command for the CLI
func (a *AgentCommand) Synopsis() string {
	return "Run dkron"
}

// Run will execute the main functions of the agent command.
// This includes the main eventloop and starting the server if enabled.
//
// The returned value is the exit code.
// protoc -I proto/ proto/executor.proto --go_out=plugins=grpc:proto/
func (a *AgentCommand) Run(args []string) int {
	// Make sure we clean up any managed plugins at the end of this
	p := &Plugins{}
	if err := p.DiscoverPlugins(); err != nil {
		log.Fatal(err)
	}
	plugins := &dkron.Plugins{
		Processors: p.Processors,
		Executors:  p.Executors,
	}

	config := dkron.NewConfig(args)

	agent := dkron.NewAgent(config, plugins)
	if err := agent.Start(); err != nil {
		a.Ui.Error(err.Error())
		return 1
	}
	a.agent = agent

	return a.handleSignals()
}

// handleSignals blocks until we get an exit-causing signal
func (a *AgentCommand) handleSignals() int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

WAIT:
	// Wait for a signal
	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	case <-a.ShutdownCh:
		sig = os.Interrupt
	}
	a.Ui.Output(fmt.Sprintf("Caught signal: %v", sig))

	// Check if this is a SIGHUP
	if sig == syscall.SIGHUP {
		a.handleReload()
		goto WAIT
	}

	// Check if we should do a graceful leave
	graceful := false
	if sig == syscall.SIGTERM || sig == os.Interrupt {
		graceful = true
	}

	// Fail fast if not doing a graceful leave
	if !graceful {
		return 1
	}

	// Attempt a graceful leave
	gracefulCh := make(chan struct{})
	a.Ui.Output("Gracefully shutting down agent...")
	log.Info("agent: Gracefully shutting down agent...")
	go func() {
		plugin.CleanupClients()
		if err := a.agent.Leave(); err != nil {
			a.Ui.Error(fmt.Sprintf("Error: %s", err))
			log.Error(fmt.Sprintf("Error: %s", err))
			return
		}
		close(gracefulCh)
	}()

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
func (a *AgentCommand) handleReload() {
	a.Ui.Output("Reloading configuration...")
	newConf := dkron.ReadConfig()
	if newConf == nil {
		a.Ui.Error(fmt.Sprintf("Failed to reload configs"))
		return
	}
	a.config = newConf

	// Reset serf tags
	if err := a.agent.SetTags(a.config.Tags); err != nil {
		a.Ui.Error(fmt.Sprintf("Failed to reload tags %v", a.config.Tags))
		return
	}
	//Config reloading will also reload Notification settings
}
