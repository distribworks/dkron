# REST API

You can communicate with dkron using a RESTful JSON API over HTTP. dkron nodes usually listen on port `8080` for API requests. All examples in this section assume that you've found a running leader at `dkron-node:8080`.

dkron implements a RESTful JSON API over HTTP to communicate with software clients. dkron listens in port `8080` by default. All examples in this section assume that you're using the default port.

Default API responses are unformatted JSON add the `pretty=true` param to format the response.

## Status

Status represents details about the node.

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **agent** | *object* | Node basic details | `{"name":"dkron2","version":"0.0.4"}` |
| **serf** | *object* | Serf status | `{"encrypted":"false","...":"..."}` |
| **tags** | *object* | Tags asociated with this node | `{"role":"web","server":"true"}` |

### Status 

Status.

```
GET /v1/
```


#### Curl Example

```bash
$ curl -n dkron-node:8080/v1/
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "agent": {
    "name": "dkron2",
    "version": "0.0.4"
  },
  "serf": {
    "encrypted": "false",
    "...": "..."
  },
  "tags": {
    "role": "web",
    "server": "true"
  }
}
```


## Job

A Job represents a scheduled task to execute.

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **name** | *string* | job name | `"cron_job"` |
| **schedule** | *string* | cron expression for the job | `"0 30 * * * *"` |
| **command** | *string* | command to run. Must be a shell command to execute | `"/usr/bin/date"` |
| **owner** | *string* | owner of the job | `"John Doe"` |
| **owner_email** | *email* | email of the owner | `"john@doe.com"` |
| **run_as_user** | *hostname* | the user to use to run the job | `"johndoe"` |
| **success_count** | *integer* | number of successful executions | `20` |
| **error_count** | *integer* | number of failed executions | `5` |
| **last_success** | *date-time* | last time this job executed successfully | `"0001-01-01T00:00:00Z"` |
| **last_error** | *date-time* | last time this job failed | `"0001-01-01T100:00:00Z"` |
| **disabled** | *boolean* | disabled state of the job | `false` |
| **tags** | *object* | tags of the target server to run this job | `{"role":"web"}` |

### Job Create or update

Create or updates a new job.

```
POST /v1/jobs
```

#### Required Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **name** | *string* | job name | `"cron_job"` |
| **schedule** | *string* | cron expression for the job | `"0 30 * * * *"` |
| **command** | *string* | command to run. Must be a shell command to execute | `"/usr/bin/date"` |


#### Optional Parameters

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **owner** | *string* | owner of the job | `"John Doe"` |
| **owner_email** | *email* | email of the owner | `"john@doe.com"` |
| **run_as_user** | *hostname* | the user to use to run the job | `"johndoe"` |
| **disabled** | *boolean* | disabled state of the job | `false` |
| **tags** | *object* | tags of the target server to run this job | `{"role":"web"}` |


#### Curl Example

```bash
$ curl -n -X POST dkron-node:8080/v1/jobs \
  -H "Content-Type: application/json" \
 \
  -d '{
  "name": "cron_job",
  "schedule": "0 30 * * * *",
  "command": "/usr/bin/date",
  "owner": "John Doe",
  "owner_email": "john@doe.com",
  "run_as_user": "johndoe",
  "disabled": false,
  "tags": {
    "role": "web"
  }
}'
```


#### Response Example

```
HTTP/1.1 201 Created
```

```json
{
  "name": "cron_job",
  "schedule": "0 30 * * * *",
  "command": "/usr/bin/date",
  "owner": "John Doe",
  "owner_email": "john@doe.com",
  "run_as_user": "johndoe",
  "success_count": 20,
  "error_count": 5,
  "last_success": "0001-01-01T00:00:00Z",
  "last_error": "0001-01-01T100:00:00Z",
  "disabled": false,
  "tags": {
    "role": "web"
  }
}
```

### Job Show

Show job.

```
GET /v1/jobs/{job_name}
```


#### Curl Example

```bash
$ curl -n dkron-node:8080/v1/jobs/$JOB_NAME
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "name": "cron_job",
  "schedule": "0 30 * * * *",
  "command": "/usr/bin/date",
  "owner": "John Doe",
  "owner_email": "john@doe.com",
  "run_as_user": "johndoe",
  "success_count": 20,
  "error_count": 5,
  "last_success": "0001-01-01T00:00:00Z",
  "last_error": "0001-01-01T100:00:00Z",
  "disabled": false,
  "tags": {
    "role": "web"
  }
}
```

### Job Delete

Delete job.

```
DELETE /v1/jobs/{job_name}
```


#### Curl Example

```bash
$ curl -n -X DELETE dkron-node:8080/v1/jobs/$JOB_NAME \
  -H "Content-Type: application/json" \
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "name": "cron_job",
  "schedule": "0 30 * * * *",
  "command": "/usr/bin/date",
  "owner": "John Doe",
  "owner_email": "john@doe.com",
  "run_as_user": "johndoe",
  "success_count": 20,
  "error_count": 5,
  "last_success": "0001-01-01T00:00:00Z",
  "last_error": "0001-01-01T100:00:00Z",
  "disabled": false,
  "tags": {
    "role": "web"
  }
}
```

