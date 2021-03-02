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
			"subject": "opcuaReadRequest",
			"url":     "localhost:4222",
			"message": "{\"hello\":11}",
			"debug":   "true",
		},
	}
	nats := &Nats{}
	output, err := nats.Execute(pa, nil)
	fmt.Println(string(output.Output))
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
	}
}
