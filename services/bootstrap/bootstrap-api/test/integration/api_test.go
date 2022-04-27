package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/iamolegga/enviper"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type TestConfig struct {
	AuthClientId     string
	AuthClientSecret string
	Auth0Endpoint    string
	BootstrapApiUrl  string
	AuthAudience     string
}

type IntegrationTestSuite struct {
	suite.Suite
	config *TestConfig
}

func (t *IntegrationTestSuite) SetupSuite() {
	t.config = loadConfig()
	logrus.Info("Running test against bootstrap url: ", t.config.BootstrapApiUrl)
}

func (i *IntegrationTestSuite) Test_Ping() {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(i.config.BootstrapApiUrl + "/ping")

	i.NoError(err)
	i.True(resp.IsSuccess())
	i.Contains(resp.String(), "pong")

}

func (i *IntegrationTestSuite) Test_BootstrapApi() {

	if len(i.config.AuthClientId) == 0 || len(i.config.AuthClientSecret) == 0 {
		i.FailNow("AuthClientId or AuthClientSecret are not initialized")
	}

	token, err := i.getToken()
	if err != nil {
		i.FailNowf("Error retreiving token.", "Err: %v", err)
	}

	client := resty.New()

	i.Run("GetOrg", func() {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+token).
			SetBody(`{	"certificate":"cert", "ip": "127.0.0.1"	}`).
			Post(i.config.BootstrapApiUrl + "/orgs/org-1")

		logrus.Info("Response: ", resp.String())
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().NotEmpty(resp.String())
	})
}

func (i *IntegrationTestSuite) getToken() (string, error) {
	url := i.config.Auth0Endpoint
	logrus.Info("Client id: ", i.config.AuthClientId)
	requestBody := fmt.Sprintf("{\"client_id\":\"%s\",\"client_secret\":\"%s\",\"audience\":\"%s\",\"grant_type\":\"client_credentials\"}",
		i.config.AuthClientId, i.config.AuthClientSecret, i.config.AuthAudience)
	payload := strings.NewReader(requestBody)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request for token failed. Error code: %d. Body %s", res.StatusCode, string(body))
	}

	respBody := map[string]interface{}{}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return "", err
	}

	return respBody["access_token"].(string), nil

}

func loadConfig() *TestConfig {
	testConf := &TestConfig{
		BootstrapApiUrl: "http://localhost:8080",
		Auth0Endpoint:   "https://ukama.us.auth0.com/oauth/token",
		AuthAudience:    "bootstrap.dev.ukama.com",
	}

	CommonLoadConfig("integration", testConf)
	return testConf
}

func CommonLoadConfig(configFileName string, config interface{}) {

	e := enviper.New(viper.New())
	e.SetConfigType("yaml")

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	e.AddConfigPath(home)
	e.AddConfigPath("")
	e.SetConfigName(configFileName + ".yaml")

	e.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err = e.ReadInConfig()
	if err == nil {
		logrus.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		logrus.Infof("Config file was not loaded. Reason: %v\n", err)
	}

	err = e.Unmarshal(config)
	if err != nil {
		logrus.Fatalf("Unable to decode into struct, %v", err)
	}
}
