# REST API

You can communicate with Dcron using a RESTful JSON API over HTTP. Dcron nodes usually listen on port 8080 for API requests. All examples in this section assume that you've found a running leader at dcron-node:8080.

Dcron implements a RESTful JSON API over HTTP to communicate with software clients. Dcron listens in port 8080 by default. All examples in this section assume that you're using the default port.

[Leaders](#leaders)

## Leaders

When you have multiple Dcron nodes in server mode running, only one of them will be elected as the leader. In Dcron you can talk to any node running in server mode all of them could handle your request but only the leader will actually run the scheduler.

## Index

- Endpoint: `/`
- Method: `GET`
- Example: `curl -XGET dcron-node:8080`

Returns info about the agent queried.

## Members

- Endpoint: `/members`
- Method: `GET`
- Example: `curl -XGET dcron-node:8080/members`

Returns the cluster member list.

## Leader

- Endpoint: `/leader`
- Method: `GET`
- Example: `curl -XGET dcron-node:8080/leader`

Returns details about the current leader.

## Get Jobs

- Endpoint: `/jobs`
- Method: `GET`
- Example: `curl -L -X GET dcron-node:8080/jobs`

## Create or update a Job

Creates a new job or updates an exsiting job based on it's `name`. The `schedule` can be any valid cron expression or an interval using the `@every Xs` format.

- Endpoint: `/jobs/`
- Method: `POST` or `PUT`
- Example: `curl -X POST dcron-node:8080/jobs/ -d @jobs.json`

Sample job:

```
{
   "name":"cron_job",
   "schedule":"@every 2s",
   "command":"date",
   "owner":"foo",
   "owner_email":"foo@bar.com"
}
```

## Deleting a Job

Get a job name from the job listing above. Then:

- Endpoint: `/jobs/jobName`
- Method: `DELETE`
- Example: `curl -L -X DELETE dcron-node:8080/job/aggregate_stats`

Delete a job definition.

## Manually Starting a Job

You can manually start a job by issuing an HTTP request.

- Endpoint: `/jobs/job_name`
- Method: `POST` or `PUT`
- Example: `curl -L -X POST dcron-node:8080/jobs/aggregate_stats`

Will run `aggregate_stats` job.

## Get job executions

- Endpoint: `/executions/job`
- Method: `GET`
- Example: `curl -L -X GET dcron-node:8080/executions/aggregate_stats`

Get a list with the job executions.
