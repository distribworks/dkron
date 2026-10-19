package plugin

import (
	"errors"
	"net"
	"net/rpc"
	"testing"

	types "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type recordingProcessor struct {
	begun    []string
	chunks   []string
	writeErr error
}

func (p *recordingProcessor) Begin(args *ProcessorArgs) error {
	p.begun = append(p.begun, args.Execution.Key())
	return nil
}

func (p *recordingProcessor) Write(args *ProcessorOutputArgs) error {
	if p.writeErr != nil {
		return p.writeErr
	}
	p.chunks = append(p.chunks, args.ExecutionKey+":"+string(args.Output))
	return nil
}

func (p *recordingProcessor) Process(args *ProcessorArgs) *types.Execution {
	args.Execution.Output = []byte("processed")
	return args.Execution
}

// newTestProcessorClient serves processor over net/rpc the way go-plugin does, and returns a client to it.
func newTestProcessorClient(t *testing.T, processor Processor) *ProcessorClient {
	t.Helper()
	server := rpc.NewServer()
	require.NoError(t, server.RegisterName("Plugin", &ProcessorServer{Processor: processor}))

	serverConn, clientConn := net.Pipe()
	go server.ServeConn(serverConn)
	client := rpc.NewClient(clientConn)
	t.Cleanup(func() { client.Close() })

	return &ProcessorClient{Client: client}
}

func TestProcessorClient_StreamOverRPC(t *testing.T) {
	processor := &recordingProcessor{}
	client := newTestProcessorClient(t, processor)

	ex := &types.Execution{JobName: "job", NodeName: "node", StartedAt: timestamppb.Now()}
	key := ex.Key()

	require.NoError(t, client.Begin(&ProcessorArgs{Execution: ex, Config: Config{"a": "b"}}))
	require.NoError(t, client.Write(&ProcessorOutputArgs{JobName: "job", ExecutionKey: key, Output: []byte("one")}))
	require.NoError(t, client.Write(&ProcessorOutputArgs{JobName: "job", ExecutionKey: key, Output: []byte("two"), Error: true}))
	out := client.Process(&ProcessorArgs{Execution: ex})

	assert.Equal(t, []string{key}, processor.begun)
	assert.Equal(t, []string{key + ":one", key + ":two"}, processor.chunks)
	assert.Equal(t, "processed", string(out.Output))
}

func TestProcessorClient_WriteError(t *testing.T) {
	client := newTestProcessorClient(t, &recordingProcessor{writeErr: errors.New("disk full")})

	err := client.Write(&ProcessorOutputArgs{ExecutionKey: "k", Output: []byte("x")})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "disk full")
}

func TestNopStream(t *testing.T) {
	var p NopStream
	assert.NoError(t, p.Begin(&ProcessorArgs{}))
	assert.NoError(t, p.Write(&ProcessorOutputArgs{}))
}
