include ../../config.mk

.PHONY: integration.test test build lint deps 
build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-X github.com/ukama/ukama/systems/subscriber/api-gateway/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/package cmd/server/main.go

lint:
	golangci-lint run

test:
	go test -v ./...

integration.test:
	go test ./test/integration -tags integration  -v -count=1

integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration
	
gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
 			   package.proto
	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

clean:
	rm pb/gen/*.go

docs:
	protoc --doc_out=. --doc_opt=./../../template.tmpl,README.md pb/*.proto & sh ./../../generate-dir-tree.sh $(PWD)

server:
	go run cmd/server/main.go