# Mailpit Integration with GitHub Actions

This document provides a comprehensive guide on how Mailpit is integrated with GitHub Actions for automated email testing in the Dkron project.

## Overview

Mailpit is used in GitHub Actions as a service container to enable automated testing of email notifications without requiring external email services or credentials. This provides a reliable, fast, and secure way to test email functionality in CI/CD pipelines.

## Why Mailpit for GitHub Actions?

### Advantages

1. **No Secrets Required**: Mailpit doesn't need authentication, eliminating the need to store SMTP credentials in GitHub Secrets
2. **Fast**: Runs locally within the GitHub Actions runner with minimal overhead
3. **Reliable**: No dependency on external services that might have downtime or rate limits
4. **Free**: No costs associated with email testing in CI
5. **Consistent**: Same behavior in local development and CI environments
6. **Secure**: Emails never leave the CI runner, no risk of accidentally sending test emails
7. **Modern**: Active maintenance and regular updates
8. **Lightweight**: Small Docker image (~15MB) for faster CI runs

## GitHub Actions Configuration

### Service Container Setup

In `.github/workflows/test.yml`, Mailpit is configured as a service container:

```yaml
jobs:
    test:
        name: Test
        runs-on: ubuntu-latest
        services:
            mailpit:
                image: axllent/mailpit
                ports:
                    - 1025:1025  # SMTP port
                    - 8025:8025  # Web UI port (not accessible in CI)
        steps:
            # ... test steps
```

### How It Works

1. **Container Start**: When the workflow runs, GitHub Actions automatically starts the Mailpit container
2. **Port Mapping**: Ports 1025 (SMTP) and 8025 (Web UI) are mapped to the runner
3. **Network Access**: Tests can access Mailpit at `localhost:1025`
4. **Automatic Cleanup**: The container is automatically removed when the workflow completes

## Test Configuration

### Email Test Setup

The email notification tests in `dkron/notifier_test.go` are configured to use Mailpit:

```go
func TestNotifier_sendExecutionEmail(t *testing.T) {
    mailHost := "localhost"
    mailPort := 1025  // Mailpit SMTP port
    
    // Check if Mailpit is available
    checkMailpitAvailable(t, mailHost, mailPort)
    
    c := &Config{
        MailHost:          mailHost,
        MailPort:          uint16(mailPort),
        MailFrom:          "dkron@dkron.io",
        MailSubjectPrefix: "[Test]",
    }
    
    // ... test implementation
}
```

### Availability Check

The `checkMailpitAvailable` helper function:
- Attempts to connect to Mailpit on the specified port
- **In CI**: Mailpit is always available, tests run normally
- **Locally**: If Mailpit isn't running, test is skipped (not failed)

```go
func checkMailpitAvailable(t *testing.T, host string, port int) {
    t.Helper()
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", address, 2*time.Second)
    if err != nil {
        t.Skipf("Mailpit is not available at %s. Start it with: docker run -p 8025:8025 -p 1025:1025 axllent/mailpit", address)
        return
    }
    conn.Close()
}
```

## Workflow Execution Flow

### 1. Workflow Trigger

Tests run on:
- Every push to any branch
- Every pull request (excluding documentation-only changes)

### 2. Environment Setup

```
┌─────────────────────────────────────┐
│   GitHub Actions Runner (Ubuntu)    │
│                                      │
│  ┌────────────────────────────────┐ │
│  │   Mailpit Container            │ │
│  │   - SMTP: localhost:1025       │ │
│  │   - Web UI: localhost:8025     │ │
│  │   - Image: ~15MB               │ │
│  └────────────────────────────────┘ │
│                                      │
│  ┌────────────────────────────────┐ │
│  │   Go Test Process              │ │
│  │   - Connects to localhost:1025 │ │
│  │   - Sends test emails          │ │
│  └────────────────────────────────┘ │
└─────────────────────────────────────┘
```

### 3. Test Execution

