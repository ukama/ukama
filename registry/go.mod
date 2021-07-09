module ukamaX/registry

go 1.16

require (
	github.com/iamolegga/enviper v1.2.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/ukama/ukamaX/common v0.0.0-00010101000000-000000000000
	gorm.io/gorm v1.21.10
)

replace github.com/ukama/ukamaX/common => ../common
