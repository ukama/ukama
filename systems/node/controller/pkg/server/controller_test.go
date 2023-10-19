package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"

	"github.com/ukama/ukama/systems/node/controller/mocks"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	"github.com/ukama/ukama/systems/node/controller/pkg"
	"google.golang.org/protobuf/types/known/anypb"
)


const testOrgName = "test-org"

var orgId = uuid.NewV4()

func TestControllerServer_RestartSite(t *testing.T) {
	// Arrange
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}

	RegRepo := &mocks.RegistryProvider{}

	netId := uuid.NewV4()

	s := NewControllerServer(testOrgName,conRepo,msgclientRepo,RegRepo,pkg.IsDebugMode)
	nodeId := "uk-983794-hnode-78-7830"
	anyMsg, err := anypb.New(&pb.RestartSiteRequest{
		SiteName:  "pamoja",
		NetworkId: netId.String(),
	})
	if err != nil {
		return 
	}
	nodeLog := &db.NodeLog{
		NodeId: nodeId,
	}
	RegRepo.On("ValidateSite", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	RegRepo.On("ValidateNetwork", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	RegRepo.On("GetNodesBySite", "pamoja",mock.Anything,mock.Anything).Return([]string{nodeId}, nil).Once()
	conRepo.On("Get", nodeId,mock.Anything).Return(nodeLog,nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &cpb.NodeFeederMessage{
		Target:"test-org."+ "." + "." + nodeId,
		HTTPMethod: "POST",
		Path:       "/v1/reboot/"+nodeId,
		Msg:        anyMsg,
	}).Return(nil).Once()
	// Act
	_, err = s.RestartSite(context.TODO(), &pb.RestartSiteRequest{
		SiteName:  "pamoja",
		NetworkId: netId.String(),
	})
	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)

}
func TestControllerServer_RestartNode(t *testing.T) {
	// Arrange
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	RegRepo := &mocks.RegistryProvider{}

	nodeId := "uk-983794-hnode-78-7830"
	s := NewControllerServer(testOrgName,conRepo,msgclientRepo,RegRepo,pkg.IsDebugMode)

	anyMsg, err := anypb.New(&pb.RestartNodeRequest{
		NodeId: nodeId,
	})
	if err != nil {
		return 
	}

	NodeLog := db.NodeLog{
		NodeId: nodeId,
	}
	conRepo.On("Get", nodeId).Return(&NodeLog, nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &cpb.NodeFeederMessage{
		Target:     "test-org" + "." + "." + "." + nodeId,
		HTTPMethod: "POST",
		Path:       "/v1/reboot/"+nodeId,
		Msg:        anyMsg,
	}).Return(nil).Once()
	// Act
	_, err = s.RestartNode(context.TODO(), &pb.RestartNodeRequest{
		NodeId: nodeId,
	})
	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)

}

func TestControllerServer_RestartNodes(t *testing.T) {
	// Arrange
	msgclientRepo := &mbmocks.MsgBusServiceClient{}
	conRepo := &mocks.NodeLogRepo{}
	RegRepo := &mocks.RegistryProvider{}
	netId := uuid.NewV4()
	nodeId := "uk-983794-hnode-78-7830"
	s := NewControllerServer(testOrgName,conRepo,msgclientRepo,RegRepo,pkg.IsDebugMode)

	anyMsg, err := anypb.New(&pb.RestartNodesRequest{
		NetworkId: netId.String(),
		NodeIds: []string{nodeId},
	})
	if err != nil {
		return 
	}

	NodeLog := db.NodeLog{
		NodeId: nodeId,
	}
	conRepo.On("Get", nodeId).Return(&NodeLog, nil).Once()

	msgclientRepo.On("PublishRequest", "request.cloud.local.test-org.node.controller.nodefeeder.publish", &cpb.NodeFeederMessage{
		Target:     "test-org" + "." + "." + "." + nodeId,
		HTTPMethod: "POST",
		Path:       "/v1/reboot/"+nodeId,
		Msg:        anyMsg,
	}).Return(nil).Once()
	// Act
	_, err = s.RestartNodes(context.TODO(), &pb.RestartNodesRequest{
		NetworkId: netId.String(),
		NodeIds: []string{nodeId},
	})
	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)

}