1. Mailpit container starts automatically
2. Go sets up the test environment
3. Tests run with 200-second timeout
4. Email tests connect to Mailpit
5. Emails are captured (but not viewable in CI)
6. Tests verify email sending succeeded
7. Coverage report generated
8. Results uploaded to Codecov

### 4. Cleanup

- Mailpit container automatically stopped and removed
- No manual cleanup required

## Running Tests Locally to Match CI

### Option 1: Using the Validation Script

```bash
./scripts/test-ci-locally.sh
```

This script:
- ✅ Starts Mailpit with the same configuration as CI
- ✅ Runs tests with the same timeout (200s)
- ✅ Generates coverage report
- ✅ Provides interactive access to Mailpit web UI
- ✅ Automatically cleans up after completion

### Option 2: Manual Setup

```bash
# 1. Start Mailpit
docker run -d --rm --name mailpit -p 1025:1025 -p 8025:8025 axllent/mailpit

# 2. Run tests (matching CI configuration)
go test -v -timeout 200s -coverprofile=coverage.txt ./...

# 3. View emails (optional)
open http://localhost:8025

# 4. Cleanup
docker stop mailpit
```

### Option 3: Using Make

```bash
make test-email
```

## Viewing Email in CI

### Limitations

The Mailpit web UI (port 8025) is not accessible from outside the GitHub Actions runner. You cannot view the web interface during CI runs.

### API-Based Verification

The email notification tests automatically use Mailpit's REST API to verify email delivery. This provides comprehensive validation without requiring web UI access:

**Automatic Verification in Tests:**
1. ✅ Email was successfully received by Mailpit
2. ✅ From address matches expected sender
3. ✅ To address matches expected recipient
4. ✅ Subject line contains expected text
5. ✅ Email body contains execution output and details

**Example from `dkron/notifier_test.go`:**

```go
// Get all messages from Mailpit API
messages := getMailpitMessages(t, "http://localhost:8025")

// Find and verify the test email
for _, msg := range messages {
    if strings.Contains(msg.Subject, "[Test]") {
        // Get full message details
        detail := getMailpitMessage(t, apiURL, msg.ID)
        
        // Verify content
        assert.Equal(t, "dkron@dkron.io", detail.From.Address)
        assert.Contains(t, detail.Text, "test-output")
    }
}
```

This approach ensures reliable email testing in CI without manual inspection.

## Troubleshooting

### Tests Fail in CI But Pass Locally

**Possible causes:**
1. Different Go versions
2. Environment-specific configurations
3. Race conditions or timing issues

**Solutions:**
```bash
# Run with the same Go version as CI
go version  # Check against .github/workflows/test.yml

# Run with the same timeout as CI
go test -v -timeout 200s ./...

# Run race detector
go test -race ./...
```

### Mailpit Service Not Starting

**Error**: Tests skip with "Mailpit is not available"

**Diagnosis:**
1. Check workflow file has `services.mailpit` section
2. Verify port mappings are correct (`1025:1025`)
3. Review GitHub Actions logs for container startup errors

**Fix**:
Ensure `.github/workflows/test.yml` contains:
```yaml
services:
    mailpit:
        image: axllent/mailpit
        ports:
            - 1025:1025
            - 8025:8025
```

### Port Conflicts in CI

**Unlikely scenario**: Port 1025 or 8025 already in use on runner

**Solutions:**
- Use different ports in workflow configuration
- Update test configuration to match

### Coverage Upload Fails

**Note**: Coverage upload failures don't fail the build

**Common causes:**
- `CODECOV_TOKEN` secret not set
- Codecov service temporarily unavailable
- Network issues

**Fix**:
1. Verify secret exists in repository settings
2. Check Codecov status page
3. Update codecov-action to latest version

## Best Practices

### 1. Keep Tests Fast

✅ **Do:**
```go
func TestEmail(t *testing.T) {
    // Quick email send test
    err := sendEmail()
    assert.NoError(t, err)
}
```

❌ **Don't:**
```go
func TestEmail(t *testing.T) {
    // Don't poll or wait unnecessarily
    time.Sleep(30 * time.Second)
    // ...
}
```

### 2. Use Helper Functions

