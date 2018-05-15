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
      \"command\": \"true\"
    }
  }"
done

for i in `seq 1 10`; do
  time curl localhost:8080/v1/jobs -d "{
    \"name\": \"job_name_$i\",
    \"schedule\": \"@every `jot -r 1 3 20`s\",
    \"tags\": {
        \"dkron_server\": \"true\"
    },
    \"command\": \"false\"
  }"
done
