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
	Host     string
	Password string
	DbName   string
	Username string
	SslEnabled bool
	Port int
}

// LoadConfig loads configuration into `config` object
// It looks for env vars
// and config yml file with the same name as binary (ex registry.yml)
// ENV vars have precedence
// configName - name of config file without extension
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

	// Search config in home directory with name ".dmr" (without extension).
	e.AddConfigPath(home)
	e.SetConfigName("." + configFileName)

	e.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err = e.ReadInConfig()
	if err == nil {
		logrus.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		logrus.Debugf("Config file was not loaded. Reason: %v", err)
	}

	err = e.Unmarshal(config)
	if err != nil {
		logrus.Fatalf("Unable to decode into struct, %v", err)
	}
}
