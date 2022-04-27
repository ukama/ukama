include ../../config.mk

.PHONY: build deps lint gen clean integration.test test

build: integration.build
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build  -ldflags '-X github.com/ukama/ukamaX/cloud/device-gateway/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/device-gateway cmd/main.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run


deps:
	go get github.com/vektra/mockery/v2/.../


gen:
	mkdir ./pb/gen -p
	# generate gateway code
	protoc -I ../hss/pb/google/api -I ../hss/pb/ -I ../../common/pb/ --grpc-gateway_out ./pb/gen \
                --grpc-gateway_opt logtostderr=true \
                --grpc-gateway_opt paths=source_relative \
                --grpc-gateway_opt standalone=true \
                hss.proto
	# we need to replace importing package as we are using grpc client from other module
	sed -i -e 's/extGen "pb\/gen"/extGen "github.com\/ukama\/ukamaX\/cloud\/hss\/pb\/gen"/' ./pb/gen/hss.pb.gw.go
	# generate open api documentation
	protoc -I ../hss/pb/  -I ../../common/pb/ --openapiv2_out ./swagger-ui \
        --openapiv2_opt logtostderr=true \
        hss.proto

	mockery --all --recursive

clean:
	rm pb/generated/*.go

install-tools:
	sudo go install \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
		google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc

# integration tests

integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration


