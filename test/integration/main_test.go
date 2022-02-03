//go:build integration
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
		QueueUri:     "amqp://guest:guest@localhost:5672/",
		RegistryHost: "localhost:9090",
		NetHost:      "localhost:9090",
		DevicePort:   8080,
		WaitingTime:  10,
		DevicesCount: 3,
	}

	config.LoadConfig("integration", testConf)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", testConf)

	return testConf
}
