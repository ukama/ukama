# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../../config.mk

.PHONY: integration.test test build lint deps 

build: integration.build
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/subscriber/sim-pool/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/sim-pool cmd/server/main.go

lint:
	golangci-lint run

integration.test:
	go test ./test/integration -tags integration  -v -count=1

integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration

test:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/...
			github.com/vektra/mockery/v2/.../

gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
		sim_pool.proto
	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

docs:
	protoc --doc_out=. --doc_opt=./../../template.tmpl,README.md pb/*.proto & sh ./../../generate-dir-tree.sh $(PWD)

clean:
	rm pb/gen/*.go
	rm mocks/*.go

server:
	go run cmd/server/main.go
