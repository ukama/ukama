build:
	env CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o bin/api-gateway cmd/main.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run


deps:
	go get github.com/vektra/mockery/v2/.../


gen:
	mockery --all --recursive

clean:
	rm pb/generated/*.go
