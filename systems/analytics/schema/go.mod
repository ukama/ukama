module github.com/ukama/ukama/systems/analytics/schema

go 1.25.0

replace github.com/ukama/ukama/systems/common => ../../common

require (
	github.com/ukama/ukama/systems/common v0.0.0-00010101000000-000000000000
	gorm.io/datatypes v1.2.7
)
