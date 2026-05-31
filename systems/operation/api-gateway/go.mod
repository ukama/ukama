module github.com/ukama/ukama/systems/operation/api-gateway

go 1.24.0

toolchain go1.24.12

replace github.com/ukama/ukama/systems/common => ../../common

replace github.com/ukama/ukama/systems/operation/manager => ../manager

require (
	github.com/gin-contrib/cors v1.7.6
	github.com/gin-gonic/gin v1.10.1
	github.com/loopfz/gadgeto v0.11.5
	github.com/sirupsen/logrus v1.9.4
	github.com/ukama/ukama/systems/common v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/operation/manager v0.0.0-00010101000000-000000000000
	github.com/wI2L/fizz v0.22.0
	google.golang.org/grpc v1.80.0
)
