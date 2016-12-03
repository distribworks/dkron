#!/usr/bin/env bash

docker-compose run -e DKRON_BACKEND_MACHINE=etcd:4001 dkron go test -v $(glide novendor) $1
