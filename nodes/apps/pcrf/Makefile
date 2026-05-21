# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../../ukamaOS/config.mk

UKAMAOS_ROOT = ../../ukamaOS
VERSION_SCRIPT = $(NODES_DIR)/utils/scripts/generate_version.sh

BIN_VER ?= $(shell $(VERSION_SCRIPT) --print 2>/dev/null || echo v0.0.debug)

.PHONY: build deps lint gen clean integration.test test server

build:
	@echo Building version: \"$(BIN_VER)\"
	export CGO_LDFLAGS="-L$(UKAMAOS_ROOT)/distro/vendor/build -lsqlite3 -Wl,-rpath=./lib"
	export CGO_CFLAGS="-I$(UKAMAOS_ROOT)/distro/vendor/build/include"
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
		-ldflags ' -s -w -X github.com/ukama/ukama/nodes/apps/pcrf/cmd/version.Version=$(BIN_VER)' \
		-o bin/pcrf cmd/server/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

deps:
	go get github.com/vektra/mockery/v2/.../

gen:
	mockery --all --recursive --dir ./pkg

clean:
	rm -rf bin/

server:
	go run cmd/server/main.go
