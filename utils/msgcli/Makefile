# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ./config.mk
  
.PHONY: build msgcli tool

build: 
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/utils/msgcli/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/msgcli main.go

# Go test
test:
	go test -v ./...


# Go lint
lint:
	golangci-lint run


clean:
	rm -rf bin/
