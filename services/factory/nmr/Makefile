include ../../config.mk

.PHONY: integration.test test build lint deps 

build: integration.build
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/services/factory/nmr/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/nmr cmd/nmr/main.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

deps:
	go install


gen:
	mockery --all --recursive --dir ./internal/db

clean:
	rm -rf ./bin

# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1

integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration
