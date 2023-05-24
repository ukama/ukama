include ../../../systems/config.mk

.PHONY: integration.test test build lint deps 

build: integration.build
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/data-rate/rate/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/rate cmd/server/main.go

test:
	go test -v ./...

gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep  -I ../base-rate/pb --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
					rate.proto

	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

clean:
	rm -rf bin/
	rm -rf mocks
	rm -rf pb/gen

# Go lint
lint:
	golangci-lint run


# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration

docs:
	protoc -I ./pb -I ../../../systems/common/pb -I ../../../systems/common/pb/dep -I ../base-rate/pb --doc_out=. --doc_opt=./../../template.tmpl,README.md pb/rate.proto & sh ./../../generate-dir-tree.sh $(PWD)

server:
	go run cmd/server/main.go
