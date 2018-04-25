package cmd

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/go-plugin"
	"github.com/kardianos/osext"
	"github.com/victorcoder/dkron/dkron"
	dkplugin "github.com/victorcoder/dkron/plugin"
)

type Plugins struct {
	Processors map[string]dkron.ExecutionProcessor
	Executors  map[string]dkron.Executor
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
	p.Executors = make(map[string]dkron.Executor)

	// Look in /etc/dkron/plugins
	processors, err := plugin.Discover("dkron-processor-*", filepath.Join("/etc", "dkron", "plugins"))
	if err != nil {
		return err
	}

	// Look in /etc/dkron/plugins
	executors, err := plugin.Discover("dkron-executor-*", filepath.Join("/etc", "dkron", "plugins"))
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
		executors, err = plugin.Discover("dkron-executor-*", filepath.Dir(exePath))
		if err != nil {
			return err
		}
	}

	for _, file := range processors {
		// Look for foo-bar-baz. The plugin name is "baz"
		parts := strings.SplitN(file, "-", 3)
		if len(parts) != 3 {
			continue
		}

		raw, err := p.pluginFactory(file, dkplugin.ProcessorPluginName)
		if err != nil {
			return err
		}
		p.Processors[parts[2]] = raw.(dkron.ExecutionProcessor)
	}

	for _, file := range executors {
		// Look for foo-bar-baz. The plugin name is "baz"
		parts := strings.SplitN(file, "-", 3)
		if len(parts) != 3 {
			continue
		}

		raw, err := p.pluginFactory(file, dkplugin.ExecutorPluginName)
		if err != nil {
			return err
		}
		p.Executors[parts[2]] = raw.(dkron.Executor)
	}

	return nil
}

func (Plugins) pluginFactory(path string, pluginType string) (interface{}, error) {
	// Build the plugin client configuration and init the plugin
	var config plugin.ClientConfig
	config.Cmd = exec.Command(path)
	config.HandshakeConfig = dkplugin.Handshake
	config.Managed = true
	config.Plugins = dkplugin.PluginMap

	switch pluginType {
	case dkplugin.ProcessorPluginName:
		config.AllowedProtocols = []plugin.Protocol{plugin.ProtocolNetRPC}
	case dkplugin.ExecutorPluginName:
		config.AllowedProtocols = []plugin.Protocol{plugin.ProtocolGRPC}
	}

	client := plugin.NewClient(&config)

	// Request the RPC client so we can get the provider
	// so we can build the actual RPC-implemented provider.
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(pluginType)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
