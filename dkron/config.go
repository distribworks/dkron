package dkron

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config stores all configuration options for the dkron package.
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
	AdvertiseRPCPort      int

	MailHost          string
	MailPort          uint16
	MailUsername      string
	MailPassword      string
	MailFrom          string
	MailPayload       string
	MailSubjectPrefix string

	WebhookURL     string
	WebhookPayload string
	WebhookHeaders []string

	// DogStatsdAddr is the address of a dogstatsd instance. If provided,
	// metrics will be sent to that instance
	DogStatsdAddr string
	// DogStatsdTags are the global tags that should be sent with each packet to dogstatsd
	// It is a list of strings, where each string looks like "my_tag_name:my_tag_value"
	DogStatsdTags []string
	StatsdAddr    string
}

// DefaultBindPort is the default port that dkron will use for Serf communication
const DefaultBindPort int = 8946

func init() {
	viper.SetConfigName("dkron")        // name of config file (without extension)
	viper.AddConfigPath("/etc/dkron")   // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.dkron") // call multiple times to add many search paths
	viper.AddConfigPath("./config")     // call multiple times to add many search paths
	viper.SetEnvPrefix("dkron")         // will be uppercased automatically
	viper.AutomaticEnv()
}

// NewConfig creates a Config object and will set up the dkron configuration using
// the command line and any file configs.
func NewConfig(args []string) *Config {
	cmdFlags := ConfigFlagSet()

	ignore := args[len(args)-1] == "ignore"
	if ignore {
		args = args[:len(args)-1]
		cmdFlags.SetOutput(ioutil.Discard)
	}

	if err := cmdFlags.Parse(args); err != nil {
		if ignore {
			log.WithError(err).Error("agent: Error parsing flags")
		} else {
			log.Info("agent: Ignoring flag parse errors")
		}
	}

	cmdFlags.VisitAll(func(f *flag.Flag) {
		v := strings.Replace(f.Name, "-", "_", -1)
		if f.Value.String() != f.DefValue {
			if sliceValue, ok := f.Value.(*AppendSliceValue); ok {
				viper.Set(v, ([]string)(*sliceValue))
			} else {
				viper.Set(v, f.Value.String())
			}
		} else {
			viper.SetDefault(v, f.Value.String())
		}
	})

	return ReadConfig()
}

// configFlagSet creates all of our configuration flags.
func ConfigFlagSet() *flag.FlagSet {
	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}

	cmdFlags := flag.NewFlagSet("dkron agent [options]", flag.ContinueOnError)

	cmdFlags.Bool("server", false, "start dkron server")
	cmdFlags.String("node", hostname, "[Deprecated use node-name]")
	cmdFlags.String("node-name", hostname, "node name")
	cmdFlags.String("bind", fmt.Sprintf("0.0.0.0:%d", DefaultBindPort), "[Deprecated use bind-addr]")
	cmdFlags.String("bind-addr", fmt.Sprintf("0.0.0.0:%d", DefaultBindPort), "address to bind listeners to")
	cmdFlags.String("advertise", "", "[Deprecated use advertise-addr]")
	cmdFlags.String("advertise-addr", "", "address to advertise to other nodes")
	cmdFlags.String("http-addr", ":8080", "HTTP address")
	cmdFlags.String("discover", "dkron", "mDNS discovery name")
	cmdFlags.String("backend", "etcd", "store backend")
	cmdFlags.String("backend-machine", "127.0.0.1:2379", "store backend machines addresses")
	cmdFlags.String("profile", "lan", "timing profile to use (lan, wan, local)")
	var join []string
	cmdFlags.Var((*AppendSliceValue)(&join), "join", "address of agent to join on startup")
	var tag []string
	cmdFlags.Var((*AppendSliceValue)(&tag), "tag", "tag pair, specified as key=value")
	cmdFlags.String("keyspace", "dkron", "key namespace to use")
	cmdFlags.String("encrypt", "", "encryption key")
	cmdFlags.String("log-level", "info", "Log level (debug, info, warn, error, fatal, panic), defaults to info")
	cmdFlags.Int("rpc-port", 6868, "RPC port")
	cmdFlags.Int("advertise-rpc-port", 0, "advertise RPC port")

	// Notifications
	cmdFlags.String("mail-host", "", "notification mail server host")
	cmdFlags.String("mail-port", "", "port to use for the mail server")
	cmdFlags.String("mail-username", "", "username for the mail server")
	cmdFlags.String("mail-password", "", "password of the mail server")
	cmdFlags.String("mail-from", "", "notification emails from address")
	cmdFlags.String("mail-payload", "", "notification mail payload")
	cmdFlags.String("mail-subject-prefix", "[Dkron]", "notification mail subject prefix")

	cmdFlags.String("webhook-url", "", "notification webhook url")
	cmdFlags.String("webhook-payload", "", "notification webhook payload")
	webhookHeaders := &AppendSliceValue{}
	cmdFlags.Var(webhookHeaders, "webhook-header", "notification webhook additional header")

	cmdFlags.String("dog-statsd-addr", "", "DataDog Agent address")
	var dogStatsdTags []string
	cmdFlags.Var((*AppendSliceValue)(&dogStatsdTags), "dog-statsd-tags", "Datadog tags, specified as key:value")
	cmdFlags.String("statsd-addr", "", "Statsd Address")

	return cmdFlags
}

