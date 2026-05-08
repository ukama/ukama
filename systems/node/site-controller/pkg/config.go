package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	DB                *config.Database                    `mapstructure:"db"`
	Grpc              *grpc.Grpc                          `mapstructure:"grpc"`
	Queue             config.Queue                        `mapstructure:"queue"`
	Service           *config.Service                     `mapstructure:"service"`
	MsgClient         msgBusServiceClient.MsgClientConfig `mapstructure:"msgClient"`
	Services          GrpcEndpoints                       `mapstructure:"services"`
	OrgName           string                              `mapstructure:"orgName"`
}

type GrpcEndpoints struct {
	Timeout    time.Duration
	Controller string
}

func NewConfig(name string) *Config {
	return &Config{
		BaseConfig: config.BaseConfig{DebugMode: false},
		DB:         config.DefaultDatabase(),
		Grpc:       grpc.DefaultGrpc(),
		Queue:      *config.DefaultQueue(),
		Service:    config.LoadServiceHostConfig(name),
		MsgClient:  *msgBusServiceClient.DefaultMsgClientConfig(),
		Services:   GrpcEndpoints{Timeout: 3 * time.Second, Controller: "controller:9090"},
	}
}
