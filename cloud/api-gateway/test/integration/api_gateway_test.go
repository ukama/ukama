//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	ory "github.com/ory/kratos-client-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/ukama/ukamaX/common/config"
	"net/http"
	"time"
)

type TestConfig struct {
	BaseDomain string
}

type IntegrationTestSuite struct {
	suite.Suite
	config *TestConfig
}

func (t *IntegrationTestSuite) SetupSuite() {
	t.config = loadConfig()
}

func loadConfig() *TestConfig {
	testConf := &TestConfig{
		BaseDomain: "dev.ukama.com",
	}

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")

	config.LoadConfig("integration", testConf)
	logrus.Infof("Config: %+v", testConf)

	return testConf
}

func (i *IntegrationTestSuite) Test_RegistryApi() {
	login, err := i.Login()
	time.Sleep(3 * time.Second) //give registry some time to create a default org for account
	if err != nil {
		i.NoError(err, "Failed to login to Kratos")
		return
	}

	client := resty.New()

	PrintJSONPretty(login)

	i.Run("GetOrg", func() {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			Get(i.getApiUrl() + "/orgs/" + login.Session.Identity.GetId())

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		//i.Assert().Contains(resp.String(), "Organization not found")
	})

	i.Run("GetOrgNotFound", func() {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			Get(i.getApiUrl() + "/orgs/" + "someRandomOrgThatShouldNotExist")

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusNotFound, resp.StatusCode())
		i.Assert().Contains(resp.String(), "Organization not found")
	})

	i.Run("GetNodesUnauthorized", func() {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+"random session").
			Get(i.getApiUrl() + "/orgs/" + login.Session.Identity.GetId() + "/nodes")
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusUnauthorized, resp.StatusCode())
	})

	i.Run("GetNodes", func() {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			Get(i.getApiUrl() + "/orgs/" + login.Session.Identity.GetId() + "/nodes")
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		fmt.Println("Response: ", resp.String())
	})
}

func (i *IntegrationTestSuite) getApiUrl() string {
	return "https://api." + i.config.BaseDomain
}

func (i *IntegrationTestSuite) getKratosUrl() string {
	return "https://auth." + i.config.BaseDomain + "/.api"
}

func (i *IntegrationTestSuite) Login() (*ory.SuccessfulSelfServiceLoginWithoutBrowser, error) {
	var client = NewSDKForSelfHosted(i.getKratosUrl())

	ctx := context.Background()

	// Create a temporary user
	email, password := RandomCredentials()
	_, _, err := CreateIdentityWithSession(client, email, password)

	// Initialize the flow
	flow, res, err := client.V0alpha2Api.InitializeSelfServiceLoginFlowWithoutBrowser(ctx).Execute()
	LogKratosSdkError(err, res)
	if err != nil {
		return nil, err
	}

	// If you want, print the flow here:
	//
	PrintJSONPretty(flow)

	// Submit the form
	result, res, err := client.V0alpha2Api.SubmitSelfServiceLoginFlow(ctx).Flow(flow.Id).SubmitSelfServiceLoginFlowBody(
		ory.SubmitSelfServiceLoginFlowWithPasswordMethodBodyAsSubmitSelfServiceLoginFlowBody(&ory.SubmitSelfServiceLoginFlowWithPasswordMethodBody{
			Method:             "password",
			Password:           password,
			PasswordIdentifier: email,
		}),
	).Execute()
	LogKratosSdkError(err, res)
	if err != nil {
		return nil, err
	}

	return result, nil
}
