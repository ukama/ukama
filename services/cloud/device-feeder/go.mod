module github.com/ukama/ukamaX/cloud/device-feeder

go 1.18

require (
	github.com/coredns/coredns v1.8.7
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/rabbitmq/amqp091-go v1.2.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/ukama/ukamaX/cloud/net v0.0.0-20220422163321-18f195cf8dfa
	github.com/ukama/ukamaX/cloud/registry v0.0.0-20220422163321-18f195cf8dfa
	github.com/ukama/ukamaX/common v0.0.0-20220422163321-18f195cf8dfa
	github.com/wagslane/go-rabbitmq v0.7.0
	google.golang.org/grpc v1.45.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/iamolegga/enviper v1.2.1 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/mwitkow/go-proto-validators v0.3.2 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220422154200-b37d22cd5731 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/ukama/ukamaX/cloud/device-feeder => ./

replace github.com/ukama/ukamaX/cloud/device-feeder/mocks => ./mocks
