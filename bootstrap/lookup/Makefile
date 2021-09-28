include ../../config.mk

build:
	env CGO_ENABLED=0 go build -ldflags='-X ukamaX/bootstrap/lookup/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/lookup cmd/lookup/main.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run