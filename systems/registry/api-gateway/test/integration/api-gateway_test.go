//go:build integration
// +build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/ukama"
	userspb "github.com/ukama/ukama/systems/registry/users/pb/gen"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Before running test for the first time you have to create a test account in Identity manager and provide email and password for it

type TestConfig struct {
	ServiceHost string `default:"localhost:8080"`
}

var testConf *TestConfig

func init() {
	testConf = &TestConfig{}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
	config.LoadConfig("integration", testConf)
	logrus.Infof("Config: %+v", testConf)
}

func Test_RegistryApi(t *testing.T) {
	org := "org-name"
	nodeId := ukama.NewVirtualHomeNodeId().String()

	client := resty.New()

	t.Run("GetOrg", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/orgs/" + org)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
		}

	})

	t.Run("GetOrgNotFound", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/orgs/" + "someRandomOrgThatShouldNotExist")

		assert.NoError(t, err)
		assert.Equal(tt, http.StatusNotFound, resp.StatusCode())
		assert.Contains(tt, resp.String(), "Organization not found")
	})

	t.Run("GetNodes", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/orgs/" + org + "/nodes")
		assert.NoError(t, err)

		assert.Equal(tt, http.StatusOK, resp.StatusCode())
		fmt.Println("Response: ", resp.String())
	})

	t.Run("AddNode", func(tt *testing.T) {
		nodeName := time.Now().Format(time.RFC3339) + "testNode"
		resp, err := client.R().
			EnableTrace().
			SetBody(map[string]string{"name": nodeName}).
			Put(getApiUrl() + "/v1/orgs/" + org + "/nodes/" + nodeId)

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
			SetBody(map[string]string{"name": nodeName}).
			Put(getApiUrl() + "/v1/orgs/" + org + "/nodes/" + nodeId)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
			fmt.Println("Response: ", resp.String())
			assert.Contains(tt, resp.String(), nodeName)
		}
	})

	t.Run("DeleteNode", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Delete(getApiUrl() + "/v1/orgs/" + org + "/nodes/" + nodeId)

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusOK, resp.StatusCode())
			fmt.Println("Response: ", resp.String())
			assert.Contains(tt, resp.String(), nodeId)
		}
	})
}

func TestGetUser(b *testing.T) {
	org := "org-name"
	userId := ""

	client := resty.New()
	users := &userspb.ListResponse{}
	_, err := client.R().
		EnableTrace().
		SetResult(users).
		Get(getApiUrl() + "/v1/orgs/" + org + "/users")

	if assert.NoError(b, err, "Failed to get users") {
		return
	}

	for i := len(users.Users) - 1; i >= 0; i++ {
		if !users.Users[i].IsDeactivated {
			userId = users.Users[i].Uuid
			break
		}
	}

	var user *userspb.GetResponse
	_, err = client.R().
		EnableTrace().
		SetResult(user).
		Get(getApiUrl() + "/v1/orgs/" + org + "/users/" + userId)

	if assert.NoError(b, err, "Failed to get user") {
		assert.NotEmpty(b, user.GetUser().Uuid, "User UUID is empty")
	}

}

func getApiUrl() string {
	return testConf.ServiceHost
}
