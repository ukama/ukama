package server

import (
	"context"
	"testing"

	"gorm.io/gorm"

	mocks "github.com/ukama/ukama/systems/registry/network/mocks"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ukama/ukama/systems/common/ukama"
)

var testNodeId = ukama.NewVirtualNodeId("HomeNode")

const testOrgName = "org-1"
const testNetName = "net-1"
const testNetId = 98

func TestNetworkServer_UpdateNode(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := createNetRepoMock()

	nodeRepo.On("Update", testNodeId, mock.MatchedBy(func(ns *db.NodeAttributes) bool {
		return *ns.State == db.Onboarded
	}), mock.Anything).Return(nil).Once()

	s := NewNetworkServer(orgRepo, nodeRepo, netRepo)

	_, err := s.UpdateNode(context.TODO(), &pb.UpdateNodeRequest{
		NodeId: testNodeId.String(),
		Node: &pb.Node{
			State: pb.NodeState_ONBOARDED,
		},
	})

	// Assert
	assert.NoError(t, err)
	nodeRepo.AssertExpectations(t)
}

func TestNetworkServer_AddNode(t *testing.T) {
	// Arrange
	nodeId := testNodeId.String()
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := createNetRepoMock()

	nodeRepo.On("Add", mock.MatchedBy(func(n *db.Node) bool {
		return n.State == db.Pending && n.NodeID == nodeId && n.NetworkID == testNetId
	}), mock.Anything).Return(func(o *db.Node, f ...func() error) error {
		return f[0]()
	}).Once()

	s := NewNetworkServer(orgRepo, nodeRepo, netRepo)

	// Act
	actNode, err := s.AddNode(context.TODO(), &pb.AddNodeRequest{
		Node: &pb.Node{
			NodeId: nodeId,
			State:  pb.NodeState_PENDING,
			Name:   "node-1",
		},
		OrgName: testOrgName,
		Network: testNetName,
	})

	// Assert
	if assert.NoError(t, err) {
		assert.Equal(t, "node-1", actNode.Node.Name)
		nodeRepo.AssertExpectations(t)
	}
}

func createNetRepoMock() *mocks.NetRepo {
	netRepo := &mocks.NetRepo{}

	netRepo.On("Get", testOrgName, testNetName).
		Return(&db.Network{
			Model: gorm.Model{ID: testNetId},
			Name:  testNetName,
			Org: &db.Org{
				Name:  testOrgName,
				Model: gorm.Model{ID: 101},
			}}, nil).Once()

	return netRepo
}

func TestNetworkServer_GetNodes(t *testing.T) {
	orgName := "node-1"
	nodeUuid1 := ukama.NewVirtualNodeId(ukama.OEMCODE)
	nodeUuid2 := ukama.NewVirtualNodeId(ukama.OEMCODE)
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := &mocks.NetRepo{}

	const NodeName0 = "NodeName0"
	nodeRepo.On("GetByOrg", mock.Anything, mock.Anything).Return([]db.Node{
		{NodeID: nodeUuid1.String(), State: db.Undefined, Name: NodeName0, Network: &db.Network{Org: &db.Org{Name: orgName}}},
		{NodeID: nodeUuid2.String(), State: db.Pending, Name: "NodeNeme2", Network: &db.Network{Org: &db.Org{Name: orgName}}},
	}, nil).Once()

	s := NewNetworkServer(orgRepo, nodeRepo, netRepo)

	resp, err := s.GetNodes(context.TODO(), &pb.GetNodesRequest{
		OrgName: orgName,
	})

	if assert.NoError(t, err) {
		assert.Equal(t, pb.NodeState_UNDEFINED, resp.Nodes[0].State)
		assert.Equal(t, NodeName0, resp.Nodes[0].Name)
		assert.Equal(t, resp.OrgName, orgName)
		assert.Equal(t, pb.NodeState_PENDING, resp.Nodes[1].State)
		assert.Equal(t, nodeUuid2.String(), resp.Nodes[1].NodeId)
		nodeRepo.AssertExpectations(t)
	}
}

func TestNetworkServer_GetNodesReturnsEmptyList(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := &mocks.NetRepo{}

	nodeRepo.On("GetByOrg", mock.Anything, mock.Anything).Return([]db.Node{}, nil).Once()

	s := NewNetworkServer(orgRepo, nodeRepo, netRepo)

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

func Test_List(t *testing.T) {
	nodeRepo := &mocks.NodeRepo{}
	orgRepo := &mocks.OrgRepo{}
	netRepo := &mocks.NetRepo{}

	queryRes := map[string]map[string]map[db.NodeType]int{
		"a": {
			"n1": map[db.NodeType]int{
				db.NodeTypeAmplifier: 1,
				db.NodeTypeTower:     5,
			},
			"n2": map[db.NodeType]int{
				db.NodeTypeAmplifier: 2,
				db.NodeTypeTower:     6,
				db.NodeTypeHome:      7,
			},
		},
		"b": {
			"n1": map[db.NodeType]int{
				db.NodeTypeAmplifier: 1,
			},
		},
	}

	netRepo.On("List").Return(queryRes, nil).Once()

	s := NewNetworkServer(orgRepo, nodeRepo, netRepo)

	// act
	res, err := s.List(context.TODO(), &pb.ListRequest{})

	// assert
	if assert.NoError(t, err) && assert.NotNil(t, res.Orgs) {
		var a, b *pb.ListResponse_Org

		for _, org := range res.Orgs {
			switch org.Name {
			case "a":
				a = org
			case "b":
				b = org
			}
		}

		assert.Len(t, res.Orgs, 2)
		assert.Len(t, a.GetNetworks(), 2)

		var n1, n2 *pb.ListResponse_Network

		if a.GetNetworks()[0].GetName() == "n1" {
			n1 = a.GetNetworks()[0]
			n2 = a.GetNetworks()[1]
		} else {
			n2 = a.GetNetworks()[0]
			n1 = a.GetNetworks()[1]
		}

		assert.Equal(t, uint32(1), n1.GetNumberOfNodes()["amplifier"])

		assert.Equal(t, uint32(2), n2.GetNumberOfNodes()["amplifier"])
		assert.Equal(t, uint32(6), n2.GetNumberOfNodes()["tower"])
		assert.Equal(t, uint32(7), n2.GetNumberOfNodes()["home"])

		assert.Equal(t, "n1", b.GetNetworks()[0].GetName())
		assert.Equal(t, uint32(1), b.GetNetworks()[0].GetNumberOfNodes()["amplifier"])
	}
}
