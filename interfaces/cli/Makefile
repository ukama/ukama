# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

.PHONY: build test package

build:
	go build cmd/main.go

test:
	go test ./...

package:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/interfaces/cli/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/ukama_linux cmd/main.go
	env CGO_ENABLED=0 GOOS=windows  GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/interfaces/cli/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/ukama_windows.exe cmd/main.go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/interfaces/cli/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/ukama_darwin cmd/main.go


# Go lint
lint:
	staticcheck ./...

gen:
	mockery --all --recursive --dir ./pkg


deps:
	go install
