#!/bin/bash

date=$(gdate -d "+2 minutes" --rfc-3339=seconds | sed 's/ /T/')
for i in `seq 1 10`; do
  time curl localhost:8080/v1/jobs -d "{
    \"name\": \"job_name_$i\",
    \"schedule\": \"@every `jot -r 1 3 20`s\",
    \"tags\": {
        \"dkron_server\": \"true\"
    },
    \"executor\": \"shell\",
    \"executor_config\": {
      \"command\": \"echo 'some words'\"
    }
  }"
done

for i in `seq 10 20`; do
  time curl localhost:8080/v1/jobs -d "{
    \"name\": \"job_name_$i\",
    \"schedule\": \"@every `jot -r 1 3 20`s\",
    \"tags\": {
        \"dkron_server\": \"true\"
    },
    \"command\": \"false\"
  }"
done
curl localhost:8080/v1/jobs -d '{
    "name": "job_name_1",
    "schedule": "@every 5s",
    "tags": {
        "dkron_server": "true"
    },
    "concurrency": "forbid",
    "executor": "shell",
    "executor_config": {
      "command": "echo",
      "shell": "true"
    },
    "processors": {
      "elasticsearch": {
        "index": "foobar"
      }
    }
  }'

Parent
{
"name": "parent1",
"schedule": "@every 3s",
"command": "echo \"papa\""
}

Child job config:
{
"name": "chain1",
"parent_job": "parent1",
"command": "echo \"hijo\""
}
