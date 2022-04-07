//go:build integration
// +build integration

package integration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/common/ukama"
	"net/http"
	"testing"
	"time"

	ory "github.com/ory/kratos-client-go"
	"github.com/ukama/ukamaX/common/testing/kratos"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type TestConfig struct {
	BaseDomain       string
	DummyAuthMode    bool
	TestAccountEmail string
	TestAccountPass  string
}

var testConf *TestConfig

func init() {
	testConf = &TestConfig{
		BaseDomain:       "dev.ukama.com",
		DummyAuthMode:    false,
		TestAccountEmail: "integration-test@ukama.com",
		TestAccountPass:  "Pass2020!!",
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
	config.LoadConfig("integration", testConf)
	logrus.Infof("Config: %+v", testConf)
}

func Test_RegistryApi(t *testing.T) {
	var login *ory.SuccessfulSelfServiceLoginWithoutBrowser
	var err error

	if testConf.DummyAuthMode {
		tkn := "dummy-token"
		login = &ory.SuccessfulSelfServiceLoginWithoutBrowser{
			SessionToken: &tkn,
		}
	} else {

		login, err = kratos.Login(getKratosUrl(), testConf.TestAccountEmail, testConf.TestAccountPass)
		if err != nil {
			assert.NoError(t, err, "Failed to login to Kratos")
			assert.FailNow(t, "Failed to login to Kratos")
			return
		}
	}

	time.Sleep(3 * time.Second) //give registry some time to create a default org for account

	client := resty.New()

	kratos.PrintJSONPretty(login)

	t.Run("GetOrg", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			Get(getApiUrl() + "/orgs/" + login.Session.Identity.GetId())

		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
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

	nodeId := ukama.NewVirtualHomeNodeId().String()
	t.Run("AddNode", func(tt *testing.T) {
		nodeName := time.Now().Format(time.RFC3339) + "testNode"
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			SetBody(fmt.Sprintf("{ 'name':'%s' } ", nodeName)).
			Put(getApiUrl() + "/orgs/" + login.Session.Identity.GetId() + "/nodes/" + nodeId)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusCreated, resp.StatusCode())
			assert.Contains(tt, resp.String(), nodeName)
			fmt.Println("Response: ", resp.String())
		}
	})

	t.Run("UpdateNode", func(tt *testing.T) {
		nodeName := time.Now().Format(time.RFC3339) + "testNode"
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			SetBody(fmt.Sprintf("{ 'name':'updated-%s' } ", nodeName)).
			Put(getApiUrl() + "/orgs/" + login.Session.Identity.GetId() + "/nodes/" + nodeId)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
			fmt.Println("Response: ", resp.String())
			assert.Contains(tt, resp.String(), nodeName)
		}
	})

}

func getApiUrl() string {
	return "https://api." + testConf.BaseDomain
}

func getKratosUrl() string {
	return "https://auth." + testConf.BaseDomain + "/.api"
}
