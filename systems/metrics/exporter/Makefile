include ../../config.mk

.PHONY: integration.test test build lint deps 

build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/systems/metrics/exporter/cmd/version.Version=$(BIN_VER) -extldflags=-static' -gcflags=all="-N -l" -o bin/exporter cmd/server/main.go	

stubs:
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/systems/metrics/exporter/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/stub test/stubs/stub.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/vektra/mockery/v2/.../
	go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@latest


gen:
	protoc -I ./pb  -I ../../common/pb/dep  --go_out=. --go-grpc_out=. --govalidators_out=. \
                            --go-grpc_opt=require_unimplemented_servers=true \
                            exporter.proto
	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

clean:
	rm -r pb/gen


# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration

docs:
	protoc -I ./pb  -I ../../common/pb/dep --doc_out=. --doc_opt=./template.tmpl,README.md pb/*.proto & sh ./generate-dir-tree.sh $(PWD)
