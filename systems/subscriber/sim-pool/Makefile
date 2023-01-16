include ../../config.mk

.PHONY: integration.test test build lint deps 
build:
	@echo Building version: \"$(BIN_VER)\"
	env CGO_ENABLED=0 go build -ldflags='-X github.com/ukama/ukama/systems/subscriber/sim-pool/cmd/version.Version=$(BIN_VER) -extldflags=-static' -o bin/sim-pool cmd/server/main.go
lint:
	golangci-lint run
test:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
gen:
	protoc -I ./pb -I ../../common/pb -I ../../common/pb/dep --go_out=. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=. \
 			sim.proto
	mockery --all --recursive --dir ./internal
	mockery --all --recursive --dir ./pkg
	mockery --dir ./pb  --all --recursive --output ./pb/gen/mocks
docs:
	protoc --doc_out=. --doc_opt=./../../template.tmpl,README.md pb/*.proto & sh ./../../generate-dir-tree.sh $(PWD)
clean:
	rm pb/gen/*.go
server:
	go run cmd/server/main.go