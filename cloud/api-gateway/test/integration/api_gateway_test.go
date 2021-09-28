// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/go-yaml/yaml"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/ukama/ukamaX/common/config"
	"net/http"

	ory "github.com/ory/kratos-client-go"
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
	b, err := yaml.Marshal(testConf)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: BASEDOMAIN")
	logrus.Infoln(string(b))

	config.LoadConfig("integration", testConf)

	return testConf
}

func (i *IntegrationTestSuite) Test_RegistryApi() {
	login := i.Login()
	client := resty.New()

	PrintJSONPretty(login)

	i.Run("GetOrg", func() {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			Get(i.getApiUrl() + "/orgs/org-1")

		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		i.Assert().NotEmpty(resp.String())
	})

	i.Run("GetNodesUnauthorized", func() {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+"random session").
			Get(i.getApiUrl() + "/nodes")
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusUnauthorized, resp.StatusCode())
	})

	i.Run("GetNodes", func() {
		resp, err := client.R().
			EnableTrace().
			SetHeader("authorization", "bearer "+login.GetSessionToken()).
			Get(i.getApiUrl() + "/nodes")
		i.Assert().NoError(err)
		i.Assert().Equal(http.StatusOK, resp.StatusCode())
		fmt.Println("Response: ", resp.String())
	})
}

func (i *IntegrationTestSuite) getApiUrl() string {
	return "https://api." + i.config.BaseDomain
}

func (i *IntegrationTestSuite) getKratosUrl() string {
	return "https://" + i.config.BaseDomain + "/.ory/kratos/public"
}

func (i *IntegrationTestSuite) Login() *ory.SuccessfulSelfServiceLoginWithoutBrowser {
	var client = NewSDKForSelfHosted(i.getKratosUrl())

	ctx := context.Background()

	// Create a temporary user
	email, password := RandomCredentials()
	_, _ = CreateIdentityWithSession(client, email, password)

	// Initialize the flow
	flow, res, err := client.V0alpha1Api.InitializeSelfServiceLoginFlowWithoutBrowser(ctx).Execute()
	SDKExitOnError(err, res)

	// If you want, print the flow here:
	//
	PrintJSONPretty(flow)

	// Submit the form
	result, res, err := client.V0alpha1Api.SubmitSelfServiceLoginFlow(ctx).Flow(flow.Id).SubmitSelfServiceLoginFlowBody(
		ory.SubmitSelfServiceLoginFlowWithPasswordMethodBodyAsSubmitSelfServiceLoginFlowBody(&ory.SubmitSelfServiceLoginFlowWithPasswordMethodBody{
			Method:             "password",
			Password:           password,
			PasswordIdentifier: email,
		}),
	).Execute()
	SDKExitOnError(err, res)

	return result
}
