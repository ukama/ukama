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
	t.Run("Org exist", func(t *testing.T) {
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
	})
}

func TestNetworkServer_Get(t *testing.T) {
	t.Run("Network exists", func(t *testing.T) {
		const netID = 1
		const netName = "network-1"

		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", uint(netID)).Return(
			&db.Network{Model: gorm.Model{ID: netID},
				Name:        netName,
				OrgID:       1,
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(netRepo, nil, nil, nil)
		netResp, err := s.Get(context.TODO(), &pb.GetRequest{
			NetworkID: netID})

		assert.NoError(t, err)
		assert.Equal(t, uint64(netID), netResp.GetNetwork().GetId())
		assert.Equal(t, netName, netResp.Network.Name)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetByName(t *testing.T) {
	t.Run("Org and Network exist", func(t *testing.T) {
		const netID = 1
		const orgName = "org-1"
		const netName = "network-1"

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", orgName, netName).Return(
			&db.Network{Model: gorm.Model{ID: netID},
				Name:        netName,
				OrgID:       1,
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(netRepo, nil, nil, nil)
		netResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{
			Name: netName, OrgName: orgName})

		assert.NoError(t, err)
		assert.Equal(t, uint64(netID), netResp.GetNetwork().GetId())
		assert.Equal(t, netName, netResp.Network.Name)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetByOrg(t *testing.T) {
	t.Run("Org exist", func(t *testing.T) {
		const netID = 1
		const orgID = 1
		const orgName = "org-1"
		const netName = "network-1"

		netRepo := &mocks.NetRepo{}
		orgRepo := &mocks.OrgRepo{}

		orgRepo.On("GetByName", orgName).Return(
			&db.Org{Model: gorm.Model{ID: orgID},
				Name:        orgName,
				Deactivated: false},
			nil).Once()

		netRepo.On("GetAllByOrgId", uint(orgID)).Return(
			[]db.Network{
				db.Network{Model: gorm.Model{ID: netID},
					Name:        netName,
					OrgID:       1,
					Deactivated: false,
				}}, nil).Once()

		s := NewNetworkServer(netRepo, orgRepo, nil, nil)
		netResp, err := s.GetByOrg(context.TODO(),
			&pb.GetByOrgRequest{OrgName: orgName})

		assert.NoError(t, err)
		assert.Equal(t, uint64(netID), netResp.GetNetworks()[0].GetId())
		assert.Equal(t, orgName, netResp.Org)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_Delete(t *testing.T) {
	t.Run("Org and Network exist", func(t *testing.T) {
		const netID = 1
		const orgName = "org-1"
		const netName = "network-1"

		netRepo := &mocks.NetRepo{}

		netRepo.On("Delete", orgName, netName).Return(nil).Once()

		s := NewNetworkServer(netRepo, nil, nil, nil)
		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Name: netName, OrgName: orgName})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Network does not exist", func(t *testing.T) {
		const netID = 1
		const orgName = "org-1"
		const netName = "network-1"

		netRepo := &mocks.NetRepo{}

		netRepo.On("Delete", orgName, netName).Return(gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(netRepo, nil, nil, nil)
		_, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Name: netName, OrgName: orgName})

		assert.Error(t, err)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_AddSite(t *testing.T) {
	t.Run("Network exists", func(t *testing.T) {
		// Arrange
		const netID = uint(1)
		const orgID = uint(1)
		const netName = "network-1"
		const siteName = "site-A"

		netRepo := &mocks.NetRepo{}
		siteRepo := &mocks.SiteRepo{}

		site := &db.Site{
			Name:      siteName,
			NetworkID: netID,
		}

		netRepo.On("Get", uint(netID)).Return(
			&db.Network{Model: gorm.Model{ID: netID},
				Name:        netName,
				OrgID:       orgID,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("Add", site).Return(nil).Once()

		s := NewNetworkServer(netRepo, nil, siteRepo, nil)

		// Act
		res, err := s.AddSite(context.TODO(), &pb.AddSiteRequest{
			NetworkID: uint64(netID),
			SiteName:  siteName,
		})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, siteName, res.Site.Name)
		assert.Equal(t, uint64(netID), res.Site.NetworkID)
		netRepo.AssertExpectations(t)
	})
}
