package dkron

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"time"

	flag "github.com/spf13/pflag"
)

// Config stores all configuration options for the dkron package.
type Config struct {
	NodeName                string `mapstructure:"node-name"`
	BindAddr                string `mapstructure:"bind-addr"`
	HTTPAddr                string `mapstructure:"http-addr"`
	Profile                 string
	Interface               string
	AdvertiseAddr           string            `mapstructure:"advertise-addr"`
	Tags                    map[string]string `mapstructure:"tags"`
	SnapshotPath            string            `mapstructure:"snapshot-path"`
	ReconnectInterval       time.Duration     `mapstructure:"reconnect-interval"`
	ReconnectTimeout        time.Duration     `mapstructure:"reconnect-timeout"`
	TombstoneTimeout        time.Duration     `mapstructure:"tombstone-timeout"`
	DisableNameResolution   bool              `mapstructure:"disable-name-resolution"`
	KeyringFile             string            `mapstructure:"keyring-file"`
	RejoinAfterLeave        bool              `mapstructure:"rejoin-after-leave"`
	Server                  bool
	EncryptKey              string        `mapstructure:"encrypt"`
	StartJoin               []string      `mapstructure:"join"`
	RetryJoinLAN            []string      `mapstructure:"retry-join"`
	RetryJoinMaxAttemptsLAN int           `mapstructure:"retry-max"`
	RetryJoinIntervalLAN    time.Duration `mapstructure:"retry-interval"`
	RPCPort                 int           `mapstructure:"rpc-port"`
	AdvertiseRPCPort        int           `mapstructure:"advertise-rpc-port"`
	LogLevel                string        `mapstructure:"log-level"`
	Datacenter              string
	Region                  string

	// Bootstrap mode is used to bring up the first Dkron server.  It is
	// required so that it can elect a leader without any other nodes
	// being present
	Bootstrap bool

	// BootstrapExpect tries to automatically bootstrap the Dkron cluster,
	// by withholding peers until enough servers join.
	BootstrapExpect int `mapstructure:"bootstrap-expect"`

	// DataDir is the directory to store our state in
	DataDir string `mapstructure:"data-dir"`

	// DevMode is used for development purposes only and limits the
	// use of persistence or state.
	DevMode bool

	// ReconcileInterval controls how often we reconcile the strongly
	// consistent store with the Serf info. This is used to handle nodes
	// that are force removed, as well as intermittent unavailability during
	// leader election.
	ReconcileInterval time.Duration

	MailHost          string `mapstructure:"mail-host"`
	MailPort          uint16 `mapstructure:"mail-port"`
	MailUsername      string `mapstructure:"mail-username"`
	MailPassword      string `mapstructure:"mail-password"`
	MailFrom          string `mapstructure:"mail-from"`
	MailPayload       string `mapstructure:"mail-payload"`
	MailSubjectPrefix string `mapstructure:"mail-subject-prefix"`

	WebhookURL     string   `mapstructure:"webhook-url"`
	WebhookPayload string   `mapstructure:"webhook-payload"`
	WebhookHeaders []string `mapstructure:"webhook-headers"`

	// DogStatsdAddr is the address of a dogstatsd instance. If provided,
	// metrics will be sent to that instance
	DogStatsdAddr string `mapstructure:"dog-statsd-addr"`
	// DogStatsdTags are the global tags that should be sent with each packet to dogstatsd
	// It is a list of strings, where each string looks like "my_tag_name:my_tag_value"
	DogStatsdTags []string `mapstructure:"dog-statsd-tags"`
	StatsdAddr    string   `mapstructure:"statsd-addr"`
}

// DefaultBindPort is the default port that dkron will use for Serf communication
const DefaultBindPort int = 8946

// DefaultConfig returns a Config struct pointer with sensible
// default settings.
func DefaultConfig() *Config {
	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}

	tags := map[string]string{}

	return &Config{
		NodeName:          hostname,
		BindAddr:          fmt.Sprintf("0.0.0.0:%d", DefaultBindPort),
		HTTPAddr:          ":8080",
		Profile:           "lan",
		LogLevel:          "info",
		RPCPort:           6868,
		MailSubjectPrefix: "[Dkron]",
		Tags:              tags,
		DataDir:           "dkron.data",
		Datacenter:        "dc1",
		Region:            "global",
		ReconcileInterval: 60 * time.Second,
	}
}

