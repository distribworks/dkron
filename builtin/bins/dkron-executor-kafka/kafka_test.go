package main

import (
	"fmt"
	"testing"

	"github.com/distribworks/dkron/v2/dkron"
)

func TestPublishExecute(t *testing.T) {
	pa := &dkron.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"topic":   "test",
			"url":     "tesr",
			"message": "{\"hello\":11}",
			"debug":   "true",
		},
	}
	kafka := &Kafka{}
	output, err := kafka.Execute(pa)
	fmt.Println(string(output.Output))
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
	}
}
