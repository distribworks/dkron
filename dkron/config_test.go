package dkron

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestReadConfigTags(t *testing.T) {

	viper.Reset()
	viper.SetConfigType("json")
	var jsonConfig = []byte(`{
		"tags": {
			"foo": "bar"
		}
	}`)
	viper.ReadConfig(bytes.NewBuffer(jsonConfig))
	config := ReadConfig()
	t.Log(config.Tags)
	assert.Equal(t, "bar", config.Tags["foo"])

	viper.Set("tag", []string{"monthy=python"})
	config = ReadConfig()
	assert.NotContains(t, config.Tags, "foo")
	assert.Contains(t, config.Tags, "monthy")
	assert.Equal(t, "python", config.Tags["monthy"])

	config = NewConfig([]string{"-tag", "t1=v1", "-tag", "t2=v2"})
	assert.Equal(t, "v1", config.Tags["t1"])
	assert.Equal(t, "v2", config.Tags["t2"])
}
