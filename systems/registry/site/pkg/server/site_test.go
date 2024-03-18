/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/site/mocks"
	"github.com/ukama/ukama/systems/registry/site/pkg/db"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	pb "github.com/ukama/ukama/systems/registry/site/pb/gen"
)
 
 const OrgName = "testorg"
 
 func TestNetworkServer_AddSite(t *testing.T) {
	 t.Run("Network exists", func(t *testing.T) {
		 // Arrange
		 var netId = uuid.NewV4()
		 const netName = "network-1"
		 const siteName = "site-A"

		 msgclientRepo := &cmocks.MsgBusServiceClient{}
 
		 siteRepo := &mocks.SiteRepo{}
 		 netRepo := &mocks.NetworkClientProvider{}

		 site := &db.Site{
			ID:            uuid.NewV4(),
			Name:          siteName,
			NetworkID:     netId,
			BackhaulID:    uuid.NewV4(),
			PowerID:       uuid.NewV4(),
			AccessID:      uuid.NewV4(),
			SwitchID:      uuid.NewV4(),
			IsDeactivated: false,
			Latitude:      40.7128,   // Dummy latitude
			Longitude:     -74.0060,  // Dummy longitude
			InstallDate:   time.Now(), // Current time as install date
		 }
 
		//  netRepo.On("Get", netId).Return(
		// 	 &db.Network{Id: netId,
		// 		 Name:        netName,
		// 		 OrgId:       orgId,
		// 		 Deactivated: false,
		// 	 }, nil).Once()
 
		now := time.Now()
		timestamp, err := ptypes.TimestampProto(now)
		if err != nil {
			// Handle error
		}
				siteRepo.On("Add", site, mock.Anything).Return(nil).Once()
		 msgclientRepo.On("PublishRequest", mock.Anything, &pb.AddRequest{
			Name:          siteName,
			NetworkId:     uuid.NewV4().String(),
			BackhaulId:    uuid.NewV4().String(),
			PowerId:       uuid.NewV4().String(),
			AccessId:      uuid.NewV4().String(),
			SwitchId:      uuid.NewV4().String(),
			IsDeactivated: false,
			Latitude:      40.7128,   // Dummy latitude
			Longitude:     -74.0060,  // Dummy longitude
			InstallDate:   timestamp, // Current time as install date
		 }).Return(nil).Once()
		 siteRepo.On("GetSiteCount", netId).Return(int64(2), nil).Once()
 
		 s := NewsiteServer(OrgName, siteRepo, msgclientRepo,netRepo, "")
 
		 // Act
		 res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name:          siteName,
			NetworkId:     uuid.NewV4().String(),
			BackhaulId:    uuid.NewV4().String(),
			PowerId:       uuid.NewV4().String(),
			AccessId:      uuid.NewV4().String(),
			SwitchId:      uuid.NewV4().String(),
			IsDeactivated: false,
			Latitude:      40.7128,   // Dummy latitude
			Longitude:     -74.0060,  // Dummy longitude
			InstallDate:   timestamp, // Current time as install date
		 })
 
		 // Assert
		 assert.NoError(t, err)
		 assert.NotNil(t, res)
		 assert.Equal(t, siteName, res.Site.Name)
		 assert.Equal(t, netId.String(), res.Site.NetworkId)
		 netRepo.AssertExpectations(t)
	 })
 }
 
//  func TestNetworkServer_GetSite(t *testing.T) {
// 	 t.Run("Site exists", func(t *testing.T) {
// 		 var siteId = uuid.NewV4()
// 		 const siteName = "site-A"
// 		 msgcRepo := &cmocks.MsgBusServiceClient{}
 
// 		 siteRepo := &mocks.SiteRepo{}
 
// 		 siteRepo.On("Get", siteId).Return(
// 			 &db.Site{Id: siteId,
// 				 Name:        siteName,
// 				 NetworkId:   uuid.NewV4(),
// 				 Deactivated: false,
// 			 }, nil).Once()
 
// 		 s := NewNetworkServer(OrgName, nil, nil, siteRepo, nil, msgcRepo, "", "", "", "")
// 		 netResp, err := s.GetSite(context.TODO(), &pb.GetSiteRequest{
// 			 SiteId: siteId.String()})
 
// 		 assert.NoError(t, err)
// 		 assert.NotNil(t, netResp)
// 		 assert.Equal(t, siteId.String(), netResp.GetSite().GetId())
// 		 assert.Equal(t, siteName, netResp.GetSite().GetName())
// 		 siteRepo.AssertExpectations(t)
// 	 })
 
// 	 t.Run("Site not found", func(t *testing.T) {
// 		 var siteId = uuid.NewV4()
// 		 msgcRepo := &cmocks.MsgBusServiceClient{}
 
// 		 siteRepo := &mocks.SiteRepo{}
 
// 		 siteRepo.On("Get", siteId).Return(nil, gorm.ErrRecordNotFound).Once()
 
