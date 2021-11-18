module github.com/ukama/ukamaX/common

go 1.16

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/iamolegga/enviper v1.2.1
	github.com/jackc/pgconn v1.8.1
	github.com/lib/pq v1.3.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.5.1
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
)

replace github.com/ukama/ukamaX/common => ../../common
