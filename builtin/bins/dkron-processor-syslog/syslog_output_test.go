package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/victorcoder/dkron/dkron"
)

func TestProcess(t *testing.T) {
	now := time.Now()

	pa := &dkron.ExecutionProcessorArgs{
		Execution: dkron.Execution{
			StartedAt: now,
			NodeName:  "testNode",
			Output:    []byte("test"),
		},
		Config: dkron.PluginConfig{
			"forward": false,
		},
	}

	fo := &SyslogOutput{}
	ex := fo.Process(pa)

	assert.Equal(t, "test", string(ex.Output))
}
