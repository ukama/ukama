include ../../config.mk

.PHONY: integration.test test build lint deps 

build: integration.build
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukamaX/service-boilerplate/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/hss cmd/server/main.go	

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/vektra/mockery/v2/.../


gen:
	protoc --go_out=. --go-grpc_out=. pb/*.proto --go-grpc_opt=require_unimplemented_servers=true		
	mockery --all --recursive --dir ./pkg
	mockery -dir ./pb  -all -recursive -output ./pb/gen/mocks

clean:
	rm pb/gen/*.go


# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration