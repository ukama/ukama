//go:build integration
// +build integration

package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukamaX/common/config"
)

type TestConfig struct {
	config.BaseConfig
	ServiceHost string
}

var tConfig *TestConfig

func init() {
	tConfig = &TestConfig{
		ServiceHost: "http://localhost:8080",
	}

	config.LoadConfig("integration", tConfig)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", tConfig)
}

// Call webhost endpoint and check response
func Test_AddCApp(t *testing.T) {
	rest := resty.New().EnableTrace().SetDebug(tConfig.DebugMode)

	t.Run("Ping", func(t *testing.T) {
		r, err := rest.R().Get(tConfig.ServiceHost + "/ping")
		assert.NoError(t, err)
		assert.Equal(t, r.StatusCode(), http.StatusOK)
	})

	t.Run("Get", func(tt *testing.T) {
		r, err := rest.R().Get(fmt.Sprintf("%s/%s", tConfig.ServiceHost, "some-name"))
		if err != nil {
			assert.FailNow(tt, err.Error())
		}
		assert.Equal(tt, r.StatusCode(), http.StatusOK)
		assert.NoError(tt, err)
	})
}
