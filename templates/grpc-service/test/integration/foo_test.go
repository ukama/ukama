// +build integration

package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pb "github.com/ukama/ukamaX/foo/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	FooHost string
}

type IntegrationTestSuite struct {
	suite.Suite
	config *TestConfig
}

func NewIntegrationTestSuite(config *TestConfig) *IntegrationTestSuite {
	return &IntegrationTestSuite{config: config}
}

func (i *IntegrationTestSuite) Test_FullFlow() {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", i.config.FooHost)
	conn, err := grpc.DialContext(ctx, i.config.FooHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(i.T(), err, "did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewRegistryServiceClient(conn)

	// Contact the server and print out its response.

	r, err = c.GetFoo(ctx, &pb.GetFooRequest{Name: "some_name"})
	i.handleResponse(err, r)

	// Asserts goes here
}

func getContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	return ctx
}

func (i *IntegrationTestSuite) handleResponse(err error, r interface{}) {
	fmt.Printf("Response: %v\n", r)
	i.Assert().NoErrorf(err, "Request failed: %v\n", err)
}
