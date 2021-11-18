
gen:
	protoc --go_out=. --go-grpc_out=. pb/*.proto --go-grpc_opt=require_unimplemented_servers=true

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

