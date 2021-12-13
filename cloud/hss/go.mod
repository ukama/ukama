module github.com/ukama/ukamaX/cloud/hss

go 1.16

require (
	github.com/golang/protobuf v1.5.2
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukamaX/common v0.0.0-20211015093708-cd6e230254b5
	google.golang.org/genproto v0.0.0-20211203200212-54befc351ae9
	google.golang.org/grpc v1.42.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	gorm.io/gorm v1.21.16
)

replace github.com/ukama/ukamaX/common => ../../common

replace github.com/ukama/ukamaX/cloud/hss => ./

replace github.com/ukama/ukamaX/cloud/hss/mocks => ./mocks
