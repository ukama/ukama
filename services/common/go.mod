module github.com/ukama/ukamaX/common

go 1.16

require (
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.7.7
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0
	github.com/iamolegga/enviper v1.2.1
	github.com/jackc/pgconn v1.8.1
	github.com/lib/pq v1.3.0
	github.com/loopfz/gadgeto v0.9.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/ory/kratos-client-go v0.8.2-alpha.1
	github.com/penglongli/gin-metrics v0.1.9
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.7.0
	github.com/wI2L/fizz v0.18.1
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
)

replace github.com/ukama/ukamaX/common => ../../common
