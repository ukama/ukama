module github.com/ukama/ukama/testing/common

go 1.23.0

replace github.com/ukama/ukama/testing/common => ./

require (
	github.com/sirupsen/logrus v1.9.3
	google.golang.org/protobuf v1.36.5
	gopkg.in/yaml.v3 v3.0.1
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
