#
# Makefile for api-gateway
#

lint:
	golangci-lint run

gen:
	protoc -I ./pb -I ../../../services/common/pb/gen --go_out=. --go-grpc_out=. \
	--go-grpc_opt=require_unimplemented_servers=true pb/client/lookup.proto; \
	mv pb/client/*.go pb/gen/

build:
	go build cmd/main.go

clean:
	rm -rf main