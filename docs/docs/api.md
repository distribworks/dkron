# Dkron REST API


<a name="overview"></a>
## Overview
You can communicate with Dkron using a RESTful JSON API over HTTP. Dkron nodes usually listen on port `8080` for API requests. All examples in this section assume that you've found a running leader at `localhost:8080`.

Dkron implements a RESTful JSON API over HTTP to communicate with software clients. Dkron listens in port `8080` by default. All examples in this section assume that you're using the default port.

Default API responses are unformatted JSON add the `pretty=true` param to format the response.


### Version information
*Version* : 0.9.0


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


<a name="listexecutionsbyjob"></a>
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


<a name="createorupdatejob"></a>
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


<a name="deletejob"></a>
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


<a name="showjobbyname"></a>
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


<a name="runjob"></a>
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


<a name="getleader"></a>
### GET /leader

#### Description
List members.


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|Successful response|[member](#member)|


#### Tags

* default


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
|**name**  <br>*required*|Name for the job.|string|
|**schedule**  <br>*required*|Cron expression for the job.|string|
|**command**  <br>*required*|Command to run.|string|
|**shell**  <br>*optional*|Use shell to run the command|boolean|
|**owner**  <br>*optional*|Owner of the job|string|
|**owner_email**  <br>*optional*|Email of the owner|string|
|**success_count**  <br>*optional*  <br>*read-only*|Number of successful executions|integer|
|**error_count**  <br>*optional*  <br>*read-only*|Number of failed executions|integer|
|**last_success**  <br>*optional*  <br>*read-only*|Last time this job executed successfully|string(date-time)|
|**last_error**  <br>*optional*  <br>*read-only*|Last time this job failed|string(date-time)|
|**disabled**  <br>*optional*|Disabled state of the job|boolean|
|**tags**  <br>*optional*|Tags asociated with this node|< string, string > map|
|**retries**  <br>*optional*|Number of times to retry a failed job execution  <br>**Example** : `2`|integer|
|**parent_job**  <br>*optional*|The name/id of the job that will trigger the execution of this job  <br>**Example** : `"parent_job"`|string|
|**dependent_jobs**  <br>*optional*  <br>*read-only*|Array containing the jobs that depends on this one  <br>**Example** : `""`|string|


<a name="member"></a>
### member
A member represents a cluster member node.


|Name|Description|Schema|
|---|---|---|
|**Name**  <br>*optional*|Node name|string|
|**Addr**  <br>*optional*|IP Address|string|
|**Port**  <br>*optional*|Port number|integer|
|**Tags**  <br>*optional*|Tags asociated with this node|< string, string > map|
|**Status**  <br>*optional*|The serf status of the node see: https://godoc.org/github.com/hashicorp/serf/serf#MemberStatus|integer|
|**ProtocolMin**  <br>*optional*|Serf protocol minimum version this node can understand or speak|integer|
|**ProtocolMax**  <br>*optional*||integer|
|**ProtocolCur**  <br>*optional*|Serf protocol current version this node can understand or speak|integer|
|**DelegateMin**  <br>*optional*|Serf delegate protocol minimum version this node can understand or speak|integer|
|**DelegateMax**  <br>*optional*|Serf delegate protocol minimum version this node can understand or speak|integer|
|**DelegateCur**  <br>*optional*|Serf delegate protocol minimum version this node can understand or speak|integer|


<a name="execution"></a>
### execution
An execution represents a timed job run.


|Name|Description|Schema|
|---|---|---|
|**job_name**  <br>*optional*|job name|string|
|**started_at**  <br>*optional*|start time of the execution|string(date-time)|
|**finished_at**  <br>*optional*|when the execution finished running|string(date-time)|
|**success**  <br>*optional*|the execution run successfuly|boolean|
|**output**  <br>*optional*|partial output of the command execution|string|
|**node_name**  <br>*optional*|name of the node that executed the command|string|





