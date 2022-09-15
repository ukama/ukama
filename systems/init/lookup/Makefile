include ../../../services/config.mk

.PHONY: integration.test test build lint deps 

build: integration.build
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/systems/init/lookup/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/lookup cmd/lookup/main.go

test:
	go test -v ./...

gen:
	protoc -I ./pb -I ../../../services/common/pb -I ../../../services/common/pb/dep --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
					  lookup.proto

	protoc -I ./pb -I ../../../services/common/pb -I ../../../services/common/pb/dep --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
					  msgBus.proto
					  
	#mockery --all --recursive --dir ./pkg
	#mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

clean:
	rm pb/gen/*.go

# Go lint
lint:
	golangci-lint run


# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration




