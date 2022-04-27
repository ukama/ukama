# Bootstrap API Gateway 

API gateway for Bootstrap system that provides REST Api for Ukama-cloud instances.  
It's a [KrakenD](https://www.krakend.io/) gateway that is configure via  [kraken.json](bootstrap/bootstrap-api/config/karkend.json) file. 
It uses [flexible config feature](https://www.krakend.io/docs/configuration/flexible-config/) of krakend.

# Configuration
Default configuration of API gateway uses [JWT Validator](https://www.krakend.io/docs/authorization/jwt-validation/) for endpoints that requires authentication 
To configure default validator you will need to initialize below environment variables:
```
AUTH_JWK_URL=https://your_auth_server.com/.well-known/jwks.json
AUTH_AUDIENCE="http://api.example.com"
AUTH_ALG="RS256"
```
More info in [Krakend documentation](https://www.krakend.io/docs/authorization/jwt-validation/)

## Configuring Auth0 
Bootstap API can be integrated with Auth0 to use it as as authentication server

To configure Aut0 follow the [Krakend Auth0 integration documentation](https://www.krakend.io/docs/authorization/auth0/#the-auth0---krakend-integration)

# Using the API 

###  Health check verifies that API gateway is up and running
```
GET http://localhost:8080/__health 
{"status":"OK"} 
```

###  Ready endpoint verifies that downstream services are up and running
```
GET http://localhost:8080/ready 
```

### Add or update organization record 
```
POST http://localhost:8080/orgs/{name}
{
	"certificate":"Q2VydGlmaWNhdGUK",
	"ip": "192.124.23.1"
}
```
Response:
```
HTTP 200
{
  "status": "Organisation added or updated"
}
```



### Add or update device record
```
POST /orgs/{name}/devices/{NodeId}
```
Response:
```
HTTP 200
{
  "status": "Mapping added or updated"
}
```
