package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	dkplugin "github.com/distribworks/dkron/v3/plugin"
	dktypes "github.com/distribworks/dkron/v3/plugin/types"
)

// GCPPubSub plugin publish message to topic when Execute method is called.
type GCPPubSub struct {
}

const (
	configProjectName    = "project"
	configTopicName      = "topic"
	configDataName       = "data"
	configAttributesName = "attributes"
)

// Execute Process method of the plugin
// "executor": "gcppubsub",
// "executor_config": {
//  "project": "project-id",
//  "topic": "topic-name",
//  "data": "aGVsbG8gd29ybGQ=" // Optional. base64 encoded data
//  "attributes": "{\"hello\":\"world\",\"waka\":\"paka\"}" // JSON serialized attributes
// }
func (g *GCPPubSub) Execute(args *dktypes.ExecuteRequest, _ dkplugin.StatusHelper) (*dktypes.ExecuteResponse, error) {
	out, err := g.ExecuteImpl(args)
	resp := &dktypes.ExecuteResponse{Output: out}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

// ExecuteImpl do http request
func (g *GCPPubSub) ExecuteImpl(args *dktypes.ExecuteRequest) ([]byte, error) {
	ctx := context.Background()
	projectID := args.Config[configProjectName]
	topicName := args.Config[configTopicName]

	if projectID == "" {
		return nil, fmt.Errorf("missing project")
	}

	if topicName == "" {
		return nil, fmt.Errorf("missing topic")
	}

	msg, err := configToPubSubMessage(args.Config)
	if err != nil {
		return nil, fmt.Errorf("convert config to pubsub message: %w", err)
	}

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("create client for project %q: %w", projectID, err)
	}
	defer func() {
		_ = client.Close()
	}()

	topic := client.Topic(topicName)
	res := topic.Publish(ctx, msg)
	serverID, err := res.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("publish message to topic %q: %w", topicName, err)
	}
	return []byte(serverID), nil
}

func configToPubSubMessage(config map[string]string) (*pubsub.Message, error) {
	if config == nil {
		return nil, fmt.Errorf("invalid config")
	}

	encodedData := config[configDataName]
	attributesJSON := config[configAttributesName]

	if attributesJSON == "" && encodedData == "" {
		return nil, fmt.Errorf("at least one of these fields should be set 'attributes, data'")
	}

	msg := &pubsub.Message{}

	var attributes map[string]string
	if attributesJSON != "" {
		if err := json.Unmarshal([]byte(attributesJSON), &attributes); err != nil  {
			return nil, fmt.Errorf("invalid attributes JSON: %w", err)
		}
		msg.Attributes = attributes
	}

	if encodedData != "" {
		data, err := base64.StdEncoding.DecodeString(encodedData)
		if err != nil {
			return nil, fmt.Errorf("invalid encoded data: %w", err)
		}
		msg.Data = data
	}

	return msg, nil
}
