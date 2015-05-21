#scheduler2: bin/serf agent -node=bar -bind=127.0.0.1:5001 -rpc-addr=127.0.0.1:7374 -discover dcron
dcron2: go run *.go server -node=dcron2 -bind=127.0.0.1:5001 -rpc-addr=127.0.0.1:7374
dcron3: go run *.go server -node=dcron3 -bind=127.0.0.1:5002 -rpc-addr=127.0.0.1:7375
etc: bin/etcd -name dcron1
