package main

import (
	"errors"
	"log"

	"github.com/armon/circbuf"
	"github.com/distribworks/dkron/v2/dkron"
	"github.com/nats-io/nats.go"
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
//     "url": "http://example.com", // nats server url
//     "message": "",                  //
//     "topic": "publishTopic",             //
//     "userName":"test@hbh.dfg",
//     "password":"dfdffs"
// }
func (s *Nats) Execute(args *dkron.ExecuteRequest) (*dkron.ExecuteResponse, error) {

	out, err := s.ExecuteImpl(args)
	resp := &dkron.ExecuteResponse{Output: out}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

// ExecuteImpl do http request
func (s *Nats) ExecuteImpl(args *dkron.ExecuteRequest) ([]byte, error) {

	output, _ := circbuf.NewBuffer(maxBufSize)

	var debug bool
	if args.Config["debug"] != "" {
		debug = true
		log.Printf("config  %#v\n\n", args.Config)
	}

	if args.Config["url"] == "" {

		return output.Bytes(), errors.New("url is empty")
	}

	if args.Config["topic"] == "" {
		return output.Bytes(), errors.New("topic is empty")
	}
	nc, err := nats.Connect(args.Config["url"], nats.UserInfo(args.Config["userName"], args.Config["password"]))

	if err != nil {
		return output.Bytes(), errors.New("Error At Nats Connection")
	}

	nc.Publish(args.Config["topic"], []byte(args.Config["message"]))

	output.Write([]byte("Result: success to publish data\n"))

	if debug {
		log.Printf("request  %#v\n\n", nc)
	}

	return output.Bytes(), nil
}
