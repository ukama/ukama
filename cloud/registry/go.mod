module github.com/ukama/ukamaX/cloud/registry

go 1.16

require (
	github.com/go-resty/resty/v2 v2.6.0
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/golang-jwt/jwt/v4 v4.1.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.5.1
	github.com/ukama/ukamaX/common v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20210716203947-853a461950ff // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210719143636-1d5a45f8e492 // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
	gorm.io/gorm v1.21.15
)

replace github.com/ukama/ukamaX/common => ../../common

replace github.com/ukama/ukamaX/cloud/registry => ./

replace github.com/ukama/ukamaX/cloud/registry/mocks => ./mocks
