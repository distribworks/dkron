# AI Executor

The AI executor runs prompts against AI model providers, enabling scheduled AI-powered tasks. It supports multiple providers including OpenAI, Anthropic (Claude), and local models via OpenAI-compatible APIs (like Ollama, LM Studio, or vLLM).

## Configuration

Params:

```
provider:    Required. The AI provider to use: "openai", "anthropic", or "local"
apiKey:      Required for openai/anthropic. API key for the provider
model:       Optional. Model to use (defaults: gpt-4o-mini for OpenAI, claude-3-haiku-20240307 for Anthropic)
prompt:      Required. The prompt to send to the AI model
baseUrl:     Optional for openai/anthropic, required for local. Custom base URL for the API endpoint
maxTokens:   Optional. Maximum tokens in the response (default: 1000)
temperature: Optional. Temperature for response randomness, 0.0-2.0 (default: 0.7)
timeout:     Optional. Request timeout in seconds (default: 120, must be positive)
debug:       Optional. Enable debug logging when set to any non-empty value
```

## Providers

### OpenAI

Use with OpenAI's API (GPT-4, GPT-4o, GPT-4o-mini, etc.):

```json
{
  "executor": "ai",
  "executor_config": {
    "provider": "openai",
    "apiKey": "sk-...",
    "model": "gpt-4o-mini",
    "prompt": "Summarize the key benefits of distributed job scheduling",
    "maxTokens": "500",
    "temperature": "0.7"
  }
}
```

### Anthropic (Claude)

Use with Anthropic's Claude API:

```json
{
  "executor": "ai",
  "executor_config": {
    "provider": "anthropic",
    "apiKey": "sk-ant-...",
    "model": "claude-3-haiku-20240307",
    "prompt": "Generate a daily status report template",
    "maxTokens": "1000"
  }
}
```

Available Anthropic models include:
- `claude-3-opus-20240229`
- `claude-3-sonnet-20240229`
- `claude-3-haiku-20240307`
- `claude-3-5-sonnet-latest`

### Local Models

Use with local models via OpenAI-compatible APIs (Ollama, LM Studio, vLLM, etc.):

```json
{
  "executor": "ai",
  "executor_config": {
    "provider": "local",
    "baseUrl": "http://localhost:11434/v1",
    "model": "llama2",
    "prompt": "What is the current system status?",
    "maxTokens": "500"
  }
}
```

For Ollama, ensure it's running with:
```bash
ollama serve
```

## Example Use Cases

### Daily Report Generation

```json
{
  "name": "daily-ai-report",
  "schedule": "@daily",
  "executor": "ai",
  "executor_config": {
    "provider": "openai",
    "apiKey": "sk-...",
    "model": "gpt-4o-mini",
    "prompt": "Generate a brief daily operations checklist for a distributed system administrator.",
    "maxTokens": "300"
  }
}
```

### Scheduled Content Generation

```json
{
  "name": "weekly-summary",
  "schedule": "0 9 * * MON",
  "executor": "ai",
  "executor_config": {
    "provider": "anthropic",
    "apiKey": "sk-ant-...",
    "model": "claude-3-haiku-20240307",
    "prompt": "Write a motivational message for the start of a new work week, focusing on productivity and teamwork.",
    "maxTokens": "200",
    "temperature": "0.9"
  }
}
```

### Local Model for Privacy-Sensitive Tasks

```json
{
  "name": "local-analysis",
  "schedule": "@hourly",
  "executor": "ai",
  "executor_config": {
    "provider": "local",
    "baseUrl": "http://localhost:11434/v1",
    "model": "codellama",
    "prompt": "Analyze the following log pattern and suggest optimizations: [pattern]",
    "maxTokens": "500",
    "timeout": "180"
  }
}
```

## Security Considerations

- **API Keys**: Store API keys securely. Consider using environment variables or a secrets manager, and reference them in your job configuration.
- **Local Models**: For sensitive data processing, use local models to keep data within your infrastructure.
- **Rate Limits**: Be aware of API rate limits when scheduling frequent jobs.
- **Costs**: Monitor API usage and costs, especially for high-frequency scheduled jobs.

## Output

The executor returns the AI model's response as the job output. This can be:
- Processed by downstream processors (e.g., log, files, email)
- Used in job chains for multi-step workflows
- Stored and reviewed in the Dkron UI

## Troubleshooting

Enable debug mode to see detailed request/response information:

```json
{
  "executor": "ai",
  "executor_config": {
    "provider": "openai",
    "apiKey": "sk-...",
    "prompt": "Test prompt",
    "debug": "true"
  }
}
```

Common issues:
- **Invalid API key**: Verify your API key is correct and has appropriate permissions
- **Timeout errors**: Increase the `timeout` value for complex prompts or slow models
- **Rate limiting**: Reduce job frequency or upgrade your API plan
- **Local model connection**: Ensure the local model server is running and accessible
