//go:build integration
// +build integration

package integration

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/config"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

type TestConfig struct {
	BaseUrl string
}

func loadConfig() *TestConfig {
	testConf := &TestConfig{
		BaseUrl: "http://localhost:8080",
	}

	config.LoadConfig("integration", testConf)

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEURL")
	logrus.Infof("%+v", testConf)

	return testConf
}

type IntegrationTestSuite struct {
	suite.Suite
	config *TestConfig
}

func NewIntegrationTestSuite(config *TestConfig) *IntegrationTestSuite {
	return &IntegrationTestSuite{config: config}
}

func (i *IntegrationTestSuite) Test_Ping() {
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Get(i.config.BaseUrl + "/ping")
	i.NoError(err)
	i.Equal(200, resp.StatusCode())
}

func (i *IntegrationTestSuite) Test_SwaggerUI() {
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Get(i.config.BaseUrl + "/swagger")
	i.NoError(err)
	i.Equal(200, resp.StatusCode())
}

func (i *IntegrationTestSuite) handleResponse(err error, r interface{}) {
	fmt.Printf("Response: %v\n", r)
	i.Assert().NoErrorf(err, "Request failed: %v\n", err)
}

func (i *IntegrationTestSuite) Test_HssIntegration() {
	client := resty.New()

	imsi := "000000845466094"
	// we check for 404 not found to make sure that device-gateway hitting the hss
	i.Run("AddGuti", func() {
		resp, err := client.R().
			EnableTrace().
			SetBody(fmt.Sprintf(`{ 
							"imsi":"%s", 
							"guti":  { "plmn_id": "000001", "mmegi": "0",  "mmec":"0", "mtmsi": "%d" }, 
							"updated_at": %d }`, imsi, time.Now().Unix(), time.Now().Unix())).
			Post(i.config.BaseUrl + "/hss/guti")

		logrus.Infof("response: %s", resp.Body())
		i.NoError(err)
		i.Equal(http.StatusNotFound, resp.StatusCode())
	})

	i.Run("UpdateTai", func() {
		resp, err := client.R().
			EnableTrace().
			SetBody(fmt.Sprintf(`{ "imsi":"%s", "tac": 1234567,  "plmn_id": "000001", "updated_at": %d  }`,
				imsi, time.Now().Unix())).
			Post(i.config.BaseUrl + "/hss/tai")

		logrus.Infof("response: %s", resp.Body())
		i.NoError(err)
		i.Equal(http.StatusNotFound, resp.StatusCode())
	})
}
