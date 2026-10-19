package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	types "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/distribworks/dkron/v4/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestProcessStreamed(t *testing.T) {
	dir := t.TempDir()
	config := plugin.Config{"forward": "false", "log_dir": dir}
	ex := &types.Execution{StartedAt: timestamppb.Now(), NodeName: "testNode"}
	fo := &FilesOutput{}

	require.NoError(t, fo.Begin(&plugin.ProcessorArgs{Execution: ex, Config: config}))

	// Stream more than the 256000 bytes kept in the execution output.
	chunk := bytes.Repeat([]byte("0123456789"), 1000)
	var want []byte
	for i := 0; i < 50; i++ {
		require.NoError(t, fo.Write(&plugin.ProcessorOutputArgs{ExecutionKey: ex.Key(), Output: chunk, Config: config}))
		want = append(want, chunk...)
	}

	// The finished execution only carries the capped output, it must not replace the streamed file.
	ex.Output = []byte("capped")
	out := fo.Process(&plugin.ProcessorArgs{Execution: ex, Config: config})

	path := fmt.Sprintf("%s/%s.log", dir, ex.Key())
	assert.Equal(t, path, string(out.Output))
	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Empty(t, fo.open)
}

func TestProcessNotStreamed(t *testing.T) {
	dir := t.TempDir()
	ex := &types.Execution{StartedAt: timestamppb.Now(), NodeName: "testNode", Output: []byte("final output")}
	fo := &FilesOutput{}

	fo.Process(&plugin.ProcessorArgs{Execution: ex, Config: plugin.Config{"log_dir": dir}})

	got, err := os.ReadFile(fmt.Sprintf("%s/%s.log", dir, ex.Key()))
	require.NoError(t, err)
	assert.Equal(t, "final output", string(got))
}

func TestWriteWithoutBegin(t *testing.T) {
	fo := &FilesOutput{}
	err := fo.Write(&plugin.ProcessorOutputArgs{ExecutionKey: "unknown", Output: []byte("x")})
	assert.Error(t, err)
}
