
gen:
	protoc --go_out=. --go-grpc_out=. pb/ukama/*.proto --go-grpc_opt=require_unimplemented_servers=false
	mockery --all --recursive --dir ./


test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

