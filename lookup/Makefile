build:
	env CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o bin/lookup cmd/lookup/main.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run