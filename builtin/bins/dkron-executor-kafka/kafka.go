package main

import (
	"errors"
	"log"
	"strings"

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
//     "brokerAddress": "192.168.59.103:9092", // kafka broker url
//     "key": "",
//     "message": "",
//     "topic": "publishTopic"
// }
func (s *Kafka) Execute(args *dktypes.ExecuteRequest, cb dkplugin.StatusHelper) (*dktypes.ExecuteResponse, error) {

	out, err := s.ExecuteImpl(args)
	resp := &dktypes.ExecuteResponse{Output: out}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

// ExecuteImpl produce message on Kafka broker
func (s *Kafka) ExecuteImpl(args *dktypes.ExecuteRequest) ([]byte, error) {

	output, _ := circbuf.NewBuffer(maxBufSize)

	var debug bool
	if args.Config["debug"] != "" {
		debug = true
		log.Printf("config  %#v\n\n", args.Config)
	}

	if args.Config["brokerAddress"] == "" {

		return output.Bytes(), errors.New("brokerAddress is empty")
	}

	if args.Config["topic"] == "" {
		return output.Bytes(), errors.New("topic is empty")
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	brokers := strings.Split(args.Config["brokerAddress"], ",")
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		// Should not reach here

		if debug {
			log.Printf("sarama  %#v\n\n", config)
		}
		return output.Bytes(), err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: args.Config["topic"],
		Key:   sarama.StringEncoder(args.Config["key"]),
		Value: sarama.StringEncoder(args.Config["message"]),
	}

	_, _, err = producer.SendMessage(msg)

	if err != nil {
		return output.Bytes(), err
	}

	output.Write([]byte("Result: successfully produced the message on Kafka broker\n"))
	return output.Bytes(), nil
}
