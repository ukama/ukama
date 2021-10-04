include ../../config.mk

build: integration.build
	env CGO_ENABLED=0 go build -ldflags='-X ukamaX/bootstrap/lookup/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/lookup cmd/lookup/main.go

test:
	go test -v ./...

gen:
	mockery --all --recursive --dir ./internal ./pkg

# Go lint
lint:
	golangci-lint run


# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration




