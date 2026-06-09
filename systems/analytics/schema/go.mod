module github.com/ukama/ukama/systems/analytics/schema

go 1.25.0

replace github.com/ukama/ukama/systems/common => ../../common

require (
	github.com/ukama/ukama/systems/common v0.0.0-00010101000000-000000000000
	gorm.io/datatypes v1.2.7
	gorm.io/gorm v1.30.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gorm.io/driver/mysql v1.5.6 // indirect
)
