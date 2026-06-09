include ../../config.mk

.PHONY: build deps lint gen mocks clean integration.test test server

build: integration.build
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-X github.com/ukama/ukama/systems/operation/api-gateway/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/api-gateway cmd/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

deps:
	go install github.com/vektra/mockery/v2@latest

mocks:
	mockery --all --recursive

clean:
	rm -rf bin/

integration.test:
	go test ./test/integration -tags integration -v -count=1

integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration

server:
	go run cmd/main.go
