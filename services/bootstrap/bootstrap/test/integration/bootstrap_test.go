//go:build integration
// +build integration

package integration

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/ukama/ukama/services/common/config"
)

type TestConfig struct {
	BootstrapHost string
}

type IntegrationTestSuite struct {
	suite.Suite
	config *TestConfig
}

func (t *IntegrationTestSuite) SetupSuite() {
	t.config = loadConfig()
}

func loadConfig() *TestConfig {
	testConf := &TestConfig{
		BootstrapHost: "http://bootstrap",
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars")
	config.LoadConfig("integration", testConf)
	logrus.Infof("%+v", testConf)

	return testConf
}

func (i *IntegrationTestSuite) Test_BootstrapApi() {
	client := resty.New()

	i.Run("Ping", func() {
		resp, err := client.R().
			EnableTrace().
			Get(i.config.BootstrapHost + "/ping")

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().NotEmpty(resp.String())
	})

}
