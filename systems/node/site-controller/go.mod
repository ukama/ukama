module github.com/ukama/ukama/systems/node/site-controller

go 1.23.0

replace github.com/ukama/ukama/systems/common => ../../common
replace github.com/ukama/ukama/systems/node/controller => ../controller

require (
	github.com/num30/config v0.1.3
	github.com/sirupsen/logrus v1.9.3
	github.com/ukama/ukama/systems/common v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/node/controller v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.73.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.30.0
)
