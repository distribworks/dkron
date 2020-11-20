package main

import (
	"errors"
	"log"

	"github.com/Shopify/sarama"
	"github.com/armon/circbuf"

	dkplugin "github.com/distribworks/dkron/v3/plugin"
	dktypes "github.com/distribworks/dkron/v3/plugin/types"
)

const (
	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 500000
)

// Kafka process kafka request
type Kafka struct {
}

// Execute Process method of the plugin
// "executor": "kafka",
// "executor_config": {
//     "url": "http://example.com", // kafka server url
//     "message": "",                  //
//     "topic": "publishTopic",             //
// }
func (s *Kafka) Execute(args *dktypes.ExecuteRequest, cb dkplugin.StatusHelper) (*dktypes.ExecuteResponse, error) {

	out, err := s.ExecuteImpl(args)
	resp := &dktypes.ExecuteResponse{Output: out}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

// ExecuteImpl do http request
func (s *Kafka) ExecuteImpl(args *dktypes.ExecuteRequest) ([]byte, error) {

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
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// brokers := []string{"192.168.59.103:9092"}
	brokers := []string{args.Config["url"]}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		// Should not reach here

		if debug {
			log.Printf("request  %#v\n\n", config)
		}
		return output.Bytes(), err
	}

	topic := args.Config["topic"]
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(args.Config["message"]),
	}

	_, _, err = producer.SendMessage(msg)

	if err != nil {
		return output.Bytes(), err
	}
	defer func() {
		producer.Close()
	}()

	output.Write([]byte("Result: success to publish data\n"))
	return output.Bytes(), nil
}
