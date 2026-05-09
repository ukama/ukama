include ../../config.mk

.PHONY: build test lint clean server gen deps

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

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@latest

gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep \
		--go_out=./pb/gen --go_opt=paths=source_relative \
		--go-grpc_out=./pb/gen --go-grpc_opt=paths=source_relative,require_unimplemented_servers=true \
		--govalidators_out=./pb/gen --govalidators_opt=paths=source_relative \
		./pb/site_controller.proto
