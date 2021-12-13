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
		HssHost: "localhost:9090",
	}

	config.LoadConfig("integration", testConf)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", testConf)

	return testConf
}
