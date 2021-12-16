module github.com/ukama/ukamaX/cloud/registry

go 1.16

require (
	github.com/go-resty/resty/v2 v2.6.0
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukamaX/common v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	gorm.io/gorm v1.21.15
)

replace github.com/ukama/ukamaX/common => ../../common

replace github.com/ukama/ukamaX/cloud/registry => ./

replace github.com/ukama/ukamaX/cloud/registry/mocks => ./mocks
