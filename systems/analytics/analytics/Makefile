# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

include ../../config.mk

TARGET_EXEC := analytics
BIN_DIR := bin
CMD_DIR := cmd/server
VERSION_PKG := github.com/ukama/ukama/systems/analytics/analytics/cmd/version

GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0

.PHONY: all deps build test tidy clean

all: build

deps:
	go mod download

tidy:
	go mod tidy

build: deps
	@mkdir -p $(BIN_DIR)
	@echo Building version: \"$(VERSION)\"
	env CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build \
		-ldflags='-X $(VERSION_PKG).Version=$(VERSION) -extldflags=-static' \
		-o $(BIN_DIR)/$(TARGET_EXEC) \
		$(CMD_DIR)/main.go

test: deps
	env CGO_ENABLED=$(CGO_ENABLED) go test ./... -v

clean:
	rm -rf $(BIN_DIR)
