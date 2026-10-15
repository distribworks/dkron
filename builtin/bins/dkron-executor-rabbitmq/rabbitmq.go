package main

import (
	"encoding/base64"
	"errors"
	"strconv"

	dktypes "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	dkplugin "github.com/distribworks/dkron/v4/plugin"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

// RabbitMQ process publish rabbitmq message when Execute method is called.
type RabbitMQ struct{}

func (r *RabbitMQ) ConfigSchema() (string, error) {
	return `{
  "type": "object",
  "required": ["url", "queue.name"],
  "properties": {
    "url": {"type": "string", "title": "URL"},
    "exchange": {"type": "string", "title": "Exchange"},
    "queue.name": {"type": "string", "title": "Queue name"},
    "queue.create": {"type": "string", "title": "Create queue", "enum": ["false", "true"], "default": "false"},
    "queue.durable": {"type": "string", "title": "Durable queue", "enum": ["false", "true"], "default": "false"},
    "queue.auto_delete": {"type": "string", "title": "Auto-delete queue", "enum": ["false", "true"], "default": "false"},
    "queue.exclusive": {"type": "string", "title": "Exclusive queue", "enum": ["false", "true"], "default": "false"},
    "message.content_type": {"type": "string", "title": "Content type"},
    "message.delivery_mode": {"type": "string", "title": "Delivery mode"},
    "message.messageId": {"type": "string", "title": "Message ID"},
    "message.body": {"type": "string", "title": "Message body"},
    "message.base64Body": {"type": "string", "title": "Base64 message body"}
  }
}`, nil
}

// Execute method of the plugin
// "executor": "rabbitmq",
//
//	"executor_config": {
//			"url": "amqp://guest:guest@localhost:5672/",
//			"exchange": "amq.default",
//			"queue.name": "test",
//			"queue.create": "true",
//			"queue.durable": "true",
//			"queue.auto_delete": "false",
//			"queue.exclusive": "false",
//			"message.content_type": "application/json",
//			"message.delivery_mode": "2",
//			"message.messageId": "4373732772",
//			"message.body": "{\"key\":\"value\"}"
//			"message.base64Body": "base64encodedBody"
//	}
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
	// validate config
	cfg := args.Config
	if cfg == nil {
		return nil, errors.New("RabbitMQ config is empty")
	}

	url := cfg["url"]
	if url == "" {
		return nil, errors.New("RabbitMQ url is empty")
	}

	queueName := cfg["queue.name"]
	if queueName == "" {
		return nil, errors.New("RabbitMQ queue name is empty")
	}

	if cfg["message.body"] != "" && cfg["message.base64Body"] != "" {
		return nil, errors.New("RabbitMQ message.body and message.base64Body are both set")
	}

	// establish connection
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Error("Failed to close amqp connection", log.WithError(err))
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.Error("Failed to close channel", log.WithError(err))
		}
	}(ch)

	// create queue if necessary
	if err := createQueueIfNecessary(cfg, queueName, ch); err != nil {
		return nil, err
	}

	// publish message
	if err = publish(cfg, ch); err != nil {
		return nil, err
	}
	return nil, nil
}

func createQueueIfNecessary(cfg map[string]string, queue string, ch *amqp.Channel) error {
	if val, ok := cfg["queue.create"]; !ok || (ok && val == "false") {
		return nil
	}

	durable, _ := strconv.ParseBool(cfg["queue.durable"])
	autoDelete, _ := strconv.ParseBool(cfg["queue.auto_delete"])
	exclusive, _ := strconv.ParseBool(cfg["queue.exclusive"])

	_, err := ch.QueueDeclare(
		queue,
		durable,
		autoDelete,
		exclusive,
		false,
		nil,
	)

	return err
}

func publish(cfg map[string]string, ch *amqp.Channel) error {
	var body []byte
	b64, ok := cfg["message.base64Body"]
	if ok {
		decoded, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return err
		}
		body = decoded
	} else {
		stringBody := cfg["message.body"]
		body = []byte(stringBody)
	}

	contentType := cfg["message.content_type"]
	if contentType == "" {
		contentType = "text/plain"
	}
	messageId := cfg["message.messageId"]
	rawDeliveryMode := cfg["message.delivery_mode"]
	if rawDeliveryMode == "" {
		rawDeliveryMode = "0"
	}
	deliveryMode, err := strconv.ParseUint(rawDeliveryMode, 10, 8)
	if err != nil {
		return err
	}
	exchange, ok := cfg["exchange"]
	if !ok {
		exchange = ""
	}
	return ch.Publish(
		exchange,          // exchange
		cfg["queue.name"], // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:  contentType,
			Body:         body,
			MessageId:    messageId,
			DeliveryMode: uint8(deliveryMode),
		})
}
