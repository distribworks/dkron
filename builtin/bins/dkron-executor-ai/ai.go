package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/armon/circbuf"

	dktypes "github.com/distribworks/dkron/v4/gen/proto/types/v1"
	dkplugin "github.com/distribworks/dkron/v4/plugin"
)

const (
	// maxBufSize limits how much data we collect from a handler.
	// This is to prevent Serf's memory from growing to an enormous
	// amount due to a faulty handler.
	maxBufSize = 500000

	// defaultTimeout is the default HTTP timeout in seconds
	defaultTimeout = 120

	// default model for OpenAI
	defaultOpenAIModel = "gpt-4o-mini"

	// default model for Anthropic
	defaultAnthropicModel = "claude-3-haiku-20240307"
)

// AI executor runs prompts against AI model providers
type AI struct{}

// Execute Process method of the plugin
// "executor": "ai",
//
//	"executor_config": {
//	    "provider": "openai",                         // Provider: openai, anthropic, or local
//	    "apiKey": "sk-...",                           // API key for the provider
//	    "model": "gpt-4o-mini",                       // Model to use
//	    "prompt": "Hello, how are you?",              // The prompt to send
//	    "baseUrl": "",                                // Optional: Custom base URL for local/compatible models
//	    "maxTokens": "1000",                          // Optional: Maximum tokens in response
//	    "temperature": "0.7",                         // Optional: Temperature for response randomness
//	    "timeout": "120",                             // Optional: Request timeout in seconds
//	    "debug": "true"                               // Optional: Enable debug logging
//	}
func (a *AI) Execute(args *dktypes.ExecuteRequest, cb dkplugin.StatusHelper) (*dktypes.ExecuteResponse, error) {
	out, err := a.ExecuteImpl(args)
	resp := &dktypes.ExecuteResponse{Output: out}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, nil
}

// ExecuteImpl executes the AI prompt request
func (a *AI) ExecuteImpl(args *dktypes.ExecuteRequest) ([]byte, error) {
	output, _ := circbuf.NewBuffer(maxBufSize)

	var debug bool
	if args.Config["debug"] != "" {
		debug = true
		// Log config without sensitive API key
		safeConfig := make(map[string]string)
		for k, v := range args.Config {
			if k == "apiKey" {
				safeConfig[k] = "[REDACTED]"
			} else {
				safeConfig[k] = v
			}
		}
		log.Printf("AI executor config: %#v\n", safeConfig)
	}

	// Validate required parameters
	provider := strings.ToLower(args.Config["provider"])
	if provider == "" {
		return output.Bytes(), errors.New("provider is required (openai, anthropic, or local)")
	}

	prompt := args.Config["prompt"]
	if prompt == "" {
		return output.Bytes(), errors.New("prompt is required")
	}

	apiKey := args.Config["apiKey"]
	// API key is required for openai and anthropic, optional for local
	if provider != "local" && apiKey == "" {
		return output.Bytes(), fmt.Errorf("apiKey is required for provider %s", provider)
	}

	// Parse optional parameters
	timeout := defaultTimeout
	if args.Config["timeout"] != "" {
		t, err := strconv.Atoi(args.Config["timeout"])
		if err != nil {
			return output.Bytes(), fmt.Errorf("invalid timeout value: %s", args.Config["timeout"])
		}
		if t <= 0 {
			return output.Bytes(), fmt.Errorf("timeout must be a positive integer, got: %d", t)
		}
		timeout = t
	}

	maxTokens := 1000
	if args.Config["maxTokens"] != "" {
		mt, err := strconv.Atoi(args.Config["maxTokens"])
		if err != nil {
			return output.Bytes(), fmt.Errorf("invalid maxTokens value: %s", args.Config["maxTokens"])
		}
		if mt <= 0 {
			return output.Bytes(), fmt.Errorf("maxTokens must be a positive integer, got: %d", mt)
		}
		maxTokens = mt
	}

	temperature := 0.7
	if args.Config["temperature"] != "" {
		t, err := strconv.ParseFloat(args.Config["temperature"], 64)
		if err != nil {
			return output.Bytes(), fmt.Errorf("invalid temperature value: %s", args.Config["temperature"])
		}
		if t < 0.0 || t > 2.0 {
			return output.Bytes(), fmt.Errorf("temperature must be between 0.0 and 2.0, got: %f", t)
		}
		temperature = t
	}

	var result string
	var err error

	switch provider {
	case "openai":
		result, err = a.executeOpenAI(args.Config, prompt, apiKey, maxTokens, temperature, timeout, debug)
	case "local":
		// For local provider, baseUrl is required to avoid accidentally sending data to OpenAI
		if args.Config["baseUrl"] == "" {
			return output.Bytes(), errors.New("baseUrl is required for local provider")
		}
		result, err = a.executeOpenAI(args.Config, prompt, apiKey, maxTokens, temperature, timeout, debug)
	case "anthropic":
		result, err = a.executeAnthropic(args.Config, prompt, apiKey, maxTokens, temperature, timeout, debug)
	default:
		return output.Bytes(), fmt.Errorf("unsupported provider: %s (supported: openai, anthropic, local)", provider)
	}

	if err != nil {
		return output.Bytes(), err
	}

	output.Write([]byte(result))
	return output.Bytes(), nil
}

