# Guides

## Target nodes spec

You can choose whether a job is run on a node or nodes by specifying tags and a count of target nodes having this tag do you want a job to run.

### Examples:

Target all nodes with a tag:

```
{
    "name": "job_name",
    "command": "/bin/true",
    "schedule": "@every 2m",
    "tags": {
        "role": "web"
    }
}
```

Target only two nodes of a group of nodes with a tag:

```
{
    "name": "job_name",
    "command": "/bin/true",
    "schedule": "@every 2m",
    "tags": {
        "role": "web:2"
    }
}
```

Dkron will try to run the job in the amount of nodes indicated by that count having that tag.

## CRON Expression Format

A cron expression represents a set of times, using 6 space-separated fields.

	Field name   | Mandatory? | Allowed values  | Allowed special characters
	----------   | ---------- | --------------  | --------------------------
	Seconds      | Yes        | 0-59            | * / , -
	Minutes      | Yes        | 0-59            | * / , -
	Hours        | Yes        | 0-23            | * / , -
	Day of month | Yes        | 1-31            | * / , - ?
	Month        | Yes        | 1-12 or JAN-DEC | * / , -
	Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?

Note: Month and Day-of-week field values are case insensitive.  "SUN", "Sun",
and "sun" are equally accepted.

Special Characters

Asterisk ( * )

The asterisk indicates that the cron expression will match for all values of the
field; e.g., using an asterisk in the 5th field (month) would indicate every
month.

Slash ( / )

Slashes are used to describe increments of ranges. For example 3-59/15 in the
1st field (minutes) would indicate the 3rd minute of the hour and every 15
minutes thereafter. The form "*\/..." is equivalent to the form "first-last/...",
that is, an increment over the largest possible range of the field.  The form
"N/..." is accepted as meaning "N-MAX/...", that is, starting at N, use the
increment until the end of that specific range.  It does not wrap around.

Comma ( , )

Commas are used to separate items of a list. For example, using "MON,WED,FRI" in
the 5th field (day of week) would mean Mondays, Wednesdays and Fridays.

Hyphen ( - )

Hyphens are used to define ranges. For example, 9-17 would indicate every
hour between 9am and 5pm inclusive.

Question mark ( ? )

Question mark may be used instead of '*' for leaving either day-of-month or
day-of-week blank.

### Predefined schedules

You may use one of several pre-defined schedules in place of a cron expression.

	Entry                  | Description                                | Equivalent To
	-----                  | -----------                                | -------------
	@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
	@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
	@weekly                | Run once a week, midnight on Sunday        | 0 0 0 * * 0
	@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
	@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *
	@minutely              | Run once a minute, beginning of minute     | 0 * * * * *

### Intervals

You may also schedule a job to execute at fixed intervals.  This is supported by
formatting the cron spec like this:

    @every <duration>

where "duration" is a string accepted by time.ParseDuration
(http://golang.org/pkg/time/#ParseDuration).

For example, "@every 1h30m10s" would indicate a schedule that activates every
1 hour, 30 minutes, 10 seconds.

Note: The interval does not take the job runtime into account.  For example,
if a job takes 3 minutes to run, and it is scheduled to run every 5 minutes,
it will have only 2 minutes of idle time between each run.

### Fixed times

You may also want to schedule a job to be executed once. This is supported by
formatting the cron spec like this:

    @at <datetime>

Where "datetime" is a string accepted by time.Parse in RFC3339 format
(https://golang.org/pkg/time/#Parse).

For example, "@at 2018-01-02T15:04:00Z" would run the job on the specified date and time
assuming UTC timezone.

### Time zones

All interpretation and scheduling is done in the machine's local time zone (as
provided by the Go time package (http://www.golang.org/pkg/time).

Be aware that jobs scheduled during daylight-savings leap-ahead transitions will
not be run!

## Job chaining

You can set some jobs to run after other job is executed. To setup a job that will be executed after any other given job, just set the `parent_job` property when saving the new job.

The dependent job will be executed after the main job finished a successful execution.

Child jobs schedule property will be ignored if it's present.

## Logging output

Logging output of each job execution can be modified by using processor plugins.

Processor plugins can be used to redirect the output of a job execution to different targets.

Refer to the [plugins documentation](plugins.md) from more info

## Concurrency

Jobs can be configured to allow overlapping executions or forbid them. 

Concurrency property accepts two option: 

* **allow** (default): Allow concurrent job executions.
* **forbid**: If the job is already running don't send the execution, it will skip the executions until the next schedule.
