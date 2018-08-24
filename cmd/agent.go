package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/victorcoder/dkron/dkron"
)

var ShutdownCh chan (struct{})
var agent *dkron.Agent

const (
	// gracefulTimeout controls how long we wait before forcefully terminating
	gracefulTimeout = 3 * time.Second
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start a dkron agent",
	Long: `Start a dkron agent that schedule jobs, listen for executions and run executors.
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
	legacyConfig()

	// Make sure we clean up any managed plugins at the end of this
	p := &Plugins{}
	if err := p.DiscoverPlugins(); err != nil {
		log.Fatal(err)
	}
	plugins := &dkron.Plugins{
		Processors: p.Processors,
		Executors:  p.Executors,
	}

	agent = dkron.NewAgent(config, plugins)
	if err := agent.Start(); err != nil {
		return err
	}

	exit := handleSignals()
	if exit != 0 {
		return fmt.Errorf("Exit status: %d", exit)
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
	case <-ShutdownCh:
		sig = os.Interrupt
	}
	fmt.Printf("Caught signal: %v", sig)

	// Check if this is a SIGHUP
	if sig == syscall.SIGHUP {
		handleReload()
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
	log.Info("agent: Gracefully shutting down agent...")
	go func() {
		plugin.CleanupClients()
		if err := agent.Leave(); err != nil {
			fmt.Printf("Error: %s", err)
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
func handleReload() {
	fmt.Println("Reloading configuration...")
	initConfig()

	// Reset serf tags
	if err := agent.SetTags(agent.Config().Tags); err != nil {
		fmt.Printf("Failed to reload tags %v", agent.Config().Tags)
		return
	}
	//Config reloading will also reload Notification settings
}

// Suport legacy config files for some time
func legacyConfig() {
	s := viper.GetString("node_name")
	if s != "" && viper.GetString("node-name") == "" {
		log.WithField("param", "node_name").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.NodeName = s
	}

	s = viper.GetString("bind_addr")
	if s != "" && viper.GetString("bind-addr") == "" {
		log.WithField("param", "bind_addr").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.BindAddr = s
	}

	s = viper.GetString("http_addr")
	if s != "" && viper.GetString("http-addr") == "" {
		log.WithField("param", "http_addr").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.HTTPAddr = s
	}

	ss := viper.GetStringSlice("backend_machine")
	if ss != nil && viper.GetStringSlice("backend_machine") == nil {
		log.WithField("param", "backend_machine").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.BackendMachines = ss
	}

	s = viper.GetString("advertise_addr")
	if s != "" && viper.GetString("advertise-addr") == "" {
		log.WithField("param", "advertise_addr").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.AdvertiseAddr = s
	}

	s = viper.GetString("snapshot_path")
	if s != "" && viper.GetString("snapshot-path") == "" {
		log.WithField("param", "snapshot_path").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.SnapshotPath = s
	}

	d := viper.GetDuration("reconnect_interval")
	if d != 0 && viper.GetDuration("reconnect-interval") == 0 {
		log.WithField("param", "reconnect_interval").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.ReconnectInterval = d
	}

	d = viper.GetDuration("reconnect_timeout")
	if d != 0 && viper.GetDuration("reconnect-timeout") == 0 {
		log.WithField("param", "reconnect_timeout").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.ReconnectTimeout = d
	}

	d = viper.GetDuration("tombstone_timeout")
	if d != 0 && viper.GetDuration("tombstone-timeout") == 0 {
		log.WithField("param", "tombstone_timeout").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.TombstoneTimeout = d
	}

	b := viper.GetBool("disable_name_resolution")
	if b != false && viper.GetBool("disable-name-resolution") == false {
		log.WithField("param", "disable_name_resolution").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.DisableNameResolution = b
	}

	s = viper.GetString("keyring_file")
	if s != "" && viper.GetString("keyring-file") == "" {
		log.WithField("param", "keyring_file").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.KeyringFile = s
	}

	b = viper.GetBool("rejoin_after_leave")
	if b != false && viper.GetBool("rejoin-after-leave") == false {
		log.WithField("param", "rejoin_after_leave").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.RejoinAfterLeave = b
	}

	s = viper.GetString("encrypt_key")
	if s != "" && viper.GetString("encrypt-key") == "" {
		log.WithField("param", "encrypt_key").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.EncryptKey = s
	}

	ss = viper.GetStringSlice("start_join")
	if ss != nil && viper.GetStringSlice("start-join") == nil {
		log.WithField("param", "start_join").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.StartJoin = ss
	}

	i := viper.GetInt("rpc_port")
	if i != 0 && viper.GetInt("rpc-port") == 0 {
		log.WithField("param", "rpc_port").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.RPCPort = i
	}

	i = viper.GetInt("advertise_rpc_port")
	if i != 0 && viper.GetInt("advertise-rpc-port") == 0 {
		log.WithField("param", "advertise_rpc_port").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.AdvertiseRPCPort = i
	}

	s = viper.GetString("log_level")
	if s != "" && viper.GetString("log-level") == "" {
		log.WithField("param", "log_level").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.LogLevel = s
	}

	s = viper.GetString("mail_host")
	if s != "" && viper.GetString("mail-host") == "" {
		log.WithField("param", "mail_host").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.MailHost = s
	}

	i = viper.GetInt("mail_port")
	if i != 0 && viper.GetInt("mail-port") == 0 {
		log.WithField("param", "mail_port").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.MailPort = uint16(i)
	}

	s = viper.GetString("mail_username")
	if s != "" && viper.GetString("mail-username") == "" {
		log.WithField("param", "mail_username").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.MailUsername = s
	}

	s = viper.GetString("mail_password")
	if s != "" && viper.GetString("mail-password") == "" {
		log.WithField("param", "mail_password").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.MailPassword = s
	}

	s = viper.GetString("mail_from")
	if s != "" && viper.GetString("mail-from") == "" {
		log.WithField("param", "mail_from").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.MailFrom = s
	}

	s = viper.GetString("mail_payload")
	if s != "" && viper.GetString("mail-payload") == "" {
		log.WithField("param", "mail_payload").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.MailPayload = s
	}

	s = viper.GetString("mail_subject_prefix")
	if s != "" && viper.GetString("mail-subject-prefix") == "" {
		log.WithField("param", "mail_subject_prefix").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.MailSubjectPrefix = s
	}

	s = viper.GetString("webhook_url")
	if s != "" && viper.GetString("webhook-url") == "" {
		log.WithField("param", "webhook_url").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.WebhookURL = s
	}

	s = viper.GetString("webhook_payload")
	if s != "" && viper.GetString("webhook-payload") == "" {
		log.WithField("param", "webhook_payload").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.WebhookPayload = s
	}

	ss = viper.GetStringSlice("webhook_headers")
	if ss != nil && viper.GetStringSlice("webhook-headers") == nil {
		log.WithField("param", "webhook_headers").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.WebhookHeaders = ss
	}

	s = viper.GetString("dog_statsd_addr")
	if s != "" && viper.GetString("dog-statsd-addr") == "" {
		log.WithField("param", "dog_statsd_add").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.DogStatsdAddr = s
	}

	ss = viper.GetStringSlice("dog_statsd_tags")
	if ss != nil && viper.GetStringSlice("dog-statsd-tags") == nil {
		log.WithField("param", "dog_statsd_tags").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.DogStatsdTags = ss
	}

	s = viper.GetString("statsd_addr")
	if s != "" && viper.GetString("statsd-addr") == "" {
		log.WithField("param", "statsd_addr").Warn("Deprecation warning: Config param name is deprecated and will be removed in future versions.")
		config.StatsdAddr = s
	}

}
