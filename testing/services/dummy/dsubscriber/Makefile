# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../config.mk

build:
	@echo Building version: \"$(BIN_VER)\"
<<<<<<<< HEAD:testing/services/dummy/dsubscriber/Makefile
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/testing/services/dummy/dsubscriber/version.Version=$(BIN_VER) -extldflags=-static' -o bin/dsubscriber cmd/server/main.go
========
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/testing/services/dummy/dnode/version.Version=$(BIN_VER) -extldflags=-static' -o bin/dnode server/main.go
>>>>>>>> main:testing/services/dummy/dnode/Makefile

test:
	go test -v ./...

lint:
	golangci-lint run

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/...
			github.com/vektra/mockery/v2/.../
<<<<<<<< HEAD:testing/services/dummy/dsubscriber/Makefile


gen:
	protoc -I ./pb -I ../../../../systems/common/pb -I ../../../../systems/common/pb/dep --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
 			  dsubscriber.proto
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

clean:
	rm -rf bin
========
clean:
	rm -rf bin/

server:
	go run server/main.go
>>>>>>>> main:testing/services/dummy/dnode/Makefile