### Job List

List jobs.

```
GET /v1/jobs
```


#### Curl Example

```bash
$ curl -n dkron-node:8080/v1/jobs
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
[
  {
    "name": "cron_job",
    "schedule": "0 30 * * * *",
    "command": "/usr/bin/date",
    "owner": "John Doe",
    "owner_email": "john@doe.com",
    "run_as_user": "johndoe",
    "success_count": 20,
    "error_count": 5,
    "last_success": "0001-01-01T00:00:00Z",
    "last_error": "0001-01-01T100:00:00Z",
    "disabled": false,
    "tags": {
      "role": "web"
    }
  }
]
```

### Job Run

Run job.

```
POST /v1/jobs/{job_name}
```


#### Curl Example

```bash
$ curl -n -X POST dkron-node:8080/v1/jobs/$JOB_NAME \
  -H "Content-Type: application/json" \
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "name": "cron_job",
  "schedule": "0 30 * * * *",
  "command": "/usr/bin/date",
  "owner": "John Doe",
  "owner_email": "john@doe.com",
  "run_as_user": "johndoe",
  "success_count": 20,
  "error_count": 5,
  "last_success": "0001-01-01T00:00:00Z",
  "last_error": "0001-01-01T100:00:00Z",
  "disabled": false,
  "tags": {
    "role": "web"
  }
}
```


## Member

A member represents a cluster member node.

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **Name** | *boolean* | Node name | `"dkron1"` |
| **Addr** | *string* | IP Address | `"10.0.0.1"` |
| **Port** | *integer* | Port number | `5001` |
| **Tags** | *object* | Tags asociated with this node | `{"role":"web","server":"true"}` |
| **Status** | *integer* | The serf status of the node see: https://godoc.org/github.com/hashicorp/serf/serf#MemberStatus | `1` |
| **ProtocolMin** | *integer* | Serf protocol minimum version this node can understand or speak | `1` |
| **ProtocolMax** | *integer* | Serf protocol minimum version this node can understand or speak | `2` |
| **ProtocolCur** | *integer* | Serf protocol current version this node can understand or speak | `2` |
| **DelegateMin** | *integer* | Serf delegate protocol minimum version this node can understand or speak | `2` |
| **DelegateMax** | *integer* | Serf delegate protocol minimum version this node can understand or speak | `4` |
| **DelegateCur** | *integer* | Serf delegate protocol minimum version this node can understand or speak | `4` |

### Member List

List members.

```
GET /v1/members
```


#### Curl Example

```bash
$ curl -n dkron-node:8080/v1/members
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
[
  {
    "Name": "dkron1",
    "Addr": "10.0.0.1",
    "Port": 5001,
    "Tags": {
      "role": "web",
      "server": "true"
    },
    "Status": 1,
    "ProtocolMin": 1,
    "ProtocolMax": 2,
    "ProtocolCur": 2,
    "DelegateMin": 2,
    "DelegateMax": 4,
    "DelegateCur": 4
  }
]
```

### Member Leader

Show leader member.

```
GET /v1/leader
```


#### Curl Example

```bash
$ curl -n dkron-node:8080/v1/leader
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
{
  "Name": "dkron1",
  "Addr": "10.0.0.1",
  "Port": 5001,
  "Tags": {
    "role": "web",
    "server": "true"
  },
  "Status": 1,
  "ProtocolMin": 1,
  "ProtocolMax": 2,
  "ProtocolCur": 2,
  "DelegateMin": 2,
  "DelegateMax": 4,
  "DelegateCur": 4
}
```


## Execution

An execution represents a timed job run.

### Attributes

| Name | Type | Description | Example |
| ------- | ------- | ------- | ------- |
| **job_name** | *string* | job name | `"cron_job"` |
| **started_at** | *date-time* | start time of the execution | `"2012-01-01T12:00:00Z"` |
| **finished_at** | *date-time* | when the execution finished running | `"2012-01-01T12:00:00Z"` |
| **success** | *boolean* | the execution run successfuly | `true` |
| **output** | *string* | partial output of the command execution | `"Sat Sep  5 23:27:10 CEST 2015"` |
| **node_name** | *string* | name of the node that executed the command | `"dkron-node1"` |

### Execution List

List executions.

```
GET /v1/executions/{execution_job_name}
```


#### Curl Example

```bash
$ curl -n dkron-node:8080/v1/executions/$EXECUTION_JOB_NAME
```


#### Response Example

```
HTTP/1.1 200 OK
```

```json
[
  {
    "job_name": "cron_job",
    "started_at": "2012-01-01T12:00:00Z",
    "finished_at": "2012-01-01T12:00:00Z",
    "success": true,
    "output": "Sat Sep  5 23:27:10 CEST 2015",
    "node_name": "dkron-node1"
  }
]
```


