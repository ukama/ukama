// +build integration

package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pb "github.com/ukama/ukamaX/cloud/hss/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	HssHost string
}

type IntegrationTestSuite struct {
	suite.Suite
	config *TestConfig
}

func NewIntegrationTestSuite(config *TestConfig) *IntegrationTestSuite {
	return &IntegrationTestSuite{config: config}
}

func (i *IntegrationTestSuite) Test_ImsiService() {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", i.config.HssHost)
	conn, err := grpc.DialContext(ctx, i.config.HssHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(i.T(), err, "did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewImsiServiceClient(conn)

	// Contact the server and print out its response.
	testImsi := fmt.Sprintf("00000%010d", time.Now().Unix())
	testOrg := fmt.Sprintf("integration-test-org-imsi-service-%s", time.Now().Format("20060102150405"))
	addResp, err := c.Add(ctx, &pb.AddImsiRequest{Org: testOrg, Imsi: &pb.ImsiRecord{Imsi: testImsi, DefaultApnName: "apn-name"}})
	i.handleResponse(err, addResp)

	getResp, err := c.Get(ctx, &pb.GetImsiRequest{Imsi: testImsi})
	i.handleResponse(err, getResp)

	delResp, err := c.Delete(ctx, &pb.DeleteImsiRequest{Imsi: testImsi, Org: testOrg})
	i.handleResponse(err, delResp)
}

func (i *IntegrationTestSuite) Test_UserService() {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", i.config.HssHost)
	conn, err := grpc.DialContext(ctx, i.config.HssHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(i.T(), err, "did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)
	testOrg := fmt.Sprintf("integration-test-org-user-service-%s", time.Now().Format("20060102150405"))

	var addResp *pb.AddUserResponse

	i.Run("Add", func() {
		testImsi := fmt.Sprintf("00001%010d", time.Now().Unix())

		addResp, err = c.Add(ctx, &pb.AddUserRequest{
			User: &pb.User{
				Email:     "test@example.com",
				Imsi:      testImsi,
				LastName:  "Joe",
				FirstName: "Doe",
			},
			Org: testOrg})

		i.handleResponse(err, addResp)
		i.NotEmpty(addResp.User.Uuid)
	})
	i.Run("list", func() {

		listResp, err := c.List(ctx, &pb.ListUsersRequest{
			Org: testOrg,
		})

		i.handleResponse(err, addResp)
		i.Equal(1, len(listResp.Users))
	})

	i.Run("Delete", func() {
		getResp, err := c.Delete(ctx, &pb.DeleteUserRequest{UserUuid: addResp.User.Uuid, Org: testOrg})
		i.handleResponse(err, getResp)
	})
}

func (i *IntegrationTestSuite) handleResponse(err error, r interface{}) {
	fmt.Printf("Response: %v\n", r)
	i.Assert().NoErrorf(err, "Request failed: %v\n", err)
}
