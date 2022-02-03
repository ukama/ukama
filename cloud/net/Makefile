include ../../config.mk

.PHONY: integration.test test build lint deps 

build: integration.build
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukamaX/cloud/net/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/net cmd/server/main.go
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukamaX/cloud/net/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/listener cmd/listener/main.go

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
	mkdir ./pb/gen -p

	protoc -I ./pb -I ../../common/pb/dep --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
         			  net.proto
	# generate http
	protoc -I ./pb -I ../../common/pb/dep --grpc-gateway_out ./pb/gen \
                     --grpc-gateway_opt logtostderr=true \
                     --grpc-gateway_opt paths=source_relative \
                     net.proto

#	# generate Open API
#	protoc I ./pb -I ../../common/pb/dep --openapiv2_out ./pb/gen \
#		--openapiv2_opt logtostderr=true \
#		net.proto

	mockery --all --recursive --dir ./pkg
	mockery -dir ./pb  -all -recursive -output ./pb/gen/mocks

clean:
	rm pb/gen/*.go


# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration