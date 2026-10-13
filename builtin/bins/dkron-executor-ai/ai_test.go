package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	dktypes "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	"github.com/stretchr/testify/assert"
)

func TestAIExecute_MissingProvider(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"prompt": "Hello",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "provider is required")
}

func TestAIExecute_MissingPrompt(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "openai",
			"apiKey":   "test-key",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "prompt is required")
}

func TestAIExecute_MissingAPIKey(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "openai",
			"prompt":   "Hello",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "apiKey is required")
}

func TestAIExecute_UnsupportedProvider(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "unsupported",
			"apiKey":   "test-key",
			"prompt":   "Hello",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "unsupported provider")
}

func TestAIExecute_InvalidTimeout(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "openai",
			"apiKey":   "test-key",
			"prompt":   "Hello",
			"timeout":  "invalid",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "invalid timeout value")
}

func TestAIExecute_NegativeTimeout(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "openai",
			"apiKey":   "test-key",
			"prompt":   "Hello",
			"timeout":  "-5",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "timeout must be a positive integer")
}

func TestAIExecute_ZeroTimeout(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "openai",
			"apiKey":   "test-key",
			"prompt":   "Hello",
			"timeout":  "0",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "timeout must be a positive integer")
}

func TestAIExecute_InvalidMaxTokens(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider":  "openai",
			"apiKey":    "test-key",
			"prompt":    "Hello",
			"maxTokens": "invalid",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "invalid maxTokens value")
}

func TestAIExecute_NegativeMaxTokens(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider":  "openai",
			"apiKey":    "test-key",
			"prompt":    "Hello",
			"maxTokens": "-10",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "maxTokens must be a positive integer")
}

func TestAIExecute_ZeroMaxTokens(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider":  "openai",
			"apiKey":    "test-key",
			"prompt":    "Hello",
			"maxTokens": "0",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "maxTokens must be a positive integer")
}

func TestAIExecute_InvalidTemperature(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider":    "openai",
			"apiKey":      "test-key",
			"prompt":      "Hello",
			"temperature": "invalid",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "invalid temperature value")
}

func TestAIExecute_NegativeTemperature(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider":    "openai",
			"apiKey":      "test-key",
			"prompt":      "Hello",
			"temperature": "-0.5",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "temperature must be between 0.0 and 2.0")
}

func TestAIExecute_TemperatureTooHigh(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider":    "openai",
			"apiKey":      "test-key",
			"prompt":      "Hello",
			"temperature": "2.5",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "temperature must be between 0.0 and 2.0")
}

func TestAIExecute_LocalProviderMissingBaseUrl(t *testing.T) {
	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "local",
			"prompt":   "Hello",
			"model":    "llama2",
		},
	}
	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "baseUrl is required for local provider")
}

func TestAIExecute_OpenAISuccess(t *testing.T) {
	// Create a mock server that returns a successful OpenAI response
	mockResponse := map[string]interface{}{
		"choices": []map[string]interface{}{
			{
				"message": map[string]string{
					"content": "Hello! I'm doing well, thank you for asking.",
				},
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/chat/completions", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "openai",
			"apiKey":   "test-api-key",
			"prompt":   "Hello, how are you?",
			"baseUrl":  ts.URL,
			"model":    "gpt-4o-mini",
		},
	}

	ai := &AI{}
	output, err := ai.Execute(pa, nil)
	assert.Nil(t, err)
	assert.Equal(t, "", output.Error)
	assert.Contains(t, string(output.Output), "Hello! I'm doing well")
}

func TestAIExecute_OpenAIError(t *testing.T) {
	// Create a mock server that returns an OpenAI error response
	mockResponse := map[string]interface{}{
		"error": map[string]string{
			"message": "Invalid API key",
			"type":    "invalid_request_error",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "openai",
			"apiKey":   "invalid-key",
			"prompt":   "Hello",
			"baseUrl":  ts.URL,
		},
	}

	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "Invalid API key")
}

func TestAIExecute_AnthropicSuccess(t *testing.T) {
	// Create a mock server that returns a successful Anthropic response
	mockResponse := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "Hello! I'm Claude, nice to meet you.",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/messages", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "test-api-key", r.Header.Get("x-api-key"))
		assert.Equal(t, "2023-06-01", r.Header.Get("anthropic-version"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "anthropic",
			"apiKey":   "test-api-key",
			"prompt":   "Hello",
			"baseUrl":  ts.URL,
			"model":    "claude-3-haiku-20240307",
		},
	}

	ai := &AI{}
	output, err := ai.Execute(pa, nil)
	assert.Nil(t, err)
	assert.Equal(t, "", output.Error)
	assert.Contains(t, string(output.Output), "Hello! I'm Claude")
}

func TestAIExecute_AnthropicError(t *testing.T) {
	// Create a mock server that returns an Anthropic error response
	mockResponse := map[string]interface{}{
		"error": map[string]string{
			"type":    "authentication_error",
			"message": "Invalid API key",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "anthropic",
			"apiKey":   "invalid-key",
			"prompt":   "Hello",
			"baseUrl":  ts.URL,
		},
	}

	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "Invalid API key")
}

func TestAIExecute_LocalProvider(t *testing.T) {
	// Create a mock server that returns a successful OpenAI-compatible response
	mockResponse := map[string]interface{}{
		"choices": []map[string]interface{}{
			{
				"message": map[string]string{
					"content": "Response from local model",
				},
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "local",
			"prompt":   "Hello",
			"baseUrl":  ts.URL,
			"model":    "llama2",
		},
	}

	ai := &AI{}
	output, err := ai.Execute(pa, nil)
	assert.Nil(t, err)
	assert.Equal(t, "", output.Error)
	assert.Contains(t, string(output.Output), "Response from local model")
}

func TestAIExecute_OpenAIEmptyChoices(t *testing.T) {
	// Create a mock server that returns empty choices
	mockResponse := map[string]interface{}{
		"choices": []map[string]interface{}{},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "openai",
			"apiKey":   "test-key",
			"prompt":   "Hello",
			"baseUrl":  ts.URL,
		},
	}

	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "no response choices")
}

func TestAIExecute_AnthropicEmptyContent(t *testing.T) {
	// Create a mock server that returns empty content
	mockResponse := map[string]interface{}{
		"content": []map[string]interface{}{},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer ts.Close()

	pa := &dktypes.ExecuteRequest{
		JobName: "testJob",
		Config: map[string]string{
			"provider": "anthropic",
			"apiKey":   "test-key",
			"prompt":   "Hello",
			"baseUrl":  ts.URL,
		},
	}

	ai := &AI{}
	output, _ := ai.Execute(pa, nil)
	assert.Contains(t, output.Error, "no response content")
}
