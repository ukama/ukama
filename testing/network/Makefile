include ../config.mk
  
.PHONY: build lint deps

# Proto's
# Common
NETWORK_PROTO:= specs/controller/proto/network.proto
NETWORK_SPEC:= specs/network/spec/*

build: mock networkspec
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/testing/network/cmd/network/version.Version=$(BIN_VER) -extldflags=-static' -o bin/network cmd/network/main.go

# Event  Protobuf Spec 
networkspec: 
	protoc --go_out=./ $(NETWORK_PROTO)

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

deps:
	go install


gen:
	mockery --all --recursive --dir ./internal/db

# Build Mock
mock: networkspec
	env GO111MODULE=on CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o bin/mock test/mock/mock.go

clean:
	rm -rf ./bin

