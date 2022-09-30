//go:build integration
// +build integration

package integration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/ukama"
	"net/http"
	"testing"

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

func Test_LookupClientApi(t *testing.T) {
	nodeId := ukama.NewVirtualHomeNodeId().String()

	client := resty.New()

	t.Run("GetNodeNotFound", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/nodes/someNodeWhichDoesNotExist")

		if assert.NoError(t, err) {
			assert.Equal(tt, http.StatusNotFound, resp.StatusCode())
			assert.Contains(tt, resp.String(), "node record not found")
		}
	})

	t.Run("GetNode", func(tt *testing.T) {
		resp, err := client.R().
			EnableTrace().
			Get(getApiUrl() + "/v1/nodes/" + nodeId)
		assert.NoError(t, err)
		assert.Equal(tt, http.StatusOK, resp.StatusCode())
		fmt.Println("Response: ", resp.String())
	})
}

func getApiUrl() string {
	return testConf.ServiceHost
}
