//go:build integration
// +build integration

package integration

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

type TestConfig struct {
	BaseUrl string
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

	imsi := fmt.Sprintf("00000%010d", time.Now().Unix())

	i.Run("AddGuti", func() {
		resp, err := client.R().
			EnableTrace().
			SetBody(fmt.Sprintf(`{ "imsi":"%s", "guti":  { "plmn_id": "000001", "mmegi": "0",  "mmec":"0", "mtmsi": "%d"  }}`, imsi, time.Now().Unix())).
			Post(i.config.BaseUrl + "/hss/guti")

		logrus.Infof("response: %s", resp.Body())
		i.NoError(err)
		i.Equal(200, resp.StatusCode())
	})

	i.Run("UpdateTai", func() {
		resp, err := client.R().
			EnableTrace().
			SetBody(fmt.Sprintf(`{ "imsi":"%s", "tac": 1234567,  "plmn_id": "000001"  }`, imsi)).
			Post(i.config.BaseUrl + "/hss/tai")

		logrus.Infof("response: %s", resp.Body())
		i.NoError(err)
		i.Equal(http.StatusNotFound, resp.StatusCode())
	})
}
