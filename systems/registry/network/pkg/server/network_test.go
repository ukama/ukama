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
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"
)

const OrgName = "testorg"

func TestNetworkServer_Add(t *testing.T) {
	t.Run("networkSuccess", func(t *testing.T) {
		// Arrange
		netRepo := &mocks.NetRepo{}
		orgId := uuid.NewV4()

		const netName = "network-1"


		var netCount = int64(1)

		msgbusClient := &cmocks.MsgBusServiceClient{}
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		network := &db.Network{
			Name:       netName,
			SyncStatus: ukama.StatusTypePending,
		}

		netRepo.On("GetNetworkCount").Return(netCount, nil).Once()
		netRepo.On("Add", network, mock.Anything).Return(nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgbusClient, "", "", "", "",orgId.String())

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name:    netName,
		})
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, netName, res.Network.Name)
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
				AllowedCountries: countries,
				AllowedNetworks:  networks,
				Deactivated:      false,
			}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil,msgcRepo, "", "", "", "","")
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

		s := NewNetworkServer(OrgName, netRepo, nil,msgcRepo, "" , "", "", "", "")
		netResp, err := s.Get(context.TODO(), &pb.GetRequest{
			NetworkId: netId.String()})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetByName(t *testing.T) {
	t.Run("Network found", func(t *testing.T) {
		var netId = uuid.NewV4()
		const netName = "network-1"
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", netName).Return(
			&db.Network{Id: netId,
				Name:        netName,
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo,"" , "", "", "", "")
		netResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{
			Name: netName})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netId.String(), netResp.GetNetwork().GetId())
		assert.Equal(t, netName, netResp.Network.Name)
		netRepo.AssertExpectations(t)
	})

	t.Run("Network not found", func(t *testing.T) {
		const netName = "network-1"
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", netName).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo,"" , "", "", "", "")
		netResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{
			Name: netName})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetAll(t *testing.T) {
	t.Run("networks found", func(t *testing.T) {
		var netId = uuid.NewV4()
		const netName = "network-1"
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetAll").Return(
			[]db.Network{
				{Id: netId,
					Name:        netName,
					Deactivated: false,
				}}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "","")
		netResp, err := s.GetAll(context.TODO(),
			&pb.GetNetworksRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netId.String(), netResp.GetNetworks()[0].GetId())
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_Delete(t *testing.T) {
	t.Run("Network exist", func(t *testing.T) {
		orgId := uuid.NewV4()
		const netName = "network-1"
		msgclientRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}	
		netRepo.On("Delete", netName).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, &pb.DeleteRequest{
			Name:    netName,
		}).Return(nil).Once()
		netRepo.On("GetNetworkCount").Return(int64(2), nil).Once()
		s := NewNetworkServer(OrgName, netRepo, nil, msgclientRepo, "", "", "", "",orgId.String())
		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Name: netName})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Network does not exist", func(t *testing.T) {
		orgId := uuid.NewV4()

		const netName = "network-1"
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		
		netRepo.On("Delete",netName).Return(gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo,nil, msgcRepo,"", "", "", "", orgId.String())
		netResp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			Name: netName})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

