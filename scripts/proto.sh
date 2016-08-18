#!/usr/bin/env bash

ROOT_DIR=${PWD}
CWD=${ROOT_DIR}/dkronpb

export IMPORT_PATH=${CWD}:${ROOT_DIR}/vendor
export GENERATOR="gogofaster_out"
export OUTPUT_DIR=${CWD}
export PROTO_FILES="$CWD/*.proto"

protoc --proto_path=${IMPORT_PATH} \
       --${GENERATOR}=plugins=grpc:${OUTPUT_DIR} \
       ${PROTO_FILES}
