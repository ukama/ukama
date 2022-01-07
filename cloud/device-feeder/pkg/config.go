package pkg

import (
	"github.com/ukama/ukamaX/common/config"
)

type Registry struct {
	Host           string
	TimeoutSeconds int
}

type DeviceNetworkConfig struct {
	Port           uint // set to 0 to bypass port addition
	TimeoutSeconds uint // timeout for one request to a device
}

type ListenerConfig struct {
	ExecutionRetryCount int64 // max retries count
	RetryPeriodSec      int   // how long request waits after failure to try again
	Threads             int   // how many go routines run message processor
}

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	Queue             config.Queue
	Registry          Registry
	Device            DeviceNetworkConfig
	Listener          ListenerConfig
	Metrics           *config.Metrics
}

func NewConfig() *Config {

	return &Config{
		Queue: config.Queue{
			Uri: "amqp://guest:guest@rabbitmq:5672/",
		},
		Registry: Registry{
			Host:           "registry:9090",
			TimeoutSeconds: 3,
		},
		Device: DeviceNetworkConfig{
			Port:           0,
			TimeoutSeconds: 3,
		},
		Listener: ListenerConfig{
			ExecutionRetryCount: 3,
			RetryPeriodSec:      30,
			Threads:             3,
		},
		Metrics: config.DefaultMetrics(),
	}
}
