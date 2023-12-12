package main

import (
	"errors"
	"log"

	"github.com/armon/circbuf"
	"github.com/nats-io/nats.go"

	dkplugin "github.com/distribworks/dkron/v3/plugin"
	dktypes "github.com/distribworks/dkron/v3/plugin/types"
)

const (
	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 500000
)

// Nats process http request
type Nats struct {
}

// Execute Process method of the plugin
// "executor": "nats",
// "executor_config": {
//     "url": "tls://nats.demo.io:4443", // nats server url
//     "message": "",
//     "subject": "Subject",
//     "userName":"test@hbh.dfg",
//     "password":"dfdffs"
// }
func (s *Nats) Execute(args *dktypes.ExecuteRequest, cb dkplugin.StatusHelper) (*dktypes.ExecuteResponse, error) {

	out, err := s.ExecuteImpl(args)
	resp := &dktypes.ExecuteResponse{Output: out}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

// ExecuteImpl do http request
func (s *Nats) ExecuteImpl(args *dktypes.ExecuteRequest) ([]byte, error) {

	output, _ := circbuf.NewBuffer(maxBufSize)

	var debug bool
	if args.Config["debug"] != "" {
		debug = true
		log.Printf("config  %#v\n\n", args.Config)
	}

	if args.Config["url"] == "" {

		return output.Bytes(), errors.New("url is empty")
	}

	if args.Config["subject"] == "" {
		return output.Bytes(), errors.New("subject is empty")
	}
	nc, err := nats.Connect(args.Config["url"], nats.UserInfo(args.Config["userName"], args.Config["password"]))

	if err != nil {
		return output.Bytes(), errors.New("error connecting to NATS")
	}

	nc.Publish(args.Config["subject"], []byte(args.Config["message"]))

	output.Write([]byte("Result: Message successfully sent\n"))

	if debug {
		log.Printf("request  %#v\n\n", nc)
	}
	
	if nc.IsConnected() {
		defer nc.Flush()
		defer nc.Close()
	}

	return output.Bytes(), nil
}
