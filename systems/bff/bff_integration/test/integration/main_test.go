package integration

import (
	"testing"

	"github.com/go-yaml/yaml"
	"github.com/iamolegga/enviper"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

func TestTopLevelTestForSuite(t *testing.T) {
	// Run all tests in suite
	suite.Run(t, NewIntegrationTestSuite(loadConfig()))
}

func loadConfig() *TestConfig {
	testConf := &TestConfig{
		BaseDomain: "dev.ukama.com",
	}
	b, err := yaml.Marshal(testConf)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	CommonLoadConfig("integration", testConf)

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infoln(string(b))

	return testConf
}

func CommonLoadConfig(configFileName string, config interface{}) {

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
