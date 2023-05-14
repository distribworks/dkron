package main

import (
	"testing"

	"cloud.google.com/go/pubsub"
	dktypes "github.com/distribworks/dkron/v3/plugin/types"
	"github.com/stretchr/testify/assert"
)

func TestConfigToPubSubMessage(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]string
		expected *pubsub.Message
		wantErr  bool
	}{
		{
			name: "missing both data and attributes",
			config: map[string]string{
				"other": "key",
			},
			wantErr: true,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "only data",
			config: map[string]string{
				"data": "aGVsbG8gd29ybGQ=",
			},
			expected: &pubsub.Message{Data: []byte("hello world")},
		},
		{
			name: "only attributes",
			config: map[string]string{
				"attributes": "{\"hello\":\"world\",\"waka\":\"paka\"}",
			},
			expected: &pubsub.Message{
				Attributes: map[string]string{"hello": "world", "waka": "paka"},
			},
		},
		{
			name: "attributes and data",
			config: map[string]string{
				"data":       "aGVsbG8gd29ybGQ=",
				"attributes": "{\"hello\":\"world\",\"waka\":\"paka\"}",
			},
			expected: &pubsub.Message{
				Data:       []byte("hello world"),
				Attributes: map[string]string{"hello": "world", "waka": "paka"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := configToPubSubMessage(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.EqualValues(t, tt.expected.Attributes, got.Attributes)
				assert.Equal(t, tt.expected.Data, got.Data)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"topic":         "test",
			"brokerAddress": "testaddress",
			"message":       "{\"hello\":11}",
			"debug":         "true",
		},
	}
	output, err := (&GCPPubSub{}).Execute(pa, nil)
	assert.NoError(t, err)
	assert.NotNil(t, output)
	t.Log(string(output.Output))
}
