# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../../config.mk

.PHONY: integration.test test build lint deps 

build: integration.build
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/systems/hub/distributor/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/distributor cmd/distributor/main.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

deps:
	go install

gen:
	mockery --all --recursive --dir ./pkg

clean:
	rm -rf bin/

integration.test:
	go test ./test/integration -tags integration  -v -count=1

integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration
