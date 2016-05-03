#!/usr/bin/env bash

docker-compose up -d etcd
export COMPOSE_ETCD_PORT=`docker port dkron_etcd_1 4001/tcp | cut -d":" -f 2`
export DKRON_BACKEND_MACHINE=`docker-machine ip default`:$COMPOSE_ETCD_PORT
go test -v ./dkron $1
