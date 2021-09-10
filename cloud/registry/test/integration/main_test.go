// +build integration

package integration

import (
	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/ukama/ukamaX/common/config"
	"testing"
)

func TestExampleTestSuite(t *testing.T) {
	// Run all tests in suite
	suite.Run(t, NewIntegrationTestSuite(loadConfig()))
}

func loadConfig() *TestConfig {
	testConf := &TestConfig{
		RegistryHost: "localhost:9090",
		Rabbitmq:     "amqp://guest:guest@localhost:5672",
	}
	b, err := yaml.Marshal(testConf)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: REGISTRYHOST")
	logrus.Infoln(string(b))

	config.LoadConfig("integration", testConf)

	return testConf
}
