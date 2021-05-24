#!/bin/bash

for i in {1..10000}
do
   curl localhost:8080/v1/jobs -d "{
    \"name\": \"test_job_$i\",
    \"schedule\": \"@every $(($RANDOM%60+1))s\",
	\"executor\": \"shell\",
  	\"executor_config\": {
    	\"command\": \"echo 'run job $i'\"
  	},
	\"ephemeral\": true
}"
done
