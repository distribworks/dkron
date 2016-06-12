# Dkron REST API


<a name="overview"></a>
## Overview
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


#### Tags

* main


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


#### Tags

* executions


<a name="jobs-get"></a>
### GET /jobs

#### Description
List jobs.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|< [job](#job) > array|


#### Tags

* jobs


<a name="jobs-post"></a>
### POST /jobs

#### Description
Create or updates a new job.


#### Parameters

|Type|Name|Description|Schema|Default|
|---|---|---|---|---|
|**Body**|**body**  <br>*required*|Updated job object|[job](#job)||


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**201**|Successful response|[job](#job)|


#### Tags

* jobs


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


#### Tags

* jobs


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


#### Tags

* jobs


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


#### Tags

* jobs


<a name="leader-get"></a>
### GET /leader

#### Description
List members.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|[member](#member)|


#### Tags

* main


<a name="members-get"></a>
### GET /members

#### Description
List members.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|< [member](#member) > array|


#### Tags

* members




<a name="definitions"></a>
## Definitions

<a name="agent"></a>
### agent
Node basic details

*Type* : object


<a name="serf"></a>
### serf
Serf status

*Type* : object


<a name="tags"></a>
### tags
Tags asociated with this node

*Type* : object


<a name="job"></a>
### job
A Job represents a scheduled task to execute.


|Name|Description|Schema|
|---|---|---|
|**name**  <br>*required*|-|string|
|**schedule**  <br>*required*|Cron expression for the job.|string|
|**command**  <br>*required*|Command to run.|string|
|**shell**  <br>*optional*|-|boolean|
|**owner**  <br>*optional*|-|string|
|**owner_email**  <br>*optional*|-|string|
|**success_count**  <br>*optional*|-|integer|
|**error_count**  <br>*optional*|-|integer|
|**last_success**  <br>*optional*|-|string(date-time)|
|**last_error**  <br>*optional*|-|string(date-time)|
|**disabled**  <br>*optional*|-|boolean|
|**tags**  <br>*optional*|-|[tags](#tags)|


<a name="member"></a>
### member
A member represents a cluster member node.


|Name|Description|Schema|
|---|---|---|
|**Name**  <br>*optional*|-|string|
|**Addr**  <br>*optional*|-|string|
|**Port**  <br>*optional*|-|integer|
|**Tags**  <br>*optional*|-|[tags](#tags)|
|**Status**  <br>*optional*|-|integer|
|**ProtocolMin**  <br>*optional*|-|integer|
|**ProtocolMax**  <br>*optional*|-|integer|
|**ProtocolCur**  <br>*optional*|-|integer|
|**DelegateMin**  <br>*optional*|-|integer|
|**DelegateMax**  <br>*optional*|-|integer|
|**DelegateCur**  <br>*optional*|-|integer|


<a name="execution"></a>
### execution
An execution represents a timed job run.


|Name|Description|Schema|
|---|---|---|
|**job_name**  <br>*optional*|-|string|
|**started_at**  <br>*optional*|-|string(date-time)|
|**finished_at**  <br>*optional*|-|string(date-time)|
|**success**  <br>*optional*|-|boolean|
|**output**  <br>*optional*|-|string|
|**node_name**  <br>*optional*|-|string|





