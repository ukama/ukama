.PHONY: all build docker-build docker-run clean deps

# Default target
all: build

# Ensure bin directory exists
bin:
	mkdir -p bin

# Build the binary
build: bin deps
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/testing/services/dummy/controller/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/controller cmd/main.go


# Install dependencies
deps:
	go mod download
	go mod tidy


gen:
	protoc -I ./pb -I ../../../../systems/common/pb -I ../../../../systems/common/pb/dep --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
 			  controller.proto
	mockery --all --recursive --dir ./internal
	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks
# Clean up
clean:
	rm -rf bin/

