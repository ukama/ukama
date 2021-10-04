// +build integration

package integration

import (
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/common/ukama"
	"gopkg.in/yaml.v3"
	"net/http"
)

type TestConfig struct {
	LookupHost string
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
		LookupHost: "http://lookup",
	}
	b, err := yaml.Marshal(testConf)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.Info("Expected config ", "integration.yaml", " or env vars")
	logrus.Infoln(string(b))

	config.LoadConfig("integration", testConf)

	return testConf
}

func (i *IntegrationTestSuite) Test_LookuApi() {
	client := resty.New()
	nodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_HOMENODE)

	i.Run("Ping", func() {
		resp, err := client.R().
			EnableTrace().
			Get(i.config.LookupHost + "/ping")

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().NotEmpty(resp.String())
	})

	const orgName = "lookub-test-org-1"
	i.Run("AddOrg", func() {
		resp, err := client.R().
			EnableTrace().
			SetBody(`{	"certificate":"cert", "ip": "127.0.0.1"	}`).
			Post(i.config.LookupHost + "/orgs/" + orgName)

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().NotEmpty(resp.String())
	})

	i.Run("AddDevice", func() {
		resp, err := client.R().
			EnableTrace().
			Post(i.config.LookupHost + "/orgs/" + orgName + "/devices/" + nodeId.String())

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().NotEmpty(resp.String())
	})

	i.Run("GetDevice", func() {
		resp, err := client.R().
			EnableTrace().
			Get(i.config.LookupHost + "/orgs/" + orgName + "/devices/" + nodeId.String())

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().Contains(resp.String(), "127.0.0.1")
		i.Assert().Contains(resp.String(), "certificate")
	})

}
