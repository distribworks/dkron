# Migration from MailHog to Mailpit

This document summarizes the changes made to migrate from MailHog to Mailpit for email testing in the Dkron project.

## Overview

The project has been updated to use **Mailpit** instead of MailHog for email testing. Mailpit is a modern, actively maintained email testing tool that serves as a drop-in replacement for MailHog, providing better performance, a superior UI, and additional features.

## Why Migrate to Mailpit?

### Key Advantages

1. **Actively Maintained**: MailHog is archived/unmaintained, while Mailpit receives regular updates
2. **Better Performance**: Written in Go with optimized performance
3. **Modern UI**: Clean, responsive interface with dark mode support
4. **More Features**: Search, filtering, tagging, better API
5. **Smaller Footprint**: ~15MB compressed vs ~30MB for MailHog
6. **Drop-in Replacement**: Uses same ports (1025 SMTP, 8025 Web UI)

### Feature Comparison

| Feature | MailHog | Mailpit |
|---------|---------|---------|
| Maintenance | Archived | Active |
| Docker Image Size | ~30MB | ~15MB |
| UI | Basic | Modern, responsive |
| Search | Limited | Full-text search |
| Dark Mode | No | Yes |
| Real-time Updates | Polling | WebSocket |
| API | Basic | Comprehensive |
| Performance | Good | Excellent |

## Changes Made

### 1. Test Configuration Updated

**File**: `dkron/dkron/notifier_test.go`

The `TestNotifier_sendExecutionEmail` function has been updated to reference Mailpit:

**Before (MailHog)**:
```go
// This test requires MailHog to be running for email testing.
// Start MailHog with: docker run -p 8025:8025 -p 1025:1025 mailhog/mailhog
checkMailHogAvailable(t, mailHost, mailPort)
```

**After (Mailpit)**:
```go
// This test requires Mailpit to be running for email testing.
// Start Mailpit with: docker run -p 8025:8025 -p 1025:1025 axllent/mailpit
checkMailpitAvailable(t, mailHost, mailPort)
```

**Key Changes**:
- Function renamed from `checkMailHogAvailable` to `checkMailpitAvailable`
- Docker image changed from `mailhog/mailhog` to `axllent/mailpit`
- Updated comments and documentation references
- **No configuration changes needed** (same ports: 1025 for SMTP, 8025 for Web UI)

### 2. Docker Compose Configuration

**File**: `docker-compose.dev.yml`

Updated service definition:

**Before**:
```yaml
mailhog:
    image: mailhog/mailhog
    ports:
        - "8025:8025"
        - "1025:1025"
```

**After**:
```yaml
mailpit:
    image: axllent/mailpit
    ports:
        - "8025:8025"
        - "1025:1025"
```

### 3. GitHub Actions Workflow

**File**: `.github/workflows/test.yml`

Updated service container:

**Before**:
```yaml
services:
    mailhog:
        image: mailhog/mailhog
        ports:
            - 1025:1025
            - 8025:8025
```

**After**:
```yaml
services:
    mailpit:
        image: axllent/mailpit
        ports:
            - 1025:1025
            - 8025:8025
```

### 4. Makefile Enhancement

**File**: `Makefile`

Updated `test-email` target:

**Before**:
```makefile
test-email:
	@echo "Starting MailHog for email testing..."
	@docker run -d --rm --name dkron-mailhog -p 8025:8025 -p 1025:1025 mailhog/mailhog
	@echo "To stop MailHog, run: docker stop dkron-mailhog"
```

**After**:
```makefile
test-email:
	@echo "Starting Mailpit for email testing..."
	@docker run -d --rm --name dkron-mailpit -p 8025:8025 -p 1025:1025 axllent/mailpit
	@echo "To stop Mailpit, run: docker stop dkron-mailpit"
```

### 5. CI Testing Script

**File**: `scripts/test-ci-locally.sh`

Updated all references:
- Container name: `dkron-ci-mailhog` → `dkron-ci-mailpit`
- Docker image: `mailhog/mailhog` → `axllent/mailpit`
- All error messages and documentation references updated

### 6. Documentation Updates

**File**: `README.md`

Updated all MailHog references to Mailpit in:
- Email Testing section
- Docker commands
- Service names in docker-compose commands

**Updated Documentation Files**:
- `docs/EMAIL_TESTING.md` - Comprehensive Mailpit guide
- `docs/MAILPIT_MIGRATION.md` - This migration guide
- `.github/TESTING.md` - CI testing guide

## How to Use

### Quick Start

**Old (MailHog)**:
```bash
docker run -p 8025:8025 -p 1025:1025 mailhog/mailhog
```

**New (Mailpit)**:
```bash
docker run -p 8025:8025 -p 1025:1025 axllent/mailpit
```

**Result**: Exact same behavior, better performance!

