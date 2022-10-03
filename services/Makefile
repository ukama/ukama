SUBDIRS := $(wildcard cloud/*/.)  $(wildcard bootstrap/*/.)  $(wildcard common/.) 

build: $(SUBDIRS)	
$(SUBDIRS):
	$(MAKE) -C $@

test: $(SUBDIRS)	
$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: build test $(SUBDIRS)


go-install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@v0.3.2
	go install github.com/vektra/mockery/v2@v2.12.2	
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.46.2
