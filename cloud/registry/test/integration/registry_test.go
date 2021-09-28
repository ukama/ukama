// +build integration

package integration

import (
	"context"
	"fmt"
	uuid2 "github.com/satori/go.uuid"
	"github.com/ukama/ukamaX/common/ukama"
	"time"

	"github.com/ukama/ukamaX/cloud/registry/pb/gen/external"
	"github.com/ukama/ukamaX/common/msgbus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"google.golang.org/grpc"
)

type TestConfig struct {
	RegistryHost string
	Rabbitmq     string
}

type IntegrationTestSuite struct {
	suite.Suite
	config  *TestConfig
	orgName string
}

func NewIntegrationTestSuite(config *TestConfig) *IntegrationTestSuite {
	return &IntegrationTestSuite{config: config,
		orgName: fmt.Sprintf("registry-integration-self-test-org-%d", time.Now().Unix())}
}

func (i *IntegrationTestSuite) Test_FullFlow() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	logrus.Infoln("Connecting to registry ", i.config.RegistryHost)
	conn, err := grpc.DialContext(ctx, i.config.RegistryHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(i.T(), err, "did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewRegistryServiceClient(conn)

	// Contact the server and print out its response.
	ownerId := uuid2.NewV1()
	node := ukama.NewVirtualNodeId("HomeNode")

	var r interface{}

	r, err = c.AddOrg(ctx, &pb.AddOrgRequest{Name: i.orgName, Owner: ownerId.String()})
	i.handleResponse(err, r)

	r, err = c.GetOrg(ctx, &pb.GetOrgRequest{Name: i.orgName})
	i.handleResponse(err, r)

	r, err = c.AddNode(ctx, &pb.AddNodeRequest{
		Node: &pb.Node{
			NodeId: node.String(),
			State:  pb.NodeState_UNDEFINED,
		},
		OrgName: i.orgName,
	})
	i.handleResponse(err, r)

	r, err = c.UpdateNode(ctx, &pb.UpdateNodeRequest{NodeId: node.String(), State: pb.NodeState_ONBOARDED})
	i.handleResponse(err, r)

	nodeResp, err := c.GetNode(ctx, &pb.GetNodeRequest{NodeId: node.String()})
	i.handleResponse(err, nodeResp)
	i.Assert().Equal(pb.NodeState_ONBOARDED, nodeResp.Node.State)

	nodesResp, err := c.GetNodes(ctx, &pb.GetNodesRequest{Owner: ownerId.String()})
	i.handleResponse(err, nodesResp)
	i.Assert().Equal(1, len(nodesResp.Orgs[0].Nodes))
	i.Assert().Equal(node.String(), nodesResp.Orgs[0].Nodes[0].NodeId)
}

func getContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	return ctx
}

func (i *IntegrationTestSuite) Test_Listener() {
	// Arrange
	org := "registry-listener-integration-test-org"
	nodeId := "UK-SA2136-HNODE-A1-30DF"
	ownerId := "474bd2c3-77a0-49f2-a143-dd840dce2c91"

	conn, c, err := i.CreateRegistryClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(i.T(), err, "did not connect: %+v\n", err)
		return
	}

	_, err = c.AddOrg(getContext(), &pb.AddOrgRequest{Name: org, Owner: ownerId})
	e, ok := status.FromError(err)
	if ok && e.Code() == codes.AlreadyExists {
		logrus.Infof("org already exist, err %+v\n", err)
	} else {
		assert.NoError(i.T(), err)
		return
	}

	_, err = c.AddNode(getContext(), &pb.AddNodeRequest{Node: &pb.Node{
		NodeId: nodeId, State: pb.NodeState_UNDEFINED,
	}, OrgName: org})
	e, ok = status.FromError(err)
	if ok && e.Code() == codes.AlreadyExists {
		logrus.Infof("node already exist, err %+v\n", err)
	} else {
		assert.NoError(i.T(), err)
		return
	}

	_, err = c.UpdateNode(getContext(), &pb.UpdateNodeRequest{NodeId: nodeId, State: pb.NodeState_PENDING})
	assert.NoError(i.T(), err)

	// Act
	err = i.sendMessageToQueue(nodeId)
	assert.NoError(i.T(), err)

	// Assert
	time.Sleep(3 * time.Second)
	nodeResp, err := c.GetNode(getContext(), &pb.GetNodeRequest{NodeId: nodeId})
	i.Assert().NoError(err)
	i.Assert().Equal(pb.NodeState_ONBOARDED, nodeResp.Node.State)
}

func (i *IntegrationTestSuite) CreateRegistryClient() (*grpc.ClientConn, pb.RegistryServiceClient, error) {
	logrus.Infoln("Connecting to registry ", i.config.RegistryHost)
	conn, err := grpc.DialContext(context.Background(), i.config.RegistryHost, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewRegistryServiceClient(conn)
	return conn, c, nil
}

func (i *IntegrationTestSuite) sendMessageToQueue(nodeId string) error {
	rabbit, err := msgbus.NewPublisherClient(i.config.Rabbitmq)
	i.Assert().NoError(err)
	if err != nil {
		return err
	}

	message, err := proto.Marshal(&external.Link{Uuid: &nodeId})
	i.Assert().NoError(err)
	err = rabbit.Publish(message, "", msgbus.DeviceQ.Exchange, msgbus.DeviceConnectedRoutingKey, "direct")
	i.Assert().NoError(err)
	return err
}

func (i *IntegrationTestSuite) handleResponse(err error, r interface{}) {
	fmt.Printf("Response: %v\n", r)
	i.Assert().NoErrorf(err, "Request failed: %v\n", err)
}
