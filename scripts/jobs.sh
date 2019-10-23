#!/bin/bash

for i in {1..25}
do
   curl localhost:8080/v1/jobs -d "{
    \"name\": \"test_job_$i\",
    \"schedule\": \"@every $(($RANDOM%60+1))m\",
	\"executor\": \"shell\",
  	\"executor_config\": {
    	\"command\": \"echo 'run job $i'\"
  	}
}"
done
