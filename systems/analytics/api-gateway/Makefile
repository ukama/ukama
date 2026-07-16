include ../../config.mk

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/analytics/api-gateway/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/api-gateway cmd/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

gen:
	mockery --all --recursive --dir ./pkg

clean:
	rm -rf bin/

.PHONY: build test lint gen clean
