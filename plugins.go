package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"bitbucket.org/kardianos/osext"
	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-plugin"
	"github.com/victorcoder/dkron/dkron"
	dkplugin "github.com/victorcoder/dkron/plugin"
)

type Plugins struct {
	Outputs map[string]string
}

// Discover plugins located on disk
//
// We look in the following places for plugins:
//
// 1. Dkron configuration path
// 2. Path where Dkron is installed
//
// Whichever file is discoverd LAST wins.
func (p *Plugins) DiscoverPlugins() error {
	// Look in /etc/dkron/plugins
	outputs, err := plugin.Discover("dkron-output-*", filepath.Join("/etc", "dkron", "plugins"))
	if err != nil {
		return err
	}

	// Next, look in the same directory as the Terraform executable, usually
	// /usr/local/bin. If found, this replaces what we found in the config path.
	exePath, err := osext.Executable()
	if err != nil {
		logrus.WithError(err).Error("Error loading exe directory")
	} else {
		outputs, err = plugin.Discover("dkron-output-*", filepath.Dir(exePath))
		if err != nil {
			return err
		}
	}

	for _, file := range outputs {
		// If the filename has a ".", trim up to there
		if idx := strings.Index(file, "."); idx >= 0 {
			file = file[:idx]
		}

		// Look for foo-bar-baz. The plugin name is "baz"
		parts := strings.SplitN(file, "-", 3)
		if len(parts) != 3 {
			continue
		}

		p.Outputs[parts[2]] = file
	}

	return nil
}

// OutputterFactories returns the mapping of prefixes to
// OutputterFactory that can be used to instantiate a
// binary-based plugin.
func (p *Plugins) OutputterFactories() map[string]dkron.OutputterFactory {
	result := make(map[string]dkron.OutputterFactory)
	for k, v := range p.Outputs {
		result[k] = p.outputterFactory(v)
	}

	return result
}

func (p *Plugins) outputterFactory(path string) dkron.OutputterFactory {
	// Build the plugin client configuration and init the plugin
	var config plugin.ClientConfig
	config.Cmd = exec.Command(path)
	fmt.Println("******************: " + path)
	config.HandshakeConfig = dkplugin.Handshake
	config.Managed = true
	config.Plugins = dkplugin.PluginMap
	client := plugin.NewClient(&config)

	return func() (dkron.Outputter, error) {
		// Request the RPC client so we can get the provider
		// so we can build the actual RPC-implemented provider.
		rpcClient, err := client.Client()
		if err != nil {
			return nil, err
		}

		raw, err := rpcClient.Dispense(dkplugin.OutputterPluginName)
		if err != nil {
			return nil, err
		}

		return raw.(dkron.Outputter), nil
	}
}
