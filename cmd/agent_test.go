package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/distribworks/dkron/v3/dkron"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	logLevel = "error"
)

func getEnvWithDefault() string {
	ea := os.Getenv("DKRON_BACKEND_MACHINE")
	if ea == "" {
		return "127.0.0.1:2379"
	}
	return ea
}

func TestUnmarshalTags(t *testing.T) {
	tagPairs := []string{
		"tag1=val1",
		"tag2=val2",
	}

	tags, err := UnmarshalTags(tagPairs)

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if v, ok := tags["tag1"]; !ok || v != "val1" {
		t.Fatalf("bad: %v", tags)
	}
	if v, ok := tags["tag2"]; !ok || v != "val2" {
		t.Fatalf("bad: %v", tags)
	}
}

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
