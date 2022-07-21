module github.com/ukama/ukama/services/cloud/users

go 1.18

require (
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/stretchr/testify v1.8.0
	github.com/ukama/ukama/services/cloud/hss v0.0.0-20220421094812-6b3352c9a822
	github.com/ukama/ukama/services/common v0.0.0-20220421094812-6b3352c9a822
	github.com/vektra/mockery/v2 v2.12.2
	google.golang.org/grpc v1.46.2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0
	google.golang.org/protobuf v1.28.0
	gorm.io/gorm v1.21.16
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/iamolegga/enviper v1.4.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.8.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.0.6 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.7.0 // indirect
	github.com/jackc/pgx/v4 v4.11.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.34.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rabbitmq/amqp091-go v1.3.0 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/cobra v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.12.0 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/subosito/gotenv v1.4.0 // indirect
	github.com/wagslane/go-rabbitmq v0.8.1 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2 // indirect
	golang.org/x/sys v0.0.0-20220702020025-31831981b65f // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.12-0.20220628192153-7743d1d949f1 // indirect
	google.golang.org/genproto v0.0.0-20220519153652-3a47de7e79bd // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.1.0 // indirect
)

replace github.com/ukama/ukama/services/common => ../../common

replace github.com/ukama/ukama/services/cloud/users => ./

replace github.com/ukama/ukama/services/cloud/users/mocks => ./mocks

replace github.com/ukama/ukama/services/cloud/hss => ../hss
