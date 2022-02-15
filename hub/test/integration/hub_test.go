//go:build integration
// +build integration

package integration

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukamaX/common/config"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

type TestConfig struct {
	config.BaseConfig
	HubHost string
}

var tConfig *TestConfig

func init() {
	tConfig = &TestConfig{
		HubHost: "http://localhost:8080",
	}

	config.LoadConfig("integration", tConfig)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", tConfig)
}

// Call webhost endpoint and check response
func Test_AddCApp(t *testing.T) {
	appUrl := fmt.Sprintf("%s/capps/hub-integration-test/1.0.%d", tConfig.HubHost, time.Now().Unix())
	con := getFileContent(t)
	rest := resty.New().EnableTrace().SetDebug(tConfig.DebugMode)

	t.Run("Ping", func(t *testing.T) {
		r, err := rest.R().Get(tConfig.HubHost + "/ping")
		assert.NoError(t, err)
		assert.Equal(t, r.StatusCode(), http.StatusOK)
	})

	t.Run("Put", func(tt *testing.T) {
		r, err := rest.R().SetHeader("Content-Type", "application/octet-stream").SetBody(con).Put(appUrl)

		assert.NoError(tt, err)
		logrus.Infof("Response: '%s'", r.String())
		assert.Equal(tt, http.StatusCreated, r.StatusCode())
	})

	t.Run("Get", func(tt *testing.T) {
		r, err := rest.R().Get(appUrl)
		if err != nil {
			assert.FailNow(tt, err.Error())
		}
		assert.Equal(tt, r.StatusCode(), http.StatusOK)
		assert.NoError(tt, err)
		assert.Equal(tt, con, r.Body(), "Expected file content is not equal to actual content")
	})
}

func getFileContent(t *testing.T) []byte {
	f, err := os.Open("testdata/capp.tar.gz")
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer f.Close()

	con, err := ioutil.ReadAll(f)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	return con
}
