package main

import (
	"fmt"
	"testing"

	"github.com/victorcoder/dkron/plugintypes"
)

func TestExecute(t *testing.T) {
	pa := &plugintypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"method":     "GET",
			"url":        "https://httpbin.org/get",
			"expectCode": "200",
			"debug":      "true",
		},
	}
	http := &HTTP{}
	output, err := http.Execute(pa)
	fmt.Println(string(output.Output))
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecutePost(t *testing.T) {
	pa := &plugintypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"method":     "POST",
			"url":        "https://httpbin.org/post",
			"body":       "{\"hello\":11}",
			"headers":    "[\"Content-Type:application/json\"]",
			"expectCode": "200",
			"expectBody": "",
			"debug":      "true",
		},
	}
	http := &HTTP{}
	output, err := http.Execute(pa)
	fmt.Println(string(output.Output))
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
	}
}
