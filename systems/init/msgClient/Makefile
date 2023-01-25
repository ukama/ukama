include ../../../systems/config.mk

.PHONY: integration.test test build lint deps 

build: integration.build 
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/systems/init/msgClient/cmd/version.Version=$(BIN_VER) -extldflags=-static' -gcflags=all="-N -l" -o bin/msgClient cmd/msgClient/main.go

test: integration.test 
	go test -v ./...

gen:

	protoc -I ./pb -I ../../../systems/common/pb -I ../../../systems/common/pb/dep --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
					  msgClient.proto
					  
	mockery --all --recursive --dir ./internal
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

clean:
	rm -rf bin/

# Go lint
lint:
	golangci-lint run

# Generate Docs

docs:
	protoc --doc_out=. --doc_opt=./template.tmpl,README.md pb/*.proto & sh ./../../generate-dir-tree.sh $(PWD)


# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration

server:
	go run cmd/msgClient/main.go


