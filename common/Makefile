
gen:
	protoc --go_out=. --go-grpc_out=. pb/*.proto --go-grpc_opt=require_unimplemented_servers=true
	protoc --go_out=pb/gen --go-grpc_out=pb/gen pb/ukama/*.proto --go-grpc_opt=require_unimplemented_servers=false

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

