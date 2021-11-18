module github.com/ukama/ukamaX/cloud/hss

go 1.16

require (
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/pkg/errors v0.9.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukamaX/common v0.0.0-20211015093708-cd6e230254b5
	golang.org/x/net v0.0.0-20210716203947-853a461950ff // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210719143636-1d5a45f8e492 // indirect
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gorm.io/gorm v1.21.16
)

replace github.com/ukama/ukamaX/common => ../../common

replace github.com/ukama/ukamaX/cloud/hss => ./

replace github.com/ukama/ukamaX/cloud/hss/mocks => ./mocks
