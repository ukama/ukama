package server

import (
	"context"
	"testing"

	mocks "github.com/ukama/ukama/services/cloud/registry/mocks"
	pb "github.com/ukama/ukama/services/cloud/registry/pb/gen"
	"github.com/ukama/ukama/services/cloud/registry/pkg/db"

	uuid2 "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/services/common/ukama"
)

var testNodeId = ukama.NewVirtualNodeId("HomeNode")
var testDeviceGatewayHost = "1.1.1.1"

const testOrgName = "org-1"
const testNetName = "net-1"
const testNetId = 98

type qPubStub struct {
}

func (q qPubStub) Publish(payload any, routingKey string) error {
	return nil
}

func TestRegistryServer_GetOrg(t *testing.T) {
	orgName := "org-1"
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := &mocks.NetRepo{}
	pub := &qPubStub{}
	orgRepo.On("GetByName", mock.Anything).Return(&db.Org{Name: orgName}, nil).Once()

	s := NewRegistryServer(orgRepo, nodeRepo, netRepo, &mocks.Client{}, testDeviceGatewayHost, pub)
	org, err := s.GetOrg(context.TODO(), &pb.GetOrgRequest{Name: orgName})
	assert.NoError(t, err)
	assert.Equal(t, orgName, org.Name)
	orgRepo.AssertExpectations(t)
}

func TestRegistryServer_AddOrg_fails_without_owner_id(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := createNetRepoMock()
	pub := &qPubStub{}
	s := NewRegistryServer(orgRepo, nodeRepo, netRepo, &mocks.Client{}, testDeviceGatewayHost, pub)
	_, err := s.AddOrg(context.TODO(), &pb.AddOrgRequest{Name: testOrgName})
	assert.Error(t, err)
}

func TestRegistryServer_AddOrg_fails_with_bad_owner_id(t *testing.T) {
	orgName := "org-1"
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := createNetRepoMock()
	pub := &qPubStub{}
	s := NewRegistryServer(orgRepo, nodeRepo, netRepo, &mocks.Client{}, testDeviceGatewayHost, pub)
	_, err := s.AddOrg(context.TODO(), &pb.AddOrgRequest{Name: orgName})
	assert.Error(t, err)
}

func TestRegistryServer_AddOrg(t *testing.T) {
	// Arrange
	orgName := "org-1"
	ownerId := uuid2.NewString()
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := &mocks.NetRepo{}
	pub := &qPubStub{}

	// trick to call nested bootstrap call
	orgRepo.On("Add", mock.Anything, mock.Anything).
		Return(func(o *db.Org, f ...func() error) error {
			return f[0]()
		}).Once()

	netRepo.On("Add", mock.Anything, mock.Anything).Return(&db.Network{}, nil).Once()
	bootstrapClient := &mocks.Client{}
	bootstrapClient.On("AddOrUpdateOrg", orgName, mock.Anything, testDeviceGatewayHost).Return(nil)

	s := NewRegistryServer(orgRepo, nodeRepo, netRepo, bootstrapClient, testDeviceGatewayHost, pub)

	// Act
	res, err := s.AddOrg(context.TODO(), &pb.AddOrgRequest{Name: orgName, Owner: ownerId})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, orgName, res.Org.Name)
	assert.Equal(t, ownerId, res.Org.Owner)
	orgRepo.AssertExpectations(t)
	bootstrapClient.AssertExpectations(t)
}

func TestRegistryServer_GetNode(t *testing.T) {
	orgName := "node-1"
	ownerId := uuid2.New()
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := createNetRepoMock()
	pub := &qPubStub{}

	nodeRepo.On("Get", testNodeId).Return(&db.Node{NodeID: testNodeId.String(),
		State: db.Pending, Type: db.NodeTypeHome,
		Network: &db.Network{
			Name: testNetName,
			Org: &db.Org{
				Name:  orgName,
				Owner: ownerId,
			},
		}}, nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo, netRepo, &mocks.Client{}, testDeviceGatewayHost, pub)
	node, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{NodeId: testNodeId.String()})
	assert.NoError(t, err)
	assert.Equal(t, orgName, node.Org.Name)
	assert.Equal(t, pb.NodeState_PENDING, node.Node.State)
	assert.Equal(t, pb.NodeType_HOME, node.Node.Type)
	nodeRepo.AssertExpectations(t)
	orgRepo.AssertExpectations(t)
}

func TestRegistryServer_UpdateNodeState(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := createNetRepoMock()
	pub := qPubStub{}

	nodeRepo.On("Update", testNodeId, mock.MatchedBy(func(ns *db.NodeState) bool {
		return *ns == db.Onboarded
	}), (*string)(nil)).Return(nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo, netRepo, &mocks.Client{}, testDeviceGatewayHost, pub)
	_, err := s.UpdateNodeState(context.TODO(), &pb.UpdateNodeStateRequest{
		NodeId: testNodeId.String(),
		State:  pb.NodeState_ONBOARDED,
	})

	// Assert
	assert.NoError(t, err)
	nodeRepo.AssertExpectations(t)
}

