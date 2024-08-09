#!/bin/bash

proto_dir=internal/infrastructure/transport/rpc
proto_files=${proto_dir}/*.proto

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ${proto_files}
