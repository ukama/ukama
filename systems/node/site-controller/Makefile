include ../../config.mk

.PHONY: build test lint clean server

build:
	@echo Building version: "$(BIN_VER)"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-X github.com/ukama/ukama/systems/node/site-controller/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/site-controller cmd/server/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

server:
	go run cmd/server/main.go
