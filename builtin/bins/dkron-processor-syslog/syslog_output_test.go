package main

import (
	"testing"
	"time"

	"github.com/distribworks/dkron/v2/plugin"
	"github.com/distribworks/dkron/v2/plugin/types"
	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	now := time.Now()

	pa := &plugin.ProcessorArgs{
		Execution: types.Execution{
			Group:    now.UnixNano(),
			NodeName: "testNode",
			Output:   []byte("test"),
		},
		Config: plugin.Config{
			"forward": "true",
		},
	}

	fo := &SyslogOutput{}
	ex := fo.Process(pa)

	assert.Equal(t, "test", string(ex.Output))
}
