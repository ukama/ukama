# Rest Service Boilerplate

Includes:
- Basic configuration (pkg/config.go)[config.go]
- Rest service based on Gin, Tonik and Fizz for OpenApi spec generation
- Ping request handler (in `common` package)
- Swagger UI on `/swagger`
- serving metrics for Rest request on ":10250/metrics". 
- custom metrics declared in (pkg/metrics/metrics.go)[metrics.go]
- Integration tests in (/test/integration)(test/integration) and infrastructure to run then

## How to use 

- Copy whole directory 
- Replace `rest-service` with your service name everywhere in the project. Including Dockerfiles,  Makefile and go.mod
- Adjust paths in go.mod according to the location of your project
- Adjust import paths if needed
- run `go mod tidy`
- run `make build`
- run `make test`

