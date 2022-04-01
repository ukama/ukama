module github.com/ukama/ukamaX/cloud/device-feeder

go 1.16

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/rabbitmq/amqp091-go v1.2.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukamaX/cloud/net v0.0.0-20220128150430-55e44457630a
	github.com/ukama/ukamaX/cloud/registry v0.0.0-20211209085225-8d6ae39819e5
	github.com/ukama/ukamaX/common v0.0.0-20220322143821-0d6da632684f
	github.com/wagslane/go-rabbitmq v0.7.0
	google.golang.org/grpc v1.45.0
)

replace github.com/ukama/ukamaX/common => ../../common

replace github.com/ukama/ukamaX/cloud/device-feeder => ./

replace github.com/ukama/ukamaX/cloud/device-feeder/mocks => ./mocks

replace github.com/ukama/ukamaX/cloud/registry => ../registry
