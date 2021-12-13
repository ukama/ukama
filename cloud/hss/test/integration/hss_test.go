//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	i.Run("AddImis", func() {
		addResp, err := c.Add(ctx, &pb.AddImsiRequest{Org: testOrg, Imsi: &pb.ImsiRecord{Imsi: testImsi, Apn: &pb.Apn{Name: "test-apn-name"}}})
		i.handleResponse(err, addResp)
	})

	i.Run("GetImis", func() {
		getResp, err := c.Get(ctx, &pb.GetImsiRequest{Imsi: testImsi})
		i.handleResponse(err, getResp)
	})

	i.Run("AddGuti", func() {
		delResp, err := c.AddGuti(ctx, &pb.AddGutiRequest{Guti: &pb.Guti{
			PlmnId: "000001",
			Mmegi:  1,
			Mmec:   1,
			Mtmsi:  uint32(time.Now().Unix()),
		}, Imsi: testImsi,
			UpdatedAt: uint32(time.Now().Unix())})
		i.handleResponse(err, delResp)
	})

	i.Run("UpdateGutiAddedEarlier", func() {
		delResp, err := c.AddGuti(ctx, &pb.AddGutiRequest{Guti: &pb.Guti{
			PlmnId: "000001",
			Mmegi:  1,
			Mmec:   1,
			Mtmsi:  uint32(time.Now().Unix()) + 1,
		}, Imsi: testImsi,
			UpdatedAt: uint32(time.Now().Unix())})
		i.handleResponse(err, delResp)
	})

	i.Run("AddTai", func() {
		resp, err := c.UpdateTai(ctx, &pb.UpdateTaiRequest{Imsi: testImsi, Tac: 4654, PlmnId: "000001",
			UpdatedAt: uint32(time.Now().Unix())})
		i.handleResponse(err, resp)
	})

	i.Run("UpdateTaiAddedEarlier", func() {
		resp, err := c.UpdateTai(ctx, &pb.UpdateTaiRequest{Imsi: testImsi, Tac: 4654, PlmnId: "000001",
			UpdatedAt: uint32(time.Now().Unix())})
		i.handleResponse(err, resp)
	})

	i.Run("DeleteImis", func() {
		delResp, err := c.Delete(ctx, &pb.DeleteImsiRequest{Imsi: testImsi, Org: testOrg})
		i.handleResponse(err, delResp)
	})

	i.Run("UpdateTaiValidationFailure", func() {
		_, err := c.UpdateTai(ctx, &pb.UpdateTaiRequest{Imsi: testImsi, Tac: 4654, PlmnId: "000001",
			UpdatedAt: uint32(time.Now().Unix())})
		s, ok := status.FromError(err)
		i.True(ok, "should be a grpc error")
		i.Equal(codes.NotFound, s.Code(), "should fail with invalid argument")
	})
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
