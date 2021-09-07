module ukamaX/cloud/api-gateway

go 1.16

replace github.com/ukama/ukamaX/common => ../../common

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/go-resty/resty/v2 v2.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukamaX/cloud/registry v0.0.0-20210818102821-83b756ff1d75 // indirect
	github.com/ukama/ukamaX/common v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
)
