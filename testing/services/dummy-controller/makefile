include ../../../systems/config.mk

.PHONY: integration.test test build lint deps 

build: 
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o bin/dummy-controller


test:
	go test -v ./...

gen:
	protoc -I ./pb -I ../../../systems/common/pb -I ../../../systems/common/pb/dep --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
					  eventGenerator.proto
					  
	mockery --all --recursive --dir ./internal
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks

clean:
	rm -rf bin/

lint:
	golangci-lint run



integration.test:
	go test ./test/integration -tags integration  -v -count=1


integration.build:
	env CGO_ENABLED=0 go test ./test/integration -tags integration -v -c -o bin/integration



