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
		Execution: &dkron.Execution{
			StartedAt: now,
			Node:      "testNode",
			Output:    "test",
		},
		Config: &dkron.PluginConfig{
			"forward": false,
		},
	}

	fo := &FilesOutput{}
	ex := fo.Process(pa)

	assert.Equal(t, "./"+ex.Key(), ex.Output)
}
