# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../../../systems/config.mk

.PHONY: integration.test test build lint deps 

build: integration.build 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/services/msgClient/cmd/version.Version=$(BIN_VER) -extldflags=-static' -gcflags=all="-N -l" -o bin/msgClient cmd/msgClient/main.go

test: integration.test 
	go test -v ./...

gen:
	mockery --all --recursive --dir ./internal
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

clean:
	rm -rf bin/

lint:
	golangci-lint run

# Generate Docs

docs:
	protoc --doc_out=. --doc_opt=./template.tmpl,README.md pb/*.proto & sh ./../../generate-dir-tree.sh $(PWD)


# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration

server:
	go run cmd/msgClient/main.go
