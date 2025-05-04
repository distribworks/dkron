---
title: Cron Specification
weight: 20
---

# Cron Expression Format

Dkron uses a powerful cron expression format to define job schedules. This guide explains the syntax in detail with practical examples.

## Basic Format

A cron expression consists of 6 space-separated fields:

| Field name   | Mandatory? | Allowed values  | Allowed special characters |
|------------- | ---------- | --------------- | -------------------------- |
| Seconds      | Yes        | 0-59            | * / , - ~                  |
| Minutes      | Yes        | 0-59            | * / , - ~                  |
| Hours        | Yes        | 0-23            | * / , - ~                  |
| Day of month | Yes        | 1-31            | * / , - ? ~                |
| Month        | Yes        | 1-12 or JAN-DEC | * / , - ~                  |
| Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ? ~                |

Note: Month and Day-of-week field values are case insensitive. "SUN", "Sun", and "sun" are equally accepted.

## Special Characters

### Asterisk ( * )

The asterisk indicates that the cron expression will match for all values of the field.

Example: `* * * * * *` runs every second of every minute of every hour, every day.

### Slash ( / )

Slashes define increments within ranges.

Example: `0 */15 * * * *` runs every 15 minutes (at 0:00, 0:15, 0:30, 0:45).

### Comma ( , )

Commas are used to separate items in a list.

Example: `0 0 0 * * MON,WED,FRI` runs at midnight on Mondays, Wednesdays, and Fridays.

### Hyphen ( - )

Hyphens define ranges.

Example: `0 0 9-17 * * *` runs every hour from 9 AM to 5 PM.

### Question mark ( ? )

Question mark may be used instead of '*' for leaving either day-of-month or day-of-week blank to avoid conflicts.

Example: `0 0 0 1 * ?` runs at midnight on the first day of each month, regardless of the day of the week.

### Tilde ( ~ )

Tilde is replaced by a deterministic numeric value based on the job name, providing an even distribution of load.

Example: `0 ~ * * * *` distributes jobs evenly across different minutes within the hour.

## Predefined Schedules

For convenience, Dkron supports several predefined schedules:

| Entry                  | Description                                | Equivalent To    |
|----------------------- | ------------------------------------------ | ---------------- |
| @yearly (or @annually) | Run once a year at midnight on January 1   | 0 0 0 1 1 *      |
| @monthly               | Run once a month at midnight on first day  | 0 0 0 1 * *      |
| @weekly                | Run once a week at midnight on Sunday      | 0 0 0 * * 0      |
| @daily (or @midnight)  | Run once a day at midnight                 | 0 0 0 * * *      |
| @hourly                | Run once an hour at the beginning of hour  | 0 0 * * * *      |
| @minutely              | Run once a minute at the beginning         | 0 * * * * *      |
| @manually              | Never runs automatically (manual triggers) | N/A              |

Example: `@daily` is equivalent to `0 0 0 * * *`

## Fixed Intervals

You can schedule jobs to execute at fixed intervals using:

```
@every <duration>
```

Where "duration" is a string accepted by [time.ParseDuration](http://golang.org/pkg/time/#ParseDuration).

Example: `@every 1h30m10s` runs every 1 hour, 30 minutes, and 10 seconds.

Note: The interval does not account for job runtime. For example, if a job takes 3 minutes to run, and it is scheduled to run every 5 minutes, it will have only 2 minutes of idle time between each run.

## One-time Execution

To schedule a job to be executed just once at a specific time:

```
@at <datetime>
```

Where "datetime" is a string in [RFC3339 format](https://golang.org/pkg/time/#Parse).

Example: `@at 2023-12-31T23:59:00Z` will run the job once on December 31, 2023 at 11:59 PM UTC.

## Time Zones

Dkron supports scheduling jobs in specific time zones by specifying the `timezone` parameter in a job definition.

If no time zone is specified:
- All scheduling is done in the server's local time zone
- Jobs scheduled during daylight-savings leap-ahead transitions will not be run

When a `timezone` is specified, the job will be scheduled in the specified time zone, correctly handling daylight saving time transitions.

Example timezone values:
- "America/New_York"
- "Europe/London"
- "Asia/Tokyo"
- "UTC"

## Practical Examples

Here are some common scheduling patterns with explanations:

### Every weekday at 8 AM
```
0 0 8 * * 1-5
```
This runs at 8:00 AM Monday through Friday.

### Every 10 minutes during business hours
```
0 */10 9-17 * * 1-5
```
This runs every 10 minutes from 9 AM to 5 PM, Monday through Friday.

### First day of every month at 3 AM
```
0 0 3 1 * *
```
This runs at 3:00 AM on the first day of each month.

### Every 15 minutes with even load distribution
```
0 ~/15 * * * *
```
This distributes jobs evenly across the hour in 15-minute intervals, with the exact minute determined by the job name.

### Every quarter (Jan, Apr, Jul, Oct) on the first day at midnight
```
0 0 0 1 1,4,7,10 *
```
This runs at midnight on the first day of January, April, July, and October.

### Every weekend at 10 PM
```
0 0 22 * * 0,6
```
This runs at 10:00 PM on Saturdays and Sundays.

### Every day at 8 AM New York time
```
0 0 8 * * *
```
With `timezone` parameter set to "America/New_York" in the job configuration.

## Best Practices

1. **Avoid Running Too Frequently**: Consider resource usage when scheduling frequent jobs. Running jobs every few seconds can put unnecessary load on your system.

2. **Use Even Load Distribution**: For jobs that don't need to run at an exact time, use the tilde (~) character to distribute load evenly.

3. **Consider Time Zones Carefully**: Be explicit about time zones in distributed systems to avoid confusion, especially if your servers are in different locations.

4. **Avoid Overlapping Executions**: For long-running jobs, ensure the schedule interval is longer than the expected execution time, or enable concurrency controls in the job configuration.

5. **Use Descriptive Job Names**: With the tilde (~) feature, job names influence scheduling, so use consistent naming conventions.

6. **Test Complex Expressions**: Use tools like [crontab.guru](https://crontab.guru/) to validate your cron expressions (note: these tools typically use 5-field format, while Dkron uses 6 fields with seconds).

7. **Document Job Schedules**: Maintain documentation about why jobs are scheduled at specific times to help with maintenance and troubleshooting.
