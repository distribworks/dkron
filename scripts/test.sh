#!/usr/bin/env bash

NOVENDOR=$(docker-compose run dkron glide novendor | tr -d '\r')

docker-compose run dkron scripts/validate-gofmt
docker-compose run dkron go vet $NOVENDOR
docker-compose run -e DKRON_BACKEND_MACHINE=etcd:2389 dkron go test -v $NOVENDOR $1