// 		 s := NewNetworkServer(OrgName, nil, nil, siteRepo, nil, msgcRepo, "", "", "", "")
// 		 netResp, err := s.GetSite(context.TODO(), &pb.GetSiteRequest{
// 			 SiteId: fmt.Sprint(siteId)})
 
// 		 assert.Error(t, err)
// 		 assert.Nil(t, netResp)
// 		 siteRepo.AssertExpectations(t)
// 	 })
//  }
 
//  func TestNetworkServer_GetSiteByName(t *testing.T) {
// 	 t.Run("Site exists", func(t *testing.T) {
// 		 var siteId = uuid.NewV4()
// 		 var netId = uuid.NewV4()
// 		 var orgId = uuid.NewV4()
// 		 const siteName = "site-A"
// 		 const netName = "net-1"
// 		 msgcRepo := &cmocks.MsgBusServiceClient{}
 
// 		 siteRepo := &mocks.SiteRepo{}
// 		 netRepo := &mocks.NetRepo{}
 
// 		 netRepo.On("Get", netId).Return(
// 			 &db.Network{Id: netId,
// 				 Name:        netName,
// 				 OrgId:       orgId,
// 				 Deactivated: false,
// 			 }, nil).Once()
 
// 		 siteRepo.On("GetByName", netId, siteName).Return(
// 			 &db.Site{Id: siteId,
// 				 Name:        siteName,
// 				 NetworkId:   netId,
// 				 Deactivated: false,
// 			 }, nil).Once()
 
// 		 s := NewNetworkServer(OrgName, netRepo, nil, siteRepo, nil, msgcRepo, "", "", "", "")
// 		 netResp, err := s.GetSiteByName(context.TODO(), &pb.GetSiteByNameRequest{
// 			 NetworkId: netId.String(), SiteName: siteName})
 
// 		 assert.NoError(t, err)
// 		 assert.NotNil(t, netResp)
// 		 assert.Equal(t, siteId.String(), netResp.GetSite().GetId())
// 		 assert.Equal(t, siteName, netResp.GetSite().GetName())
// 		 siteRepo.AssertExpectations(t)
// 	 })
 
// 	 t.Run("Site not found", func(t *testing.T) {
// 		 var netId = uuid.NewV4()
// 		 var orgId = uuid.NewV4()
// 		 const siteName = "site-A"
// 		 const netName = "net-1"
// 		 msgcRepo := &cmocks.MsgBusServiceClient{}
 
// 		 siteRepo := &mocks.SiteRepo{}
// 		 netRepo := &mocks.NetRepo{}
 
// 		 netRepo.On("Get", netId).Return(
// 			 &db.Network{Id: netId,
// 				 Name:        netName,
// 				 OrgId:       orgId,
// 				 Deactivated: false,
// 			 }, nil).Once()
 
// 		 siteRepo.On("GetByName", netId, siteName).Return(nil, gorm.ErrRecordNotFound).Once()
 
// 		 s := NewNetworkServer(OrgName, netRepo, nil, siteRepo, nil, msgcRepo, "", "", "", "")
// 		 netResp, err := s.GetSiteByName(context.TODO(), &pb.GetSiteByNameRequest{
// 			 NetworkId: netId.String(), SiteName: siteName})
 
// 		 assert.Error(t, err)
// 		 assert.Nil(t, netResp)
// 		 siteRepo.AssertExpectations(t)
// 	 })
//  }
 
//  func TestNetworkServer_GetSiteByNetwork(t *testing.T) {
// 	 t.Run("Network found", func(t *testing.T) {
// 		 var netId = uuid.NewV4()
// 		 var orgId = uuid.NewV4()
// 		 const siteName = "site-A"
// 		 const netName = "network-1"
// 		 msgcRepo := &cmocks.MsgBusServiceClient{}
 
// 		 netRepo := &mocks.NetRepo{}
// 		 siteRepo := &mocks.SiteRepo{}
 
// 		 netRepo.On("Get", netId).Return(
// 			 &db.Network{Id: netId,
// 				 Name:        netName,
// 				 OrgId:       orgId,
// 				 Deactivated: false,
// 			 }, nil).Once()
 
// 		 siteRepo.On("GetByNetwork", netId).Return(
// 			 []db.Site{
// 				 {Id: netId,
// 					 Name:        siteName,
// 					 NetworkId:   netId,
// 					 Deactivated: false,
// 				 }}, nil).Once()
 
// 		 s := NewNetworkServer(OrgName, netRepo, nil, siteRepo, nil, msgcRepo, "", "", "", "")
// 		 netResp, err := s.GetSitesByNetwork(context.TODO(),
// 			 &pb.GetSitesByNetworkRequest{NetworkId: netId.String()})
 
// 		 assert.NoError(t, err)
// 		 assert.NotNil(t, netResp)
// 		 assert.Equal(t, netId.String(), netResp.GetSites()[0].GetId())
// 		 netRepo.AssertExpectations(t)
// 	 })
//  }