package server

import (
	"context"
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/registry/network/mocks"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"
	"gorm.io/gorm"
)

func TestNetworkServer_AddNetwork(t *testing.T) {
	// Arrange
	const orgID = uint(1)
	const netName = "network-1"
	const orgName = "org-1"

	netRepo := &mocks.NetRepo{}
	orgRepo := &mocks.OrgRepo{}

	net := &db.Network{
		Name:  netName,
		OrgID: orgID,
	}

	orgRepo.On("GetByName", orgName).Return(
		&db.Org{Model: gorm.Model{ID: orgID},
			Name:        orgName,
			Deactivated: false},
		nil).Once()
	netRepo.On("Add", orgID, netName).Return(net, nil).Once()

	s := NewNetworkServer(netRepo, orgRepo, nil, nil)

	// Act
	res, err := s.Add(context.TODO(), &pb.AddRequest{
		Name:    netName,
		OrgName: orgName,
	})

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, orgName, res.Org)
	assert.Equal(t, netName, res.Network.Name)
	netRepo.AssertExpectations(t)
}
