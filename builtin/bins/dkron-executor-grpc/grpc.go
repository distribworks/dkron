package main

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/armon/circbuf"
	dkplugin "github.com/distribworks/dkron/v3/plugin"
	dktypes "github.com/distribworks/dkron/v3/plugin/types"
	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

const (
	defaultTimeout = 30
	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 256000
)

type GRPC struct{}

// Execute Process method of the plugin
// "executor": "grpc",
// "executor_config": {
//     "url": "127.0.0.1:9000/demo.DemoService/Demo", // Request url
//     "body": "",                                    // POST body
//     "timeout": "30",                               // Request timeout, unit seconds
//     "expectCode": "0",                             // Expect response code, any of the described here https://grpc.github.io/grpc/core/md_doc_statuscodes.html
// }
func (g *GRPC) Execute(args *dktypes.ExecuteRequest, _ dkplugin.StatusHelper) (*dktypes.ExecuteResponse, error) {
	out, err := g.ExecuteImpl(args)
	resp := &dktypes.ExecuteResponse{Output: out}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

// ExecuteImpl do grpc request
func (g *GRPC) ExecuteImpl(args *dktypes.ExecuteRequest) ([]byte, error) {
	output, _ := circbuf.NewBuffer(maxBufSize)

	if args.Config["url"] == "" {
		return output.Bytes(), errors.New("url is empty")
	}

	segments := strings.Split(args.Config["url"], "/")
	if len(segments) < 2 {
		return output.Bytes(), errors.New("we require at least a url and a path to do a proto request")
	}

	var timeout int64 = defaultTimeout
	if args.Config["timeout"] != "" {
		t, convErr := strconv.ParseInt(args.Config["timeout"], 10, 64)
		if convErr != nil {
			return output.Bytes(), errors.New("Invalid timeout")
		}

		timeout = t
	}

	expectedStatusCode := codes.OK
	if args.Config["expectCode"] != "" {
		t, convErr := strconv.ParseUint(args.Config["expectCode"], 10, 32)
		if convErr != nil {
			return output.Bytes(), errors.New("Invalid timeout")
		}

		expectedStatusCode = codes.Code(t)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	body := strings.NewReader(args.Config["body"])

	var descSource grpcurl.DescriptorSource
	var refClient *grpcreflect.Client
	var cc *grpc.ClientConn
	var opts []grpc.DialOption
	var creds credentials.TransportCredentials

	cc, grpcDialErr := grpcurl.BlockingDial(ctx, "tcp", segments[0], creds, opts...)
	if grpcDialErr != nil {
		return output.Bytes(), grpcDialErr
	}
	defer cc.Close() // nolint:errcheck

	md := grpcurl.MetadataFromHeaders([]string{})
	refCtx := metadata.NewOutgoingContext(ctx, md)
	refClient = grpcreflect.NewClient(refCtx, reflectpb.NewServerReflectionClient(cc))
	descSource = grpcurl.DescriptorSourceFromServer(ctx, refClient)

	rf, formatter, refErr := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, descSource, body, grpcurl.FormatOptions{})
	if refErr != nil {
		return output.Bytes(), errors.Wrap(refErr, "Failed querying reflection server")
	}

	var b []byte
	out := bytes.NewBuffer(b)

	h := &grpcurl.DefaultEventHandler{
		Out:       out,
		Formatter: formatter,
	}

	rpcCallErr := grpcurl.InvokeRPC(ctx, descSource, cc, strings.Join(segments[1:], "/"), []string{}, h, rf.Next)
	if rpcCallErr != nil {
		return output.Bytes(), errors.Wrap(rpcCallErr, "Failed querying reflection server")
	}

	if h.Status.Code() != expectedStatusCode {
		return output.Bytes(), fmt.Errorf("server returned %v code, expected %v", h.Status.Code(), expectedStatusCode)
	}

	return output.Bytes(), nil
}
