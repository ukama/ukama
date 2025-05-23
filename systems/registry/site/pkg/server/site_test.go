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

	"github.com/tj/assert"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/site/mocks"
	pb "github.com/ukama/ukama/systems/registry/site/pb/gen"
	"github.com/ukama/ukama/systems/registry/site/pkg/db"
	"gorm.io/gorm"
)

const OrgName = "Ukama"

func TestSiteService_Get(t *testing.T) {
	siteRepo := &mocks.SiteRepo{}
	msgclientRepo := &cmocks.MsgBusServiceClient{}
	netRepo := &mocks.NetworkClientProvider{}

	s := NewSiteServer(OrgName, siteRepo, msgclientRepo, netRepo, "", nil)

	t.Run("SiteFound", func(t *testing.T) {
		siteId := uuid.NewV4()

		siteRepo.On("Get", siteId).Return(&db.Site{
			Id: siteId,
		}, nil)

		uResp, err := s.Get(context.TODO(), &pb.GetRequest{SiteId: siteId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, uResp)

		assert.NoError(t, err)
		assert.Equal(t, siteId.String(), uResp.GetSite().Id)
		siteRepo.AssertExpectations(t)
	})

	t.Run("SiteNotFound", func(t *testing.T) {
		siteId := uuid.NewV4()

		siteRepo.On("Get", siteId).Return(nil, gorm.ErrRecordNotFound).Once()

		uResp, err := s.Get(context.TODO(), &pb.GetRequest{SiteId: siteId.String()})

		assert.Error(t, err)
		assert.Nil(t, uResp)
		siteRepo.AssertExpectations(t)
	})
}

func TestSiteService_List(t *testing.T) {
	siteRepo := &mocks.SiteRepo{}
	msgclientRepo := &cmocks.MsgBusServiceClient{}
	netRepo := &mocks.NetworkClientProvider{}

	s := NewSiteServer(OrgName, siteRepo, msgclientRepo, netRepo, "", nil)

	t.Run("ValidRequest", func(t *testing.T) {
		netId := uuid.NewV4()

		mockSites := []*db.Site{
			{
				Id:        uuid.NewV4(),
				NetworkId: netId,
			},
		}

		var mockSitesConverted []db.Site
		for _, site := range mockSites {
			mockSitesConverted = append(mockSitesConverted, *site)
		}

		siteRepo.On("List", &netId, false).Return(mockSitesConverted, nil)

		req := &pb.ListRequest{
			NetworkId:     netId.String(),
			IsDeactivated: false,
		}

		resp, err := s.List(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, len(mockSites), len(resp.Sites))
		siteRepo.AssertExpectations(t)
	})
}
