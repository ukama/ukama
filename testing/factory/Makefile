include ../config.mk
  
.PHONY: build lint deps

# Proto's
# Common
FACTORY_PROTO:= specs/factory/proto/factory.proto
FACTORY_SPEC:= specs/factory/spec/*

build: factoryspec
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/testing/factory/cmd/factory/version.Version=$(BIN_VER) -extldflags=-static' -o bin/factory cmd/factory/main.go

# Event  Protobuf Spec 
factoryspec: 
	protoc --go_out=./ $(FACTORY_PROTO)

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

