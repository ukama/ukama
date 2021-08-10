build:
	env CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o bin/registry cmd/server/main.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run



postgres:
	mkdir ${HOME}/postgres-data/
	# user - postgres
	docker run -d --name dev-postgres -e POSTGRES_PASSWORD=Pass2020! \
            -v ${HOME}/postgres-data/:/var/lib/postgresql/data \
            -p 5432:5432 \
            postgres

deps:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc \
			github.com/ysugimoto/grpc-graphql-gateway/protoc-gen-graphql/...
			github.com/vektra/mockery/v2/.../


gen:
	protoc --go_out=. --go-grpc_out=. pb/*.proto --go-grpc_opt=require_unimplemented_servers=false
	mockery --all --recursive

clean:
	rm pb/generated/*.go

server: gen
	go run cmd/server/main.go -port 55555

client:
	go run cmd/client/main.go -address 0.0.0.0:55555
