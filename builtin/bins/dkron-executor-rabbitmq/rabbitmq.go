package main

import (
	"encoding/base64"
	"errors"

	dkplugin "github.com/distribworks/dkron/v3/plugin"
	dktypes "github.com/distribworks/dkron/v3/plugin/types"
	"github.com/streadway/amqp"
)

// RabbitMQ process publish rabbitmq message when Execute method is called.
type RabbitMQ struct {
}

// Execute method of the plugin
// "executor": "rabbitmq",
// "executor_config": {
//     "url": "amqp://guest:guest@localhost:5672/", // rabbitmq server url
//     "text": "hello world!",                  				// or "base64" to send bytes as rabbitmq message
//     "queue": "test",             				//
// }
func (r *RabbitMQ) Execute(args *dktypes.ExecuteRequest, cb dkplugin.StatusHelper) (*dktypes.ExecuteResponse, error) {
	out, err := r.ExecuteImpl(args, cb)
	resp := &dktypes.ExecuteResponse{Output: out}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

// ExecuteImpl do rabbitmq publish
func (r *RabbitMQ) ExecuteImpl(args *dktypes.ExecuteRequest, cb dkplugin.StatusHelper) ([]byte, error) {

	if args.Config["url"] == "" {
		return nil, errors.New("url is empty")
	}

	if args.Config["queue"] == "" {
		return nil, errors.New("queue is empty")
	}

	// broker := "amqp://guest:guest@localhost:5672/"
	broker := args.Config["url"]
	conn, err := amqp.Dial(broker)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	queue := args.Config["queue"]
	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	var body []byte
	b64, ok := args.Config["base64"]
	if ok {
		decoded, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return nil, err
		}
		body = decoded
	} else {
		text := args.Config["text"]
		body = []byte(text)
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
