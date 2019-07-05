package main

import (
	"testing"
	"time"

	"github.com/distribworks/dkron/dkron"
	"github.com/stretchr/testify/assert"
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
			"forward": true,
		},
	}

	fo := &SyslogOutput{}
	ex := fo.Process(pa)

	assert.Equal(t, "test", string(ex.Output))
}
