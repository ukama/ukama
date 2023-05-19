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

replace github.com/ukama/ukama/systems/registry/org => ../../systems/registry/org

replace github.com/ukama/ukama/systems/registry/users => ../../systems/registry/users

replace github.com/ukama/ukama/systems/registry/network => ../../systems/registry/network

replace github.com/ukama/ukama/systems/registry/node => ../../systems/registry/node

require (
	github.com/bxcodec/faker/v4 v4.0.0-beta.3
	github.com/go-resty/resty/v2 v2.7.0
	github.com/goombaio/namegenerator v0.0.0-20181006234301-989e774b106e
	github.com/rabbitmq/amqp091-go v1.8.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.9.2
	github.com/stretchr/testify v1.8.3
	github.com/tj/assert v0.0.3
	github.com/ukama/ukama/systems/common v0.0.0-20230208235400-d17899b75cbb
	github.com/ukama/ukama/systems/data-plan/api-gateway v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/data-plan/base-rate v0.0.0-20221114075906-a5be6bf1d178
	github.com/ukama/ukama/systems/data-plan/package v0.0.0-20230208235400-d17899b75cbb
	github.com/ukama/ukama/systems/data-plan/rate v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/init/api-gateway v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/init/lookup v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/registry/api-gateway v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/registry/network v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/registry/org v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/subscriber/api-gateway v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/subscriber/registry v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/subscriber/sim-manager v0.0.0-00010101000000-000000000000
	github.com/ukama/ukama/systems/subscriber/sim-pool v0.0.0-00010101000000-000000000000
	github.com/wagslane/go-rabbitmq v0.12.3
	google.golang.org/protobuf v1.30.0
	k8s.io/apimachinery v0.27.2
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bytedance/sonic v1.8.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gin-contrib/cors v1.4.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.9.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.11.2 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/iamolegga/enviper v1.4.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.2.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/loopfz/gadgeto v0.11.2 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mwitkow/go-proto-validators v0.3.2 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/penglongli/gin-metrics v0.1.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.13.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.13.0 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	github.com/subosito/gotenv v1.4.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.9 // indirect
	github.com/ukama/ukama/systems/registry/users v0.0.0-00010101000000-000000000000 // indirect
	github.com/wI2L/fizz v0.22.0 // indirect
	github.com/willf/bitset v1.1.11 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230202175211-008b39050e57 // indirect
	google.golang.org/grpc v1.53.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.4.6 // indirect
	gorm.io/gorm v1.24.3 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
)
