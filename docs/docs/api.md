# Dkron


<a name="overview"></a>
## Overview
# REST API

You can communicate with Dkron using a RESTful JSON API over HTTP. Dkron nodes usually listen on port `8080` for API requests. All examples in this section assume that you've found a running leader at `dkron-node:8080`.

Dkron implements a RESTful JSON API over HTTP to communicate with software clients. Dkron listens in port `8080` by default. All examples in this section assume that you're using the default port.

Default API responses are unformatted JSON add the `pretty=true` param to format the response.


### Version information
*Version* : 0.7.2


### URI scheme
*Host* : localhost:8080  
*BasePath* : /v1


### Consumes

* `application/json`


### Produces

* `application/json`




<a name="paths"></a>
## Paths

<a name="get"></a>
### GET /

#### Description
Gets `Status` object.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|[Response 200](#get-response-200)|

<a name="get-response-200"></a>
**Response 200**

|Name|Description|Schema|
|---|---|---|
|**agent**  <br>*optional*||[agent](#agent)|
|**serf**  <br>*optional*||[serf](#serf)|
|**tags**  <br>*optional*||[tags](#tags)|


#### Example HTTP response

##### Response 200
```
json :
{
  "application/json" : "{\n     \"agent\": {\n        \"name\": \"dkron2\",\n        \"version\": \"0.7.2\"\n      },\n      \"serf\": {\n        \"encrypted\": \"false\",\n        \"...\": \"...\"\n      },\n      \"tags\": {\n        \"role\": \"web\",\n        \"dkron_server\": true\n      }\n}"
}
```


<a name="executions-job_name-get"></a>
### GET /executions/{job_name}

#### Description
List executions.


#### Parameters

|Type|Name|Description|Schema|Default|
|---|---|---|---|---|
|**Path**|**job_name**  <br>*required*|The job that owns the executions to be fetched.|string||


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|< [execution](#execution) > array|


<a name="jobs-post"></a>
### POST /jobs

#### Description
Create or updates a new job.
# Expected responses for this operation:


#### Parameters

|Type|Name|Description|Schema|Default|
|---|---|---|---|---|
|**Body**|**body**  <br>*required*|Updated job object|[job](#job)||


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**201**|Successful response|[job](#job)|


<a name="jobs-get"></a>
### GET /jobs

#### Description
List jobs.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|< [job](#job) > array|


<a name="jobs-job_name-post"></a>
### POST /jobs/{job_name}

#### Description
Executes a job.


#### Parameters

|Type|Name|Description|Schema|Default|
|---|---|---|---|---|
|**Path**|**job_name**  <br>*required*|The job that needs to be run.|string||


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|[job](#job)|


<a name="jobs-job_name-get"></a>
### GET /jobs/{job_name}

#### Description
Show a job.


#### Parameters

|Type|Name|Description|Schema|Default|
|---|---|---|---|---|
|**Path**|**job_name**  <br>*required*|The job that needs to be fetched.|string||


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|[job](#job)|


<a name="jobs-job_name-delete"></a>
### DELETE /jobs/{job_name}

#### Description
Delete a job.


#### Parameters

|Type|Name|Description|Schema|Default|
|---|---|---|---|---|
|**Path**|**job_name**  <br>*required*|The job that needs to be deleted.|string||


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|[job](#job)|


<a name="leader-get"></a>
### GET /leader

#### Description
List members.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|[member](#member)|


<a name="members-get"></a>
### GET /members

#### Description
List members.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|< [member](#member) > array|




<a name="definitions"></a>
## Definitions

<a name="agent"></a>
### agent
Node basic details

*Type* : object


<a name="execution"></a>
### execution
An execution represents a timed job run.


|Name|Description|Schema|
|---|---|---|
|**finished_at**  <br>*optional*||string(date-time)|
|**job_name**  <br>*optional*||string|
|**node_name**  <br>*optional*||string|
|**output**  <br>*optional*||string|
|**started_at**  <br>*optional*||string(date-time)|
|**success**  <br>*optional*||boolean|


<a name="job"></a>
### job
A Job represents a scheduled task to execute.


|Name|Description|Schema|
|---|---|---|
|**command**  <br>*required*|Command to run.|string|
|**disabled**  <br>*optional*||boolean|
|**error_count**  <br>*optional*||integer|
|**last_error**  <br>*optional*||string(date-time)|
|**last_success**  <br>*optional*||string(date-time)|
|**name**  <br>*required*||string|
|**owner**  <br>*optional*||string|
|**owner_email**  <br>*optional*||string|
|**schedule**  <br>*required*|Cron expression for the job. <br>0 0 0 * * *|string|
|**shell**  <br>*optional*||boolean|
|**success_count**  <br>*optional*||integer|
|**tags**  <br>*optional*||[tags](#tags)|


<a name="member"></a>
### member
A member represents a cluster member node.


|Name|Description|Schema|
|---|---|---|
|**Addr**  <br>*optional*||string|
|**DelegateCur**  <br>*optional*||integer|
|**DelegateMax**  <br>*optional*||integer|
|**DelegateMin**  <br>*optional*||integer|
|**Name**  <br>*optional*||string|
|**Port**  <br>*optional*||integer|
|**ProtocolCur**  <br>*optional*||integer|
|**ProtocolMax**  <br>*optional*||integer|
|**ProtocolMin**  <br>*optional*||integer|
|**Status**  <br>*optional*||integer|
|**Tags**  <br>*optional*||[tags](#tags)|


<a name="serf"></a>
### serf
Serf status

*Type* : object


<a name="tags"></a>
### tags
Tags asociated with this node

*Type* : object





