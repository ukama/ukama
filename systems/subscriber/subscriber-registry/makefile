include ../../config.mk

.PHONY: integration.test test build lint deps 

build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/systems/subscriber/subscriber/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/package cmd/server/main.go
lint:
	golangci-lint run
test:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
	# integration tests
integration.test:
	go test ./test/integration -tags integration  -v -count=1
integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration
gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
 			  subscriber.proto
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
 			events.proto
	mockery --all --recursive --dir ./internal
	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks
clean:
	rm pb/gen/*.go
server:
	go run cmd/server/main.go