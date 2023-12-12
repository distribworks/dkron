package main

import (
	"context"
	"github.com/distribworks/dkron/v3/builtin/bins/dkron-executor-grpc/test"
	dktypes "github.com/distribworks/dkron/v3/plugin/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"testing"
)

type DemoServer struct {
	test.UnimplementedTestServiceServer
}

func (d DemoServer) Test(_ context.Context, request *test.TestRequest) (*test.TestRequest, error) {
	return request, nil
}

func serverSetup() *grpc.Server {
	lis, _ := net.Listen("tcp", ":9000")
	grpcServer := grpc.NewServer()

	d := &DemoServer{}

	test.RegisterTestServiceServer(grpcServer, d)
	reflection.Register(grpcServer)
	go func() {
		grpcServer.Serve(lis)
	}()

	return grpcServer
}

func TestGRPC_ExecuteImpl(t *testing.T) {
	type args struct {
		config map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "It passes if correct data is provided",
			args: args{
				config: map[string]string{
					"url":  "127.0.0.1:9000/test.TestService/Test",
					"body": `{"body": "test"}`,
				},
			},
			wantErr: false,
		},
		{
			name: "it fails if bad address is provided",
			args: args{
				config: map[string]string{
					"url":  "127.0.0.1:9000",
					"body": `{"body": "test"}`,
				},
			},
			wantErr: true,
		},
		{
			name: "it fails if service didn't returned expected code",
			args: args{
				config: map[string]string{
					"url":        "127.0.0.1:9000/test.TestService/Test",
					"body":       `{"body": "test"}`,
					"expectCode": "1",
				},
			},
			wantErr: true,
		},
	}

	srv := serverSetup()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GRPC{}
			_, err := g.ExecuteImpl(&dktypes.ExecuteRequest{Config: tt.args.config})
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteImpl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

	srv.Stop()
}
