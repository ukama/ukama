build:
	env CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o bin/registry cmd/registry/main.go

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


