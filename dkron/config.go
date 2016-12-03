package dkron

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	NodeName              string
	BindAddr              string
	HTTPAddr              string
	Discover              string
	Backend               string
	BackendMachines       []string
	Profile               string
	Interface             string
	AdvertiseAddr         string
	Tags                  map[string]string
	SnapshotPath          string
	ReconnectInterval     time.Duration
	ReconnectTimeout      time.Duration
	TombstoneTimeout      time.Duration
	DisableNameResolution bool
	KeyringFile           string
	RejoinAfterLeave      bool
	Server                bool
	EncryptKey            string
	StartJoin             AppendSliceValue
	Keyspace              string
	UIDir                 string
	RPCPort               int

	MailHost     string
	MailPort     uint16
	MailUsername string
	MailPassword string
	MailFrom     string
	MailPayload  string

	WebhookURL     string
	WebhookPayload string
	WebhookHeaders []string
}

// This is the default port that we use for Serf communication
const DefaultBindPort int = 8946

func init() {
	viper.SetConfigName("dkron")        // name of config file (without extension)
	viper.AddConfigPath("/etc/dkron")   // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.dkron") // call multiple times to add many search paths
	viper.AddConfigPath("./config")     // call multiple times to add many search paths
	viper.SetEnvPrefix("dkron")         // will be uppercased automatically
	viper.AutomaticEnv()
}

// readConfig is responsible for setup of our configuration using
// the command line and any file configs
func NewConfig(args []string, agent *AgentCommand) *Config {
	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}

	cmdFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	cmdFlags.Usage = func() { agent.Ui.Output(agent.Help()) }
	cmdFlags.String("node", hostname, "node name")
	viper.SetDefault("node_name", cmdFlags.Lookup("node").Value)
	cmdFlags.String("bind", fmt.Sprintf("0.0.0.0:%d", DefaultBindPort), "address to bind listeners to")
	viper.SetDefault("bind_addr", cmdFlags.Lookup("bind").Value)
	cmdFlags.String("advertise", "", "address to advertise to other nodes")
	viper.SetDefault("advertise_addr", cmdFlags.Lookup("advertise").Value)
	cmdFlags.String("http-addr", ":8080", "HTTP address")
	viper.SetDefault("http_addr", cmdFlags.Lookup("http-addr").Value)
	cmdFlags.String("discover", "dkron", "mDNS discovery name")
	viper.SetDefault("discover", cmdFlags.Lookup("discover").Value)
	cmdFlags.String("backend", "etcd", "store backend")
	viper.SetDefault("backend", cmdFlags.Lookup("backend").Value)
	cmdFlags.String("backend-machine", "127.0.0.1:2379", "store backend machines addresses")
	viper.SetDefault("backend_machine", cmdFlags.Lookup("backend-machine").Value)
	cmdFlags.String("profile", "lan", "timing profile to use (lan, wan, local)")
	viper.SetDefault("profile", cmdFlags.Lookup("profile").Value)
	viper.SetDefault("server", cmdFlags.Bool("server", false, "start dkron server"))
	var startJoin []string
	cmdFlags.Var((*AppendSliceValue)(&startJoin), "join", "address of agent to join on startup")
	var tag []string
	cmdFlags.Var((*AppendSliceValue)(&tag), "tag", "tag pair, specified as key=value")
	cmdFlags.String("keyspace", "dkron", "key namespace to use")
	viper.SetDefault("keyspace", cmdFlags.Lookup("keyspace").Value)
	cmdFlags.String("encrypt", "", "encryption key")
	viper.SetDefault("encrypt", cmdFlags.Lookup("encrypt").Value)

	cmdFlags.String("log-level", "info", "Log level (debug, info, warn, error, fatal, panic), defaults to info")
	viper.SetDefault("log_level", cmdFlags.Lookup("log-level").Value)

	cmdFlags.String("ui-dir", ".", "directory to serve web UI")
	viper.SetDefault("ui_dir", cmdFlags.Lookup("ui-dir").Value)
	viper.SetDefault("rpc_port", cmdFlags.Int("rpc-port", 6868, "RPC port"))

	// Notifications
	cmdFlags.String("mail-host", "", "notification mail server host")
	viper.SetDefault("mail_host", cmdFlags.Lookup("mail-host").Value)
	cmdFlags.String("mail-port", "", "port to use for the mail server")
	viper.SetDefault("mail_port", cmdFlags.Lookup("mail-port").Value)
	cmdFlags.String("mail-username", "", "username for the mail server")
	viper.SetDefault("mail_username", cmdFlags.Lookup("mail-username").Value)
	cmdFlags.String("mail-password", "", "password of the mail server")
	viper.SetDefault("mail_password", cmdFlags.Lookup("mail-password").Value)
	cmdFlags.String("mail-from", "", "notification emails from address")
	viper.SetDefault("mail_from", cmdFlags.Lookup("mail-from").Value)
	cmdFlags.String("mail-payload", "", "notification mail payload")
	viper.SetDefault("mail_payload", cmdFlags.Lookup("mail-payload").Value)

	cmdFlags.String("webhook-url", "", "notification webhook url")
	viper.SetDefault("webhook_url", cmdFlags.Lookup("webhook-url").Value)
	cmdFlags.String("webhook-payload", "", "notification webhook payload")
	viper.SetDefault("webhook_payload", cmdFlags.Lookup("webhook-payload").Value)
	webhookHeaders := &AppendSliceValue{}
	cmdFlags.Var(webhookHeaders, "webhook-header", "notification webhook additional header")

	if err := cmdFlags.Parse(args); err != nil {
		log.Fatal(err)
	}

	ut, err := UnmarshalTags(tag)
	if err != nil {
		log.Fatal(err)
	}
	viper.SetDefault("tags", ut)
	viper.SetDefault("join", startJoin)
	viper.SetDefault("webhook_headers", webhookHeaders)

	return ReadConfig(agent)
}

