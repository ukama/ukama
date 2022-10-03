package server

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/init/lookup/internal/db"
	mocks "github.com/ukama/ukama/systems/init/lookup/mocks"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	mb "github.com/ukama/ukama/systems/init/lookup/pkg/msgBusClient"

	"github.com/stretchr/testify/assert"
)

var testNodeId = ukama.NewVirtualNodeId("HomeNode")

func TestLookupServer_AddOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	msgbusClient := &mb.MsgBusClient{}

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
	}

	orgRepo.On("Add", org).Return(nil).Once()

	s := NewLookupServer(nil, orgRepo, nil, msgbusClient)
	_, err = s.AddOrg(context.TODO(), &pb.AddOrgRequest{OrgName: "ukama", Certificate: "ukama_certs", Ip: "0.0.0.0"})

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_UpdateOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	msgbusClient := &mb.MsgBusClient{}

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
	orgRepo.On("Update", org).Return(nil).Once()

	s := NewLookupServer(nil, orgRepo, nil, msgbusClient)
	_, err = s.UpdateOrg(context.TODO(), &pb.UpdateOrgRequest{OrgName: "ukama", Certificate: "ukama_certs", Ip: "0.0.0.0"})

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_GetOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	msgbusClient := &mb.MsgBusClient{}

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

	s := NewLookupServer(nil, orgRepo, nil, msgbusClient)
	resp, err := s.GetOrg(context.TODO(), &pb.GetOrgRequest{OrgName: "ukama"})

	assert.NoError(t, err)
	assert.Equal(t, org.Name, resp.GetOrgName())
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_AddNodeForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	nodeRepo := &mocks.NodeRepo{}
	msgbusClient := &mb.MsgBusClient{}

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
	nodeRepo.On("AddOrUpdate", node).Return(nil).Once()

	s := NewLookupServer(nodeRepo, orgRepo, nil, msgbusClient)
	_, err = s.AddNodeForOrg(context.TODO(), &pb.AddNodeRequest{NodeId: nodeStr, OrgName: "ukama"})

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_GetNode(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	nodeRepo := &mocks.NodeRepo{}
	msgbusClient := &mb.MsgBusClient{}

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

	s := NewLookupServer(nodeRepo, orgRepo, nil, msgbusClient)
	resp, err := s.GetNode(context.TODO(), &pb.GetNodeRequest{NodeId: nodeStr})

	assert.NoError(t, err)
	assert.Equal(t, nodeStr, resp.NodeId)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_GetNodeForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	nodeRepo := &mocks.NodeRepo{}
	msgbusClient := &mb.MsgBusClient{}

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

	s := NewLookupServer(nodeRepo, orgRepo, nil, msgbusClient)
	resp, err := s.GetNodeForOrg(context.TODO(), &pb.GetNodeForOrgRequest{NodeId: nodeStr, OrgName: "ukama"})

	assert.NoError(t, err)
	assert.Equal(t, nodeStr, resp.NodeId)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_DeleteNodeForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	nodeRepo := &mocks.NodeRepo{}
	msgbusClient := &mb.MsgBusClient{}

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

	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	nodeRepo.On("Delete", testNodeId).Return(nil).Once()

	s := NewLookupServer(nodeRepo, orgRepo, nil, msgbusClient)
	_, err = s.DeleteNodeForOrg(context.TODO(), &pb.DeleteNodeRequest{NodeId: nodeStr, OrgName: "ukama"})

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)
}

func TestLookupServer_GetSystemForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	systemRepo := &mocks.SystemRepo{}
	msgbusClient := &mb.MsgBusClient{}

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
	}

	system := &db.System{
		Name:        "sys",
		Uuid:        uuid.New().String(),
		Certificate: "ukama_certs",
		Ip:          orgIp,
		Port:        100,
	}

	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	systemRepo.On("GetByName", system.Name).Return(system, nil).Once()

	s := NewLookupServer(nil, orgRepo, systemRepo, msgbusClient)
	resp, err := s.GetSystemForOrg(context.TODO(), &pb.GetSystemRequest{SystemName: system.Name, OrgName: "ukama"})

	assert.NoError(t, err)
	assert.Equal(t, system.Name, resp.SystemName)
	orgRepo.AssertExpectations(t)

}

func TestLookupServer_UpdateSystemForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	systemRepo := &mocks.SystemRepo{}
	msgbusClient := &mb.MsgBusClient{}

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
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

	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	systemRepo.On("GetByName", system.Name).Return(system, nil).Once()
	systemRepo.On("Update", system).Return(nil).Once()
	s := NewLookupServer(nil, orgRepo, systemRepo, msgbusClient)
	_, err = s.UpdateSystemForOrg(context.TODO(), &pb.UpdateSystemRequest{SystemName: system.Name, OrgName: "ukama", Certificate: "ukama_certs", Ip: ip, Port: 100})

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)

}

func TestLookupServer_DeleteSystemForOrg(t *testing.T) {
	orgRepo := &mocks.OrgRepo{}
	systemRepo := &mocks.SystemRepo{}
	msgbusClient := &mb.MsgBusClient{}

	var orgIp pgtype.Inet
	const ip = "0.0.0.0"
	err := orgIp.Set(ip)
	assert.NoError(t, err)

	org := &db.Org{
		Name:        "ukama",
		Certificate: "ukama_certs",
		Ip:          orgIp,
	}

	system := &db.System{
		Name:        "sys",
		Uuid:        uuid.New().String(),
		Certificate: "ukama_certs",
		Ip:          orgIp,
		Port:        100,
	}

	orgRepo.On("GetByName", org.Name).Return(org, nil).Once()
	systemRepo.On("Delete", system.Name).Return(nil).Once()
	s := NewLookupServer(nil, orgRepo, systemRepo, msgbusClient)
	_, err = s.DeleteSystemForOrg(context.TODO(), &pb.DeleteSystemRequest{SystemName: system.Name, OrgName: "ukama"})

	assert.NoError(t, err)
	orgRepo.AssertExpectations(t)

}
