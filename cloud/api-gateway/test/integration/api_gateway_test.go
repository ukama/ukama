//go:build integration
// +build integration

package integration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukamaX/common/config"
	"net/http"
	"testing"
	"time"

	"github.com/ukama/ukamaX/common/testing/kratos"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type TestConfig struct {
	BaseDomain string
}

var testConf *TestConfig

func init() {
	testConf = &TestConfig{
		BaseDomain: "dev.ukama.com",
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")

	config.LoadConfig("integration", testConf)
	logrus.Infof("Config: %+v", testConf)
}

func Test_RegistryApi(t *testing.T) {
	login, err := kratos.Login(getKratosUrl())
	time.Sleep(3 * time.Second) //give registry some time to create a default org for account
	if err != nil {
		assert.NoError(t, err, "Failed to login to Kratos")
		return
	}

	client := resty.New()

	kratos.PrintJSONPretty(login)

	t.Run("GetOrg", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			Get(getApiUrl() + "/orgs/" + login.Session.Identity.GetId())

		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
		//i.Assert().Contains(resp.String(), "Organization not found")
	})

	t.Run("GetOrgNotFound", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			Get(getApiUrl() + "/orgs/" + "someRandomOrgThatShouldNotExist")

		assert.NoError(t, err)
		assert.Equal(tt, http.StatusNotFound, resp.StatusCode())
		assert.Contains(tt, resp.String(), "Organization not found")
	})

	t.Run("GetNodesUnauthorized", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+"random session").
			Get(getApiUrl() + "/orgs/" + login.Session.Identity.GetId() + "/nodes")
		assert.NoError(t, err)
		assert.Equal(tt, http.StatusUnauthorized, resp.StatusCode())
	})

	t.Run("GetNodes", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			Get(getApiUrl() + "/orgs/" + login.Session.Identity.GetId() + "/nodes")
		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
		fmt.Println("Response: ", resp.String())
	})
}

func getApiUrl() string {
	return "https://api." + testConf.BaseDomain
}

func getKratosUrl() string {
	return "https://auth." + testConf.BaseDomain + "/.api"
}
