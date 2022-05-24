//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/ukama/ukama/services/common/config"
	"testing"
	"time"

	uuid2 "github.com/google/uuid"
	"github.com/ukama/ukama/services/common/ukama"

	"github.com/ukama/ukama/services/common/msgbus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	pb "github.com/ukama/ukama/services/cloud/registry/pb/gen"
	commonpb "github.com/ukama/ukama/services/common/pb/gen/ukamaos/mesh"
	"google.golang.org/grpc"
)

var tConfig *TestConfig
var orgName string

func init() {
	// set org name
	orgName = fmt.Sprintf("registry-integration-self-test-org-%d", time.Now().Unix())

	// load config
	tConfig = &TestConfig{
		RegistryHost: "localhost:9090",
		Rabbitmq:     "amqp://guest:guest@localhost:5672",
	}

	config.LoadConfig("integration", tConfig)

	logrus.Info("Expected config ", "integration.yaml", " or env vars for ex: REGISTRYHOST")
	logrus.Infof("Config: %+v\n", tConfig)
}

type TestConfig struct {
	RegistryHost string
	Rabbitmq     string
}

func Test_FullFlow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	logrus.Infoln("Connecting to registry ", tConfig.RegistryHost)
	conn, err := grpc.DialContext(ctx, tConfig.RegistryHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)
		return
	}

	c := pb.NewRegistryServiceClient(conn)

	// Contact the server and print out its response.
	ownerId := uuid2.NewString()
	node := ukama.NewVirtualNodeId("HomeNode")

	var r interface{}

	t.Run("AddOrg", func(tt *testing.T) {
		r, err = c.AddOrg(ctx, &pb.AddOrgRequest{Name: orgName, Owner: ownerId})
		handleResponse(tt, err, r)
	})

	t.Run("GetOrg", func(tt *testing.T) {
		r, err = c.GetOrg(ctx, &pb.GetOrgRequest{Name: orgName})
		handleResponse(tt, err, r)
	})

	t.Run("AddAndUpdateNode", func(tt *testing.T) {
		nodeName := "HomeNodeX"
		addResp, err := c.AddNode(ctx, &pb.AddNodeRequest{
			Node: &pb.Node{
				NodeId: node.String(),
				State:  pb.NodeState_UNDEFINED,
				Name:   nodeName,
			},
			OrgName: orgName,
		})
		handleResponse(tt, err, addResp)
		assert.NotNil(tt, addResp.Node)
		assert.Equal(tt, nodeName, addResp.Node.Name)

		r, err = c.UpdateNodeState(ctx, &pb.UpdateNodeStateRequest{NodeId: node.String(), State: pb.NodeState_ONBOARDED})
		handleResponse(tt, err, r)

		nodeResp, err := c.GetNode(ctx, &pb.GetNodeRequest{NodeId: node.String()})
		handleResponse(tt, err, nodeResp)
		assert.Equal(tt, pb.NodeState_ONBOARDED, nodeResp.Node.State)
		assert.Equal(tt, pb.NodeType_HOME, nodeResp.Node.Type)
		assert.NotNil(tt, nodeResp.Network)
	})

	t.Run("AddTowerNodeWithAmplifiers", func(tt *testing.T) {
		tNodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)
		_, err := c.AddNode(ctx, &pb.AddNodeRequest{
			Node: &pb.Node{
				NodeId: tNodeId.String(),
				State:  pb.NodeState_UNDEFINED,
			},
			OrgName: orgName,
		})
		if err != nil {
			assert.FailNow(tt, "AddNode failed", err.Error())
		}

		aNodeId := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_AMPNODE)
		_, err = c.AddNode(ctx, &pb.AddNodeRequest{
			Node: &pb.Node{
				NodeId: aNodeId.String(),
				State:  pb.NodeState_UNDEFINED,
			},
			OrgName: orgName,
		})
		if err != nil {
			assert.FailNow(tt, "AddNode failed", err.Error())
		}

		_, err = c.AttachNodes(ctx, &pb.AttachNodesRequest{
			ParentNodeId:    tNodeId.String(),
			AttachedNodeIds: []string{aNodeId.String()},
		})
		if err != nil {
			assert.FailNow(tt, "AttachNodes failed", err.Error())
		}

		resp, err := c.GetNode(ctx, &pb.GetNodeRequest{NodeId: tNodeId.String()})
		if assert.NoError(tt, err, "GetNode failed") {
			assert.NotNil(tt, resp.Node.Attached)
			assert.Equal(tt, 1, len(resp.Node.Attached))
			assert.Equal(tt, aNodeId.StringLowercase(), resp.Node.Attached[0].NodeId)
		}
	})

	t.Run("DeleteNode", func(tt *testing.T) {
		r, err = c.DeleteNode(ctx, &pb.DeleteNodeRequest{NodeId: node.String()})
		handleResponse(tt, err, r)

		r, err = c.AddNode(ctx, &pb.AddNodeRequest{
			Node: &pb.Node{
				NodeId: node.String(),
				State:  pb.NodeState_UNDEFINED,
			},
			OrgName: orgName,
		})
		handleResponse(tt, err, r)
	})

	t.Run("GetNodes", func(tt *testing.T) {
		nodesResp, err := c.GetNodes(ctx, &pb.GetNodesRequest{OrgName: orgName})
		handleResponse(t, err, nodesResp)
		if assert.Equal(tt, 3, len(nodesResp.Nodes)) {
			cont := false
			for _, n := range nodesResp.Nodes {
				if node.String() == n.NodeId {
					cont = true
					break
				}
			}
			assert.True(tt, cont, "Can't find added node %s", node.String())
		}

	})
}

