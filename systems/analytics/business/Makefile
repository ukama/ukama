# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

include ../../config.mk

.PHONY: integration.test test build lint deps gen clean server

build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/analytics/business/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/business cmd/server/main.go

lint:
	golangci-lint run

test:
	go test -v ./...

integration.test:
	go test ./test/integration -tags integration -v -count=1

integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration

deps:
	go get github.com/vektra/mockery/v2/.../

gen:
	@echo "proto generation placeholder for business"

clean:
	rm -rf bin/

server:
	go run cmd/server/main.go
