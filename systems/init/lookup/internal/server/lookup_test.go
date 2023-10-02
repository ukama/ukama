package server

import (
	"context"
	"testing"

	"github.com/jackc/pgtype"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/init/lookup/internal/db"
	mocks "github.com/ukama/ukama/systems/init/lookup/mocks"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testNodeId = ukama.NewVirtualNodeId("HomeNode")
var orgName = "test-org"

func TestLookupServer_AddOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	id := uuid.NewV4()
	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
		OrgId:       id,
	}

	porg := &pb.AddOrgRequest{OrgName: "ukama", Certificate: "ukama_certs", Ip: "0.0.0.0", OrgId: id.String()}

	orgRepo.On("Add", org).Return(nil).Once()
	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, porg).Return(nil).Once()

	s := NewLookupServer(nil, orgRepo, nil, msgbusClient, orgName)
	_, err = s.AddOrg(context.TODO(), porg)

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_UpdateOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
	}

	porg := &pb.UpdateOrgRequest{OrgName: "ukama", Certificate: "ukama_certs", Ip: "0.0.0.0"}

	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	orgRepo.On("Update", org).Return(nil).Once()
	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, porg).Return(nil).Once()

	s := NewLookupServer(nil, orgRepo, nil, msgbusClient, orgName)
	_, err = s.UpdateOrg(context.TODO(), porg)

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_GetOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
	}

	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()

	s := NewLookupServer(nil, orgRepo, nil, msgbusClient, orgName)
	resp, err := s.GetOrg(context.TODO(), &pb.GetOrgRequest{OrgName: "ukama"})

	assert.NoError(t, err)
	assert.Equal(t, org.Name, resp.GetOrgName())
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_AddNodeForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	nodeRepo := &mocks.NodeRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	nodeStr := testNodeId.StringLowercase()

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
	}

	node := &db.Node{
		NodeID: nodeStr,
		OrgID:  org.ID,
	}

	pnode := &pb.AddNodeRequest{NodeId: nodeStr, OrgName: "ukama"}
	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	nodeRepo.On("AddOrUpdate", node).Return(nil).Once()
	nodeRepo.On("Get", testNodeId).Return(node, nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, pnode).Return(nil).Once()

	s := NewLookupServer(nodeRepo, orgRepo, nil, msgbusClient, orgName)
	_, err = s.AddNodeForOrg(context.TODO(), pnode)

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_GetNode(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	nodeRepo := &mocks.NodeRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	nodeStr := testNodeId.StringLowercase()

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
	}

	node := &db.Node{
		NodeID: nodeStr,
		OrgID:  org.ID,
	}

	nodeRepo.On("Get", testNodeId).Return(node, nil).Once()

	s := NewLookupServer(nodeRepo, orgRepo, nil, msgbusClient, orgName)
	resp, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{NodeId: nodeStr})

	assert.NoError(t, err)
	assert.Equal(t, nodeStr, resp.NodeId)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_GetNodeForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	nodeRepo := &mocks.NodeRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	nodeStr := testNodeId.StringLowercase()

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
	}

	node := &db.Node{
		NodeID: nodeStr,
		OrgID:  org.ID,
	}

	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	nodeRepo.On("Get", testNodeId).Return(node, nil).Once()

	s := NewLookupServer(nodeRepo, orgRepo, nil, msgbusClient, orgName)
	resp, err := s.GetNodeForOrg(context.TODO(), &pb.GetNodeForOrgRequest{NodeId: nodeStr, OrgName: "ukama"})

	assert.NoError(t, err)
	assert.Equal(t, nodeStr, resp.NodeId)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_DeleteNodeForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	nodeRepo := &mocks.NodeRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	nodeStr := testNodeId.StringLowercase()

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
	}

	pnode := &pb.DeleteNodeRequest{NodeId: nodeStr, OrgName: "ukama"}

	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	nodeRepo.On("Delete", testNodeId).Return(nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, pnode).Return(nil).Once()

	s := NewLookupServer(nodeRepo, orgRepo, nil, msgbusClient, orgName)
	_, err = s.DeleteNodeForOrg(context.TODO(), pnode)

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_GetSystemForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	systemRepo := &mocks.SystemRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Model: gorm.Model{
			ID: 1,
		},
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
		OrgId:       uuid.NewV4(),
	}

	system := &db.System{
		Name:        "sys",
		Uuid:        uuid.NewV4().String(),
		Certificate: "ukama_certs",
		Ip:          orgIp,
		Port:        100,
	}

	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	systemRepo.On("GetByName", system.Name, org.ID).Return(system, nil).Once()

	s := NewLookupServer(nil, orgRepo, systemRepo, msgbusClient, orgName)
	resp, err := s.GetSystemForOrg(context.TODO(), &pb.GetSystemRequest{SystemName: system.Name, OrgName: "ukama"})

	assert.NoError(t, err)
	assert.Equal(t, system.Name, resp.SystemName)
	orgRepo.AssertExpectations(t)

}

func TestLookupServer_UpdateSystemForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	systemRepo := &mocks.SystemRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Model: gorm.Model{
			ID: 1,
		},
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
	}

	system := &db.System{
		Name:        "sys",
		Certificate: "ukama_certs",
		Ip:          orgIp,
		Port:        100,
	}

	psys := &pb.UpdateSystemRequest{SystemName: system.Name, OrgName: "ukama", Certificate: "ukama_certs", Ip: ip, Port: 100}

	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	systemRepo.On("GetByName", system.Name, org.ID).Return(system, nil).Once()
	systemRepo.On("Update", system, org.ID).Return(nil).Once()
	systemRepo.On("GetByName", system.Name, org.ID).Return(system, nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, psys).Return(nil).Once()

	s := NewLookupServer(nil, orgRepo, systemRepo, msgbusClient, orgName)
	_, err = s.UpdateSystemForOrg(context.TODO(), psys)

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)

}

func TestLookupServer_DeleteSystemForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	systemRepo := &mocks.SystemRepo{}
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Model: gorm.Model{
			ID: 1,
		},
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
	}

	system := &db.System{
		Name:        "sys",
		Uuid:        uuid.NewV4().String(),
		Certificate: "ukama_certs",
		Ip:          orgIp,
		Port:        100,
	}

	psys := &pb.DeleteSystemRequest{SystemName: system.Name, OrgName: "ukama"}

	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	systemRepo.On("Delete", system.Name, org.ID).Return(nil).Once()
	msgbusClient.On("PublishRequest", mock.Anything, psys).Return(nil).Once()

	s := NewLookupServer(nil, orgRepo, systemRepo, msgbusClient, orgName)
	_, err = s.DeleteSystemForOrg(context.TODO(), psys)

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)

}
