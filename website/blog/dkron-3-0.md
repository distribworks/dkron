---
Description: "Dkron release 3.0 public release"
Keywords:    [ "Development", "OpenSource", "Distributed systems", "cron" ]
Tags:        [ "Development", "OpenSource", "Distributed systems", "cron" ]
date:        "2020-05-12"
Topics:      [ "Development", "OpenSource", "Distributed Systems" ]
Slug:        "dkron-3-0"
---
# Dkron 3.0 Release

I'm thrilled to announce that Dkron/Pro v3.0 is here!

This release brings a big internal refactor in the job execution engine, incorporating breaking changes but ensuring no more missing executions.

## What's new

### Job execution engine

Refactored the job execution engine for proper synchronization of executions, no more missing executions under normal conditions, and if there is one, Dkron will report the issue clearly in the logs.

New node targeting algorithm, transparent for the user.

### UI Improvements

Change the notification JS code to a pop-up like system that provides better comfort in using the UI, previously causing some weird effects on certain job operations like Run, Toggle, and Delete.

## Wrap-up

This update brings no public API changes, and no changes in storage format so your upgrade path should be easy if you follow the [rolling upgrade notes](/docs/usage/upgrade/#rolling-upgrade).

Download and install from [here](/docs/basics/installation/)

*Thank you to all my Dkron Pro customers for ensuring the long-term support and maintenance of Dkron. Support OSS software and your infrastructure vendors so we can support you!*
