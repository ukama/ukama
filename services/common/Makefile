
gen:
	protoc --go_out=. --go-grpc_out=. pb/ukama/*.proto --go-grpc_opt=require_unimplemented_servers=false

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

