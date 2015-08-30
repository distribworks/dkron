package dkron

import (
	"encoding/base64"
	"fmt"
	"net"
	"time"

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
}

// This is the default port that we use for Serf communication
const DefaultBindPort int = 8946

func init() {
	viper.SetConfigName("dkron")        // name of config file (without extension)
	viper.AddConfigPath("/etc/dkron")   // call multiple times to add many search paths
	viper.AddConfigPath("$HOME/.dkron") // call multiple times to add many search paths
	viper.AddConfigPath("./config")     // call multiple times to add many search paths
	err := viper.ReadInConfig()         // Find and read the config file
	if err != nil {                     // Handle errors reading the config file
		log.Infof("No valid config found: %s \n Applying default values.", err)
	}

	viper.SetEnvPrefix("dcr") // will be uppercased automatically
	viper.AutomaticEnv()
}

// BindAddrParts returns the parts of the BindAddr that should be
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
