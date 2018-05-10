---
title: Dkron REST API
toc: true
---

<a name="overview"></a>
## Overview
You can communicate with Dkron using a RESTful JSON API over HTTP. Dkron nodes usually listen on port `8080` for API requests. All examples in this section assume that you've found a running leader at `localhost:8080`.

Dkron implements a RESTful JSON API over HTTP to communicate with software clients. Dkron listens in port `8080` by default. All examples in this section assume that you're using the default port.

Default API responses are unformatted JSON add the `pretty=true` param to format the response.


### Version information
*Version* : 0.9.5


### URI scheme
*Host* : localhost:8080  
*BasePath* : /v1  
*Schemes* : HTTP


### Consumes

* `application/json`


### Produces

* `application/json`




<a name="paths"></a>
## Paths

<a name="status"></a>
### GET /

#### Description
Gets `Status` object.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|[status](#status)|


#### Tags

* default


#### Example HTTP request

##### Request path
```
/
```


#### Example HTTP response

##### Response 200
```
json :
"{ }"
```


<a name="getjobs"></a>
### GET /jobs

#### Description
List jobs.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|< [job](#job) > array|


#### Tags

* jobs


#### Example HTTP request

##### Request path
```
/jobs
```


#### Example HTTP response

##### Response 200
```
json :
"array"
```


<a name="createorupdatejob"></a>
### POST /jobs

#### Description
Create or updates a new job.


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Body**|**body**  <br>*required*|Updated job object|[job](#job)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**201**|Successful response|[job](#job)|


#### Tags

* jobs


#### Example HTTP request

##### Request path
```
/jobs
```


##### Request body
```
json :
{
  "name" : "string",
  "schedule" : "string",
  "command" : "string",
  "shell" : true,
  "owner" : "string",
  "owner_email" : "string",
  "success_count" : 0,
  "error_count" : 0,
  "last_success" : "string",
  "last_error" : "string",
  "disabled" : true,
  "tags" : {
    "string" : "string"
  },
  "retries" : 2,
  "parent_job" : "parent_job",
  "dependent_jobs" : [ "string" ],
  "processors" : {
    "string" : "string"
  },
  "concurrency" : "allow",
  "timezone": "string"
}
```


#### Example HTTP response

##### Response 201
```
json :
{
  "name" : "string",
  "schedule" : "string",
  "command" : "string",
  "shell" : true,
  "owner" : "string",
  "owner_email" : "string",
  "success_count" : 0,
  "error_count" : 0,
  "last_success" : "string",
  "last_error" : "string",
  "disabled" : true,
  "tags" : {
    "string" : "string"
  },
  "retries" : 2,
  "parent_job" : "parent_job",
  "dependent_jobs" : [ "string" ],
  "processors" : {
    "string" : "string"
  },
  "concurrency" : "allow",
  "timezone": "string"
}
```


<a name="showjobbyname"></a>
### GET /jobs/{job_name}

#### Description
Show a job.


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Path**|**job_name**  <br>*required*|The job that needs to be fetched.|string|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|[job](#job)|


#### Tags

* jobs


#### Example HTTP request

##### Request path
```
/jobs/string
```


#### Example HTTP response

##### Response 200
```
json :
{
  "name" : "string",
  "schedule" : "string",
  "command" : "string",
  "shell" : true,
  "owner" : "string",
  "owner_email" : "string",
  "success_count" : 0,
  "error_count" : 0,
  "last_success" : "string",
  "last_error" : "string",
  "disabled" : true,
  "tags" : {
    "string" : "string"
  },
  "retries" : 2,
  "parent_job" : "parent_job",
  "dependent_jobs" : [ "string" ],
  "processors" : {
    "string" : "string"
  },
  "concurrency" : "allow",
  "timezone": "string"
}
```


<a name="runjob"></a>
### POST /jobs/{job_name}

#### Description
Executes a job.


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Path**|**job_name**  <br>*required*|The job that needs to be run.|string|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**202**|Successful response|[job](#job)|


#### Tags

* jobs


#### Example HTTP request

##### Request path
```
/jobs/string
```


#### Example HTTP response

##### Response 202
```
json :
{
  "name" : "string",
  "schedule" : "string",
  "command" : "string",
  "shell" : true,
  "owner" : "string",
  "owner_email" : "string",
  "success_count" : 0,
  "error_count" : 0,
  "last_success" : "string",
  "last_error" : "string",
  "disabled" : true,
  "tags" : {
    "string" : "string"
  },
  "retries" : 2,
  "parent_job" : "parent_job",
  "dependent_jobs" : [ "string" ],
  "processors" : {
    "string" : "string"
  },
  "concurrency" : "allow",
  "timezone": "string"
}
```


<a name="deletejob"></a>
### DELETE /jobs/{job_name}

#### Description
Delete a job.


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Path**|**job_name**  <br>*required*|The job that needs to be deleted.|string|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|[job](#job)|


#### Tags

* jobs


#### Example HTTP request

##### Request path
```
/jobs/string
```


#### Example HTTP response

##### Response 200
```
json :
{
  "name" : "string",
  "schedule" : "string",
  "command" : "string",
  "shell" : true,
  "owner" : "string",
  "owner_email" : "string",
  "success_count" : 0,
  "error_count" : 0,
  "last_success" : "string",
  "last_error" : "string",
  "disabled" : true,
  "tags" : {
    "string" : "string"
  },
  "retries" : 2,
  "parent_job" : "parent_job",
  "dependent_jobs" : [ "string" ],
  "processors" : {
    "string" : "string"
  },
  "concurrency" : "allow",
  "timezone": "string"
}
```


<a name="listexecutionsbyjob"></a>
### GET /jobs/{job_name}/executions

#### Description
List executions.


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Path**|**job_name**  <br>*required*|The job that owns the executions to be fetched.|string|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|< [execution](#execution) > array|


#### Tags

* executions


#### Example HTTP request

##### Request path
```
/jobs/string/executions
```


#### Example HTTP response

##### Response 200
```
json :
"array"
```


<a name="getleader"></a>
### GET /leader

#### Description
List leader of cluster.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|[member](#member)|


#### Tags

* default


#### Example HTTP request

##### Request path
```
/leader
```


#### Example HTTP response

##### Response 200
```
json :
{
  "Name" : "string",
  "Addr" : "string",
  "Port" : 0,
  "Tags" : {
    "string" : "string"
  },
  "Status" : 0,
  "ProtocolMin" : 0,
  "ProtocolMax" : 0,
  "ProtocolCur" : 0,
  "DelegateMin" : 0,
  "DelegateMax" : 0,
  "DelegateCur" : 0
}
```


<a name="leave"></a>
### GET /leave

#### Description
Force the node to leave the cluster.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|< [member](#member) > array|


#### Tags

* default


#### Example HTTP request

##### Request path
```
/leave
```


#### Example HTTP response

##### Response 200
```
json :
"array"
```


<a name="getmember"></a>
### GET /members

#### Description
List members.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|< [member](#member) > array|


#### Tags

* members


#### Example HTTP request

##### Request path
```
/members
```


#### Example HTTP response

##### Response 200
```
json :
"array"
```




<a name="definitions"></a>
## Definitions

<a name="status"></a>
### status
Status represents details about the node.

*Type* : object


<a name="job"></a>
### job
A Job represents a scheduled task to execute.


|Name|Description|Schema|
|---|---|---|
|**name**  <br>*required*|Name for the job.  <br>**Example** : `"string"`|string|
|**schedule**  <br>*required*|Cron expression for the job.  <br>**Example** : `"string"`|string|
|**command**  <br>*required*|Command to run.  <br>**Example** : `"string"`|string|
|**timezone**  <br>*optional*|Timezone where the schedule will be executed. <br>**Example** : `"Europe/Paris"`|string|
|**shell**  <br>*optional*|Use shell to run the command  <br>**Example** : `true`|boolean|
|**owner**  <br>*optional*|Owner of the job  <br>**Example** : `"string"`|string|
|**owner_email**  <br>*optional*|Email of the owner  <br>**Example** : `"string"`|string|
|**success_count**  <br>*optional*  <br>*read-only*|Number of successful executions  <br>**Example** : `0`|integer|
|**error_count**  <br>*optional*  <br>*read-only*|Number of failed executions  <br>**Example** : `0`|integer|
|**last_success**  <br>*optional*  <br>*read-only*|Last time this job executed successfully  <br>**Example** : `"string"`|string(date-time)|
|**last_error**  <br>*optional*  <br>*read-only*|Last time this job failed  <br>**Example** : `"string"`|string(date-time)|
|**disabled**  <br>*optional*|Disabled state of the job  <br>**Example** : `true`|boolean|
|**tags**  <br>*optional*|Target nodes tags of this job  <br>**Example** : `{<br>  "string" : "string"<br>}`|< string, string > map|
|**retries**  <br>*optional*|Number of times to retry a failed job execution  <br>**Example** : `2`|integer|
|**parent_job**  <br>*optional*|The name/id of the job that will trigger the execution of this job  <br>**Example** : `"parent_job"`|string|
|**dependent_jobs**  <br>*optional*|Array containing the jobs that depends on this one  <br>**Example** : `[ "string" ]`|< string > array|
|**processors**  <br>*optional*|Processor plugins used to process executions results of this job  <br>**Example** : `{<br>  "string" : "string"<br>}`|< string, string > map|
|**concurrency**  <br>*optional*|Concurrency policy for the job allow/forbid  <br>**Example** : `"allow"`|string|


<a name="member"></a>
### member
A member represents a cluster member node.


|Name|Description|Schema|
|---|---|---|
|**Name**  <br>*optional*|Node name  <br>**Example** : `"string"`|string|
|**Addr**  <br>*optional*|IP Address  <br>**Example** : `"string"`|string|
|**Port**  <br>*optional*|Port number  <br>**Example** : `0`|integer|
|**Tags**  <br>*optional*|Tags asociated with this node  <br>**Example** : `{<br>  "string" : "string"<br>}`|< string, string > map|
|**Status**  <br>*optional*|The serf status of the node see: https://godoc.org/github.com/hashicorp/serf/serf#MemberStatus  <br>**Example** : `0`|integer|
|**ProtocolMin**  <br>*optional*|Serf protocol minimum version this node can understand or speak  <br>**Example** : `0`|integer|
|**ProtocolMax**  <br>*optional*|Serf protocol maximum version this node can understand or speak  <br>**Example** : `0`|integer|
|**ProtocolCur**  <br>*optional*|Serf protocol current version this node can understand or speak  <br>**Example** : `0`|integer|
|**DelegateMin**  <br>*optional*|Serf delegate protocol minimum version this node can understand or speak  <br>**Example** : `0`|integer|
|**DelegateMax**  <br>*optional*|Serf delegate protocol maximum version this node can understand or speak  <br>**Example** : `0`|integer|
|**DelegateCur**  <br>*optional*|Serf delegate protocol current version this node can understand or speak  <br>**Example** : `0`|integer|


<a name="execution"></a>
### execution
An execution represents a timed job run.


|Name|Description|Schema|
|---|---|---|
|**job_name**  <br>*optional*|job name  <br>**Example** : `"string"`|string|
|**started_at**  <br>*optional*|start time of the execution  <br>**Example** : `"string"`|string(date-time)|
|**finished_at**  <br>*optional*|when the execution finished running  <br>**Example** : `"string"`|string(date-time)|
|**success**  <br>*optional*|the execution run successfuly  <br>**Example** : `true`|boolean|
|**output**  <br>*optional*|partial output of the command execution  <br>**Example** : `"string"`|string|
|**node_name**  <br>*optional*|name of the node that executed the command  <br>**Example** : `"string"`|string|





