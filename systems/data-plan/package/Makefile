include ../../config.mk

.PHONY: integration.test test build lint deps 
build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/subscriber/api-gateway/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/package cmd/server/main.go

lint:
	golangci-lint run

test:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
		package.proto
	mockery --all --recursive --dir ./pkg

clean:
	rm pb/gen/*.go

server:
	go run cmd/server/main.go