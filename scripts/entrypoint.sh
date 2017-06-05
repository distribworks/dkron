#!/usr/bin/env bash

GOBIN=`pwd` go clean -i ./builtin/...
GOBIN=`pwd` go clean
GOBIN=`pwd` go install ./builtin/...
go build -o main
exec ./main $@
