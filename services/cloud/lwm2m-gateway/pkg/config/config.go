package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	LogLevel      string
	SchemaDir     string
	DevConfigType string
	Lwm2mServer   ServerConfig
	Lwm2mGateway  ServerConfig
	MsgBus        ServerConfig
}

// Server Configuration
type ServerConfig struct {
	Address string
	Port    string
}

var Config Configuration

// LoadConfig loads config from files
func LoadConfig(cfgName string, cfgType string, path string) error {

	v := viper.New()

	// Set the file name of the configurations file
	v.SetConfigName(cfgName)

	// Set the file type of the configurations file
	v.SetConfigType(cfgType)

	// Set the path to look for the configurations file
	v.AddConfigPath(path)

	// Enable VIPER to read Environment Variables
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Errorf("Config:: Failed to read the configuration file: %s. Err %s", cfgName, err)
		return err
	}

	err := v.Unmarshal(&Config)
	if err != nil {
		log.Errorf("Config:: Unable to decode into struct, %v", err)
		return err
	}

	log.Debugf("Config:: Loaded service config form %s. Config %+v.", path, Config)
	return nil
}
