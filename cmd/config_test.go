package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/victorcoder/dkron/dkron"
)

func TestReadConfigTags(t *testing.T) {
	viper.Reset()
	viper.SetConfigType("yaml")
	var yamlConfig = []byte(`
tags:
  - foo: bar
`)
	if err := viper.ReadConfig(bytes.NewBuffer(yamlConfig)); err != nil {
		t.Fatal(err)
	}
	config := dkron.DefaultConfig()
	viper.Unmarshal(config)
	assert.Equal(t, "bar", config.Tags["foo"])

	config = dkron.DefaultConfig()
	viper.Set("tags", map[string]string{"monthy": "python"})
	viper.Unmarshal(config)
	assert.NotContains(t, config.Tags, "foo")
	assert.Contains(t, config.Tags, "monthy")
	assert.Equal(t, "python", config.Tags["monthy"])

	config = &dkron.Config{Tags: map[string]string{"t1": "v1", "t2": "v2"}}
	assert.Equal(t, "v1", config.Tags["t1"])
	assert.Equal(t, "v2", config.Tags["t2"])
}
