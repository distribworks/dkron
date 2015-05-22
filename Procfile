#scheduler2: bin/serf agent -node=bar -bind=127.0.0.1:5001 -rpc-addr=127.0.0.1:7374 -discover dcron
etc: bin/etcd -name dcron1
dcron2: go run *.go server -node=dcron2 -bind=127.0.0.1:5001 -rpc-addr=127.0.0.1:7374 -http-addr=:8081
dcron3: go run *.go server -node=dcron3 -bind=127.0.0.1:5002 -rpc-addr=127.0.0.1:7375 -http-addr=:8082
