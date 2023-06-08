include ../../../systems/config.mk

.PHONY: build deps lint gen clean integration.test test

build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build  -ldflags '-X github.com/ukama/ukama/systems/messaging/api-gateway/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/api-gateway cmd/main.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

deps:
	go get github.com/vektra/mockery/v2/.../

gen:
	mockery --all --recursive --dir ./pkg

clean:
	rm -rf ./bin/ 

# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration


