include ../../config.mk

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/analytics/aggregator/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/aggregator cmd/server/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/vektra/mockery/v2/.../

gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
 			  aggregator.proto
	mockery --all --recursive --dir ./internal
	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks


clean:
	rm -rf bin/

.PHONY: build test lint deps gen clean
