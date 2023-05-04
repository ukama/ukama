
gen:
	protoc --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src pb/ukama/*.proto --go-grpc_opt=require_unimplemented_servers=false
	protoc --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src pb/*.proto --go-grpc_opt=require_unimplemented_servers=false
	protoc -I ./pb/dep -I ./pb  --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
 			   pb/events/*.proto
	mockery --all --recursive --dir ./


test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

