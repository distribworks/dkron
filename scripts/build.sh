#!/usr/bin/env bash

set -e

GOARCH=amd64 GOOS=linux godep go build -o bin/dcron
