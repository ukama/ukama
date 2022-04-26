module ukamaX/bootstrap/lookup

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/gin-gonic/gin v1.7.2
	github.com/go-playground/validator/v10 v10.6.1 // indirect
	github.com/go-resty/resty/v2 v2.6.0
	github.com/jackc/pgtype v1.7.0
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukamaX/common v0.0.0-00010101000000-000000000000
	github.com/ysugimoto/grpc-graphql-gateway v0.22.0 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
)

replace github.com/ukama/ukamaX/common => ../../common
