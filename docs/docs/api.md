# REST API

You can communicate with dkron using a RESTful JSON API over HTTP. dkron nodes usually listen on port 8080 for API requests. All examples in this section assume that you've found a running leader at dkron-node:8080.

dkron implements a RESTful JSON API over HTTP to communicate with software clients. dkron listens in port 8080 by default. All examples in this section assume that you're using the default port.

[Leaders](#leaders)

## Leaders

When you have multiple dkron nodes in server mode running, only one of them will be elected as the leader. In dkron you can talk to any node running in server mode all of them could handle your request but only the leader will actually run the scheduler.

## Index

- Endpoint: `/v1`
- Method: `GET`
- Example: `curl -XGET dkron-node:8080/v1`

Returns info about the agent queried.

## Members

- Endpoint: `/v1/members`
- Method: `GET`
- Example: `curl -XGET dkron-node:8080/v1/members`

Returns the cluster member list.

## Leader

- Endpoint: `/v1/leader`
- Method: `GET`
- Example: `curl -XGET dkron-node:8080/v1/leader`

Returns details about the current leader.

## Get Jobs

- Endpoint: `/v1/jobs`
- Method: `GET`
- Example: `curl -L -X GET dkron-node:8080/v1/jobs`

## Create or update a Job

Creates a new job or updates an exsiting job based on it's `name`. The `schedule` can be any valid cron expression or an interval using the `@every Xs` format.

- Endpoint: `/v1/jobs/`
- Method: `POST` or `PUT`
- Example: `curl -X POST dkron-node:8080/v1/jobs/ -d @jobs.json`

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

- Endpoint: `/v1/jobs/jobName`
- Method: `DELETE`
- Example: `curl -L -X DELETE dkron-node:8080/v1/job/aggregate_stats`

Delete a job definition.

## Manually Starting a Job

You can manually start a job by issuing an HTTP request.

- Endpoint: `/v1/jobs/job_name`
- Method: `POST` or `PUT`
- Example: `curl -L -X POST dkron-node:8080/v1/jobs/aggregate_stats`

Will run `aggregate_stats` job.

## Get job executions

- Endpoint: `/v1/executions/job`
- Method: `GET`
- Example: `curl -L -X GET dkron-node:8080/v1/executions/aggregate_stats`

Get a list with the job executions.
