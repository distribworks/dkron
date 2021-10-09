package main

import (
	"fmt"
	"testing"

	dktypes "github.com/distribworks/dkron/v3/plugin/types"
)

func TestProduceExecuteWithKey(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJobWithKey",
		Config: map[string]string{
			"topic":         "test",
			"brokerAddress": "testaddress",
			"key":           "testkey",
			"message":       "{\"hello\":11}",
			"debug":         "true",
		},
	}
	kafka := &Kafka{}
	output, err := kafka.Execute(pa, nil)
	fmt.Println(string(output.Output))
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestProduceExecuteWithoutKey(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJobWithoutKey",
		Config: map[string]string{
			"topic":         "test",
			"brokerAddress": "testaddress",
			"message":       "{\"hello\":11}",
			"debug":         "true",
		},
	}
	kafka := &Kafka{}
	output, err := kafka.Execute(pa, nil)
	fmt.Println(string(output.Output))
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
	}
}