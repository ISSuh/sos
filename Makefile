# Makefile for sos

message_proto_dir := internal/domain/model/message
message_proto_files := $(wildcard $(message_proto_dir)/*.proto)

gprc_proto_dir := internal/infrastructure/transport/rpc/message
gprc_service_dir := internal/infrastructure/transport/rpc
gprc_proto_files := $(wildcard $(gprc_proto_dir)/*.proto)

RELEASE ?= 0
ifeq ($(RELEASE), 1)
	build_options := -ldflags="-s -w" -trimpath
else
	build_options :=
endif

generate: 
	protoc -I=$(gprc_proto_dir) \
		--go_out=${gprc_proto_dir} \
		--go-grpc_out=${gprc_service_dir} \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(gprc_proto_files)

	protoc -I=$(message_proto_dir) \
		--go_out=${message_proto_dir} \
		--go_opt=paths=source_relative \
		$(message_proto_files)

	mv internal/infrastructure/transport/rpc/message/metadata_registry.pb.go internal/infrastructure/transport/rpc

.PHONY: vendor
vendor:
	go mod tidy && go mod vendor

.PHONY: all
all: vendor generate build_api

build_api:
	go build -o bin/api ${build_options} cmd/api/main.go 

.PHONY: test
test:
	
clean:
	rm -rf vendor
	rm -rf bin
	rm ${message_proto_dir}/*.pb.go
	rm ${gprc_proto_dir}/*.pb.go
	rm ${gprc_service_dir}/*.pb.go
	
