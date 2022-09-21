//go:build integration
// +build integration

package integration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/ukama"
	"net/http"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Before running test for the first time you have to create a test account in Identity manager and provide email and password for it

type TestConfig struct {
	ApiUrl string
}

var testConf *TestConfig

func init() {
	testConf = &TestConfig{
		ApiUrl: "http://0.0.0.0:8080",
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
	config.LoadConfig("integration", testConf)
	logrus.Infof("Config: %+v", testConf)
}

func Test_LookupClientApi(t *testing.T) {

	org := "org-name"
	nodeId := ukama.NewVirtualHomeNodeId().String()
	system := "sys-name"

	client := resty.New()

	t.Run("GetOrgNotFound", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/lookup/orgs/someOrgWhichDoesnotExist")

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusNotFound, resp.StatusCode())
			assert.Contains(tt, resp.String(), "org record not found")
		}
	})

	t.Run("AddOrg", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetBody(strings.NewReader(`{"Certificate": "helloOrg","Ip": "0.0.0.0"}`)).
			Put(getApiUrl() + "/lookup/orgs/" + org)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusCreated, resp.StatusCode())
		}

	})

	t.Run("GetOrg", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/lookup/orgs/" + org)

		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
	})

	t.Run("AddNode", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Put(getApiUrl() + "/lookup/orgs/" + org + "/nodes/" + nodeId)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusCreated, resp.StatusCode())
		}
	})

	t.Run("GetNodes", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/lookup/orgs/" + org + "/nodes/" + nodeId)
		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
		fmt.Println("Response: ", resp.String())
	})

	t.Run("DeleteNode", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Delete(getApiUrl() + "/lookup/orgs/" + org + "/nodes/" + nodeId)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
		}
	})

	t.Run("AddSystem", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetBody(strings.NewReader(`{ "org":"org-name", "system":"sys", "ip":"0.0.0.0", "certificate":"certs", "port":100}`)).
			Put(getApiUrl() + "/lookup/orgs/" + org + "/systems/" + system)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusCreated, resp.StatusCode())
			fmt.Println("Response: ", resp.String())
		}
	})

	t.Run("GetSystems", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/lookup/orgs/" + org + "/systems/" + system)
		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
		fmt.Println("Response: ", resp.String())
	})

	t.Run("DeleteSystems", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Delete(getApiUrl() + "/lookup/orgs/" + org + "/systems/" + system)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
		}
	})
}

func getApiUrl() string {
	return testConf.ApiUrl
}
