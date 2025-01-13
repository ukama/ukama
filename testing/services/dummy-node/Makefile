.PHONY: test build deps 

build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o bin/dummy-node

lint:
	golangci-lint run

test:
	echo not tests