func ReadConfig(agent *AgentCommand) *Config {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logrus.WithError(err).Info("No valid config found: Applying default values.")
	}

	tags := viper.GetStringMapString("tags")
	server := viper.GetBool("server")
	nodeName := viper.GetString("node_name")

	if server {
		tags["dkron_server"] = "true"
	}
	tags["dkron_version"] = agent.Version

	InitLogger(viper.GetString("log_level"), nodeName)

	return &Config{
		NodeName:        nodeName,
		BindAddr:        viper.GetString("bind_addr"),
		AdvertiseAddr:   viper.GetString("advertise_addr"),
		HTTPAddr:        viper.GetString("http_addr"),
		Discover:        viper.GetString("discover"),
		Backend:         viper.GetString("backend"),
		BackendMachines: viper.GetStringSlice("backend_machine"),
		Server:          server,
		Profile:         viper.GetString("profile"),
		StartJoin:       viper.GetStringSlice("join"),
		Tags:            tags,
		Keyspace:        viper.GetString("keyspace"),
		EncryptKey:      viper.GetString("encrypt"),
		UIDir:           viper.GetString("ui_dir"),
		RPCPort:         viper.GetInt("rpc_port"),

		MailHost:     viper.GetString("mail_host"),
		MailPort:     uint16(viper.GetInt("mail_port")),
		MailUsername: viper.GetString("mail_username"),
		MailPassword: viper.GetString("mail_password"),
		MailFrom:     viper.GetString("mail_from"),
		MailPayload:  viper.GetString("mail_payload"),

		WebhookURL:     viper.GetString("webhook_url"),
		WebhookPayload: viper.GetString("webhook_payload"),
		WebhookHeaders: viper.GetStringSlice("webhook_headers"),
	}
}

// AddrParts returns the parts of the BindAddr that should be
// used to configure Serf.
func (c *Config) AddrParts(address string) (string, int, error) {
	checkAddr := address

START:
	_, _, err := net.SplitHostPort(checkAddr)
	if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
		checkAddr = fmt.Sprintf("%s:%d", checkAddr, DefaultBindPort)
		goto START
	}
	if err != nil {
		return "", 0, err
	}

	// Get the address
	addr, err := net.ResolveTCPAddr("tcp", checkAddr)
	if err != nil {
		return "", 0, err
	}

	return addr.IP.String(), addr.Port, nil
}

// Networkinterface is used to get the associated network
// interface from the configured value
func (c *Config) NetworkInterface() (*net.Interface, error) {
	if c.Interface == "" {
		return nil, nil
	}
	return net.InterfaceByName(c.Interface)
}

// EncryptBytes returns the encryption key configured.
func (c *Config) EncryptBytes() ([]byte, error) {
	return base64.StdEncoding.DecodeString(c.EncryptKey)
}
