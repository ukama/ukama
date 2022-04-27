include ../../config.mk

.PHONY: integration.test test build lint deps 

build: 
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/services/bootstrap/lookup/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/lookup cmd/lookup/main.go

test:
	go test -v ./...

gen:
	mockery --all --recursive --dir ./internal

# Go lint
lint:
	golangci-lint run


# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration




