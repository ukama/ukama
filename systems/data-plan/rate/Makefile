# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../../../systems/config.mk

.PHONY: integration.test test build lint deps 

build: integration.build
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/data-rate/rate/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/rate cmd/server/main.go

test:
	go test -v ./...

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/...
			github.com/vektra/mockery/v2/.../

gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep -I ../base-rate/pb --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
 			  rate.proto
	mockery --all --recursive --dir ./internal
	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks
	
clean:
	rm -rf bin/

lint:
	golangci-lint run

integration.test:
	go test ./test/integration -tags integration  -v -count=1

integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration

docs:
	protoc -I ./pb -I ../../../systems/common/pb -I ../../../systems/common/pb/dep -I ../base-rate/pb --doc_out=. --doc_opt=./../../template.tmpl,README.md pb/rate.proto & sh ./../../generate-dir-tree.sh $(PWD)

server:
	go run cmd/server/main.go
