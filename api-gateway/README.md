# Api-gateway

Api-gateway is a Rest interface of Ukama cloud for user facing apps: frontend, mobile or API consumers


## Run Locally

- Start docker-compose with dependencies. Make sure that all service images are up-to-date 

`docker-compose -f docker-compose.deps.yaml up`

- enable debug mode and bypass authentication. You still have to provide empty token or authorization header   
```
export DEBUGMODE=true
export BYPASSAUTHMODE=true
```
- Run service

`go run ./cmd/main.go`

- Test request
`curl -s -X GET  -H "token: bearer token"  http://localhost:8080/orgs/org-1`

