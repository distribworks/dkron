package main

import (
	"fmt"
	"testing"

	dktypes "github.com/distribworks/dkron/v3/plugin/types"
)

func TestPublishExecute(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"queue": "test",
			"url":   "amqp://guest:guest@localhost:5672",
			"text":  "{\"hello\":11}",
		},
	}
	rabbitmq := &RabbitMQ{}
	output, err := rabbitmq.Execute(pa, nil)
	fmt.Println(string(output.Output))
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
	}
}
