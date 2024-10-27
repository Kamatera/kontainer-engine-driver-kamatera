#!/usr/bin/env bash

docker run --network=host fullstorydev/grpcurl -plaintext -d "$(cat ${1:-tests/grpc_postcheck.json})" 127.0.0.1:8888 types.Driver.PostCheck
