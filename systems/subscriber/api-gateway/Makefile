include ../../config.mk

.PHONY: integration.test test build lint deps 
build: gen
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/systems/subscriber/api-gateway/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/package cmd/server/main.go

lint:
	golangci-lint run

test:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

gen:
	mockery --all --recursive

clean:
	rm pb/gen/*.go

server:
	go run cmd/server/main.go