### Using Make

```bash
make test-email
```

### Using Docker Compose

**Old**:
```bash
docker compose -f docker-compose.dev.yml up mailhog
```

**New**:
```bash
docker compose -f docker-compose.dev.yml up mailpit
```

### Validate CI Setup

```bash
./scripts/test-ci-locally.sh
```

## Migration is Seamless

### What Stays the Same

✅ **Ports**: Still uses 1025 (SMTP) and 8025 (Web UI)
✅ **Configuration**: Same SMTP settings in tests
✅ **API**: Compatible API for programmatic access
✅ **Docker Usage**: Same deployment pattern
✅ **No Authentication**: Still doesn't require credentials

### What Gets Better

🚀 **Performance**: Faster email processing and UI rendering
🎨 **UI/UX**: Modern, responsive interface with better usability
🔍 **Search**: Full-text search across all emails
🌙 **Dark Mode**: Toggle between light and dark themes
📊 **Features**: Better filtering, tagging, and message management
🔄 **Updates**: Active maintenance and regular improvements

## Benefits

### For Developers

1. **Better DX**: Improved UI makes debugging easier
2. **Faster**: Quicker startup and email processing
3. **More Reliable**: Active maintenance ensures bug fixes
4. **Better Search**: Find test emails more easily

### For CI/CD

1. **Faster Builds**: Smaller Docker image downloads faster
2. **More Stable**: Active project with bug fixes
3. **Same Integration**: No CI configuration changes needed
4. **Future-Proof**: Won't become unmaintained

### For the Project

1. **Modern Stack**: Using actively maintained tools
2. **Better Support**: Community support and updates
3. **Enhanced Features**: Can leverage new Mailpit features
4. **No Breaking Changes**: Drop-in replacement

## Migration Checklist

If you're updating an existing environment:

- [x] Update Docker Compose service name (`mailhog` → `mailpit`)
- [x] Update Docker image reference (`mailhog/mailhog` → `axllent/mailpit`)
- [x] Update GitHub Actions workflow service
- [x] Update test helper function names
- [x] Update Makefile targets
- [x] Update documentation and comments
- [x] Update CI testing scripts
- [x] Stop any running MailHog containers
- [x] Start Mailpit containers
- [x] Verify tests pass with Mailpit

## Testing the Migration

### Step 1: Remove Old MailHog Container

```bash
docker stop dkron-mailhog 2>/dev/null
docker stop dkron-ci-mailhog 2>/dev/null
docker rm dkron-mailhog 2>/dev/null
docker rm dkron-ci-mailhog 2>/dev/null
```

### Step 2: Start Mailpit

```bash
docker run -d --rm --name mailpit -p 8025:8025 -p 1025:1025 axllent/mailpit
```

### Step 3: Run Tests

```bash
go test -v -run TestNotifier_sendExecutionEmail ./dkron
```

### Step 4: Verify in Browser

Open http://localhost:8025 and verify the email appears in the modern Mailpit UI.

## Troubleshooting

### Old MailHog Container Still Running

```bash
# List all containers
docker ps -a | grep mailhog

# Stop and remove
docker stop $(docker ps -a -q --filter="ancestor=mailhog/mailhog")
docker rm $(docker ps -a -q --filter="ancestor=mailhog/mailhog")
```

### Port Conflicts

If ports are in use:

```bash
# Check what's using the ports
lsof -i :1025
lsof -i :8025

# Or use different ports for Mailpit
docker run -p 8026:8025 -p 1026:1025 axllent/mailpit
```

### Tests Still Reference MailHog

All references should be updated. If you see MailHog mentioned:

```bash
# Search for any remaining references
grep -r "mailhog" . --exclude-dir=.git
```

## Rollback (If Needed)

If you need to rollback to MailHog temporarily:

```bash
# Stop Mailpit
docker stop mailpit

# Start MailHog
docker run -d --rm --name mailhog -p 8025:8025 -p 1025:1025 mailhog/mailhog
```

However, we recommend staying with Mailpit for the benefits listed above.

## Additional Resources

- [Mailpit GitHub Repository](https://github.com/axllent/mailpit)
- [Mailpit Documentation](https://mailpit.axllent.org/)
- [Email Testing Guide](EMAIL_TESTING.md)
- [GitHub Actions Integration](GITHUB_ACTIONS_MAILPIT.md)
- [CI Testing Guide](../.github/TESTING.md)

## Summary

The migration from MailHog to Mailpit is:

- ✅ **Complete**: All references updated
- ✅ **Seamless**: Drop-in replacement with same ports
- ✅ **Tested**: Works in local dev and GitHub Actions
- ✅ **Beneficial**: Better performance, UI, and features
- ✅ **Future-proof**: Actively maintained project

The implementation is production-ready and provides immediate improvements to the email testing workflow.