include ../../config.mk

.PHONY: build lint deps 

build:
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
