include ../../config.mk

.PHONY: integration.test test build lint deps 

build: integration.build
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/systems/data-plan/base-rate/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/rate cmd/server/main.go
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/systems/data-plan/base-rate/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/listener cmd/listener/main.go
lint:
	golangci-lint run
gen:
	protoc -I ./pb -I ../../common/pb/dep  --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true  \
					  rate.proto
	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks
clean:
	rm pb/gen/*.go
server:
	go run cmd/server/main.go