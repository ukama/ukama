//go:build integration
// +build integration

package integration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	userspb "github.com/ukama/ukama/services/cloud/users/pb/gen"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/ukama"
	"net/http"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	ory "github.com/ory/kratos-client-go"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/common/testing/kratos"
)

// Before running test for the first time you have to create a test account in Identity manager and provide email and password for it

type TestConfig struct {
	ApiUrl           string
	KratosUrl        string
	DummyAuthMode    bool
	TestAccountEmail string
	TestAccountPass  string
}

var testConf *TestConfig

func init() {
	testConf = &TestConfig{
		DummyAuthMode:    false,
		TestAccountEmail: "integration-test@ukama.com",
		TestAccountPass:  "Pass2020!!",
		KratosUrl:        "https://auth.dev.ukama.com/.api",
		ApiUrl:           "https://api.dev.ukama.com",
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
	config.LoadConfig("integration", testConf)
	logrus.Infof("Config: %+v", testConf)
}

func Test_RegistryApi(t *testing.T) {
	var login *ory.SuccessfulSelfServiceLoginWithoutBrowser
	var err error

	if testConf.DummyAuthMode {
		fmt.Println("Dummy auth mode enabled")
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

	time.Sleep(3 * time.Second) //give network some time to create a default org for account

	client := resty.New()

	kratos.PrintJSONPretty(login)

	t.Run("GetOrg", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			Get(getApiUrl() + "/orgs/" + login.Session.Identity.GetId())

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
		}

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
			SetBody(map[string]string{"name": nodeName}).
			Put(getApiUrl() + "/orgs/" + login.Session.Identity.GetId() + "/nodes/" + nodeId)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusCreated, resp.StatusCode())
			assert.Contains(tt, resp.String(), nodeName)
			fmt.Println("Response: ", resp.String())
		}
	})

	t.Run("UpdateNode", func(tt *testing.T) {
		nodeName := "updated-testNode-" + time.Now().Format(time.RFC3339)
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			SetBody(map[string]string{"name": nodeName}).
			Put(getApiUrl() + "/orgs/" + login.Session.Identity.GetId() + "/nodes/" + nodeId)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
			fmt.Println("Response: ", resp.String())
			assert.Contains(tt, resp.String(), nodeName)
		}
	})

	t.Run("DeleteNode", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			Delete(getApiUrl() + "/orgs/" + login.Session.Identity.GetId() + "/nodes/" + nodeId)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
			fmt.Println("Response: ", resp.String())
			assert.Contains(tt, resp.String(), nodeId)
		}
	})
}

func TestGetUser(b *testing.T) {
	session, err := kratos.Login(getKratosUrl(), testConf.TestAccountEmail, testConf.TestAccountPass)
	if err != nil {
		assert.NoError(b, err, "Failed to login to Kratos")
		assert.FailNow(b, "Failed to login to Kratos")
		return
	}
	orgId := session.Session.Identity.GetId()

	client := resty.New()
	users := &userspb.ListResponse{}
	_, err = client.R().
		EnableTrace().
		SetHeader("authorization", "bearer "+session.GetSessionToken()).
		SetResult(users).
		Get(getApiUrl() + "/orgs/" + orgId + "/users")

	if assert.NoError(b, err, "Failed to get users") {
		return
	}

	userId := ""
	for i := len(users.Users) - 1; i >= 0; i++ {
		if !users.Users[i].IsDeactivated {
			userId = users.Users[i].Uuid
			break
		}
	}

	var user *userspb.GetResponse
	_, err = client.R().
		EnableTrace().
		SetHeader("authorization", "bearer "+session.GetSessionToken()).
		SetResult(user).
		Get(getApiUrl() + "/orgs/" + orgId + "/users/" + userId)

	if assert.NoError(b, err, "Failed to get user") {
		assert.NotEmpty(b, user.GetUser().Uuid, "User UUID is empty")
	}

}

func getApiUrl() string {
	return testConf.ApiUrl
}

func getKratosUrl() string {
	return testConf.KratosUrl
}
