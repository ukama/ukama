module github.com/ukama/ukama/testing/integration

go 1.20

replace github.com/ukama/ukama/systems/common => ../../systems/common

replace github.com/ukama/ukama/systems/init/api-gateway => ../../systems/init/api-gateway

replace github.com/ukama/ukama/systems/init/lookup => ../../systems/init/lookup

replace github.com/ukama/ukama/systems/subscriber/api-gateway => ../../systems/subscriber/api-gateway

replace github.com/ukama/ukama/systems/subscriber/registry => ../../systems/subscriber/registry

replace github.com/ukama/ukama/systems/subscriber/sim-pool => ../../systems/subscriber/sim-pool

replace github.com/ukama/ukama/systems/subscriber/sim-manager => ../../systems/subscriber/sim-manager

replace github.com/ukama/ukama/systems/data-plan/api-gateway => ../../systems/data-plan/api-gateway

replace github.com/ukama/ukama/systems/data-plan/rate => ../../systems/data-plan/rate

replace github.com/ukama/ukama/systems/data-plan/base-rate => ../../systems/data-plan/base-rate

replace github.com/ukama/ukama/systems/data-plan/package => ../../systems/data-plan/package

replace github.com/ukama/ukama/systems/registry/api-gateway => ../../systems/registry/api-gateway

replace github.com/ukama/ukama/systems/nucleus/org => ../../systems/nucleus/org

replace github.com/ukama/ukama/systems/nucleus/user => ../../systems/nucleus/user

replace github.com/ukama/ukama/systems/registry/network => ../../systems/registry/network

replace github.com/ukama/ukama/systems/registry/node => ../../systems/registry/node

replace github.com/ukama/ukama/systems/billing/invoice => ../../systems/billing/invoice

require (
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/iamolegga/enviper v1.4.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.12.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.11.0 // indirect
	github.com/jackc/pgx/v4 v4.16.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.2 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.12.0 // indirect
	github.com/subosito/gotenv v1.4.0 // indirect
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.4.0 // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.3.5 // indirect
	gorm.io/gorm v1.24.3 // indirect
)
