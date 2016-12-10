#!/usr/bin/env bash

docker-compose run dkron scripts/validate-gofmt
docker-compose run dkron go vet ./...
docker-compose run -e DKRON_BACKEND_MACHINE=etcd:4001 dkron go test -v $(glide novendor) $1
