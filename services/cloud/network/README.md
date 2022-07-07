# Node Network

Node Network service is responsible for storing user->org->nodes structure
It exposes gRPC [interface](/network/pb/network.proto#L5).

### Structure
- Every node belongs to an organization
- Every organization belongs to a user (owner)

## Run Locally

- Start postgres

`make postgres`
  
- Run network service

`go run .\cmd\server\main.go`

- use [Evans](https://github.com/ktr0731/evans) to connect to service 

`evans repl --proto .\pb\health.proto,.\pb\network.proto --port 9090`

- run set of requests 

` go run .\cmd\test\main.go`
