module github.com/ukama/ukama/systems/analytics/analytics

go 1.25.0

replace github.com/ukama/ukama/systems/common => ../../common
replace github.com/ukama/ukama/systems/analytics/business => ../business
replace github.com/ukama/ukama/systems/analytics/customer => ../customer
replace github.com/ukama/ukama/systems/analytics/network => ../network
replace github.com/ukama/ukama/systems/analytics/schema => ../schema
replace github.com/ukama/ukama/systems/analytics/analytics => ./

require (
	github.com/num30/config v0.1.3
	github.com/sirupsen/logrus v1.9.4
	github.com/ukama/ukama/systems/analytics/business v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/analytics/customer v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/analytics/network v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/analytics/schema v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/common v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.81.1
	gopkg.in/yaml.v2 v2.4.0
)
