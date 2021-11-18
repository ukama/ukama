module github.com/ukama/ukamaX/cloud/api-gateway

go 1.16

replace github.com/ukama/ukamaX/common => ../../common

replace github.com/ukama/ukamaX/cloud/api-gateway => ./

replace github.com/ukama/ukamaX/cloud/registry => ../registry

replace github.com/ukama/ukamaX/cloud/hss => ../hss

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.2
	github.com/go-resty/resty/v2 v2.6.0
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/ory/kratos-client-go v0.8.0-alpha.2
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukamaX/cloud/hss v0.0.0-00010101000000-000000000000
	github.com/ukama/ukamaX/cloud/registry v0.0.0-20210818102821-83b756ff1d75
	github.com/ukama/ukamaX/common v0.0.0-20211015093708-cd6e230254b5
	golang.org/x/net v0.0.0-20211104161157-7594919abe59 // indirect
	golang.org/x/oauth2 v0.0.0-20211028175245-ba495a64dcb5 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
)
