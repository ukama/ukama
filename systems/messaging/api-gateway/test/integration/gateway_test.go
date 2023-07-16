//go:build integration
// +build integration

package integration

import (
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"

	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Before running test for the first time you have to create a test account in Identity manager and provide email and password for it

type TestConfig struct {
	DebugMode        bool
	DummyAuthMode    bool
	TestAccountEmail string
	TestAccountPass  string
	ApiUrl           string
}

var testConf *TestConfig

func init() {
	testConf = &TestConfig{
		DebugMode:        true,
		DummyAuthMode:    false,
		TestAccountEmail: "integration-test@ukama.com",
		TestAccountPass:  "Pass2020!!",
		ApiUrl:           "https://api.dev.ukama.com",
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
	config.LoadConfig("integration", testConf)
	logrus.Infof("Config: %+v", testConf)
}

func Test_MessagingApi(t *testing.T) {
	client := resty.New().EnableTrace().SetDebug(testConf.DebugMode)
	t.Run("Prometheus", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/prometheus")

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
		}

	})

	t.Run("GetNodeIPList", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/nns/list")

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
		}

	})
}

func getApiUrl() string {
	return testConf.ApiUrl
}
