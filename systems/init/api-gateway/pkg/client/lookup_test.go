package client

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	//"github.com/ukama/ukama/systems/init/lookup/gen/mocks"
	"github.com/ukama/ukama/services/common/ukama"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	mocks "github.com/ukama/ukama/systems/init/lookup/pb/gen/mocks"
)

const sys = "sys"
const org = "org-name"

func TestLookupClient_AddOrg(t *testing.T) {
	lc := &mocks.LookupServiceClient{}
	orgReq := &pb.AddOrgRequest{
		OrgName:     org,
		Certificate: "certs",
		Ip:          "0.0.0.0",
	}

	lc.On("AddOrg", mock.Anything, orgReq).Return(&pb.AddOrgResponse{}, nil)

	l := &Lookup{
		client: lc,
	}

	_, err := l.AddOrg(orgReq)
	assert.NoError(t, err)
}

func TestLookupClient_UpdateOrg(t *testing.T) {
	lc := &mocks.LookupServiceClient{}
	orgReq := &pb.UpdateOrgRequest{
		OrgName:     org,
		Certificate: "updated_certs",
	}

	lc.On("UpdateOrg", mock.Anything, orgReq).Return(&pb.UpdateOrgResponse{}, nil)

	l := &Lookup{
		client: lc,
	}

	_, err := l.UpdateOrg(orgReq)
	assert.NoError(t, err)
}

func TestLookupClient_GetOrg(t *testing.T) {
	lc := &mocks.LookupServiceClient{}
	orgReq := &pb.GetOrgRequest{
		OrgName: org,
	}

	orgResp := &pb.GetOrgResponse{
		OrgName:     org,
		Certificate: "certs",
		Ip:          "0.0.0.0",
	}

	lc.On("GetOrg", mock.Anything, orgReq).Return(orgResp, nil)

	l := &Lookup{
		client: lc,
	}

	resp, err := l.GetOrg(orgReq)
	if assert.NoError(t, err) {
		lc.AssertExpectations(t)
		assert.Contains(t, resp.OrgName, org)
	}
}

func TestLookupClient_AddSystemForOrg(t *testing.T) {
	lc := &mocks.LookupServiceClient{}
	sysReq := &pb.AddSystemRequest{
		SystemName:  sys,
		OrgName:     org,
		Certificate: "certs",
		Ip:          "0.0.0.0",
		Port:        100,
	}

	lc.On("AddSystemForOrg", mock.Anything, sysReq).Return(&pb.AddSystemResponse{}, nil)

	l := &Lookup{
		client: lc,
	}

	_, err := l.AddSystemForOrg(sysReq)
	assert.NoError(t, err)
}

func TestLookupClient_UpdateSystemForOrg(t *testing.T) {
	lc := &mocks.LookupServiceClient{}
	sysReq := &pb.UpdateSystemRequest{
		SystemName:  sys,
		OrgName:     org,
		Certificate: "update_certs",
		Ip:          "127.0.0.1",
		Port:        101,
	}

	lc.On("UpdateSystemForOrg", mock.Anything, sysReq).Return(&pb.UpdateSystemResponse{}, nil)

	l := &Lookup{
		client: lc,
	}

	_, err := l.UpdateSystemForOrg(sysReq)
	assert.NoError(t, err)
}

func TestLookupClient_GetSystemForOrg(t *testing.T) {
	lc := &mocks.LookupServiceClient{}
	sysId := uuid.New().String()
	sysReq := &pb.GetSystemRequest{
		OrgName:    org,
		SystemName: sys,
	}
	sysResp := &pb.GetSystemResponse{
		SystemName:  sys,
		SystemId:    sysId,
		Certificate: "certs",
		Ip:          "0.0.0.0",
		Port:        100,
	}

	lc.On("GetSystemForOrg", mock.Anything, sysReq).Return(sysResp, nil)

	l := &Lookup{
		client: lc,
	}

	resp, err := l.GetSystemForOrg(sysReq)
	if assert.NoError(t, err) {
		lc.AssertExpectations(t)
		assert.Contains(t, resp.SystemId, sysId)
		assert.Contains(t, resp.SystemName, sys)
	}

}

func TestLookupClient_DeleteSystemForOrg(t *testing.T) {
	lc := &mocks.LookupServiceClient{}

	sysReq := &pb.DeleteSystemRequest{
		OrgName:    org,
		SystemName: sys,
	}

	lc.On("DeleteSystemForOrg", mock.Anything, sysReq).Return(&pb.DeleteSystemResponse{}, nil)

	l := &Lookup{
		client: lc,
	}

	_, err := l.DeleteSystemForOrg(sysReq)
	assert.NoError(t, err)

}

func TestLookupClient_AddNodeForOrg(t *testing.T) {
	lc := &mocks.LookupServiceClient{}
	nodeId := ukama.NewVirtualNodeId("homenode").String()

	nodeReq := &pb.AddNodeRequest{
		NodeId:  nodeId,
		OrgName: "org-name",
	}

	lc.On("AddNodeForOrg", mock.Anything, nodeReq).Return(&pb.AddNodeResponse{}, nil)

	l := &Lookup{
		client: lc,
	}

	_, err := l.AddNodeForOrg(nodeReq)
	assert.NoError(t, err)
}

func TestLookupClient_GetNodeForOrg(t *testing.T) {
	lc := &mocks.LookupServiceClient{}
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	nodeReq := &pb.GetNodeForOrgRequest{
		OrgName: org,
		NodeId:  nodeId,
	}
	nodeResp := &pb.GetNodeResponse{
		NodeId:      nodeId,
		OrgName:     "org-name",
		Certificate: "certs",
		Ip:          "0.0.0.0",
	}

	lc.On("GetNodeForOrg", mock.Anything, nodeReq).Return(nodeResp, nil)

	l := &Lookup{
		client: lc,
	}

	resp, err := l.GetNodeForOrg(nodeReq)
	if assert.NoError(t, err) {
		lc.AssertExpectations(t)
		assert.Contains(t, resp.NodeId, nodeId)
		assert.Contains(t, resp.OrgName, org)
	}

}

func TestLookupClient_DeleteNodeForOrg(t *testing.T) {
	lc := &mocks.LookupServiceClient{}
	nodeId := ukama.NewVirtualNodeId("homenode").String()
	sysReq := &pb.DeleteNodeRequest{
		OrgName: org,
		NodeId:  nodeId,
	}

	lc.On("DeleteNodeForOrg", mock.Anything, sysReq).Return(&pb.DeleteNodeResponse{}, nil)

	l := &Lookup{
		client: lc,
	}

	_, err := l.DeleteNodeForOrg(sysReq)
	assert.NoError(t, err)

}
