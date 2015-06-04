etc: bin/etcd -name dcron1
dcron2: go run *.go agent -server -node=dcron2 -join=127.0.0.1:5002 -bind=127.0.0.1:5001 -http-addr=:8081
dcron3: go run *.go agent -server -node=dcron3 -join=127.0.0.1:5001 -bind=127.0.0.1:5002 -http-addr=:8082
