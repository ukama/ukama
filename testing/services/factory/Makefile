include ../../config.mk
  
.PHONY: test build lint deps

# Proto's
# Common
FACTORY_PROTO:= factory.proto
FACTORY_SPEC:= ./specs/factory/spec

build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/testing/services/factory/cmd/factory/version.Version=$(BIN_VER) -extldflags=-static' -o bin/factory cmd/factory/main.go

# Event  Protobuf Spec 
factoryspec:
	protoc -I ./specs/factory/proto/ --go_out=. --go-grpc_out=. --govalidators_out=.  $(FACTORY_PROTO)

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc	
	go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@latest


gen:
	mockery --all --recursive --dir ./internal/db

# Build Mock
mock: factoryspec
	env GO111MODULE=on CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o bin/mock test/mock/mock.go

clean:
	rm -rf ./bin
	rm -rf $(FACTORY_SPEC)/*

