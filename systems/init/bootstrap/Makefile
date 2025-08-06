include ../../config.mk

.PHONY: integration.test test build lint deps 

build: integration.build 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/init/bootstrap/cmd/version.Version=$(BIN_VER) -extldflags=-static' -gcflags=all="-N -l" -o bin/bootstrap cmd/server/main.go

lint:
	golangci-lint run

test:
	go test -v ./...

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/...
			github.com/vektra/mockery/v2/...

gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
		bootstrap.proto
	mockery --all --recursive
	mockery --all --recursive --dir ./internal
	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

clean:
	rm -rf bin/

docs:
	protoc --doc_out=. --doc_opt=./../../template.tmpl,README.md pb/*.proto & sh ./../../generate-dir-tree.sh $(PWD)


integration.test:
	go test ./test/integration -tags integration  -v -count=1

integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration

server:
	go run cmd/server/main.go
