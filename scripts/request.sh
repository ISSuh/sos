#!/bin/bash

file=$1

curl -X POST -F "upload=@./sample/sample.mp3" 127.0.0.1:33221/v1/1/2/test
