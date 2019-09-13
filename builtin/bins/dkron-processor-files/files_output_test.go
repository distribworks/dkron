package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/distribworks/dkron/v2/dkron"
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
			"forward": "false",
			"log_dir": "/tmp",
		},
	}

	fo := &FilesOutput{}
	ex := fo.Process(pa)

	assert.Equal(t, fmt.Sprintf("/tmp/%s.log", ex.Key()), string(ex.Output))
}
