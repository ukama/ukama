module ukamaX/api-gateway

go 1.16

replace github.com/ukama/ukamaX/common => ../common

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/go-resty/resty/v2 v2.6.0
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukamaX/common v0.0.0-00010101000000-000000000000
	github.com/ukama/ukamaX/registry v0.0.0-20210806133156-c318edfb57c3 // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1 // indirect
)
