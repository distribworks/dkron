#!/bin/bash

date=$(gdate -d "+2 minutes" --rfc-3339=seconds | sed 's/ /T/')
for i in `seq 1 5`; do
  time curl localhost:8080/v1/jobs -d "{
    \"name\": \"job_name_$i\",
    \"schedule\": \"@every 4s\",
    \"tags\": {
        \"dkron_server\": \"true\"
    },
    \"executor\": \"shell\",
    \"executor_config\": {
      \"command\": \"echo 'hallo'\"
    }
  }"
done
