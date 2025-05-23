# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../../config.mk
  
.PHONY: integration.test test

build: integration.build
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/node/notify/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/notify cmd/server/main.go

# Go test
test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/vektra/mockery/v2/.../

gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
		notify.proto
	mockery --all --recursive --dir ./internal
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks


clean:
	rm -rf ./bin/

integration.test:
	go test ./test/integration -tags integration  -v -count=1

integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration

mock: 
	env GO111MODULE=on CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o bin/mock test/mock/listener.go

server:
	go run cmd/server/main.go
