// +build integration

package integration

import (
	"testing"

	"github.com/ukama/ukamaX/common/config"

	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

func TestTopLevelTestForSuite(t *testing.T) {
	// Run all tests in suite
	suite.Run(t, NewIntegrationTestSuite(loadConfig()))
}

func loadConfig() *TestConfig {
	testConf := &TestConfig{
		FooHost: "localhost:9090",
	}
	b, err := yaml.Marshal(testConf)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	config.LoadConfig("integration", testConf)

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infoln(string(b))

	return testConf
}
