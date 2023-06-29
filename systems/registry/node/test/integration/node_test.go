//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/num30/config"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
	uconf "github.com/ukama/ukama/systems/common/config"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
)

var tConfig *TestConfig
var orgName string

func init() {
	// set org name
	orgName = fmt.Sprintf("node-integration-self-test-org-%d", time.Now().Unix())

	// load config
	tConfig = &TestConfig{}

	err := config.NewConfReader("integration").Read(tConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if tConfig.DebugMode {
		b, err := yaml.Marshal(tConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}

	log.Info("Expected config ", "integration.yaml", " or env vars for ex: SERVICEHOST")
	log.Infof("Config: %+v\n", tConfig)
}

type TestConfig struct {
	uconf.BaseConfig `mapstructure:",squash"`
	ServiceHost      string `default:"localhost:9090"`
}

func Test_FullFlow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Infoln("Connecting to network ", tConfig.ServiceHost)
	conn, err := grpc.DialContext(ctx, tConfig.ServiceHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		assert.NoError(t, err, "did not connect: %v", err)

		return
	}

	c := pb.NewNodeServiceClient(conn)
	// keep all used nodes here so we could delete them after test
	ndToClean := []ukama.NodeID{}

	// Contact the server and print out its response.
	node := ukama.NewVirtualHomeNodeId()
	orgId := uuid.NewV4()

	ndToClean = append(ndToClean, node)

	defer cleanupNodes(t, c, ndToClean)

	var r interface{}

	t.Run("AddAndUpdateNode", func(tt *testing.T) {
		nodeName := "HomeNodeX"
		addResp, err := c.AddNode(ctx, &pb.AddNodeRequest{
			NodeId: node.String(),
			OrgId:  orgId.String(),
			State:  db.Undefined.String(),
			Name:   nodeName,
		})

		handleResponse(tt, err, addResp)
		assert.NotNil(tt, addResp.Node)
		assert.Equal(tt, nodeName, addResp.Node.Name)

		r, err = c.UpdateNodeState(ctx, &pb.UpdateNodeStateRequest{
			NodeId: node.String(),
			State:  db.Offline.String(),
		})

		handleResponse(tt, err, r)

		nodeResp, err := c.GetNode(ctx, &pb.GetNodeRequest{
			NodeId: node.String()})

		handleResponse(tt, err, nodeResp)
		assert.Equal(tt, db.Onboarded.String(), nodeResp.Node.State)
		assert.Equal(tt, node.GetNodeType(), nodeResp.Node.Type)
	})

	tNodeID := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE)
	aNodeID := ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_AMPNODE)

	t.Run("AddTowerNodeWithAmplifiers", func(tt *testing.T) {
		ndToClean = append(ndToClean, tNodeID)

		_, err := c.AddNode(ctx, &pb.AddNodeRequest{
			NodeId: tNodeID.String(),
			State:  db.Undefined.String(),
			OrgId:  orgId.String(),
		})

		if err != nil {
			assert.FailNow(tt, "AddNode failed", err.Error())
		}

		ndToClean = append(ndToClean, aNodeID)

		_, err = c.AddNode(ctx, &pb.AddNodeRequest{
			NodeId: aNodeID.String(),
			OrgId:  orgId.String(),
			State:  db.Undefined.String(),
		})

		if err != nil {
			assert.FailNow(tt, "AddNode failed", err.Error())
		}

		_, err = c.AttachNodes(ctx, &pb.AttachNodesRequest{
			NodeId:        tNodeID.String(),
			AttachedNodes: []string{aNodeID.String()},
		})

		if err != nil {
			assert.FailNow(tt, "AttachNodes failed", err.Error())
		}

		resp, err := c.GetNode(ctx, &pb.GetNodeRequest{
			NodeId: tNodeID.String()})

		if assert.NoError(tt, err, "GetNode failed") {
			assert.NotNil(tt, resp.Node.Attached)
			assert.Equal(tt, 1, len(resp.Node.Attached))
			assert.Equal(tt, aNodeID.StringLowercase(), resp.Node.Attached[0].Id)
		}
	})

	t.Run("DetachNode", func(tt *testing.T) {
		_, err := c.DetachNode(ctx, &pb.DetachNodeRequest{
			NodeId: aNodeID.String(),
		})

		if assert.NoError(t, err) {
			resp, err := c.GetNode(ctx, &pb.GetNodeRequest{NodeId: tNodeID.String()})
			if assert.NoError(t, err) {
				assert.Nil(t, resp.Node.Attached)
			}
		}
	})
}

func cleanupNodes(tt *testing.T, c pb.NodeServiceClient, nodes []ukama.NodeID) {
	for _, node := range nodes {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		_, err := c.DeleteNode(ctx, &pb.DeleteNodeRequest{NodeId: node.String()})
		if err != nil {
			assert.FailNow(tt, "DeleteNode failed", err.Error())
		}

		_, err = c.GetNode(ctx, &pb.GetNodeRequest{NodeId: node.String()})
		if assert.Error(tt, err) {
			assert.Equal(tt, codes.NotFound, status.Code(err))
		}
	}
}

func Test_Listener(t *testing.T) {
	// Arrange
	nodeID := "UK-SA2136-HNODE-A1-30DF"
	orgId := uuid.NewV4()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, c, err := CreateRegistryClient()
	if err != nil {
		assert.NoErrorf(t, err, "did not connect: %+v\n", err)

		return
	}
	defer conn.Close()

	_, err = c.AddNode(ctx, &pb.AddNodeRequest{
		NodeId: nodeID,
		OrgId:  orgId.String(),
		State:  db.Undefined.String(),
	})

	e, ok := status.FromError(err)
	if ok && e.Code() == codes.AlreadyExists {
		log.Infof("node already exist, err %+v\n", err)
	} else {
		assert.NoError(t, err)

		return
	}

	_, err = c.UpdateNodeState(ctx, &pb.UpdateNodeStateRequest{
		NodeId: nodeID,
		State:  db.Offline.String()})

	if err != nil {
		assert.FailNow(t, "Failed to update node. Error: %s", err.Error())
	}

	// Act

	// Assert
	assert.NoError(t, err)

	time.Sleep(3 * time.Second)

	nodeResp, err := c.GetNode(ctx, &pb.GetNodeRequest{NodeId: nodeID})
	assert.NoError(t, err)

	if err != nil {
		assert.Equal(t, db.Onboarded.String(), nodeResp.Node.State)
	}
}

func CreateRegistryClient() (*grpc.ClientConn, pb.NodeServiceClient, error) {
	log.Infoln("Connecting to network ", tConfig.ServiceHost)

	context, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	conn, err := grpc.DialContext(context, tConfig.ServiceHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewNodeServiceClient(conn)

	return conn, c, nil
}

func handleResponse(t *testing.T, err error, r interface{}) {
	t.Helper()

	log.Printf("Response: %v\n", r)

	if err != nil {
		assert.FailNow(t, "Request failed: %v\n", err)
	}
}
