include ../../config.mk

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/analytics/analysis/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/analysis cmd/server/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

deps:
	go get github.com/vektra/mockery/v2/.../

gen:
	mockery --all --recursive --dir ./pkg

clean:
	rm -rf bin/

.PHONY: build test lint deps gen clean
