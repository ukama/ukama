include ../../config.mk

.PHONY: integration.test test

build: integration.build
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukamaX/cloud/registry/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/registry cmd/server/main.go
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukamaX/cloud/registry/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/listener cmd/listener/main.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/...
			github.com/vektra/mockery/v2/.../


gen:
	protoc -I ./pb -I ../../common/pb --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
 			  registry.proto
	protoc --go_out=. --go-grpc_out=./pb/external pb/external/*.proto --go-grpc_opt=require_unimplemented_servers=false --go_opt=Mpb/external/link.proto=github.com/ukama/ukamaos/mesh
	mockery --all --recursive --dir ./internal
	mockery --all --recursive --dir ./pkg
	mockery -dir ./pb  -all -recursive -output ./pb/gen/mocks

clean:
	rm pb/generated/*.go


# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration




