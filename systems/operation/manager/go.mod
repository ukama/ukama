module github.com/ukama/ukama/systems/operation/manager

go 1.24.0

toolchain go1.24.12

replace github.com/ukama/ukama/systems/common => ../../common

replace github.com/ukama/ukama/systems/services/msgClient => ../../services/msgClient

require (
	github.com/golang/protobuf v1.5.4
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/num30/config v0.1.3
	github.com/sirupsen/logrus v1.9.4
	github.com/ukama/ukama/systems/common v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
)
