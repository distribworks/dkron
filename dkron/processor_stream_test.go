package dkron

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	typesv1 "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/distribworks/dkron/v4/plugin"
	"github.com/distribworks/dkron/v4/plugin/shell"
	"github.com/hashicorp/serf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type streamingProcessor struct {
	beginErr error
	writeErr error
	configs  []plugin.Config
	chunks   []string
}

func (p *streamingProcessor) Begin(args *plugin.ProcessorArgs) error {
	p.configs = append(p.configs, args.Config)
	return p.beginErr
}

func (p *streamingProcessor) Write(args *plugin.ProcessorOutputArgs) error {
	if p.writeErr != nil {
		return p.writeErr
	}
	p.chunks = append(p.chunks, string(args.Output))
	return nil
}

func (p *streamingProcessor) Process(args *plugin.ProcessorArgs) *typesv1.Execution {
	return args.Execution
}

func streamTestJob(processors ...string) *typesv1.Job {
	job := &typesv1.Job{Name: "stream-job", Processors: map[string]*typesv1.PluginConfig{}}
	for _, name := range processors {
		job.Processors[name] = &typesv1.PluginConfig{Config: map[string]string{"log_dir": "/logs/" + name}}
	}
	return job
}

func TestProcessorStream_WritesChunksInOrder(t *testing.T) {
	files := &streamingProcessor{}
	other := &streamingProcessor{}
	plugins := map[string]plugin.Processor{"files": files, "other": other}
	execution := &typesv1.Execution{JobName: "stream-job", NodeName: "agent-1", StartedAt: timestamppb.Now()}

	ps := newProcessorStream(plugins, streamTestJob("files", "other", "missing"), execution, "server-1", getTestLogger())
	ps.write([]byte("one"), false)
	ps.write(nil, false)
	ps.write([]byte("two"), true)

	assert.Equal(t, []plugin.Config{{"log_dir": "/logs/files", "reporting_node": "server-1"}}, files.configs)
	assert.Equal(t, []string{"one", "two"}, files.chunks)
	assert.Equal(t, []string{"one", "two"}, other.chunks)
	assert.Equal(t, execution.Key(), ps.key)
}

func TestProcessorStream_SkipsFailingProcessors(t *testing.T) {
	failsBegin := &streamingProcessor{beginErr: errors.New("no dir")}
	failsWrite := &streamingProcessor{writeErr: errors.New("disk full")}
	ok := &streamingProcessor{}
	plugins := map[string]plugin.Processor{"begin": failsBegin, "write": failsWrite, "ok": ok}
	execution := &typesv1.Execution{JobName: "stream-job", NodeName: "agent-1", StartedAt: timestamppb.Now()}

	ps := newProcessorStream(plugins, streamTestJob("begin", "write", "ok"), execution, "server-1", getTestLogger())
	ps.write([]byte("one"), false)
	ps.write([]byte("two"), false)

	assert.Empty(t, failsBegin.chunks)
	assert.NotContains(t, ps.targets, "write")
	assert.Equal(t, []string{"one", "two"}, ok.chunks)
}

// recordingProcessor collects everything a processor sees from a real execution.
type recordingProcessor struct {
	mu        sync.Mutex
	begins    int
	output    bytes.Buffer
	processed *typesv1.Execution
}

func (p *recordingProcessor) Begin(args *plugin.ProcessorArgs) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.begins++
	return nil
}

func (p *recordingProcessor) Write(args *plugin.ProcessorOutputArgs) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.output.Write(args.Output)
	return nil
}

func (p *recordingProcessor) Process(args *plugin.ProcessorArgs) *typesv1.Execution {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.processed = args.Execution
	return args.Execution
}

func TestProcessorStream_StreamsCompleteOutputOfRunningJob(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses seq")
	}

	ip1, returnFn1 := testutil.TakeIP()
	defer returnFn1()

	c := DefaultConfig()
	c.BindAddr = ip1.String()
	c.NodeName = "test-stream"
	c.Server = true
	c.LogLevel = logLevel
	c.BootstrapExpect = 1
	c.DevMode = true
	c.DataDir = t.TempDir()
	c.HTTPAddr = "127.0.0.1:0"

	processor := &recordingProcessor{}
	a := NewAgent(c, WithPlugins(Plugins{
		Processors: map[string]plugin.Processor{"recording": processor},
		Executors:  map[string]plugin.Executor{"shell": &shell.Shell{}},
	}))
	require.NoError(t, a.Start())
	defer func() { _ = a.Stop() }()

	require.Eventually(t, a.IsLeader, 10*time.Second, 10*time.Millisecond)

	ctx := context.Background()
	parent := &Job{Name: "stream-parent", Schedule: "@manually", Executor: "shell",
		ExecutorConfig: map[string]string{"command": "/bin/true"}, Disabled: true}
	require.NoError(t, a.Store.SetJob(ctx, parent, true))

	// Output well over the 256000 bytes kept in the stored execution.
	job := &Job{Name: "stream-child", ParentJob: parent.Name, Executor: "shell",
		ExecutorConfig: map[string]string{"command": "seq 1 100000", "shell": "true"},
		Processors:     map[string]plugin.Config{"recording": {}}}
	require.NoError(t, a.Store.SetJob(ctx, job, true))

	_, err := a.Run(ctx, job.Name, &Execution{JobName: job.Name, Group: time.Now().UnixNano(), Attempt: 1})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		processor.mu.Lock()
		defer processor.mu.Unlock()
		return processor.processed != nil
	}, 30*time.Second, 50*time.Millisecond)

	var want bytes.Buffer
	for i := 1; i <= 100000; i++ {
		fmt.Fprintf(&want, "%d\n", i)
	}

	processor.mu.Lock()
	defer processor.mu.Unlock()
	assert.Equal(t, 1, processor.begins)
	assert.Equal(t, want.Len(), processor.output.Len())
	assert.Equal(t, want.String(), processor.output.String())
	// The stored output stays capped, and ends with the end of the stream.
	assert.LessOrEqual(t, len(processor.processed.Output), maxBufSize)
	assert.True(t, bytes.HasSuffix(want.Bytes(), processor.processed.Output))
}
