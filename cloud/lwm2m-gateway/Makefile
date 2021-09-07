# A very simple makefile. Just for initial testing and build.
.PHONY: clean

PROJECT = "lwm2m-gateway"
# Proto's
# Common
COMMON_PROTO:= specs/common/proto/statuscode.proto
COMMON_SPEC:= specs/common/spec/*

#LWM2M GATEWAY Interface
LWM2MGATEWAY_PROTO:=specs/lwm2mIface/proto/lwm2m.proto
LWM2MGATEWAY_SPEC:=specs/lwm2mIface/spec*

#LwM2M Gateway to Server interface
LWM2MSERVER_PROTO:=specs/lwm2mServer/proto/lwm2mserver.proto
LWM2MSERVER_SPEC:=specs/lwm2mServer/spec*

# Build all
all: speclIface specServer gateway testService

# Common Protobuf Spec 
commonproto: 
	protoc --go_out=../  $(COMMON_PROTO)

# LWM2M GATEWAY Protobuf Spec 
speclIface: commonproto
	protoc --go_out=../  $(LWM2MGATEWAY_PROTO)

# LwM2M Interface Protobuf Spec 
specServer: commonproto
	protoc --go_out=../  $(LWM2MSERVER_PROTO)

# Gateway
gateway: speclIface
	env CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o bin/lwm2mGateway cmd/gateway/main.go	

#test service
testService:
	env CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o bin/testService test/service/main.go
	
# Remove the build binaries 
clean:
	rm -rf bin/*
	rm -rf $(LWM2MGATEWAY_SPEC) $(LWM2MSERVER_SPEC) $(COMMON_SPEC)


# CI Commands

build: gateway

test: 
	go test ./...

lint:
	golangci-lint run