include ../../config.mk

.PHONY: integration.test test build lint deps 

build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/data-plan/base-rate/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/base-rate cmd/server/main.go	

test:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

lint:
	golangci-lint run

gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
 			rate.proto
	mockery --all --recursive --dir ./internal
	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks
	
clean:
	rm pb/gen/*.go

server:
	go run cmd/server/main.go