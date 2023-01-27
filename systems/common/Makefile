
gen:
	protoc --go_out=. --go-grpc_out=. pb/ukama/*.proto --go-grpc_opt=require_unimplemented_servers=false
	protoc --go_out=. --go-grpc_out=. pb/*.proto --go-grpc_opt=require_unimplemented_servers=false
	protoc --go_out=. -I ./pb/dep -I ./pb  --go-grpc_out=. pb/events/*.proto --go-grpc_opt=require_unimplemented_servers=false --govalidators_out=.
	
	mockery --all --recursive --dir ./


test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

