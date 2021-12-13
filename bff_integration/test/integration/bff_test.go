package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	graphql "github.com/machinebox/graphql"
	"github.com/stretchr/testify/suite"
)

type TestConfig struct {
	BFFHost string
}

type IntegrationTestSuite struct {
	suite.Suite
	config        *TestConfig
	graphqlClient *graphql.Client
}

func NewIntegrationTestSuite(config *TestConfig) *IntegrationTestSuite {
	return &IntegrationTestSuite{config: config, graphqlClient: graphql.NewClient(config.BFFHost)}
}

func (i *IntegrationTestSuite) Test_Ping() {

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 15,
		Transport: netTransport,
	}
	response, error := netClient.Get("https://bff.dev.ukama.com/ping")
	bodyBytes, _ := ioutil.ReadAll(response.Body)

	fmt.Println("Response of Ping Service: ", string(bodyBytes))
	i.Assert().NoError(error)
	i.Assert().Equal(http.StatusOK, response.StatusCode)
	i.Assert().Equal("pong", string(bodyBytes))

}

func (i *IntegrationTestSuite) Test_GetConnectedUsers() {

	_, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	graphqlRequest := graphql.NewRequest(GetConnectedUsers)

	graphqlRequest.Header.Set("csrf-token", "authorization")
	graphqlRequest.Header.Set("ukama-session", "test")

	var res GetConnectedUsersResponse

	err := i.graphqlClient.Run(context.Background(), graphqlRequest, &res)
	fmt.Println("Response of Test_GetConnectedUsers Query: ", "%+v", res)

	i.Assert().NoError(err)
	i.Assert().GreaterOrEqual(res.ConnectedUser.TotalUsers, 0)

}
