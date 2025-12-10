# Testing in CI/CD

This document describes how automated testing works in GitHub Actions for the Dkron project.

## Overview

The Dkron project uses GitHub Actions for continuous integration. Tests run automatically on:
- Every push to any branch
- Every pull request (excluding documentation-only changes)

## Test Workflow

The main test workflow is defined in `.github/workflows/test.yml`.

### Services

The following services are automatically started for testing:

#### Mailpit (Email Testing)

Mailpit runs as a service container during all test runs to enable email notification testing without external dependencies.

**Configuration:**
- Image: `axllent/mailpit`
- SMTP Port: `1025`
- Web UI Port: `8025` (not accessible in CI)

**Usage in Tests:**
Tests can send emails to `localhost:1025` without authentication. The email notification tests in `dkron/notifier_test.go` are configured to use this by default.

## Running Tests Locally to Match CI

To ensure your tests will pass in GitHub Actions, run them with the same environment:

### 1. Start Required Services

```bash
# Start Mailpit for email testing
docker run -d --rm --name mailpit -p 1025:1025 -p 8025:8025 axllent/mailpit
```

### 2. Run Tests

```bash
# Run all tests with the same timeout as CI
go test -v -timeout 200s -coverprofile=coverage.txt ./...
```

### 3. Cleanup

```bash
# Stop MailHog
docker stop mailhog
```

Or use the Makefile:

```bash
make test-email  # Runs email tests with MailHog
```

## Test Configuration

### Go Version

The CI uses Go version `1.23.1`. Ensure you're testing with the same version locally:

```bash
go version
```

### Test Timeout

Tests have a 200-second timeout in CI:

```bash
go test -v -timeout 200s ./...
```

### Coverage

Test coverage is automatically uploaded to Codecov after successful test runs.

## Email Testing

Email notification tests require Mailpit to be running. In GitHub Actions, this is handled automatically via service containers.

**Key Points:**
- ✅ Mailpit runs automatically in CI
- ✅ No credentials or secrets required
- ✅ Tests use `localhost:1025` for SMTP
- ✅ If Mailpit is not available locally, tests are skipped (not failed)

See [docs/EMAIL_TESTING.md](../docs/EMAIL_TESTING.md) for detailed email testing documentation.

## Skipped Tests

Some tests may be skipped in certain conditions:

- **Email tests**: Skipped if Mailpit is not available (only in local development, not in CI)

## Troubleshooting CI Failures

### Email Tests Failing

If email tests fail in CI:

1. **Check Mailpit Service**: Ensure the service is defined in `.github/workflows/test.yml`
2. **Verify Port Configuration**: SMTP should be on port `1025`
3. **Check Test Configuration**: Verify `dkron/notifier_test.go` uses `localhost:1025`

### Timeout Failures

If tests timeout:

1. **Check Test Duration**: Individual tests should complete quickly
2. **Review Logs**: Look for hanging operations or infinite loops
3. **Increase Timeout**: If legitimate, adjust the timeout in the workflow

### Coverage Upload Failures

Coverage upload failures don't fail the build. If you need coverage data:

1. **Check Codecov Token**: Verify `CODECOV_TOKEN` secret is set in repository settings
2. **Review codecov-action**: Ensure the action version is up to date

## Adding New Tests

When adding new tests that require external services:

1. **Update Workflow**: Add service containers to `.github/workflows/test.yml`
2. **Document Requirements**: Update this file with service details
3. **Add Local Setup**: Update `docs/` with local development instructions
4. **Make Tests Skippable**: Use `t.Skip()` if service is unavailable locally

Example:

```go
func TestNewFeature(t *testing.T) {
    // Check if required service is available
    if !isServiceAvailable("localhost", 1234) {
        t.Skip("Service not available. Start with: docker run ...")
        return
    }
    
    // Test implementation
}
```

## Workflow Triggers

Tests run on:

```yaml
on:
  push:
    branches:
      - '**'        # All branches
    tags-ignore:
      - '**'        # Ignore tags
      
  pull_request:
    paths-ignore:   # Skip if only these files changed
      - '**.md'
      - 'website/**'
      - 'docs/**'
      - 'examples/**'
      - 'ui'
```

## Best Practices

1. **Keep Tests Fast**: CI has a 200-second timeout for all tests
2. **Use Service Containers**: Don't install services in workflow steps
3. **Match Local and CI**: Use the same configuration locally as in CI
4. **Provide Clear Errors**: Use descriptive error messages and skip conditions
5. **Clean Up Resources**: Ensure tests clean up after themselves
6. **Avoid External Dependencies**: Use local services (like MailHog) instead of external APIs

## Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Service Containers](https://docs.github.com/en/actions/using-containerized-services/about-service-containers)
- [Mailpit Documentation](https://github.com/axllent/mailpit)
- [Email Testing Guide](../docs/EMAIL_TESTING.md)