package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/network/mocks"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"

	"github.com/lib/pq"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"gorm.io/gorm"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"
)

const OrgName = "testorg"

func TestNetworkServer_Get(t *testing.T) {
	t.Run("Network found", func(t *testing.T) {
		const netName = "network-1"
		var netId = uuid.NewV4()

		networks := pq.StringArray{"Verizon"}
		countries := pq.StringArray{"USA"}

		netRepo := &mocks.NetRepo{}
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo.On("Get", netId).Return(
			&db.Network{Id: netId,
				Name:             netName,
				OrgId:            uuid.NewV4(),
				AllowedCountries: countries,
				AllowedNetworks:  networks,
				Deactivated:      false,
			}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, nil, nil, msgcRepo, "")
		netResp, err := s.Get(context.TODO(), &pb.GetRequest{
			NetworkId: netId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netId.String(), netResp.GetNetwork().GetId())
		assert.Equal(t, netName, netResp.Network.Name)
		assert.Equal(t, len(networks), len(netResp.Network.AllowedNetworks))
		assert.Equal(t, len(countries), len(netResp.Network.AllowedCountries))
		netRepo.AssertExpectations(t)
	})

	t.Run("Network not found", func(t *testing.T) {
		var netId = uuid.NewV4()
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", netId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, nil, nil, msgcRepo, "")
		netResp, err := s.Get(context.TODO(), &pb.GetRequest{
			NetworkId: netId.String()})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetByName(t *testing.T) {
	t.Run("Org and Network found", func(t *testing.T) {
		var netId = uuid.NewV4()
		const orgName = "org-1"
		const netName = "network-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", orgName, netName).Return(
			&db.Network{Id: netId,
				Name:        netName,
				OrgId:       uuid.NewV4(),
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, nil, nil, msgcRepo, "")
		netResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{
			Name: netName, OrgName: orgName})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netId.String(), netResp.GetNetwork().GetId())
		assert.Equal(t, netName, netResp.Network.Name)
		netRepo.AssertExpectations(t)
	})

	t.Run("Org or Network not found", func(t *testing.T) {
		const orgName = "org-1"
		const netName = "network-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", orgName, netName).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, nil, nil, msgcRepo, "")
		netResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{
			Name: netName, OrgName: orgName})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetByOrg(t *testing.T) {
	t.Run("Org found", func(t *testing.T) {
		var netId = uuid.NewV4()
		var orgId = uuid.NewV4()
		const netName = "network-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		orgRepo := &mocks.OrgRepo{}

		netRepo.On("GetByOrg", orgId).Return(
			[]db.Network{
				{Id: netId,
					Name:        netName,
					OrgId:       orgId,
					Deactivated: false,
				}}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, orgRepo, nil, nil, msgcRepo, "")
		netResp, err := s.GetByOrg(context.TODO(),
			&pb.GetByOrgRequest{OrgId: orgId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netId.String(), netResp.GetNetworks()[0].GetId())
		assert.Equal(t, orgId.String(), netResp.OrgId)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_Delete(t *testing.T) {
	t.Run("Org and Network exist", func(t *testing.T) {
		const orgName = "org-1"
		orgId := uuid.NewV4()
		const netName = "network-1"
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		orgRepo := &mocks.OrgRepo{}
		orgRepo.On("GetByName", orgName).Return(&db.Org{Id: orgId}, nil).Once()
		netRepo.On("Delete", orgName, netName).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, &pb.DeleteRequest{
			Name:    netName,
			OrgName: orgName,
		}).Return(nil).Once()
		netRepo.On("GetNetworkCount", orgId).Return(int64(2), nil).Once()
		s := NewNetworkServer(OrgName, netRepo, orgRepo, nil, nil, msgclientRepo, "")
		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Name: netName, OrgName: orgName})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Network does not exist", func(t *testing.T) {
		// const netId = 1
		const orgName = "org-1"
		orgId := uuid.NewV4()
		const netName = "network-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		orgRepo := &mocks.OrgRepo{}
		orgRepo.On("GetByName", orgName).Return(&db.Org{Id: orgId}, nil).Once()
		netRepo.On("Delete", orgName, netName).Return(gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, orgRepo, nil, nil, msgcRepo, "")
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
		var netId = uuid.NewV4()
		var orgId = uuid.NewV4()
		const netName = "network-1"
		const siteName = "site-A"
		msgclientRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		siteRepo := &mocks.SiteRepo{}

		site := &db.Site{
			Name:      siteName,
			NetworkId: netId,
		}

		netRepo.On("Get", netId).Return(
			&db.Network{Id: netId,
				Name:        netName,
				OrgId:       orgId,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("Add", site, mock.Anything).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, &pb.AddSiteRequest{
			NetworkId: netId.String(),
			SiteName:  siteName,
		}).Return(nil).Once()
		siteRepo.On("GetSiteCount", netId).Return(int64(2), nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, siteRepo, nil, msgclientRepo, "")

		// Act
		res, err := s.AddSite(context.TODO(), &pb.AddSiteRequest{
			NetworkId: netId.String(),
			SiteName:  siteName,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, siteName, res.Site.Name)
		assert.Equal(t, netId.String(), res.Site.NetworkId)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetSite(t *testing.T) {
	t.Run("Site exists", func(t *testing.T) {
		var siteId = uuid.NewV4()
		const siteName = "site-A"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		siteRepo := &mocks.SiteRepo{}

		siteRepo.On("Get", siteId).Return(
			&db.Site{Id: siteId,
				Name:        siteName,
				NetworkId:   uuid.NewV4(),
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(OrgName, nil, nil, siteRepo, nil, msgcRepo, "")
		netResp, err := s.GetSite(context.TODO(), &pb.GetSiteRequest{
			SiteId: siteId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, siteId.String(), netResp.GetSite().GetId())
		assert.Equal(t, siteName, netResp.GetSite().GetName())
		siteRepo.AssertExpectations(t)
	})

	t.Run("Site not found", func(t *testing.T) {
		var siteId = uuid.NewV4()
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		siteRepo := &mocks.SiteRepo{}

		siteRepo.On("Get", siteId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, nil, nil, siteRepo, nil, msgcRepo, "")
		netResp, err := s.GetSite(context.TODO(), &pb.GetSiteRequest{
			SiteId: fmt.Sprint(siteId)})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		siteRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetSiteByName(t *testing.T) {
	t.Run("Site exists", func(t *testing.T) {
		var siteId = uuid.NewV4()
		var netId = uuid.NewV4()
		var orgId = uuid.NewV4()
		const siteName = "site-A"
		const netName = "net-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		siteRepo := &mocks.SiteRepo{}
		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", netId).Return(
			&db.Network{Id: netId,
				Name:        netName,
				OrgId:       orgId,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("GetByName", netId, siteName).Return(
			&db.Site{Id: siteId,
				Name:        siteName,
				NetworkId:   netId,
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, siteRepo, nil, msgcRepo, "")
		netResp, err := s.GetSiteByName(context.TODO(), &pb.GetSiteByNameRequest{
			NetworkId: netId.String(), SiteName: siteName})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, siteId.String(), netResp.GetSite().GetId())
		assert.Equal(t, siteName, netResp.GetSite().GetName())
		siteRepo.AssertExpectations(t)
	})

	t.Run("Site not found", func(t *testing.T) {
		var netId = uuid.NewV4()
		var orgId = uuid.NewV4()
		const siteName = "site-A"
		const netName = "net-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		siteRepo := &mocks.SiteRepo{}
		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", netId).Return(
			&db.Network{Id: netId,
				Name:        netName,
				OrgId:       orgId,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("GetByName", netId, siteName).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, siteRepo, nil, msgcRepo, "")
		netResp, err := s.GetSiteByName(context.TODO(), &pb.GetSiteByNameRequest{
			NetworkId: netId.String(), SiteName: siteName})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		siteRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetSiteByNetwork(t *testing.T) {
	t.Run("Network found", func(t *testing.T) {
		var netId = uuid.NewV4()
		var orgId = uuid.NewV4()
		const siteName = "site-A"
		const netName = "network-1"
		msgcRepo := &mbmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		siteRepo := &mocks.SiteRepo{}

		netRepo.On("Get", netId).Return(
			&db.Network{Id: netId,
				Name:        netName,
				OrgId:       orgId,
				Deactivated: false,
			}, nil).Once()

		siteRepo.On("GetByNetwork", netId).Return(
			[]db.Site{
				{Id: netId,
					Name:        siteName,
					NetworkId:   netId,
					Deactivated: false,
				}}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, siteRepo, nil, msgcRepo, "")
		netResp, err := s.GetSitesByNetwork(context.TODO(),
			&pb.GetSitesByNetworkRequest{NetworkId: netId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netId.String(), netResp.GetSites()[0].GetId())
		netRepo.AssertExpectations(t)
	})
}
