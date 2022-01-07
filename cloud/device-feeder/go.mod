module github.com/ukama/ukamaX/cloud/device-feeder

go 1.16

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	github.com/rabbitmq/amqp091-go v1.2.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.7.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukamaX/cloud/registry v0.0.0-20211209085225-8d6ae39819e5
	github.com/ukama/ukamaX/common v0.0.0-20211209085225-8d6ae39819e5
	github.com/wagslane/go-rabbitmq v0.7.0 // indirect
	google.golang.org/grpc v1.42.0
)

replace github.com/ukama/ukamaX/common => ../../common

replace github.com/ukama/ukamaX/cloud/device-feeder => ./

replace github.com/ukama/ukamaX/cloud/device-feeder/mocks => ./mocks

replace github.com/ukama/ukamaX/cloud/registry => ../registry
