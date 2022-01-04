//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"github.com/ukama/ukamaX/common/ukama"
	"google.golang.org/grpc"
	"time"
)

func (i *IntegrationTestSuite) Test_NetworkingService() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to registry ", i.config.RegistryHost)
	conn, err := grpc.DialContext(ctx, i.config.RegistryHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(i.T(), err, "did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewNetworkingClient(conn)

	// Contact the server and print out its response.
	node := ukama.NewVirtualNodeId("HomeNode")

	var r interface{}

	i.Run("SetNodeIp", func() {
		r, err = c.SetNodeIp(ctx, &pb.SetNodeIpRequest{
			NodeId: node.String(),
			Ip:     "8.8.8.8"})
		i.handleResponse(err, r)
	})

	i.Run("GetNodeIp", func() {
		r, err = c.ResolveNodeIp(ctx, &pb.ResolveNodeIpRequest{
			NodeId: node.String(),
		})
		i.handleResponse(err, r)
	})
}
