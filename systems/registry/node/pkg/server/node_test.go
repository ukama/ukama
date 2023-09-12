package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/node/mocks"
	"github.com/ukama/ukama/systems/registry/node/pkg/db"
	"github.com/ukama/ukama/systems/registry/node/pkg/server"
	"gorm.io/gorm"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	opb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	omocks "github.com/ukama/ukama/systems/nucleus/org/pb/gen/mocks"
	pb "github.com/ukama/ukama/systems/registry/node/pb/gen"
)

var testNode = ukama.NewVirtualNodeId("HomeNode")
var orgId = uuid.NewV4()

const OrgName = "testorg"

func TestNodeServer_Add(t *testing.T) {
	nodeId := testNode.String()

	msgbusClient := &mbmocks.MsgBusServiceClient{}
	nodeRepo := &mocks.NodeRepo{}
	nodeStatusRepo := &mocks.NodeStatusRepo{}
	orgService := &mocks.OrgClientProvider{}
	networkService := &mocks.NetworkClientProvider{}

	const nodeName = "node-A"
	const nodeType = "hnode"

	s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", msgbusClient, orgService, networkService, orgId)

	node := &db.Node{
		Id:    nodeId,
		Name:  nodeName,
		OrgId: orgId,
		Type:  testNode.GetNodeType(),
		Status: db.NodeStatus{
			NodeId: nodeId,
			Conn:   db.Unknown,
			State:  db.Undefined,
		},
	}

	nodeRepo.On("Add", node, mock.Anything).Return(nil).Once()
	nodeRepo.On("GetNodeCount").Return(int64(1), int64(1), int64(0), nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	t.Run("NodeStateValid", func(t *testing.T) {
		// Arrange
		orgClient := orgService.On("GetClient").
			Return(&omocks.OrgServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*omocks.OrgServiceClient)

		orgClient.On("Get", mock.Anything,
			&opb.GetRequest{Id: orgId.String()}).
			Return(&opb.GetResponse{
				Org: &opb.Organization{
					Id: orgId.String(),
				},
			}, nil).Once()

		// Act
		res, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
			NodeId: nodeId,
			Name:   nodeName,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, nodeName, res.Node.Name)
		assert.Equal(t, nodeType, res.Node.Type)
		nodeRepo.AssertExpectations(t)
	})

}

func TestNodeServer_Get(t *testing.T) {
	t.Run("NodeFound", func(t *testing.T) {
		const nodeName = "node-A"
		const nodeType = ukama.NODE_ID_TYPE_HOMENODE
		var nodeId = ukama.NewVirtualNodeId(nodeType)

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).Return(
			&db.Node{Id: nodeId.StringLowercase(),
				Name: nodeName,
				Type: ukama.NODE_ID_TYPE_HOMENODE,
			}, nil).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, nil, orgId)

		resp, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{
			NodeId: nodeId.StringLowercase()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, nodeId.String(), resp.GetNode().GetId())
		assert.Equal(t, nodeName, resp.GetNode().Name)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		var nodeId = ukama.NewVirtualAmplifierNodeId()

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}

		nodeRepo.On("Get", nodeId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, nil, orgId)

		resp, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{
			NodeId: nodeId.StringLowercase()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})

	t.Run("NodeIdInvalid", func(t *testing.T) {
		var nodeId = uuid.NewV4()

		nodeRepo := &mocks.NodeRepo{}
		nodeStatusRepo := &mocks.NodeStatusRepo{}
		s := server.NewNodeServer(OrgName, nodeRepo, nil, nodeStatusRepo, "", nil, nil, nil, orgId)

		resp, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{
			NodeId: nodeId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		nodeRepo.AssertExpectations(t)
	})
}
