package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	uuid "github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/network/mocks"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"
	"gorm.io/gorm"
)

func TestNetworkServer_Get(t *testing.T) {
	t.Run("Network found", func(t *testing.T) {
		var netID = uuid.NewV4()
		const netName = "network-1"

		netRepo := &mocks.NetRepo{}
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo.On("Get", netID).Return(
			&db.Network{ID: netID,
				Name:        netName,
				OrgID:       uuid.NewV4(),
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(netRepo, nil, nil, nil, msgcRepo)
		netResp, err := s.Get(context.TODO(), &pb.GetRequest{
			NetworkID: netID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netID.String(), netResp.GetNetwork().GetId())
		assert.Equal(t, netName, netResp.Network.Name)
		netRepo.AssertExpectations(t)
	})

	t.Run("Network not found", func(t *testing.T) {
		var netID = uuid.NewV4()
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", netID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(netRepo, nil, nil, nil, msgcRepo)
		netResp, err := s.Get(context.TODO(), &pb.GetRequest{
			NetworkID: netID.String()})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetByName(t *testing.T) {
	t.Run("Org and Network found", func(t *testing.T) {
		var netID = uuid.NewV4()
		const orgName = "org-1"
		const netName = "network-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", orgName, netName).Return(
			&db.Network{ID: netID,
				Name:        netName,
				OrgID:       uuid.NewV4(),
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(netRepo, nil, nil, nil, msgcRepo)
		netResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{
			Name: netName, OrgName: orgName})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netID.String(), netResp.GetNetwork().GetId())
		assert.Equal(t, netName, netResp.Network.Name)
		netRepo.AssertExpectations(t)
	})

	t.Run("Org or Network not found", func(t *testing.T) {
		const orgName = "org-1"
		const netName = "network-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", orgName, netName).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(netRepo, nil, nil, nil, msgcRepo)
		netResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{
			Name: netName, OrgName: orgName})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetByOrg(t *testing.T) {
	t.Run("Org found", func(t *testing.T) {
		var netID = uuid.NewV4()
		var orgID = uuid.NewV4()
		const netName = "network-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		orgRepo := &mocks.OrgRepo{}

		netRepo.On("GetByOrg", orgID).Return(
			[]db.Network{
				{ID: netID,
					Name:        netName,
					OrgID:       orgID,
					Deactivated: false,
				}}, nil).Once()

		s := NewNetworkServer(netRepo, orgRepo, nil, nil, msgcRepo)
		netResp, err := s.GetByOrg(context.TODO(),
			&pb.GetByOrgRequest{OrgID: orgID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netID.String(), netResp.GetNetworks()[0].GetId())
		assert.Equal(t, orgID.String(), netResp.OrgID)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_Delete(t *testing.T) {
	t.Run("Org and Network exist", func(t *testing.T) {
		const orgName = "org-1"
		const netName = "network-1"
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("Delete", orgName, netName).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, &pb.DeleteRequest{
			Name:    netName,
			OrgName: orgName,
		}).Return(nil).Once()
		s := NewNetworkServer(netRepo, nil, nil, nil, msgclientRepo)
		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Name: netName, OrgName: orgName})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Network does not exist", func(t *testing.T) {
		// const netID = 1
		const orgName = "org-1"
		const netName = "network-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("Delete", orgName, netName).Return(gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(netRepo, nil, nil, nil, msgcRepo)
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
		var netID = uuid.NewV4()
		var orgID = uuid.NewV4()
		const netName = "network-1"
		const siteName = "site-A"
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		siteRepo := &mocks.SiteRepo{}

		site := &db.Site{
			Name:      siteName,
			NetworkID: netID,
		}

		netRepo.On("Get", netID).Return(
			&db.Network{ID: netID,
				Name:        netName,
				OrgID:       orgID,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("Add", site, mock.Anything).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, &pb.AddSiteRequest{
			NetworkID: netID.String(),
			SiteName:  siteName,
		}).Return(nil).Once()
		s := NewNetworkServer(netRepo, nil, siteRepo, nil, msgclientRepo)

		// Act
		res, err := s.AddSite(context.TODO(), &pb.AddSiteRequest{
			NetworkID: netID.String(),
			SiteName:  siteName,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, siteName, res.Site.Name)
		assert.Equal(t, netID.String(), res.Site.NetworkID)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetSite(t *testing.T) {
	t.Run("Site exists", func(t *testing.T) {
		var siteID = uuid.NewV4()
		const siteName = "site-A"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		siteRepo := &mocks.SiteRepo{}

		siteRepo.On("Get", siteID).Return(
			&db.Site{ID: siteID,
				Name:        siteName,
				NetworkID:   uuid.NewV4(),
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(nil, nil, siteRepo, nil, msgcRepo)
		netResp, err := s.GetSite(context.TODO(), &pb.GetSiteRequest{
			SiteID: siteID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, siteID.String(), netResp.GetSite().GetId())
		assert.Equal(t, siteName, netResp.GetSite().GetName())
		siteRepo.AssertExpectations(t)
	})

	t.Run("Site not found", func(t *testing.T) {
		var siteID = uuid.NewV4()
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		siteRepo := &mocks.SiteRepo{}

		siteRepo.On("Get", siteID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(nil, nil, siteRepo, nil, msgcRepo)
		netResp, err := s.GetSite(context.TODO(), &pb.GetSiteRequest{
			SiteID: fmt.Sprint(siteID)})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		siteRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetSiteByName(t *testing.T) {
	t.Run("Site exists", func(t *testing.T) {
		var siteID = uuid.NewV4()
		var netID = uuid.NewV4()
		var orgID = uuid.NewV4()
		const siteName = "site-A"
		const netName = "net-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		siteRepo := &mocks.SiteRepo{}
		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", netID).Return(
			&db.Network{ID: netID,
				Name:        netName,
				OrgID:       orgID,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("GetByName", netID, siteName).Return(
			&db.Site{ID: siteID,
				Name:        siteName,
				NetworkID:   netID,
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(netRepo, nil, siteRepo, nil, msgcRepo)
		netResp, err := s.GetSiteByName(context.TODO(), &pb.GetSiteByNameRequest{
			NetworkID: netID.String(), SiteName: siteName})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, siteID.String(), netResp.GetSite().GetId())
		assert.Equal(t, siteName, netResp.GetSite().GetName())
		siteRepo.AssertExpectations(t)
	})

	t.Run("Site not found", func(t *testing.T) {
		var netID = uuid.NewV4()
		var orgID = uuid.NewV4()
		const siteName = "site-A"
		const netName = "net-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		siteRepo := &mocks.SiteRepo{}
		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", netID).Return(
			&db.Network{ID: netID,
				Name:        netName,
				OrgID:       orgID,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("GetByName", netID, siteName).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(netRepo, nil, siteRepo, nil, msgcRepo)
		netResp, err := s.GetSiteByName(context.TODO(), &pb.GetSiteByNameRequest{
			NetworkID: netID.String(), SiteName: siteName})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		siteRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetSiteByNetwork(t *testing.T) {
	t.Run("Network found", func(t *testing.T) {
		var netID = uuid.NewV4()
		var orgID = uuid.NewV4()
		const siteName = "site-A"
		const netName = "network-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		siteRepo := &mocks.SiteRepo{}

		netRepo.On("Get", netID).Return(
			&db.Network{ID: netID,
				Name:        netName,
				OrgID:       orgID,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("GetByNetwork", netID).Return(
			[]db.Site{
				{ID: netID,
					Name:        siteName,
					NetworkID:   netID,
					Deactivated: false,
				}}, nil).Once()

		s := NewNetworkServer(netRepo, nil, siteRepo, nil, msgcRepo)
		netResp, err := s.GetSitesByNetwork(context.TODO(),
			&pb.GetSitesByNetworkRequest{NetworkID: netID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netID.String(), netResp.GetSites()[0].GetId())
		netRepo.AssertExpectations(t)
	})
}
