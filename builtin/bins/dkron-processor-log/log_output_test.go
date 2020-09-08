package main

import (
	"testing"

	"github.com/distribworks/dkron/v3/plugin"
	"github.com/distribworks/dkron/v3/plugin/types"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	now := ptypes.TimestampNow()

	pa := &plugin.ProcessorArgs{
		Execution: types.Execution{
			StartedAt: now,
			NodeName:  "testNode",
			Output:    []byte("test"),
		},
		Config: plugin.Config{
			"forward": "false",
		},
	}

	fo := &LogOutput{}
	ex := fo.Process(pa)

	assert.Equal(t, "Output in dkron log", string(ex.Output))
}