func TestRegistryServer_AddNode(t *testing.T) {
	// Arrange
	nodeId := testNodeId.String()
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := createNetRepoMock()
	pub := &qPubStub{}

	nodeRepo.On("Add", mock.MatchedBy(func(n *db.Node) bool {
		return n.State == db.Pending && n.NodeID == nodeId && n.NetworkID == testNetId
	}), mock.Anything).Return(func(o *db.Node, f ...func() error) error {
		return f[0]()
	}).Once()
	bootstrapClient := &mocks.Client{}
	bootstrapClient.On("AddNode", testOrgName, nodeId).Return(nil)
	s := NewRegistryServer(orgRepo, nodeRepo, netRepo, bootstrapClient, testDeviceGatewayHost, pub)

	// Act
	actNode, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
		Node: &pb.Node{
			NodeId: nodeId,
			State:  pb.NodeState_PENDING,
		},
		OrgName: testOrgName,
		Network: testNetName,
	})

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, actNode.Node.Name)
	nodeRepo.AssertExpectations(t)
	bootstrapClient.AssertExpectations(t)
}

func createNetRepoMock() *mocks.NetRepo {
	netRepo := &mocks.NetRepo{}
	netRepo.On("Get", testOrgName, testNetName).
		Return(&db.Network{
			BaseModel: db.BaseModel{ID: testNetId},
			Name:      testNetName,
			Org: &db.Org{
				Name: testOrgName,
				BaseModel: db.BaseModel{
					ID: 101,
				},
			}}, nil).Once()
	return netRepo
}

func TestRegistryServer_GetNodes(t *testing.T) {
	orgName := "node-1"
	nodeUuid1 := ukama.NewVirtualNodeId(ukama.OEMCODE)
	nodeUuid2 := ukama.NewVirtualNodeId(ukama.OEMCODE)
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := &mocks.NetRepo{}
	pub := &qPubStub{}

	const NodeName0 = "NodeName0"
	nodeRepo.On("GetByOrg", mock.Anything, mock.Anything).Return([]db.Node{
		{NodeID: nodeUuid1.String(), State: db.Undefined, Name: NodeName0, Network: &db.Network{Org: &db.Org{Name: orgName}}},
		{NodeID: nodeUuid2.String(), State: db.Pending, Name: "NodeNeme2", Network: &db.Network{Org: &db.Org{Name: orgName}}},
	}, nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo, netRepo, &mocks.Client{}, testDeviceGatewayHost, pub)
	resp, err := s.GetNodes(context.TODO(), &pb.GetNodesRequest{
		OrgName: orgName,
	})

	assert.NoError(t, err)
	assert.Equal(t, pb.NodeState_UNDEFINED, resp.Nodes[0].State)
	assert.Equal(t, NodeName0, resp.Nodes[0].Name)
	assert.Equal(t, resp.OrgName, orgName)
	assert.Equal(t, resp.Nodes[1].State, pb.NodeState_PENDING)
	assert.Equal(t, resp.Nodes[1].NodeId, nodeUuid2.String())
	nodeRepo.AssertExpectations(t)
}

func TestRegistryServer_GetNodesReturnsEmptyList(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := &mocks.NetRepo{}
	pub := &qPubStub{}

	nodeRepo.On("GetByOrg", mock.Anything, mock.Anything).Return([]db.Node{}, nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo, netRepo, &mocks.Client{}, testDeviceGatewayHost, pub)

	// act
	res, err := s.GetNodes(context.TODO(), &pb.GetNodesRequest{
		OrgName: "org-test",
	})

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, res.Nodes)
	assert.Equal(t, 0, len(res.Nodes))
}

func Test_toDbNodeType(t *testing.T) {

	tests := []struct {
		nodeId ukama.NodeID
		want   db.NodeType
	}{
		{
			nodeId: ukama.NewVirtualHomeNodeId(),
			want:   db.NodeTypeHome,
		},
		{
			nodeId: ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_TOWERNODE),
			want:   db.NodeTypeTower,
		},
		{
			nodeId: ukama.NewVirtualNodeId(ukama.NODE_ID_TYPE_AMPNODE),
			want:   db.NodeTypeAmplifier,
		},
		{
			nodeId: ukama.NewVirtualNodeId("unknown"),
			want:   db.NodeTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.nodeId.String(), func(t *testing.T) {

			got := toDbNodeType(tt.nodeId.GetNodeType())
			assert.Equal(t, tt.want, got)
		})
	}
}
