# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../../config.mk

.PHONY: integration.test test build lint deps 
build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/subscriber/api-gateway/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/api-gateway cmd/server/main.go

lint:
	golangci-lint run

test:
	go test -v ./...

deps:
	go get github.com/vektra/mockery/v2/.../

gen:
	mockery --all --recursive

clean:
	rm -rf bin/

server:
	go run cmd/server/main.go
