package cmd

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/distribworks/dkron/v3/dkron"
	dkplugin "github.com/distribworks/dkron/v3/plugin"
	"github.com/hashicorp/go-plugin"
	"github.com/kardianos/osext"
	"github.com/sirupsen/logrus"
)

type Plugins struct {
	Processors map[string]dkplugin.Processor
	Executors  map[string]dkplugin.Executor
	LogLevel   string
	NodeName   string
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
	p.Processors = make(map[string]dkplugin.Processor)
	p.Executors = make(map[string]dkplugin.Executor)

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
		p, err := plugin.Discover("dkron-processor-*", filepath.Dir(exePath))
		if err != nil {
			return err
		}
		processors = append(processors, p...)
		e, err := plugin.Discover("dkron-executor-*", filepath.Dir(exePath))
		if err != nil {
			return err
		}
		executors = append(executors, e...)
	}

	for _, file := range processors {

		pluginName, ok := getPluginName(file)
		if !ok {
			continue
		}

		raw, err := p.pluginFactory(file, dkplugin.ProcessorPluginName)
		if err != nil {
			return err
		}
		p.Processors[pluginName] = raw.(dkplugin.Processor)
	}

	for _, file := range executors {

		pluginName, ok := getPluginName(file)
		if !ok {
			continue
		}

		raw, err := p.pluginFactory(file, dkplugin.ExecutorPluginName)
		if err != nil {
			return err
		}
		p.Executors[pluginName] = raw.(dkplugin.Executor)
	}

	return nil
}

func getPluginName(file string) (string, bool) {
	// Look for foo-bar-baz. The plugin name is "baz"
	base := path.Base(file)
	parts := strings.SplitN(base, "-", 3)
	if len(parts) != 3 {
		return "", false
	}

	// This cleans off the .exe for windows plugins
	name := strings.TrimSuffix(parts[2], ".exe")
	return name, true
}

func (p *Plugins) pluginFactory(path string, pluginType string) (interface{}, error) {
	// Build the plugin client configuration and init the plugin
	var config plugin.ClientConfig
	config.Cmd = exec.Command(path)
	config.HandshakeConfig = dkplugin.Handshake
	config.Managed = true
	config.Plugins = dkplugin.PluginMap
	config.SyncStdout = os.Stdout
	config.SyncStderr = os.Stderr
	config.Logger = &dkron.HCLogAdapter{Logger: dkron.InitLogger(p.LogLevel, p.NodeName), LoggerName: "plugins"}

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