// OpenAI API structures
type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// executeOpenAI handles OpenAI and OpenAI-compatible API requests
func (a *AI) executeOpenAI(config map[string]string, prompt, apiKey string, maxTokens int, temperature float64, timeout int, debug bool) (string, error) {
	model := config["model"]
	if model == "" {
		model = defaultOpenAIModel
	}

	baseURL := config["baseUrl"]
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	reqBody := openAIRequest{
		Model: model,
		Messages: []openAIMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	if debug {
		log.Printf("OpenAI request URL: %s/chat/completions\n", baseURL)
		log.Printf("OpenAI request body: %s\n", string(jsonBody))
	}

	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if debug {
		log.Printf("OpenAI response status: %d\n", resp.StatusCode)
		log.Printf("OpenAI response body: %s\n", string(body))
	}

	var openAIResp openAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		// If JSON parsing fails and we have a non-2xx status, report the status
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("API error: %s", openAIResp.Error.Message)
	}

	// Check for non-2xx status codes that didn't include an error object
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	if len(openAIResp.Choices) == 0 {
		return "", errors.New("no response choices returned from API")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// Anthropic API structures
type anthropicRequest struct {
	Model       string             `json:"model"`
	Messages    []anthropicMessage `json:"messages"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature float64            `json:"temperature,omitempty"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// executeAnthropic handles Anthropic Claude API requests
func (a *AI) executeAnthropic(config map[string]string, prompt, apiKey string, maxTokens int, temperature float64, timeout int, debug bool) (string, error) {
	model := config["model"]
	if model == "" {
		model = defaultAnthropicModel
	}

	baseURL := config["baseUrl"]
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}
	baseURL = strings.TrimSuffix(baseURL, "/")

	reqBody := anthropicRequest{
		Model: model,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	if debug {
		log.Printf("Anthropic request URL: %s/messages\n", baseURL)
		log.Printf("Anthropic request body: %s\n", string(jsonBody))
	}

	req, err := http.NewRequest("POST", baseURL+"/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if debug {
		log.Printf("Anthropic response status: %d\n", resp.StatusCode)
		log.Printf("Anthropic response body: %s\n", string(body))
	}

	var anthropicResp anthropicResponse
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		// If JSON parsing fails and we have a non-2xx status, report the status
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if anthropicResp.Error != nil {
		return "", fmt.Errorf("API error: %s", anthropicResp.Error.Message)
	}

	// Check for non-2xx status codes that didn't include an error object
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	if len(anthropicResp.Content) == 0 {
		return "", errors.New("no response content returned from API")
	}

	// Combine all text content
	var result strings.Builder
	for i, content := range anthropicResp.Content {
		if content.Type == "text" {
			if i > 0 {
				result.WriteString("\n")
			}
			result.WriteString(content.Text)
		}
	}

	return result.String(), nil
}
