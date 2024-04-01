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

	"github.com/lib/pq"
	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/network/mocks"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"
)

const OrgName = "testorg"

func TestNetworkServer_Add(t *testing.T) {
	t.Run("OrgFound", func(t *testing.T) {
		// Arrange
		netRepo := &mocks.NetRepo{}
		orgRepo := &mocks.OrgRepo{}

		const netName = "network-1"
		const orgName = "org-1"

		var orgId = uuid.NewV4()
		var netCount = int64(1)

		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		network := &db.Network{
			Name:       netName,
			OrgId:      orgId,
			SyncStatus: ukama.StatusTypePending,
		}

		orgRepo.On("GetByName", orgName).Return(
			&db.Org{
				Id:          orgId,
				Name:        orgName,
				Deactivated: false,
			}, nil).Once()

		netRepo.On("GetNetworkCount", mock.Anything).Return(netCount, nil).Once()
		netRepo.On("Add", network, mock.Anything).Return(nil).Once()

		s := NewNetworkServer(OrgName, netRepo, orgRepo, nil, msgbusClient, "", "", "", "")

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name:    netName,
			OrgName: orgName,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, netName, res.Network.Name)
		assert.Equal(t, orgId.String(), res.Network.OrgId)
		netRepo.AssertExpectations(t)
	})

	t.Run("OrgNotFound", func(t *testing.T) {
		// Arrange
		netRepo := &mocks.NetRepo{}
		orgRepo := &mocks.OrgRepo{}
		orgClient := &cmocks.OrgClient{}

		const netName = "network-1"
		const orgName = "org-1"

		var orgId = uuid.NewV4()
		var netCount = int64(1)

		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		netRepo.On("GetNetworkCount", mock.Anything).Return(netCount, nil).Once()

		org := &db.Org{
			Id:          orgId,
			Name:        orgName,
			Deactivated: false,
		}

		network := &db.Network{
			Name:       netName,
			OrgId:      orgId,
			SyncStatus: ukama.StatusTypePending,
		}

		orgRepo.On("GetByName", orgName).Return(nil, gorm.ErrRecordNotFound).Once()

		orgClient.On("Get", orgName).Return(
			&cnucl.OrgInfo{
				Id:            orgId.String(),
				Name:          orgName,
				IsDeactivated: false,
			}, nil).Once()

		orgRepo.On("Add", org, mock.Anything).Return(nil).Once()

		netRepo.On("Add", network, mock.Anything).Return(nil).Once()

		s := NewNetworkServer(OrgName, netRepo, orgRepo, orgClient, msgbusClient, "", "", "", "")

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name:    netName,
			OrgName: orgName,
		})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, netName, res.Network.Name)
		assert.Equal(t, orgId.String(), res.Network.OrgId)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_Get(t *testing.T) {
	t.Run("Network found", func(t *testing.T) {
		const netName = "network-1"
		var netId = uuid.NewV4()

		networks := pq.StringArray{"Verizon"}
		countries := pq.StringArray{"USA"}

		netRepo := &mocks.NetRepo{}
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo.On("Get", netId).Return(
			&db.Network{Id: netId,
				Name:             netName,
				OrgId:            uuid.NewV4(),
				AllowedCountries: countries,
				AllowedNetworks:  networks,
				Deactivated:      false,
			}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, nil, msgcRepo, "", "", "", "")
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
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", netId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, nil, msgcRepo, "", "", "", "")
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
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", orgName, netName).Return(
			&db.Network{Id: netId,
				Name:        netName,
				OrgId:       uuid.NewV4(),
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, nil, msgcRepo, "", "", "", "")
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
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", orgName, netName).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, nil, msgcRepo, "", "", "", "")
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
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		orgRepo := &mocks.OrgRepo{}

		netRepo.On("GetByOrg", orgId).Return(
			[]db.Network{
				{Id: netId,
					Name:        netName,
					OrgId:       orgId,
					Deactivated: false,
				}}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, orgRepo, nil, msgcRepo, "", "", "", "")
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
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		orgRepo := &mocks.OrgRepo{}
		orgRepo.On("GetByName", orgName).Return(&db.Org{Id: orgId}, nil).Once()
		netRepo.On("Delete", orgName, netName).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, &pb.DeleteRequest{
			Name:    netName,
			OrgName: orgName,
		}).Return(nil).Once()
		netRepo.On("GetNetworkCount", orgId).Return(int64(2), nil).Once()
		s := NewNetworkServer(OrgName, netRepo, orgRepo, nil, msgclientRepo, "", "", "", "")
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
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		orgRepo := &mocks.OrgRepo{}
		orgRepo.On("GetByName", orgName).Return(&db.Org{Id: orgId}, nil).Once()
		netRepo.On("Delete", orgName, netName).Return(gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, orgRepo, nil, msgcRepo, "", "", "", "")
		netResp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Name: netName, OrgName: orgName})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

