module github.com/ukama/ukama/systems/hub/distributor

go 1.20

replace github.com/ukama/ukama/systems/common => ../../common

replace github.com/ukama/ukama/systems/services/msgClient => ../../services/msgClient

replace github.com/ukama/ukama/systems/hub/hub => ../hub

replace github.com/ukama/ukama/systems/hub/distributor => ./

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/folbricht/desync v0.9.5
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.1
	github.com/go-resty/resty/v2 v2.7.0
	github.com/loopfz/gadgeto v0.11.2
	github.com/mholt/archiver/v3 v3.5.1
	github.com/minio/minio-go/v6 v6.0.57
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
	github.com/ukama/ukama/systems/common v0.0.0-20230831053137-df467d889284
	github.com/ukama/ukama/systems/hub/hub v0.0.0-00010101000000-000000000000
	github.com/wI2L/fizz v0.22.0
)

require (
	cloud.google.com/go v0.110.0 // indirect
	cloud.google.com/go/compute v1.19.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v0.13.0 // indirect
	cloud.google.com/go/storage v1.30.1 // indirect
	github.com/DataDog/zstd v1.5.2 // indirect
	github.com/andybalholm/brotli v1.0.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boljen/go-bitmap v0.0.0-20151001105940-23cd2fb0ce7d // indirect
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dchest/siphash v1.2.3 // indirect
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/folbricht/tempfile v0.0.1 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.8.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0 // indirect
	github.com/hanwen/go-fuse/v2 v2.2.0 // indirect
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
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.61 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mwitkow/go-proto-validators v0.3.2 // indirect
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/penglongli/gin-metrics v0.1.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/sftp v1.13.5 // indirect
	github.com/pkg/xattr v0.4.9 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.12.2 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.34.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rabbitmq/amqp091-go v1.7.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.12.0 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.4.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	github.com/ulikunitz/xz v0.5.9 // indirect
	github.com/wagslane/go-rabbitmq v0.12.4 // indirect
	github.com/willf/bitset v1.1.11 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.12.0 // indirect
	golang.org/x/oauth2 v0.7.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/term v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.116.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230403163135-c38d8f061ccd // indirect
	google.golang.org/grpc v1.54.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/cheggaaa/pb.v1 v1.0.28 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.3.5 // indirect
	gorm.io/gorm v1.24.3 // indirect
)
