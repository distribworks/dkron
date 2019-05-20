package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/victorcoder/dkron/plugintypes"
)

func TestProcess(t *testing.T) {
	now := time.Now()

	pa := &plugintypes.ExecutionProcessorArgs{
		Execution: plugintypes.Execution{
			StartedAt: now,
			NodeName:  "testNode",
			Output:    []byte("test"),
		},
		Config: plugintypes.PluginConfig{
			"forward": false,
		},
	}

	fo := &LogOutput{}
	ex := fo.Process(pa)

	assert.Equal(t, "Output in dkron log", string(ex.Output))
}