func Test_Listener(t *testing.T) {
	// Arrange
	org := "registry-listener-integration-test-org"
	nodeId := "UK-SA2136-HNODE-A1-30DF"
	ownerId := "474bd2c3-77a0-49f2-a143-dd840dce2c91"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, c, err := CreateRegistryClient()
	defer conn.Close()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)
		return
	}

	_, err = c.AddOrg(ctx, &pb.AddOrgRequest{Name: org, Owner: ownerId})
	e, ok := status.FromError(err)
	if ok && e.Code() == codes.AlreadyExists {
		logrus.Infof("org already exist, err %+v\n", err)
	} else {
		assert.NoError(t, err)
		return
	}

	_, err = c.AddNode(ctx, &pb.AddNodeRequest{Node: &pb.Node{
		NodeId: nodeId, State: pb.NodeState_UNDEFINED,
	}, OrgName: org})
	e, ok = status.FromError(err)
	if ok && e.Code() == codes.AlreadyExists {
		logrus.Infof("node already exist, err %+v\n", err)
	} else {
		assert.NoError(t, err)
		return
	}

	_, err = c.UpdateNodeState(ctx, &pb.UpdateNodeStateRequest{NodeId: nodeId, State: pb.NodeState_PENDING})
	if err != nil {
		assert.FailNow(t, "Failed to update node. Error: %s", err.Error())
	}

	// Act
	err = sendMessageToQueue(nodeId)

	// Assert
	assert.NoError(t, err)
	time.Sleep(3 * time.Second)
	nodeResp, err := c.GetNode(ctx, &pb.GetNodeRequest{NodeId: nodeId})
	assert.NoError(t, err)
	if err != nil {
		assert.Equal(t, pb.NodeState_ONBOARDED, nodeResp.Node.State)
	}
}

func CreateRegistryClient() (*grpc.ClientConn, pb.RegistryServiceClient, error) {
	logrus.Infoln("Connecting to registry ", tConfig.RegistryHost)
	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(context, tConfig.RegistryHost, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewRegistryServiceClient(conn)
	return conn, c, nil
}

func sendMessageToQueue(nodeId string) error {
	rabbit, err := msgbus.NewPublisherClient(tConfig.Rabbitmq)

	if err != nil {
		return err
	}

	message, err := proto.Marshal(&commonpb.Link{NodeId: &nodeId, Ip: proto.String("1.1.1.1")})
	if err != nil {
		return err
	}
	err = rabbit.Publish(message, "", msgbus.DeviceQ.Exchange, msgbus.DeviceConnectedRoutingKey, "topic")
	return err
}

func handleResponse(t *testing.T, err error, r interface{}) {
	fmt.Printf("Response: %v\n", r)
	if err != nil {
		assert.FailNow(t, "Request failed: %v\n", err)
	}
}
