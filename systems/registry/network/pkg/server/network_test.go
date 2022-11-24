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
	t.Run("Org exists", func(t *testing.T) {
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
		assert.NotNil(t, res)
		assert.Equal(t, orgName, res.Org)
		assert.Equal(t, netName, res.Network.Name)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_Get(t *testing.T) {
	t.Run("Network found", func(t *testing.T) {
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
		assert.NotNil(t, netResp)
		assert.Equal(t, uint64(netID), netResp.GetNetwork().GetId())
		assert.Equal(t, netName, netResp.Network.Name)
		netRepo.AssertExpectations(t)
	})

	t.Run("Network not found", func(t *testing.T) {
		const netID = 1
		const netName = "network-1"

		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", uint(netID)).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(netRepo, nil, nil, nil)
		netResp, err := s.Get(context.TODO(), &pb.GetRequest{
			NetworkID: netID})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetByName(t *testing.T) {
	t.Run("Org and Network found", func(t *testing.T) {
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
		assert.NotNil(t, netResp)
		assert.Equal(t, uint64(netID), netResp.GetNetwork().GetId())
		assert.Equal(t, netName, netResp.Network.Name)
		netRepo.AssertExpectations(t)
	})

	t.Run("Org or Network not found", func(t *testing.T) {
		const netID = 1
		const orgName = "org-1"
		const netName = "network-1"

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", orgName, netName).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(netRepo, nil, nil, nil)
		netResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{
			Name: netName, OrgName: orgName})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetByOrg(t *testing.T) {
	t.Run("Org found", func(t *testing.T) {
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
		assert.NotNil(t, netResp)
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
		netResp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Name: netName, OrgName: orgName})

		assert.Error(t, err)
		assert.Nil(t, netResp)
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
		assert.NotNil(t, res)
		assert.Equal(t, siteName, res.Site.Name)
		assert.Equal(t, uint64(netID), res.Site.NetworkID)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetSite(t *testing.T) {
	t.Run("Site exists", func(t *testing.T) {
		const siteID = 1
		const siteName = "site-A"

		siteRepo := &mocks.SiteRepo{}

		siteRepo.On("Get", uint(siteID)).Return(
			&db.Site{Model: gorm.Model{ID: siteID},
				Name:        siteName,
				NetworkID:   1,
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(nil, nil, siteRepo, nil)
		netResp, err := s.GetSite(context.TODO(), &pb.GetSiteRequest{
			SiteID: siteID})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, uint64(siteID), netResp.GetSite().GetId())
		assert.Equal(t, siteName, netResp.GetSite().GetName())
		siteRepo.AssertExpectations(t)
	})

	t.Run("Site not found", func(t *testing.T) {
		const siteID = 1
		const siteName = "site-A"

		siteRepo := &mocks.SiteRepo{}

		siteRepo.On("Get", uint(siteID)).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(nil, nil, siteRepo, nil)
		netResp, err := s.GetSite(context.TODO(), &pb.GetSiteRequest{
			SiteID: siteID})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		siteRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetSiteByName(t *testing.T) {
	t.Run("Site exists", func(t *testing.T) {
		const siteID = 1
		const netID = 1
		const orgID = 1
		const siteName = "site-A"
		const netName = "net-1"

		siteRepo := &mocks.SiteRepo{}
		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", uint(netID)).Return(
			&db.Network{Model: gorm.Model{ID: netID},
				Name:        netName,
				OrgID:       orgID,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("GetByName", uint(netID), siteName).Return(
			&db.Site{Model: gorm.Model{ID: siteID},
				Name:        siteName,
				NetworkID:   1,
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(netRepo, nil, siteRepo, nil)
		netResp, err := s.GetSiteByName(context.TODO(), &pb.GetSiteByNameRequest{
			NetworkID: netID, SiteName: siteName})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, uint64(siteID), netResp.GetSite().GetId())
		assert.Equal(t, siteName, netResp.GetSite().GetName())
		siteRepo.AssertExpectations(t)
	})

	t.Run("Site not found", func(t *testing.T) {
		const siteID = 1
		const netID = 1
		const orgID = 1
		const siteName = "site-A"
		const netName = "net-1"

		siteRepo := &mocks.SiteRepo{}
		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", uint(netID)).Return(
			&db.Network{Model: gorm.Model{ID: netID},
				Name:        netName,
				OrgID:       orgID,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("GetByName", uint(netID), siteName).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(netRepo, nil, siteRepo, nil)
		netResp, err := s.GetSiteByName(context.TODO(), &pb.GetSiteByNameRequest{
			NetworkID: netID, SiteName: siteName})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		siteRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetSiteByNetwork(t *testing.T) {
	t.Run("Network found", func(t *testing.T) {
		const netID = 1
		const orgID = 1
		const siteName = "site-A"
		const netName = "network-1"

		netRepo := &mocks.NetRepo{}
		siteRepo := &mocks.SiteRepo{}

		netRepo.On("Get", uint(netID)).Return(
			&db.Network{Model: gorm.Model{ID: netID},
				Name:        netName,
				OrgID:       1,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("GetByNetwork", uint(orgID)).Return(
			[]db.Site{
				db.Site{Model: gorm.Model{ID: netID},
					Name:        siteName,
					NetworkID:   1,
					Deactivated: false,
				}}, nil).Once()

		s := NewNetworkServer(netRepo, nil, siteRepo, nil)
		netResp, err := s.GetSiteByNetwork(context.TODO(),
			&pb.GetSiteByNetworkRequest{NetworkID: netID})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, uint64(netID), netResp.GetSites()[0].GetId())
		assert.Equal(t, uint64(netID), netResp.NetworkID)
		netRepo.AssertExpectations(t)
	})
}
