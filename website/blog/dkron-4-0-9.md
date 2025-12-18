---
title: Dkron v4.0.9 Release - Enhanced Scheduling and Control
Keywords:    [ "Development", "OpenSource", "Distributed systems", "cron", "job scheduling" ]
Tags:        [ "Development", "OpenSource", "Distributed systems", "cron" ]
date: "2025-12-16"
author: The Dkron Team
---

We're pleased to announce the upcoming release of Dkron v4.0.9! This release brings several powerful new features and important fixes that enhance job scheduling control, resource management, and system reliability.

## What's New in Dkron v4.0.9?

### Schedule Jobs with Precise Start Times

One of the most requested features is now here! The new `starts_at` field allows you to schedule jobs to begin execution at a specific date and time. This is perfect for one-time events, maintenance windows, or any scenario where you need precise control over when a job becomes active.

```json
{
  "name": "maintenance-job",
  "schedule": "@daily",
  "starts_at": "2025-12-20T02:00:00Z",
  "executor": "shell",
  "executor_config": {
    "command": "maintenance.sh"
  }
}
```

The UI has been updated to support this new field, making it easy to configure start times directly from the dashboard.

### Pause and Resume Job Submissions

Gain better control during maintenance windows or system updates with the new API endpoints to pause and unpause job submissions. When paused, no new jobs will be accepted, but existing scheduled jobs continue to run normally.

```bash
# Pause new job submissions
curl -X POST http://dkron-server:8080/v1/pause

# Resume job submissions
curl -X POST http://dkron-server:8080/v1/unpause
```

### Memory Limits for Shell Executor Jobs

Prevent runaway processes from consuming all system resources! You can now set memory limits for shell executor jobs using the new `memory_limit` configuration option. This helps ensure system stability and predictable resource usage.

```json
{
  "executor": "shell",
  "executor_config": {
    "command": "data-processing.sh",
    "memory_limit": "512m"
  }
}
```

### Enhanced Metrics and Monitoring

Job execution metrics are now tracked using the go-metrics package, providing better insights into your job performance. Monitor execution times, success rates, and system health with greater precision.

### Improved Concurrency Handling

The "forbid" concurrency policy now correctly survives node restarts. Previously, jobs with concurrency set to "forbid" could incorrectly block after a node restart. This fix ensures your concurrency policies work reliably across the entire cluster lifecycle.

## Bug Fixes and Improvements

### Critical Fixes

- **Fixed nil pointer panic during startup**: Resolved issues with Raft-dependent methods being called before initialization was complete
- **Fixed Docker startup with custom address pools**: Dkron now works properly with custom Docker network configurations
- **Fixed mutex copy in processor plugin interface**: Improved thread safety in plugin communications
- **Fixed dependent_jobs preservation**: Job updates now correctly preserve existing dependent_jobs values

### Developer Experience

- **Mailpit integration**: Replaced the unmaintained MailHog with actively maintained Mailpit for email testing
- **Documentation search**: The documentation site now includes full-text search functionality
- **Architecture diagrams**: Added comprehensive diagrams documenting the job execution flow
- **UI improvements**: Job IDs are now visible on small screens for better mobile experience

## Upgrading to v4.0.9

Dkron v4.0.9 is designed to be a straightforward upgrade from v4.0.8. As always, we recommend:

1. Backing up your data store before upgrading
2. Testing the upgrade in a non-production environment first
3. Following a rolling upgrade pattern for production clusters

## What's Next?

We continue to invest in making Dkron the most reliable and feature-rich distributed job scheduler. Future releases will focus on:

- Enhanced observability with OpenTelemetry integration improvements
- Advanced job dependencies and workflows
- Performance optimizations for large-scale deployments

## Get Dkron v4.0.9

- [Download Dkron](https://github.com/distribworks/dkron/releases)
- [Read the Full Documentation](https://dkron.io/docs/)
- [View the Changelog](https://github.com/distribworks/dkron/blob/main/CHANGELOG.md)

## Community and Support

Join our growing community:

- **GitHub**: [github.com/distribworks/dkron](https://github.com/distribworks/dkron)
- **Discussions**: Share your use cases and get help
- **Issues**: Report bugs and request features

Thank you to all our contributors who made this release possible, especially [@NAlexandrov](https://github.com/NAlexandrov) for the starts_at feature and [@indeedhat](https://github.com/indeedhat) for the dependent_jobs fix!

Happy scheduling!

The Dkron Team
