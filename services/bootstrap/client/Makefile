include ../../config.mk

.PHONY: test lint build deps gen clean

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

gen:
	mockery --all --recursive --dir ./	
