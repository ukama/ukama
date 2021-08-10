package server

import (
	"context"
	uuid2 "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukamaX/common/ukama"
	"testing"
	"ukamaX/registry/internal/db"
	mocks "ukamaX/registry/mocks"
	pb "ukamaX/registry/pb/generated"
)

var nodeUuid = ukama.NewVirtualNodeId("HomeNode")

func TestRegistryServer_GetOrg(t *testing.T) {
	orgName := "org-1"
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	orgRepo.On("GetByName", mock.Anything).Return(&db.Org{Name: orgName}, nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo)
	org, err := s.GetOrg(context.TODO(), &pb.GetOrgRequest{Name: orgName})
	assert.NoError(t, err)
	assert.Equal(t, orgName, org.Name)
	orgRepo.AssertExpectations(t)
}

func TestRegistryServer_AddOrg_fails_without_owner_id(t *testing.T) {
	orgName := "org-1"
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	s := NewRegistryServer(orgRepo, nodeRepo)
	_, err := s.AddOrg(context.TODO(), &pb.AddOrgRequest{Name: orgName})
	assert.Error(t, err)
}

func TestRegistryServer_AddOrg_fails_with_bad_owner_id(t *testing.T) {
	orgName := "org-1"
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	s := NewRegistryServer(orgRepo, nodeRepo)
	_, err := s.AddOrg(context.TODO(), &pb.AddOrgRequest{Name: orgName})
	assert.Error(t, err)
}

func TestRegistryServer_AddOrg(t *testing.T) {
	orgName := "org-1"
	ownerId := uuid2.NewV1().String()
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	orgRepo.On("Add", mock.Anything).Return(nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo)
	res, err := s.AddOrg(context.TODO(), &pb.AddOrgRequest{Name: orgName, Owner: ownerId})
	assert.NoError(t, err)
	assert.Equal(t, orgName, res.Org.Name)
	assert.Equal(t, ownerId, res.Org.Owner)
	orgRepo.AssertExpectations(t)
}

func TestRegistryServer_GetNode(t *testing.T) {
	orgName := "node-1"
	ownerId := uuid2.NewV1()
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	nodeRepo.On("Get", nodeUuid).Return(&db.Node{NodeID: nodeUuid.String(), State: db.Pending,
		Org: &db.Org{
			Name:  orgName,
			Owner: ownerId,
		}}, nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo)
	node, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{NodeId: nodeUuid.String()})
	assert.NoError(t, err)
	assert.Equal(t, orgName, node.Org.Name)
	assert.Equal(t, pb.NodeState_PENDING, node.Node.State)
	nodeRepo.AssertExpectations(t)
	orgRepo.AssertExpectations(t)
}

func TestRegistryServer_UpdateNode(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	nodeRepo.On("Update", nodeUuid, db.Onboarded).Return(nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo)
	_, err := s.UpdateNode(context.TODO(), &pb.UpdateNodeRequest{
		NodeId: nodeUuid.String(),
		State:  pb.NodeState_ONBOARDED,
	})
	assert.NoError(t, err)
	nodeRepo.AssertExpectations(t)
}

func TestRegistryServer_AddNode(t *testing.T) {
	orgName := "node-1"
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	orgRepo.On("GetByName", mock.Anything).Return(&db.Org{BaseModel: db.BaseModel{ID: 1}}, nil).Once()
	nodeRepo.On("Add", mock.MatchedBy(func(n *db.Node) bool {
		return n.State == db.Pending && n.NodeID == nodeUuid.String() && n.OrgID == 1
	})).Return(nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo)
	_, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
		Node: &pb.Node{
			NodeId: nodeUuid.String(),
			State:  pb.NodeState_PENDING,
		},
		OrgName: orgName,
	})
	assert.NoError(t, err)
	nodeRepo.AssertExpectations(t)
}

func TestRegistryServer_GetNodes(t *testing.T) {
	orgName := "node-1"
	nodeUuid1 := ukama.NewVirtualNodeId(ukama.OEMCODE)
	nodeUuid2 := ukama.NewVirtualNodeId(ukama.OEMCODE)
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	ownerId := uuid2.NewV1()

	nodeRepo.On("GetByOrg", mock.Anything, mock.Anything).Return([]db.Node{
		{NodeID: nodeUuid1.String(), State: db.Undefined},
		{NodeID: nodeUuid2.String(), State: db.Pending},
	}, nil).Once()
	s := NewRegistryServer(orgRepo, nodeRepo)
	res, err := s.GetNodes(context.TODO(), &pb.GetNodesRequest{
		OrgName: orgName,
		Owner:   ownerId.String(),
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.GetNodes()))
	assert.Equal(t, res.GetNodes()[0].State, pb.NodeState_UNDEFINED)
	assert.Equal(t, res.GetNodes()[1].State, pb.NodeState_PENDING)
	assert.Equal(t, res.GetNodes()[1].NodeId, nodeUuid2.String())
	nodeRepo.AssertExpectations(t)
}
