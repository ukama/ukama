package config

import (
	"github.com/iamolegga/enviper"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/ukama/ukama/services/common/rest"
	cors "github.com/gin-contrib/cors"
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

/*
Route are key value pair
*/
type Route map[string]string

/*
Service would a have some finite number of routes
thorugh which it could be reached.
*/
type Pattern struct {
	Routes []Route
}

/*
Service would be listing on this
IP and Port for incoming messages
*/
type Forward struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	Path string `json:"default_path"`
}

/*
Service API interface is registered to Router service.
So that Router service know when and where to reach service.
*/
type ServiceApiIf struct {
	Name string  `json:"name"`
	P    []Route `json:"patterns"`
	F    Forward `json:"forward"`
}

type Queue struct {
	Uri string // Env var name: QUEUE_URI or in file Queue: { Uri: "" }. Example: QUEUE_URI=amqp://guest:guest@localhost:5672/
}

type Grpc struct {
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
	e.AddConfigPath("./")
	e.SetConfigName(configFileName + ".yaml")

	e.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err = e.ReadInConfig()
	if err == nil {
		logrus.Info("Using config file:", e.ConfigFileUsed())
	} else {
		logrus.Infof("Config file was not loaded. Reason: %v\n", err)
	}

	err = e.Unmarshal(config)
	if err != nil {
		logrus.Fatalf("Unable to decode into struct, %v", err)
	}

	if e.GetBool("DebugMode") {
		logrus.Infoln("Debug mode is enabled")
		logrus.Infoln("vvvv Config file vvvv")
		logrus.Infof("%+v", config)
	}
}

func DefaultMetrics() *Metrics {
	return &Metrics{
		Enabled: true,
		Port:    10250,
	}
}

func DefaultDatabase() Database {
	return Database{
		Host:       "localhost",
		Password:   "Pass2020!",
		DbName:     "registry",
		Username:   "postgres",
		Port:       5432,
		SslEnabled: false,
	}
}

func DefaultHTTPConfig() rest.HttpConfig {
	return rest.HttpConfig{
		Port: 8080,
		Cors: cors.Config{
			AllowOrigins: []string{"http://localhost", "https://localhost", "*"},
		},
	}
}