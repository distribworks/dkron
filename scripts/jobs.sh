#!/bin/bash

for i in {1..2000}
do
   curl test.dkron.io:8080/v1/jobs -d "{
    \"name\": \"test_job_$i\",
    \"schedule\": \"@every $(($RANDOM%10+1))s\",
	\"concurrency\": \"forbid\",
	\"tags\": {
            \"role\": \"dkron:1\"
          },
	\"executor\": \"http\",
  	\"executor_config\": {
            \"method\": \"GET\",
            \"url\": \"https://httpbin.org/get\"
          }
  	},
	\"ephemeral\": false,
}"
done
