#!/bin/bash

file=$1

curl -v -N -X PUT -H "Transfer-Encoding: chunked" -d ~/Downloads/online.pdf 127.0.0.1:33221/v1/1/2/test
