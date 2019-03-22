package dkron

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/abronan/valkeyrie/store"

	flag "github.com/spf13/pflag"
)

// Config stores all configuration options for the dkron package.
type Config struct {
	NodeName              string `mapstructure:"node-name"`
	BindAddr              string `mapstructure:"bind-addr"`
	HTTPAddr              string `mapstructure:"http-addr"`
	Backend               store.Backend
	BackendMachines       []string `mapstructure:"backend-machine"`
	BackendPassword       string   `mapstructure:"backend-password"`
	Profile               string
	Interface             string
	AdvertiseAddr         string            `mapstructure:"advertise-addr"`
	Tags                  map[string]string `mapstructure:"tags"`
	SnapshotPath          string            `mapstructure:"snapshot-path"`
	ReconnectInterval     time.Duration     `mapstructure:"reconnect-interval"`
	ReconnectTimeout      time.Duration     `mapstructure:"reconnect-timeout"`
	TombstoneTimeout      time.Duration     `mapstructure:"tombstone-timeout"`
	DisableNameResolution bool              `mapstructure:"disable-name-resolution"`
	KeyringFile           string            `mapstructure:"keyring-file"`
	RejoinAfterLeave      bool              `mapstructure:"rejoin-after-leave"`
	Server                bool
	EncryptKey            string   `mapstructure:"encrypt"`
	StartJoin             []string `mapstructure:"join"`
	Keyspace              string
	RPCPort               int    `mapstructure:"rpc-port"`
	AdvertiseRPCPort      int    `mapstructure:"advertise-rpc-port"`
	LogLevel              string `mapstructure:"log-level"`

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

func DefaultConfig() *Config {
	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}

	tags := map[string]string{"dkron_version": Version}

	return &Config{
		NodeName:          hostname,
		BindAddr:          fmt.Sprintf("0.0.0.0:%d", DefaultBindPort),
		HTTPAddr:          ":8080",
		Backend:           "boltdb",
		BackendMachines:   []string{"./dkron.db"},
		BackendPassword:   "",
		Profile:           "lan",
		Keyspace:          "dkron",
		LogLevel:          "info",
		RPCPort:           6868,
		MailSubjectPrefix: "[Dkron]",
		Tags:              tags,
	}
}

// configFlagSet creates all of our configuration flags.
func ConfigFlagSet() *flag.FlagSet {
	c := DefaultConfig()
	cmdFlags := flag.NewFlagSet("agent flagset", flag.ContinueOnError)

	cmdFlags.Bool("server", false, "This node is running in server mode")
	cmdFlags.String("node-name", c.NodeName, "Name of this node. Must be unique in the cluster")
	cmdFlags.String("bind-addr", c.BindAddr, "Address to bind network listeners to")
	cmdFlags.String("advertise-addr", "", "Address used to advertise to other nodes in the cluster. By default, the bind address is advertised")
	cmdFlags.String("http-addr", c.HTTPAddr, "Address to bind the UI web server to. Only used when server")
	cmdFlags.String("backend", string(c.Backend), "Store backend (etcd|etcdv3|consul|zk|redis|boltdb|dynamodb)")
	cmdFlags.StringSlice("backend-machine", c.BackendMachines, "Store backend machines addresses")
	cmdFlags.String("backend-password", c.BackendPassword, "Store backend machines password or token, only REDIS/CONSUL")
	cmdFlags.String("profile", c.Profile, "Profile is used to control the timing profiles used")
	cmdFlags.StringSlice("join", []string{}, "An initial agent to join with. This flag can be specified multiple times")
	cmdFlags.StringSlice("tag", []string{}, "Tag can be specified multiple times to attach multiple key/value tag pairs to the given node, specified as key=value")
	cmdFlags.String("keyspace", c.Keyspace, "The keyspace to use. A prefix under all data is stored for this instance")
	cmdFlags.String("encrypt", "", "Key for encrypting network traffic. Must be a base64-encoded 16-byte key")
	cmdFlags.String("log-level", c.LogLevel, "Log level (debug|info|warn|error|fatal|panic)")
	cmdFlags.Int("rpc-port", c.RPCPort, "RPC Port used to communicate with clients. Only used when server. The RPC IP Address will be the same as the bind address")
	cmdFlags.Int("advertise-rpc-port", 0, "Use the value of rpc-port by default")

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
