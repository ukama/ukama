module github.com/ukama/ukama/services/cloud/registry

go 1.16

require (
	github.com/go-resty/resty/v2 v2.6.0
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/goombaio/namegenerator v0.0.0-20181006234301-989e774b106e
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/jackc/pgconn v1.8.1
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukama/services/common v0.0.0-20220322143821-0d6da632684f
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.27.1
	gorm.io/gorm v1.21.15
)

replace github.com/ukama/ukama/services/common => ../../common

replace github.com/ukama/ukama/services/cloud/registry => ./

replace github.com/ukama/ukama/services/cloud/registry/mocks => ./mocks
