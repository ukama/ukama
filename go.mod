module github.com/ukama/ukamaX/common

go 1.16

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/iamolegga/enviper v1.2.1
	github.com/jackc/pgconn v1.8.1
	github.com/lib/pq v1.3.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/ory/kratos-client-go v0.8.0-alpha.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20211205041911-012df41ee64c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211203200212-54befc351ae9 // indirect
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
)

replace github.com/ukama/ukamaX/common => ../../common
