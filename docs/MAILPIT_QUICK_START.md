# Mailpit Quick Start Guide

Quick reference for getting started with Mailpit email testing in Dkron.

## TL;DR

```bash
# Start Mailpit
docker run -p 8025:8025 -p 1025:1025 axllent/mailpit

# Run tests
make test-email

# View emails
open http://localhost:8025
```

## What is Mailpit?

Mailpit is a modern email testing tool that captures outgoing emails without sending them. It's the actively maintained successor to MailHog.

**Key Features:**
- 🚀 Fast and lightweight (~15MB Docker image)
- 🎨 Modern UI with dark mode
- 🔍 Full-text search
- 📱 Responsive design
- ⚡ Real-time updates via WebSocket
- 🔓 No authentication required

## Installation

### Option 1: Docker (Recommended)

```bash
docker run -p 8025:8025 -p 1025:1025 axllent/mailpit
```

### Option 2: Docker Compose

```bash
docker compose -f docker-compose.dev.yml up mailpit
```

### Option 3: Using Make

```bash
make test-email
```

## Ports

- **1025** - SMTP server (send emails here)
- **8025** - Web UI (view emails here)

## Running Tests

### Run Email Tests

```bash
# Ensure Mailpit is running, then:
go test -v -run TestNotifier_sendExecutionEmail ./dkron
```

### Run All Tests (CI Mode)

```bash
./scripts/test-ci-locally.sh
```

This simulates GitHub Actions CI environment.

## Web Interface

Open http://localhost:8025 in your browser to:

- View all captured emails
- Search emails by subject, sender, recipient
- Toggle dark mode
- Inspect headers and raw MIME
- Delete or mark emails

## Configuration

Dkron is pre-configured to use Mailpit:

```go
MailHost: "localhost"
MailPort: 1025
// No authentication required
```

## GitHub Actions

Mailpit runs automatically in CI. No configuration needed!

The workflow includes Mailpit as a service container:

```yaml
services:
    mailpit:
        image: axllent/mailpit
        ports:
            - 1025:1025
            - 8025:8025
```

## Common Commands

```bash
# Start Mailpit
docker run -d --name mailpit -p 8025:8025 -p 1025:1025 axllent/mailpit

# Stop Mailpit
docker stop mailpit

# Remove Mailpit
docker rm mailpit

# View logs
docker logs mailpit

# Restart Mailpit
docker restart mailpit
```

## Troubleshooting

### Port already in use

```bash
# Find what's using the port
lsof -i :1025
lsof -i :8025

# Use different ports
docker run -p 8026:8025 -p 1026:1025 axllent/mailpit
```

### Mailpit not starting

```bash
# Check if Docker is running
docker ps

# Pull latest image
docker pull axllent/mailpit

# Check container logs
docker logs mailpit
```

### Tests skip with "Mailpit not available"

```bash
# Verify Mailpit is running
docker ps | grep mailpit

# Check connection
curl -I http://localhost:8025
nc -zv localhost 1025
```

## Migrating from MailHog

Mailpit is a drop-in replacement:

1. Stop MailHog containers:
   ```bash
   docker stop $(docker ps -a -q --filter="ancestor=mailhog/mailhog")
   ```

2. Start Mailpit:
   ```bash
   docker run -p 8025:8025 -p 1025:1025 axllent/mailpit
   ```

3. Run tests - everything works the same!

**Note:** Same ports, same configuration, better features.

## API Usage (Optional)

Mailpit provides a REST API for automation:

```bash
# List all messages
curl http://localhost:8025/api/v1/messages

# Search messages
curl http://localhost:8025/api/v1/search?query=test

# Get specific message
curl http://localhost:8025/api/v1/message/{id}

# Delete all messages
curl -X DELETE http://localhost:8025/api/v1/messages
```

See [API docs](https://mailpit.axllent.org/docs/api-v1/) for more.

## Comparison: MailHog vs Mailpit

| Feature | MailHog | Mailpit |
|---------|---------|---------|
| Status | 🔴 Archived | ✅ Active |
| Size | 30MB | 15MB |
| UI | Basic | Modern |
| Search | Limited | Full-text |
| Dark Mode | ❌ | ✅ |
| Performance | Good | Excellent |

## Next Steps

- 📖 Read full guide: [docs/EMAIL_TESTING.md](EMAIL_TESTING.md)
- 🔄 Migration details: [docs/MAILPIT_MIGRATION.md](MAILPIT_MIGRATION.md)
- 🤖 CI integration: [docs/GITHUB_ACTIONS_MAILPIT.md](GITHUB_ACTIONS_MAILPIT.md)
- 🧪 CI testing: [.github/TESTING.md](../.github/TESTING.md)

## Resources

- [Mailpit GitHub](https://github.com/axllent/mailpit)
- [Mailpit Docs](https://mailpit.axllent.org/)
- [API Reference](https://mailpit.axllent.org/docs/api-v1/)

## Support

If you encounter issues:

1. Check this guide
2. Review [troubleshooting](#troubleshooting) section
3. Run `./scripts/test-ci-locally.sh` to verify setup
4. Check [full documentation](EMAIL_TESTING.md)

---

**Quick Help:**

```bash
# Fresh start
docker stop mailpit && docker rm mailpit
docker run -d --name mailpit -p 8025:8025 -p 1025:1025 axllent/mailpit
make test-email
open http://localhost:8025
```
