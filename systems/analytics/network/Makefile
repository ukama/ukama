# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

include ../../config.mk

GOPATH ?= $(shell go env GOPATH)
PROTO_FILE := $(shell find ./pb -maxdepth 1 -name '*.proto' -printf '%f\n' | head -n 1)

.PHONY: test lint deps gen mocks clean info

test:
	go test -v ./pkg/...

lint:
	golangci-lint run ./pkg/...

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

gen:
	@test -n "$(PROTO_FILE)" || (echo "missing proto file under ./pb"; exit 1)
	@mkdir -p $(GOPATH)/src
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep \
		--go_out=$(GOPATH)/src \
		--go-grpc_out=$(GOPATH)/src \
		--go-grpc_opt=require_unimplemented_servers=true \
		--govalidators_out=$(GOPATH)/src \
		$(PROTO_FILE)

mocks:
	mockery --all --recursive --dir ./pkg

clean:
	rm -rf bin/

info:
	@echo "network is a domain module only."
	@echo "Runtime service is analytics."