```go
// Reusable helper for email testing
func setupEmailTest(t *testing.T) *Config {
    t.Helper()
    checkMailpitAvailable(t, "localhost", 1025)
    return &Config{
        MailHost: "localhost",
        MailPort: 1025,
        // ...
    }
}
```

### 3. Clear Error Messages

```go
if err != nil {
    t.Fatalf("Failed to send email: %v\nCheck that Mailpit is running", err)
}
```

### 4. Skip Gracefully

```go
// Skip test if Mailpit unavailable (local dev)
// Test runs normally in CI (Mailpit always available)
checkMailpitAvailable(t, host, port)
```

## Monitoring and Maintenance

### Workflow Health

Monitor test runs in the GitHub Actions tab:
- Check for consistent pass rates
- Watch for timeout increases
- Monitor for flaky tests

### Mailpit Updates

The `axllent/mailpit` image is pulled automatically by GitHub Actions:
- No manual updates needed
- GitHub caches images for performance
- Updates happen when new workflow runs pull newer images

### Go Version Updates

When updating Go version in workflow:

1. Update `.github/workflows/test.yml`:
   ```yaml
   - uses: actions/setup-go@v5
     with:
       go-version: 1.23.1  # Update this
   ```

2. Test locally with the same version:
   ```bash
   go version
   ```

3. Update documentation if needed

## Mailpit Features in CI

### Performance Benefits

- **Fast Startup**: Container starts in ~1-2 seconds
- **Low Memory**: Minimal memory footprint during tests
- **Quick Processing**: Emails processed instantly
- **Efficient Storage**: In-memory storage for test duration

### API Usage in Tests

The Dkron tests use Mailpit's API for comprehensive email verification. This is the recommended approach for CI testing:

**Messages API:**
```go
// List all messages
resp, err := http.Get("http://localhost:8025/api/v1/messages")

// Response includes: total count, message list with ID, From, To, Subject
var result mailpitMessagesResponse
json.NewDecoder(resp.Body).Decode(&result)
```

**Message Detail API:**
```go
// Get specific message with full content
resp, err := http.Get("http://localhost:8025/api/v1/message/MESSAGE_ID")

// Response includes: full Text/HTML content, headers, attachments
var detail mailpitMessageDetail
json.NewDecoder(resp.Body).Decode(&detail)
```

**Cleanup API:**
```go
// Delete all messages (useful for test isolation)
req, err := http.NewRequest("DELETE", "http://localhost:8025/api/v1/messages", nil)
http.DefaultClient.Do(req)
```

See [Mailpit API documentation](https://mailpit.axllent.org/docs/api-v1/) for complete API reference.

**Real Example:** Check `dkron/notifier_test.go` for a complete implementation showing how to verify email delivery, content, and metadata using the API.

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Service Containers Guide](https://docs.github.com/en/actions/using-containerized-services/about-service-containers)
- [Mailpit GitHub Repository](https://github.com/axllent/mailpit)
- [Mailpit Documentation](https://mailpit.axllent.org/)
- [Email Testing Guide](EMAIL_TESTING.md)
- [CI Testing Guide](../.github/TESTING.md)

## Support

If you encounter issues with Mailpit in GitHub Actions:

1. Check this documentation first
2. Review the [troubleshooting section](#troubleshooting)
3. Run tests locally with `./scripts/test-ci-locally.sh`
4. Check GitHub Actions logs for detailed error messages
5. Review recent workflow changes
6. Open an issue with reproduction steps

## Summary

Mailpit integration with GitHub Actions provides:

- ✅ **Automated email testing** without external dependencies
- ✅ **No secrets or credentials** required in CI
- ✅ **Fast and reliable** test execution with modern tooling
- ✅ **Consistent behavior** between local and CI environments
- ✅ **Zero cost** for email testing
- ✅ **Simple setup** with minimal configuration
- ✅ **Active maintenance** ensuring long-term reliability
- ✅ **Better performance** with smaller Docker image and faster processing

The integration is production-ready and has been tested to work reliably in the Dkron CI/CD pipeline.