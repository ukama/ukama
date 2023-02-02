//go:build integration
// +build integration

package integration

import (
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	"net/http"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Before running test for the first time you have to create a test account in Identity manager and provide email and password for it

type TestConfig struct {
	ServiceHost string `default:"localhost:8080"`
}

var testConf *TestConfig

var iccid = "012345678901234567891"
var network = "40987edb-ebb6-4f84-a27c-99db7c136127"

// var orgId = "880f7c63-eb57-461a-b514-248ce91e9b3e"
var packageId = "8adcdfb4-ed30-405d-b32f-d0b2dda4a1e0"

func init() {
	testConf = &TestConfig{}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
	config.LoadConfig("integration", testConf)
	logrus.Infof("Config: %+v", testConf)
}

func Test_UkamaAgentClientApi(t *testing.T) {

	client := resty.New()

	t.Run("Activate", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetBody(strings.NewReader(`{"network":"` + network + `","packageId":"` + packageId + `"}`)).
			Put(getApiUrl() + "/v1/subscriber/" + iccid)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusCreated, resp.StatusCode())
		}
	})

	t.Run("UpdatePackage", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetBody(strings.NewReader(`{"packageId":"` + packageId + `"}`)).
			Patch(getApiUrl() + "/v1/subscriber/" + iccid)

		if err != nil {
			if assert.Error(t, err) {
				assert.Equal(tt, http.StatusOK, resp.StatusCode())
			}
		}
	})

	t.Run("Read", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/subscriber/" + iccid)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
			assert.Contains(tt, iccid, resp.String())
		}
	})

	t.Run("Inactivate", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Delete(getApiUrl() + "/v1/subscriber/" + iccid)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
		}
	})

}

func getApiUrl() string {
	return testConf.ServiceHost
}
