etc: bin/etcd -name dcron1
dkron2: godep go run *.go agent -server -node=dkron2 -join=127.0.0.1:5002 -bind=127.0.0.1:5001 -http-addr=:8080
dkron3: godep go run *.go agent -server -node=dkron3 -join=127.0.0.1:5001 -bind=127.0.0.1:5002 -http-addr=:8081
