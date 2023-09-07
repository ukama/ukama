include ../../config.mk
  
.PHONY: build lint deps

# Proto's
# Common
NETWORK_PROTO:= network.proto
NETWORK_SPEC:= specs/controller/spec

build: mock
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/testing/services/network/cmd/network/version.Version=$(BIN_VER) -extldflags=-static' -o bin/network cmd/network/main.go

# Event  Protobuf Spec 
networkspec: 
	protoc -I ./specs/controller/proto/ --go_out=. --go-grpc_out=. --govalidators_out=.  $(NETWORK_PROTO)

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
mock: 
	env GO111MODULE=on CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o bin/mock test/mock/mock.go

clean:
	rm -rf ./bin
	rm -rf $(NETWORK_SPEC)/*

