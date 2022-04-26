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
	DistributionHost string
}

var (
	cappname    = "ukamaos"
	cappversion = "1.0.1"
)

var tConfig *TestConfig

func init() {
	tConfig = &TestConfig{
		DistributionHost: "http://localhost:8098",
	}

	config.LoadConfig("integration", tConfig)
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	logrus.Infof("%+v", tConfig)
}

// Call webhost endpoint and check response
func Test_PutChunks(t *testing.T) {
	appUrl := fmt.Sprintf("%s/chunk", tConfig.DistributionHost)

	rest := resty.New().EnableTrace().SetDebug(tConfig.DebugMode)

	t.Run("Ping", func(t *testing.T) {
		r, err := rest.R().Get(tConfig.DistributionHost + "/ping")
		assert.NoError(t, err)
		assert.Equal(t, r.StatusCode(), http.StatusOK)
	})

	t.Run("Put", func(tt *testing.T) {
		r, err := rest.R().SetBody(map[string]interface{}{
			"store": "testdata/art"}).Put(appUrl + "/" + cappname + "/" + cappversion)

		assert.NoError(tt, err)
		logrus.Infof("Response: '%d'", r.StatusCode())

		assert.Equal(tt, r.StatusCode(), http.StatusOK)

	})

}
