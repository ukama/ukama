module github.com/ukama/ukamaX/cloud/api-gateway

go 1.16

replace github.com/ukama/ukamaX/common => ../../common

replace github.com/ukama/ukamaX/cloud/api-gateway => ./

replace github.com/ukama/ukamaX/cloud/registry => ../registry

require (
	github.com/gin-contrib/cors v1.3.1 // indirect
	github.com/gin-gonic/gin v1.7.2
	github.com/go-resty/resty/v2 v2.6.0
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/jarcoal/httpmock v1.0.8 // indirect
	github.com/ory/kratos-client-go v0.7.3-alpha.8
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukamaX/cloud/registry v0.0.0-20210818102821-83b756ff1d75
	github.com/ukama/ukamaX/common v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
)
