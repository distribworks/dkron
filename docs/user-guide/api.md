# REST API

You can communicate with Dcron using a RESTful JSON API over HTTP. Dcron nodes usually listen on port 8080 for API requests. All examples in this section assume that you've found a running leader at dcron-node:8080.

Dcron implements a RESTful JSON API over HTTP to communicate with software clients. Dcron listens in port 8080 by default. All examples in this section assume that you're using the default port.

[Leaders](#leaders)

## Leaders

When you have multiple Dcron nodes in server mode running, only one of them will be elected as the leader. In Dcron you can talk to any node running in server mode all of them could handle your request but only the leader will actually run the scheduler.

## Index

- Endpoint: /
- Method: GET
- Example: curl -XGET dcron-node:8080

It returns info about the agent queried.

## Get Jobs

Endpoint: /jobs
Method: GET
Example: curl -L -X GET dcron-node:8080/jobs
A job listing returns a JSON list containing all of the jobs. Each job is a JSON hash. Interesting fields in the hashes are:

invocationCount: the number of times the job completed
executor: auto-determined by Chronos, but will usually be "" for non-async jobs
parents: for dependent jobs, a list of all other jobs that must run before this job will do so
If there is a parents field there will be no schedule field and vice-versa.

## Manually Starting a Job
You can manually start a job by issuing an HTTP request.

Endpoint: /job
Method: PUT
Query string parameters: arguments - optional string with a list of command line arguments that is appended to job's command
Example: curl -L -X PUT chronos-node:8080/scheduler/job/request_event_counter_hourly
Example: curl -L -X PUT chronos-node:8080/scheduler/job/job_name?arguments=-debug
Response: HTTP 204

