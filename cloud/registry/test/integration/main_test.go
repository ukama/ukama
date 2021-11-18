// +build integration

package integration

import (
	"testing"

	"github.com/ukama/ukamaX/common/config"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

func TestTopLevelTestForSuite(t *testing.T) {
	// Run all tests in suite
	suite.Run(t, NewIntegrationTestSuite(loadConfig()))
}

func loadConfig() *TestConfig {
	testConf := &TestConfig{
		RegistryHost: "localhost:9090",
		Rabbitmq:     "amqp://guest:guest@localhost:5672",
	}

	config.LoadConfig("integration", testConf)

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: REGISTRYHOST")
	//b, err := yaml.Marshal(testConf)
	logrus.Infof("Config: %+v\n", testConf)

	return testConf
}
