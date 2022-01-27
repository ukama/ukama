package config

import (
	"github.com/iamolegga/enviper"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Common properties for all configs.
// Don't forget to use `mapstructure:",squash"`. See unittest for example
type BaseConfig struct {
	DebugMode bool
}

type Database struct {
	Host       string
	Password   string
	DbName     string
	Username   string
	SslEnabled bool
	Port       int
}

type Queue struct {
	Uri string // Env var name: QUEUE_URI or in file Queue: { Uri: "" }. Example: QUEUE_URI=amqp://guest:guest@localhost:5672/
}

type Grpc struct {
	Port int
}

type Http struct {
	Port int
}

type Metrics struct {
	Port    int
	Enabled bool
}

// LoadConfig loads configuration into `config` object
// Pulls configuration from env vars and config file
// Config should be a yaml file with `configFileName` and 'yaml' extension
// for example: `.registry`
// Available paths: $HOME, same as binary
// ENV vars have precedence
// configFileName - name of config file without extension.
// Config file should have yaml format and property Names should start with lowercase latter
// Evn var should be uppercased
func LoadConfig(configFileName string, config interface{}) {

	e := enviper.New(viper.New())
	e.SetConfigType("yaml")

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	e.AddConfigPath(home)
	e.AddConfigPath("")
	e.SetConfigName(configFileName + ".yaml")

	e.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err = e.ReadInConfig()
	if err == nil {
		logrus.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		logrus.Infof("Config file was not loaded. Reason: %v\n", err)
	}

	err = e.Unmarshal(config)
	if err != nil {
		logrus.Fatalf("Unable to decode into struct, %v", err)
	}
}

func DefaultMetrics() *Metrics {
	return &Metrics{
		Enabled: true,
		Port:    10250,
	}
}
