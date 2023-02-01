//go:build integration
// +build integration

package integration

import (
	"github.com/ukama/ukama/systems/common/config"

	rconf "github.com/num30/config"
	"github.com/sirupsen/logrus"
)

var tConfig *TestConfig

func init() {
	// load config
	tConfig = &TestConfig{}

	reader := rconf.NewConfReader("integration")

	err := reader.Read(tConfig)
	if err != nil {
		logrus.Fatalf("Failed to read config: %v", err)
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("Config: %+v\n", tConfig)
}

type TestConfig struct {
	ServiceHost string        `default:"localhost:9090"`
	Queue       *config.Queue `default:"{}"`
}
