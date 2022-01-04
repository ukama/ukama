package server

import (
	"context"
	"testing"

	mocks "github.com/ukama/ukamaX/cloud/registry/mocks"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"

	"github.com/ukama/ukamaX/cloud/registry/internal/db"

	uuid2 "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukamaX/common/ukama"
)

var testNodeId = ukama.NewVirtualNodeId("HomeNode")
var testDeviceGatewayHost = "1.1.1.1"

func TestRegistryServer_GetOrg(t *testing.T) {
	orgName := "org-1"
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	orgRepo.On("GetByName", mock.Anything).Return(&db.Org{Name: orgName}, nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo, &mocks.Client{}, testDeviceGatewayHost)
	org, err := s.GetOrg(context.TODO(), &pb.GetOrgRequest{Name: orgName})
	assert.NoError(t, err)
	assert.Equal(t, orgName, org.Name)
	orgRepo.AssertExpectations(t)
}

func TestRegistryServer_AddOrg_fails_without_owner_id(t *testing.T) {
	orgName := "org-1"
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	s := NewRegistryServer(orgRepo, nodeRepo, &mocks.Client{}, testDeviceGatewayHost)
	_, err := s.AddOrg(context.TODO(), &pb.AddOrgRequest{Name: orgName})
	assert.Error(t, err)
}

func TestRegistryServer_AddOrg_fails_with_bad_owner_id(t *testing.T) {
	orgName := "org-1"
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	s := NewRegistryServer(orgRepo, nodeRepo, &mocks.Client{}, testDeviceGatewayHost)
	_, err := s.AddOrg(context.TODO(), &pb.AddOrgRequest{Name: orgName})
	assert.Error(t, err)
}

func TestRegistryServer_AddOrg(t *testing.T) {
	// Arrange
	orgName := "org-1"
	ownerId := uuid2.NewV1().String()
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	// trick to call nested bootstrap call
	orgRepo.On("Add", mock.Anything, mock.Anything).
		Return(func(o *db.Org, f ...func() error) error {
			return f[0]()
		}).Once()
	bootstrapClient := &mocks.Client{}
	bootstrapClient.On("AddOrUpdateOrg", orgName, mock.Anything, testDeviceGatewayHost).Return(nil)

	s := NewRegistryServer(orgRepo, nodeRepo, bootstrapClient, testDeviceGatewayHost)

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
	ownerId := uuid2.NewV1()
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	nodeRepo.On("Get", testNodeId).Return(&db.Node{NodeID: testNodeId.String(), State: db.Pending,
		Org: &db.Org{
			Name:  orgName,
			Owner: ownerId,
		}}, nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo, &mocks.Client{}, testDeviceGatewayHost)
	node, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{NodeId: testNodeId.String()})
	assert.NoError(t, err)
	assert.Equal(t, orgName, node.Org.Name)
	assert.Equal(t, pb.NodeState_PENDING, node.Node.State)
	nodeRepo.AssertExpectations(t)
	orgRepo.AssertExpectations(t)
}

func TestRegistryServer_UpdateNode(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	nodeRepo.On("Update", testNodeId, db.Onboarded).Return(nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo, &mocks.Client{}, testDeviceGatewayHost)
	_, err := s.UpdateNode(context.TODO(), &pb.UpdateNodeRequest{
		NodeId: testNodeId.String(),
		State:  pb.NodeState_ONBOARDED,
	})
	assert.NoError(t, err)
	nodeRepo.AssertExpectations(t)
}

func TestRegistryServer_AddNode(t *testing.T) {
	// Arrange
	orgName := "node-1"
	nodeId := testNodeId.String()
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	orgRepo.On("GetByName", mock.Anything).Return(&db.Org{BaseModel: db.BaseModel{ID: 1}, Name: orgName}, nil).Once()
	nodeRepo.On("Add", mock.MatchedBy(func(n *db.Node) bool {
		return n.State == db.Pending && n.NodeID == nodeId && n.OrgID == 1
	}), mock.Anything).Return(func(o *db.Node, f ...func() error) error {
		return f[0]()
	}).Once()
	bootstrapClient := &mocks.Client{}
	bootstrapClient.On("AddDevice", orgName, nodeId).Return(nil)
	s := NewRegistryServer(orgRepo, nodeRepo, bootstrapClient, testDeviceGatewayHost)

	// Act
	_, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
		Node: &pb.Node{
			NodeId: nodeId,
			State:  pb.NodeState_PENDING,
		},
		OrgName: orgName,
	})

	// Assert
	assert.NoError(t, err)
	nodeRepo.AssertExpectations(t)
	bootstrapClient.AssertExpectations(t)
}

func TestRegistryServer_GetNodes(t *testing.T) {
	orgName := "node-1"
	nodeUuid1 := ukama.NewVirtualNodeId(ukama.OEMCODE)
	nodeUuid2 := ukama.NewVirtualNodeId(ukama.OEMCODE)
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}

	nodeRepo.On("GetByOrg", mock.Anything, mock.Anything).Return([]db.Node{
		{NodeID: nodeUuid1.String(), State: db.Undefined, Org: &db.Org{Name: orgName}},
		{NodeID: nodeUuid2.String(), State: db.Pending, Org: &db.Org{Name: orgName}},
	}, nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo, &mocks.Client{}, testDeviceGatewayHost)
	res, err := s.GetNodes(context.TODO(), &pb.GetNodesRequest{
		OrgName: orgName,
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, len(res.GetOrgs()))
	assert.Equal(t, res.Orgs[0].GetNodes()[0].State, pb.NodeState_UNDEFINED)
	assert.Equal(t, res.Orgs[0].OrgName, orgName)
	assert.Equal(t, res.Orgs[0].GetNodes()[1].State, pb.NodeState_PENDING)
	assert.Equal(t, res.Orgs[0].GetNodes()[1].NodeId, nodeUuid2.String())
	nodeRepo.AssertExpectations(t)
}

func TestRegistryServer_GetNodesReturnsEmptyList(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}

	nodeRepo.On("GetByOrg", mock.Anything, mock.Anything).Return([]db.Node{}, nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo, &mocks.Client{}, testDeviceGatewayHost)

	// act
	res, err := s.GetNodes(context.TODO(), &pb.GetNodesRequest{
		OrgName: "org-test",
	})

	// assert
	assert.NoError(t, err)
	assert.NotNil(t, res.Orgs)
	assert.Equal(t, 1, len(res.Orgs))
	assert.Equal(t, 0, len(res.Orgs[0].Nodes))
}
