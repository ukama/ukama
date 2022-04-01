module github.com/ukama/ukamaX/cloud/device-gateway

go 1.17

require (
	github.com/go-resty/resty/v2 v2.7.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.1
	github.com/ukama/ukamaX/cloud/hss v0.0.0-20211125132149-eb580caa29dd
	github.com/ukama/ukamaX/common v0.0.0-20211015093708-cd6e230254b5
	google.golang.org/grpc v1.45.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/iamolegga/enviper v1.2.1 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/afero v1.1.2 // indirect
	github.com/spf13/cast v1.3.0 // indirect
	github.com/spf13/jwalterweatherman v1.0.0 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

require (
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/mwitkow/go-proto-validators v0.3.2 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220317150908-0efb43f6373e // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/ukama/ukamaX/cloud/hss => ../hss
