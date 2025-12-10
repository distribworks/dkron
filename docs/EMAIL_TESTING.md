# Email Testing with Mailpit

This document describes how to test email notifications in Dkron using Mailpit.

## What is Mailpit?

[Mailpit](https://github.com/axllent/mailpit) is a modern email testing tool for developers. It runs a fake SMTP server to capture outgoing emails instead of sending them to real recipients. This allows you to:

- Test email functionality without sending real emails
- Inspect email content, headers, and formatting via a web UI
- Verify emails are being sent correctly during development and testing
- View emails in a modern, responsive interface with search capabilities

## Why Mailpit?

Mailpit offers several advantages over older tools:

- **Actively Maintained**: Regular updates and bug fixes
- **Modern UI**: Clean, responsive interface with better UX
- **Fast**: Written in Go with excellent performance
- **Feature-Rich**: Search, filtering, tagging, and more
- **Small Footprint**: Lightweight Docker image (~15MB compressed)
- **API Support**: RESTful API for automated testing

## Running Mailpit

### Using Docker (Recommended)

The easiest way to run Mailpit is using Docker:

```bash
docker run -p 8025:8025 -p 1025:1025 axllent/mailpit
```

This exposes:
- **Port 1025**: SMTP server (for your application to send emails)
- **Port 8025**: Web UI (to view captured emails)

### Using Docker Compose

Mailpit is included in the development docker-compose configuration. Start it with:

```bash
docker compose -f docker-compose.dev.yml up mailpit
```

Or start the entire development environment including Mailpit:

```bash
docker compose -f docker-compose.dev.yml up
```

### Standalone Installation

You can also install Mailpit directly on your system. See the [Mailpit documentation](https://github.com/axllent/mailpit#installation) for installation instructions.

## Configuring Dkron to Use Mailpit

To configure Dkron to send emails through Mailpit, use the following SMTP settings:

```yaml
MailHost: localhost
MailPort: 1025
MailFrom: dkron@example.com
```

**Note**: Mailpit does not require authentication, so you don't need to set `MailUsername` or `MailPassword`.

## Running Email Tests

The email notification tests in `dkron/notifier_test.go` are configured to use Mailpit by default. The tests automatically verify that emails are correctly sent and received using Mailpit's REST API.

To run the email notification test:

```bash
# Make sure Mailpit is running first
docker run -d -p 8025:8025 -p 1025:1025 axllent/mailpit

# Run the test
go test -v -run TestNotifier_sendExecutionEmail ./dkron
```

Or use the convenient Makefile target:

```bash
make test-email
```

## Viewing Captured Emails

1. Open your browser and navigate to: http://localhost:8025
2. The Mailpit web UI will show all emails sent during testing
3. Click on any email to view its details, including:
   - Subject line
   - Recipients (To, Cc, Bcc)
   - Email body (HTML and plain text)
   - Attachments
   - Raw MIME content
   - Headers
   - Source code

### Web UI Features

Mailpit's web interface includes:
- **Search**: Full-text search across all emails
- **Filtering**: Filter by sender, recipient, subject
- **Responsive Design**: Works on desktop and mobile
- **Dark Mode**: Toggle between light and dark themes
- **Real-time Updates**: Emails appear instantly via WebSocket

## Test Configuration

The test in `dkron/notifier_test.go` uses the following configuration:

```go
c := &Config{
    MailHost:          "localhost",
    MailPort:          1025,
    MailFrom:          "dkron@dkron.io",
    MailSubjectPrefix: "[Test]",
}
```

This configuration:
- Connects to Mailpit's SMTP server at `localhost:1025`
- Sets the sender address to `dkron@dkron.io`
- Adds a `[Test]` prefix to all email subjects
- Uses Mailpit API at `http://localhost:8025` to verify email delivery

### API Verification

The test automatically verifies email delivery using Mailpit's REST API:

```go
// The test performs the following verifications:
// 1. Checks that email was received
// 2. Verifies From address is correct
// 3. Verifies To address matches recipient
// 4. Confirms subject contains expected text
// 5. Validates email body contains execution output
```

This ensures that emails are not just sent, but actually captured correctly by Mailpit.

## Troubleshooting

### Connection Refused Errors

If you see connection errors when running tests:

1. Verify Mailpit is running:
   ```bash
   docker ps | grep mailpit
   ```

2. Check that port 1025 is accessible:
   ```bash
   nc -zv localhost 1025
   ```

3. If using Docker, ensure the ports are properly mapped

### Emails Not Appearing in Web UI

1. Check the Mailpit web UI at http://localhost:8025
2. Verify the test completed successfully without errors
3. Check Mailpit logs for any issues:
   ```bash
   docker logs <mailpit-container-id>
   ```

### Port Already in Use

If port 1025 or 8025 is already in use:

```bash
# Find what's using the port
lsof -i :1025
lsof -i :8025

# Or use different ports
docker run -p 8026:8025 -p 1026:1025 axllent/mailpit
# Then update test configuration to use port 1026
```

## GitHub Actions CI/CD Integration

Mailpit is automatically configured in the GitHub Actions test workflow. The workflow includes Mailpit as a service container, making it available during test execution.

### How It Works

The `.github/workflows/test.yml` file includes Mailpit as a service:

```yaml
services:
    mailpit:
        image: axllent/mailpit
        ports:
            - 1025:1025
            - 8025:8025
```

When tests run in GitHub Actions:
1. Mailpit starts automatically as a service container
2. The SMTP server is available at `localhost:1025`
3. Tests can send emails without any additional configuration
4. No secrets or credentials are needed

### Viewing Emails in CI

Since GitHub Actions runners don't expose web UIs, you cannot view the Mailpit web interface during CI runs. However, you can:

1. Verify emails are sent successfully by checking test results
2. Tests automatically use Mailpit's API to verify email content
3. Email assertions include subject, sender, recipient, and body content

### Testing Locally Before CI

To ensure your tests will pass in GitHub Actions:

```bash
# Start Mailpit locally with the same configuration as CI
docker run -p 1025:1025 -p 8025:8025 axllent/mailpit

# Run the full test suite
go test -v -timeout 200s ./...
```

Or use the validation script:

```bash
./scripts/test-ci-locally.sh
```

## Advantages Over External Services

Mailpit provides several advantages for local development and testing:

- **Free and Open Source**: No account required, runs completely locally
- **No External Dependencies**: Doesn't require internet connection
- **Simple Setup**: Single Docker command to get started
- **Fast**: No network latency from external services
- **Privacy**: All emails stay on your local machine
- **CI/CD Friendly**: Easy to integrate into automated testing pipelines
- **No Secrets Required**: No API keys or credentials needed in CI

## Production Email Configuration

Remember that Mailpit is for **testing only**. For production, configure Dkron with your actual SMTP server credentials:

```yaml
MailHost: smtp.your-provider.com
MailPort: 587
MailUsername: your-username
MailPassword: your-password
MailFrom: noreply@yourdomain.com
```

## Using Mailpit API in Tests

The email notification tests demonstrate how to use Mailpit's API for verification:

### Get All Messages

```go
resp, err := http.Get("http://localhost:8025/api/v1/messages")
```

### Get Specific Message

```go
resp, err := http.Get(fmt.Sprintf("http://localhost:8025/api/v1/message/%s", messageID))
```

### Delete All Messages

```go
req, err := http.NewRequest("DELETE", "http://localhost:8025/api/v1/messages", nil)
resp, err := http.DefaultClient.Do(req)
```

See `dkron/notifier_test.go` for a complete example of API usage in tests.

## Additional Resources

- [Mailpit GitHub Repository](https://github.com/axllent/mailpit)
- [Mailpit Documentation](https://mailpit.axllent.org/)
- [API Documentation](https://mailpit.axllent.org/docs/api-v1/)
- [GitHub Actions Integration Guide](GITHUB_ACTIONS_MAILPIT.md)
- [CI Testing Guide](../.github/TESTING.md)