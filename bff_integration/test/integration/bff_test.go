package integration

import (
	"context"
	"fmt"
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

func (i *IntegrationTestSuite) handleResponse(err error, r interface{}) {
	fmt.Printf("Response: %v\n", r)
	i.Assert().NoErrorf(err, "Request failed: %v\n", err)
}

func (i *IntegrationTestSuite) Test_GetConnectedUsers() {

	_, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	graphqlRequest := graphql.NewRequest(`{
		getConnectedUsers(filter:WEEK){
			totalUser
			residentUsers
			guestUsers
		}
	}`)

	graphqlRequest.Header.Set("csrf-token", "authorization")
	graphqlRequest.Header.Set("ukama-session", "test")

	var graphqlResponse interface{}

	err := i.graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse)

	i.handleResponse(err, graphqlResponse)

}