// ConfigFlagSet creates all of our configuration flags.
func ConfigFlagSet() *flag.FlagSet {
	c := DefaultConfig()
	cmdFlags := flag.NewFlagSet("agent flagset", flag.ContinueOnError)

	cmdFlags.Bool("server", false, "This node is running in server mode")
	cmdFlags.String("node-name", c.NodeName, "Name of this node. Must be unique in the cluster")
	cmdFlags.String("bind-addr", c.BindAddr, "Address to bind network listeners to")
	cmdFlags.String("advertise-addr", "", "Address used to advertise to other nodes in the cluster. By default, the bind address is advertised")
	cmdFlags.String("http-addr", c.HTTPAddr, "Address to bind the UI web server to. Only used when server")
	cmdFlags.String("profile", c.Profile, "Profile is used to control the timing profiles used")
	cmdFlags.StringSlice("join", []string{}, "An initial agent to join with. This flag can be specified multiple times")
	cmdFlags.StringSlice("retry-join", []string{}, "Address of an agent to join at start time with retries enabled. Can be specified multiple times.")
	cmdFlags.Int("retry-max", 0, "Maximum number of join attempts. Defaults to 0, which will retry indefinitely.")
	cmdFlags.String("retry-interval", "0", "Time to wait between join attempts.")
	cmdFlags.StringSlice("tag", []string{}, "Tag can be specified multiple times to attach multiple key/value tag pairs to the given node, specified as key=value")
	cmdFlags.String("encrypt", "", "Key for encrypting network traffic. Must be a base64-encoded 16-byte key")
	cmdFlags.String("log-level", c.LogLevel, "Log level (debug|info|warn|error|fatal|panic)")
	cmdFlags.Int("rpc-port", c.RPCPort, "RPC Port used to communicate with clients. Only used when server. The RPC IP Address will be the same as the bind address")
	cmdFlags.Int("advertise-rpc-port", 0, "Use the value of rpc-port by default")
	cmdFlags.Int("bootstrap-expect", 0, "Provides the number of expected servers in the datacenter. Either this value should not be provided or the value must agree with other servers in the cluster. When provided, Dkron waits until the specified number of servers are available and then bootstraps the cluster. This allows an initial leader to be elected automatically. This flag requires server mode.")
	cmdFlags.String("data-dir", c.DataDir, "Specifies the directory to use for server-specific data, including the replicated log. By default, this is the top-level data-dir, like [/var/lib/dkron]")
	cmdFlags.String("datacenter", c.Datacenter, "Specifies the data center of the local agent. All members of a datacenter should share a local LAN connection.")
	cmdFlags.String("region", c.Region, "Specifies the region the Dkron agent is a member of. A region typically maps to a geographic region, for example us, with potentially multiple zones, which map to datacenters such as us-west and us-east")

	// Notifications
	cmdFlags.String("mail-host", "", "Mail server host address to use for notifications")
	cmdFlags.Uint16("mail-port", 0, "Mail server port")
	cmdFlags.String("mail-username", "", "Mail server username used for authentication")
	cmdFlags.String("mail-password", "", "Mail server password to use")
	cmdFlags.String("mail-from", "", "From email address to use")
	cmdFlags.String("mail-payload", "", "Notification mail payload")
	cmdFlags.String("mail-subject-prefix", c.MailSubjectPrefix, "Notification mail subject prefix")

	cmdFlags.String("webhook-url", "", "Webhook url to call for notifications")
	cmdFlags.String("webhook-payload", "", "Body of the POST request to send on webhook call")
	cmdFlags.StringSlice("webhook-header", []string{}, "Headers to use when calling the webhook URL. Can be specified multiple times")

	cmdFlags.String("dog-statsd-addr", "", "DataDog Agent address")
	cmdFlags.StringSlice("dog-statsd-tags", []string{}, "Datadog tags, specified as key:value")
	cmdFlags.String("statsd-addr", "", "Statsd address")

	return cmdFlags
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
