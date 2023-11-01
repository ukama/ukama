# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../../../systems/config.mk

.PHONY: integration.test test build lint deps 

build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/services/eventGenerator/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/eventGenerator cmd/eventGenerator/main.go

test:
	go test -v ./...

gen:
	protoc -I ./pb -I ../../../systems/common/pb -I ../../../systems/common/pb/dep --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
					  eventGenerator.proto
					  
	mockery --all --recursive --dir ./internal
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

clean:
	rm -rf bin/

# Go lint
lint:
	golangci-lint run


# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration




