# Makefile for sos

proto_dir := internal/rpc
proto_files := $(wildcard $(proto_dir)/*.proto)

generate:
	protoc -I=$(proto_dir) --go_out=${proto_dir} --go-grpc_out=${proto_dir} --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(proto_files)

.PHONY: vendor
vendor:
	go mod tidy && go mod vendor

build: vendor generate

.PHONY: test
test:
	
clean:
	rm -rf vendor
	rm ${proto_dir}/*.pb.go
