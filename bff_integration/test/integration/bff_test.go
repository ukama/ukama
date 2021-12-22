package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	graphql "github.com/machinebox/graphql"
	"github.com/stretchr/testify/suite"
	"github.com/ukama/ukamaX/common/testing/kratos"
)

type TestConfig struct {
	BaseDomain string
}

type IntegrationTestSuite struct {
	suite.Suite
	config        *TestConfig
	graphqlClient *graphql.Client
}

func NewIntegrationTestSuite(config *TestConfig) *IntegrationTestSuite {
	return &IntegrationTestSuite{config: config, graphqlClient: graphql.NewClient("https://bff." + config.BaseDomain + "/graphql")}
}

func (i *IntegrationTestSuite) TestPing() {
	var netClient = &http.Client{}
	response, error := netClient.Get("https://bff." + i.config.BaseDomain + "/ping")
	bodyBytes, _ := ioutil.ReadAll(response.Body)

	i.Assert().NoError(error)
	i.Assert().Equal(http.StatusOK, response.StatusCode)
	i.Assert().Equal("pong", string(bodyBytes))
}

func (i *IntegrationTestSuite) TestGetConnectedUsers() {
	login, error := kratos.Login("https://auth." + i.config.BaseDomain + "/.api")
	i.Assert().NoError(error)

	graphqlRequest := graphql.NewRequest(GetConnectedUsers)
	graphqlRequest.Header.Set("authorization", "Bearer "+login.GetSessionToken())

	var res GetConnectedUsersResponse
	err := i.graphqlClient.Run(context.Background(), graphqlRequest, &res)

	fmt.Println("TestGetConnectedUsers Response: ", res.ConnectedUser.TotalUsers)
	i.Assert().NoError(err)
	i.Assert().GreaterOrEqual(res.ConnectedUser.TotalUsers, 0)
}

func (i *IntegrationTestSuite) TestGetNodesByOrg() {
	login, error := kratos.Login("https://auth." + i.config.BaseDomain + "/.api")
	i.Assert().NoError(error)

	graphqlRequest := graphql.NewRequest(fmt.Sprintf(GetNodesByOrg, login.Session.Identity.GetId()))
	graphqlRequest.Header.Set("authorization", "Bearer "+login.GetSessionToken())

	var res GetNodesByOrgResponse
	err := i.graphqlClient.Run(context.Background(), graphqlRequest, &res)

	fmt.Println("TestGetNodesByOrg Response: ", res.GetNodesByOrg.OrgName, res.GetNodesByOrg.TotalNodes)
	i.Assert().NoError(err)
}
