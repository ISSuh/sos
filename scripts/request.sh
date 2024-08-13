#!/bin/bash

file=$1

# curl -X POST -F "upload=@./sample/sample.mp3" 127.0.0.1:33221/v1/1/2/test

go run cmd/tools/uploader/main.go -s localhost:33221 -f $1 -g 1 -p 2 -o 3
