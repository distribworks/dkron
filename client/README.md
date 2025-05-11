# Dkron Go client 

## Installation

```bash
go get github.com/distribworks/dkron/v4/client
```

## Usage

```go
package main

import (
	"context"

	"github.com/distribworks/dkron/v4/client"
	log "github.com/sirupsen/logrus"
)

func main() {
	newClient, err := client.NewClient("http://localhost:9090")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new job
	params := &client.CreateOrUpdateJobParams{
		Runoncreate: nil,
	}
	body := client.CreateOrUpdateJobJSONRequestBody{
		Name:       "testjob1",
		Timezone:   "Europe/Paris",
		Schedule:   "@every 1m",
		Owner:      "test",
		OwnerEmail: "",
		Tags: map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		},
		Retries:  3,
		Executor: "Shell",
		ExecutorConfig: map[string]string{
			"command": "echo hello world!",
		},
		Displayname: "Job",
		Ephemeral:   false,
	}

	_, err = newClient.CreateOrUpdateJob(context.Background(), params, body)
	if err != nil {
		log.WithError(err).Error("Failed to create job")
	}
}

```