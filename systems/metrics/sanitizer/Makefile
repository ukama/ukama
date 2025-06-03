# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

include ../../config.mk

.PHONY: integration.test test build lint deps

build: integration.build
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/metrics/sanitizer/cmd/version.Version=$(BIN_VER) -extldflags=-static' -gcflags=all="-N -l" -o bin/sanitizer cmd/server/main.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/vektra/mockery/v2/.../
	go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@latest


gen:
	protoc -I ./pb  -I ../../common/pb/dep  --go_out=. --go-grpc_out=. --govalidators_out=. \
                            --go-grpc_opt=require_unimplemented_servers=true \
                            sanitizer.proto

	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

clean:
	rm -rf bin/

integration.test:
	go test ./test/integration -tags integration  -v -count=1

integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration

docs:
	protoc -I ./pb  -I ../../common/pb/dep --doc_out=. --doc_opt=./template.tmpl,README.md pb/*.proto & sh ./generate-dir-tree.sh $(PWD)
