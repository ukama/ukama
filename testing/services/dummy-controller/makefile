.PHONY: all build docker-build docker-run clean deps

all: build

bin:
	mkdir -p bin

build: bin deps
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/dummy-controller main.go


clean:
	rm -rf bin/

# Start everything
start: build docker-build docker-run


