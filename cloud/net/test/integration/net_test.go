//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	pb "github.com/ukama/ukamaX/cloud/net/pb/gen"
	"github.com/ukama/ukamaX/common/ukama"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	//pb "github.com/ukama/ukamaX/cloud/net/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	NetHost string
}

type IntegrationTestSuite struct {
	suite.Suite
	config *TestConfig
}

func NewIntegrationTestSuite(config *TestConfig) *IntegrationTestSuite {
	return &IntegrationTestSuite{config: config}
}

func (is *IntegrationTestSuite) Test_FullFlow() {
	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to service ", is.config.NetHost)
	conn, err := grpc.DialContext(ctx, is.config.NetHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(is.T(), err, "did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewNnsClient(conn)
	nodeId := ukama.NewVirtualHomeNodeId().String()
	const ip = "1.1.1.1"

	is.Run("SetIp", func() {
		r, err := c.Set(ctx, &pb.SetRequest{NodeId: nodeId, Ip: ip})
		is.handleResponse(err, r)
	})

	is.Run("ResolevIp", func() {
		r, err := c.Get(ctx, &pb.GetRequest{NodeId: nodeId})
		is.handleResponse(err, r)
		is.Equal(ip, r.Ip)
	})

	is.Run("ResolevMissinIp", func() {
		_, err := c.Get(ctx, &pb.GetRequest{NodeId: ukama.NewVirtualHomeNodeId().String()})
		s, ok := status.FromError(err)
		is.True(ok)
		is.Equal(codes.NotFound, s.Code())
	})

	is.Run("GetIpList", func() {
		_, err := c.Set(ctx, &pb.SetRequest{NodeId: ukama.NewVirtualHomeNodeId().String(), Ip: ip})
		is.handleResponse(err, nil)
		r, err := c.List(ctx, &pb.ListRequest{})
		is.NoError(err)
		// just make sure it's unique list
		is.Greater(len(r.Ips), 2)
		un := make(map[string]bool)
		for _, i := range r.Ips {
			if _, ok := un[i]; ok {
				is.Fail("Duplicate ip")
			}
			un[i] = true
		}
	})

	is.Run("Delete", func() {
		_, err := c.Delete(ctx, &pb.DeleteRequest{NodeId: nodeId})
		is.NoError(err)
		_, err = c.Get(ctx, &pb.GetRequest{NodeId: nodeId})
		e, ok := status.FromError(err)
		is.True(ok)
		is.Equal(codes.NotFound, e.Code())
	})
}

func getContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	return ctx
}

func (is *IntegrationTestSuite) handleResponse(err error, r interface{}) {
	fmt.Printf("Response: %v\n", r)
	is.Assert().NoErrorf(err, "Request failed: %v\n", err)
	if err != nil {
		is.FailNow("Unexpected response")
	}
}
