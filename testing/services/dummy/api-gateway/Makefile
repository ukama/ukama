# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../config.mk

build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -ldflags '-X github.com/ukama/ukama/testing/services/dummy/api-gateway/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/api-gateway cmd/main.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run


deps:
	go get github.com/vektra/mockery/v2/.../


gen:
	mockery --all --recursive

clean:
	rm -rf bin/

server:
	go run cmd/main.go