module lwm2m-gateway

go 1.16

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/pelletier/go-toml v1.8.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	github.com/theherk/viper v0.0.0-20171202031228-e0502e82247d
	github.com/ukama/ukamaX/common v0.0.0-20210910150531-bb65155448ea
	google.golang.org/protobuf v1.27.1

)

replace github.com/ukama/ukamaX/common => ../../common
