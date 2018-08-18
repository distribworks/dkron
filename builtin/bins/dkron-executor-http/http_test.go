package main

import (
	"fmt"
	"testing"

	"github.com/victorcoder/dkron/dkron"
)

func TestExecute(t *testing.T) {
	pa := &dkron.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"method":     "GET",
			"url":        "https://webhook.site/2a570753-499b-4c08-a8d5-d20f9b625ea8",
			"expectCode": "200",
			"debug":      "true",
		},
	}
	http := &HTTP{}
	output, err := http.Execute(pa)
	fmt.Println(string(output))
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecutePost(t *testing.T) {
	pa := &dkron.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"method":     "POST",
			"url":        "https://webhook.site/2a570753-499b-4c08-a8d5-d20f9b625ea8",
			"body":       "{\"hello\":11}",
			"headers":    "[\"Content-Type:application/json\"]",
			"expectCode": "200",
			"expectBody": "",
			"debug":      "true",
		},
	}
	http := &HTTP{}
	output, err := http.Execute(pa)
	fmt.Println(string(output))
	fmt.Println(err)
	if err != nil {
		t.Fatal(err)
	}
}
