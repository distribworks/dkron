package main

import (
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
	Processors map[string]dkron.ExecutionProcessor
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
	p.Processors = make(map[string]dkron.ExecutionProcessor)

	// Look in /etc/dkron/plugins
	processors, err := plugin.Discover("dkron-processor-*", filepath.Join("/etc", "dkron", "plugins"))
	if err != nil {
		return err
	}

	// Next, look in the same directory as the Dkron executable, usually
	// /usr/local/bin. If found, this replaces what we found in the config path.
	exePath, err := osext.Executable()
	if err != nil {
		logrus.WithError(err).Error("Error loading exe directory")
	} else {
		processors, err = plugin.Discover("dkron-processor-*", filepath.Dir(exePath))
		if err != nil {
			return err
		}
	}

	for _, file := range processors {
		// If the filename has a ".", trim up to there
		// if idx := strings.Index(file, "."); idx >= 0 {
		// 	file = file[:idx]
		// }

		// Look for foo-bar-baz. The plugin name is "baz"
		parts := strings.SplitN(file, "-", 3)
		if len(parts) != 3 {
			continue
		}

		processor, _ := p.processorFactory(file)
		p.Processors[parts[2]] = processor
	}

	return nil
}

func (p *Plugins) processorFactory(path string) (dkron.ExecutionProcessor, error) {
	// Build the plugin client configuration and init the plugin
	var config plugin.ClientConfig
	config.Cmd = exec.Command(path)
	config.HandshakeConfig = dkplugin.Handshake
	config.Managed = true
	config.Plugins = dkplugin.PluginMap
	client := plugin.NewClient(&config)

	// Request the RPC client so we can get the provider
	// so we can build the actual RPC-implemented provider.
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(dkplugin.ProcessorPluginName)
	if err != nil {
		return nil, err
	}

	return raw.(dkron.ExecutionProcessor), nil
}
