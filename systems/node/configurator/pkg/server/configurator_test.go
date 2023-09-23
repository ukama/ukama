package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/configurator/mocks"
	pb "github.com/ukama/ukama/systems/node/configurator/pb/gen"
	"github.com/ukama/ukama/systems/node/configurator/pkg"
)

const testOrgName = "test-org"

var orgId = uuid.NewV4()

func TestConfiguratorServer_RestartSite(t *testing.T) {
	// Arrange
	msgclientRepo := &mbmocks.MsgBusServiceClient{}

	RegRepo := &mocks.RegistryProvider{}

	netId := uuid.NewV4()

	s := NewConfiguratorServer(msgclientRepo, RegRepo, pkg.IsDebugMode, testOrgName)

	RegRepo.On("ValidateSite", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	msgclientRepo.On("PublishRequest", mock.Anything, &pb.RestartSiteRequest{
		SiteName:  "pamoja",
		NetworkId: netId.String(),
	}).Return(nil).Once()
	// Act
	_, err := s.RestartSite(context.TODO(), &pb.RestartSiteRequest{
		SiteName:  "pamoja",
		NetworkId: netId.String(),
	})
	// Assert
	msgclientRepo.AssertExpectations(t)
	assert.NoError(t, err)

}
