include ../../config.mk

.PHONY: integration.test test build lint deps 

build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/systems/data-plan/base-rate/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/rate cmd/server/main.go
lint:
	golangci-lint run
test:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
gen:
	protoc --go_out=. --go_opt=paths=source_relative \
   --go-grpc_out=. --go-grpc_opt=paths=source_relative \
   pb/package.proto

	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks
clean:
	rm pb/gen/*.go
server:
	go run cmd/server/main.go