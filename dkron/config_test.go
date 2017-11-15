package dkron

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/mitchellh/cli"
)

func TestReadConfigTags(t *testing.T) {
	shutdownCh := make(chan struct{})
	defer close(shutdownCh)

	ui := new(cli.MockUi)
	a := &AgentCommand{
		Ui:         ui,
		ShutdownCh: shutdownCh,
	}

	viper.Reset()
	viper.SetConfigType("json")
	var jsonConfig = []byte(`{
		"tags": {
			"foo": "bar"
		}
	}`)
	viper.ReadConfig(bytes.NewBuffer(jsonConfig))
	config := ReadConfig(a)
	t.Log(config.Tags)
	assert.Equal(t, "bar", config.Tags["foo"])

	viper.Set("tag", []string{"monthy=python"})
	config = ReadConfig(a)
	assert.NotContains(t, config.Tags, "foo")
	assert.Contains(t, config.Tags, "monthy")
	assert.Equal(t, "python", config.Tags["monthy"])
}