// readConfig from file and create the actual config object.
func ReadConfig() *Config {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logrus.WithError(err).Info("No valid config found: Applying default values.")
	}

	cliTags := viper.GetStringSlice("tag")
	var tags map[string]string

	if len(cliTags) > 0 {
		tags, err = UnmarshalTags(cliTags)
		if err != nil {
			logrus.Fatal("config: Error unmarshaling cli tags")
		}
	} else {
		tags = viper.GetStringMapString("tags")
	}

	server := viper.GetBool("server")
	nodeName := viper.GetString("node_name")

	if server {
		tags["dkron_server"] = "true"
	} else {
		tags["dkron_server"] = "false"
	}
	tags["dkron_version"] = Version

	InitLogger(viper.GetString("log_level"), nodeName)

	return &Config{
		NodeName:         nodeName,
		BindAddr:         viper.GetString("bind_addr"),
		AdvertiseAddr:    viper.GetString("advertise_addr"),
		HTTPAddr:         viper.GetString("http_addr"),
		Discover:         viper.GetString("discover"),
		Backend:          viper.GetString("backend"),
		BackendMachines:  viper.GetStringSlice("backend_machine"),
		Server:           server,
		Profile:          viper.GetString("profile"),
		StartJoin:        viper.GetStringSlice("join"),
		Tags:             tags,
		Keyspace:         viper.GetString("keyspace"),
		EncryptKey:       viper.GetString("encrypt"),
		UIDir:            viper.GetString("ui_dir"),
		RPCPort:          viper.GetInt("rpc_port"),
		AdvertiseRPCPort: viper.GetInt("advertise_rpc_port"),

		MailHost:          viper.GetString("mail_host"),
		MailPort:          uint16(viper.GetInt("mail_port")),
		MailUsername:      viper.GetString("mail_username"),
		MailPassword:      viper.GetString("mail_password"),
		MailFrom:          viper.GetString("mail_from"),
		MailPayload:       viper.GetString("mail_payload"),
		MailSubjectPrefix: viper.GetString("mail_subject_prefix"),

		WebhookURL:     viper.GetString("webhook_url"),
		WebhookPayload: viper.GetString("webhook_payload"),
		WebhookHeaders: viper.GetStringSlice("webhook_headers"),

		DogStatsdAddr: viper.GetString("dog_statsd_addr"),
		DogStatsdTags: viper.GetStringSlice("dog_statsd_tags"),
		StatsdAddr:    viper.GetString("statsd_addr"),
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

// NetworkInterface is used to get the associated network
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
