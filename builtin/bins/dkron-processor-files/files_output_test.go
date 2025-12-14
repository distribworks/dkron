package main

import (
	"fmt"
	"testing"

	"github.com/distribworks/dkron/v4/plugin"
	"github.com/distribworks/dkron/v4/types"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestProcess(t *testing.T) {
	now := timestamppb.Now()

	pa := &plugin.ProcessorArgs{
		Execution: &types.Execution{
			StartedAt: now,
			NodeName:  "testNode",
			Output:    []byte("test"),
		},
		Config: plugin.Config{
			"forward": "false",
			"log_dir": "/tmp",
		},
	}

	fo := &FilesOutput{}
	ex := fo.Process(pa)

	assert.Equal(t, fmt.Sprintf("/tmp/%s.log", ex.Key()), string(ex.Output))
}
