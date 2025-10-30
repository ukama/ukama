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

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	cnucl "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/registry/network/mocks"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	pb "github.com/ukama/ukama/systems/registry/network/pb/gen"
)

const (
	OrgName = "testorg"

	// Test network data
	TestNetworkName1 = "network-1"
	TestNetworkName2 = "network-2"

	// Test organization data
	TestOrgName = "org-1"

	// Test network properties
	TestNetworkProvider = "Verizon"
	TestCountry         = "USA"

	// Test counts
	TestNetworkCount1 = int64(1)
	TestNetworkCount2 = int64(2)
)

func TestNetworkServer_Add(t *testing.T) {
	t.Run("networkSuccess", func(t *testing.T) {
		// Arrange
		netRepo := &mocks.NetRepo{}
		orgId := uuid.NewV4()
		orgClient := &cmocks.OrgClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		netName := TestNetworkName1
		netCount := TestNetworkCount1

		network := &db.Network{
			Name:       netName,
			SyncStatus: ukama.StatusTypePending,
		}
		orgClient.On("Get", orgName).Return(
			&cnucl.OrgInfo{
				Id:            orgId.String(),
				Name:          orgName,
				IsDeactivated: false,
			}, nil).Once()
		netRepo.On("Add", network, mock.Anything).Return(nil).Once()

		netRepo.On("GetNetworkCount").Return(netCount, nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewNetworkServer(orgName, netRepo, orgClient, msgbusClient, "", "", "", "", orgId.String())

		// Act
		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name: netName,
		})
		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, netName, res.Network.Name)
		netRepo.AssertExpectations(t)
	})

	t.Run("org client error", func(t *testing.T) {
		netRepo := &mocks.NetRepo{}
		orgId := uuid.NewV4()
		orgClient := &cmocks.OrgClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		netName := TestNetworkName1

		orgClient.On("Get", orgName).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(orgName, netRepo, orgClient, msgbusClient, "", "", "", "", orgId.String())

		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name: netName,
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		netRepo.AssertExpectations(t)
	})

	t.Run("invalid network name", func(t *testing.T) {
		netRepo := &mocks.NetRepo{}
		orgId := uuid.NewV4()
		orgClient := &cmocks.OrgClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		invalidNetName := "Invalid_Network_Name_With_Uppercase_And_Underscores_That_Exceeds_253_Characters_Limit_And_Contains_Invalid_Characters_Like_Spaces_And_Special_Symbols_That_Are_Not_Allowed_In_DNS_Labels_According_To_RFC_1123_Standards_Which_Require_Lowercase_Alphanumeric_Characters_And_Hyphens_Only_For_Valid_DNS_Label_Names"

		orgClient.On("Get", orgName).Return(
			&cnucl.OrgInfo{
				Id:            orgId.String(),
				Name:          orgName,
				IsDeactivated: false,
			}, nil).Once()

		s := NewNetworkServer(orgName, netRepo, orgClient, msgbusClient, "", "", "", "", orgId.String())

		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name: invalidNetName,
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		netRepo.AssertExpectations(t)
	})

	t.Run("deactivated org", func(t *testing.T) {
		netRepo := &mocks.NetRepo{}
		orgId := uuid.NewV4()
		orgClient := &cmocks.OrgClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		netName := TestNetworkName1

		orgClient.On("Get", orgName).Return(
			&cnucl.OrgInfo{
				Id:            orgId.String(),
				Name:          orgName,
				IsDeactivated: true,
			}, nil).Once()

		s := NewNetworkServer(orgName, netRepo, orgClient, msgbusClient, "", "", "", "", orgId.String())

		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name: netName,
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		netRepo.AssertExpectations(t)
	})

	t.Run("database error during add", func(t *testing.T) {
		netRepo := &mocks.NetRepo{}
		orgId := uuid.NewV4()
		orgClient := &cmocks.OrgClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		netName := TestNetworkName1

		network := &db.Network{
			Name:       netName,
			SyncStatus: ukama.StatusTypePending,
		}

		orgClient.On("Get", orgName).Return(
			&cnucl.OrgInfo{
				Id:            orgId.String(),
				Name:          orgName,
				IsDeactivated: false,
			}, nil).Once()
		netRepo.On("Add", network, mock.Anything).Return(gorm.ErrInvalidDB).Once()

		s := NewNetworkServer(orgName, netRepo, orgClient, msgbusClient, "", "", "", "", orgId.String())

		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name: netName,
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		netRepo.AssertExpectations(t)
	})

	t.Run("add with message bus failure", func(t *testing.T) {
		netRepo := &mocks.NetRepo{}
		orgId := uuid.NewV4()
		orgClient := &cmocks.OrgClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		netName := TestNetworkName1
		netCount := TestNetworkCount1

		network := &db.Network{
			Name:       netName,
			SyncStatus: ukama.StatusTypePending,
		}

		orgClient.On("Get", orgName).Return(
			&cnucl.OrgInfo{
				Id:            orgId.String(),
				Name:          orgName,
				IsDeactivated: false,
			}, nil).Once()
		netRepo.On("Add", network, mock.Anything).Return(nil).Once()
		netRepo.On("GetNetworkCount").Return(netCount, nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(gorm.ErrInvalidDB).Once()

		s := NewNetworkServer(orgName, netRepo, orgClient, msgbusClient, "", "", "", "", orgId.String())

		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name: netName,
		})

		// Add should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, netName, res.Network.Name)
		netRepo.AssertExpectations(t)
	})

	t.Run("add with nil message bus", func(t *testing.T) {
		netRepo := &mocks.NetRepo{}
		orgId := uuid.NewV4()
		orgClient := &cmocks.OrgClient{}

		orgName := TestOrgName
		netName := TestNetworkName1
		netCount := TestNetworkCount1

		network := &db.Network{
			Name:       netName,
			SyncStatus: ukama.StatusTypePending,
		}

		orgClient.On("Get", orgName).Return(
			&cnucl.OrgInfo{
				Id:            orgId.String(),
				Name:          orgName,
				IsDeactivated: false,
			}, nil).Once()
		netRepo.On("Add", network, mock.Anything).Return(nil).Once()
		netRepo.On("GetNetworkCount").Return(netCount, nil).Once()

		s := NewNetworkServer(orgName, netRepo, orgClient, nil, "", "", "", "", orgId.String())

		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name: netName,
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, netName, res.Network.Name)
		netRepo.AssertExpectations(t)
	})

	t.Run("add with network count error", func(t *testing.T) {
		netRepo := &mocks.NetRepo{}
		orgId := uuid.NewV4()
		orgClient := &cmocks.OrgClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		netName := TestNetworkName1

		network := &db.Network{
			Name:       netName,
			SyncStatus: ukama.StatusTypePending,
		}

		orgClient.On("Get", orgName).Return(
			&cnucl.OrgInfo{
				Id:            orgId.String(),
				Name:          orgName,
				IsDeactivated: false,
			}, nil).Once()
		netRepo.On("Add", network, mock.Anything).Return(nil).Once()
		netRepo.On("GetNetworkCount").Return(int64(0), gorm.ErrInvalidDB).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewNetworkServer(orgName, netRepo, orgClient, msgbusClient, "", "", "", "", orgId.String())

		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name: netName,
		})

		// Add should still succeed even if network count fails
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, netName, res.Network.Name)
		netRepo.AssertExpectations(t)
	})

	t.Run("add with callback function that sets network ID", func(t *testing.T) {
		netRepo := &mocks.NetRepo{}
		orgId := uuid.NewV4()
		orgClient := &cmocks.OrgClient{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		orgName := TestOrgName
		netName := TestNetworkName1
		netCount := TestNetworkCount1

		network := &db.Network{
			Name:       netName,
			SyncStatus: ukama.StatusTypePending,
		}

		// Create a custom mock that actually executes the callback
		netRepo.On("Add", network, mock.MatchedBy(func(callback func(*db.Network, *gorm.DB) error) bool {
			// Execute the callback to ensure the network ID is set
			err := callback(network, nil)
			return err == nil
		})).Return(nil).Once()

		orgClient.On("Get", orgName).Return(
			&cnucl.OrgInfo{
				Id:            orgId.String(),
				Name:          orgName,
				IsDeactivated: false,
			}, nil).Once()
		netRepo.On("GetNetworkCount").Return(netCount, nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewNetworkServer(orgName, netRepo, orgClient, msgbusClient, "", "", "", "", orgId.String())

		res, err := s.Add(context.TODO(), &pb.AddRequest{
			Name: netName,
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, netName, res.Network.Name)
		assert.NotEqual(t, "", res.Network.Id) // Verify ID was set
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_SetDefault(t *testing.T) {
	t.Run("Set default network", func(t *testing.T) {
		var netId = uuid.NewV4()

		netRepo := &mocks.NetRepo{}
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo.On("SetDefault", netId, true).Return(
			&db.Network{}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.SetDefault(context.TODO(), &pb.SetDefaultRequest{
			NetworkId: netId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Invalid UUID format", func(t *testing.T) {
		netRepo := &mocks.NetRepo{}
		msgcRepo := &cmocks.MsgBusServiceClient{}

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.SetDefault(context.TODO(), &pb.SetDefaultRequest{
			NetworkId: "invalid-uuid",
		})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Database error during set default", func(t *testing.T) {
		netId := uuid.NewV4()
		netRepo := &mocks.NetRepo{}
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo.On("SetDefault", netId, true).Return(
			nil, gorm.ErrInvalidDB).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.SetDefault(context.TODO(), &pb.SetDefaultRequest{
			NetworkId: netId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Network not found during set default", func(t *testing.T) {
		netId := uuid.NewV4()
		netRepo := &mocks.NetRepo{}
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo.On("SetDefault", netId, true).Return(
			nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.SetDefault(context.TODO(), &pb.SetDefaultRequest{
			NetworkId: netId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetDefault(t *testing.T) {
	t.Run("Get default network", func(t *testing.T) {
		var netId = uuid.NewV4()

		netName := TestNetworkName1
		networks := pq.StringArray{TestNetworkProvider}
		countries := pq.StringArray{TestCountry}

		netRepo := &mocks.NetRepo{}
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo.On("GetDefault").Return(
			&db.Network{Id: netId,
				Name:             netName,
				AllowedCountries: countries,
				AllowedNetworks:  networks,
				Deactivated:      false,
				IsDefault:        true,
			}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.GetDefault(context.TODO(), &pb.GetDefaultRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netId.String(), netResp.GetNetwork().GetId())
		assert.Equal(t, true, netResp.GetNetwork().GetIsDefault())
		assert.Equal(t, netName, netResp.Network.Name)
		assert.Equal(t, len(networks), len(netResp.Network.AllowedNetworks))
		assert.Equal(t, len(countries), len(netResp.Network.AllowedCountries))
		netRepo.AssertExpectations(t)
	})
	t.Run("No default network", func(t *testing.T) {

		netRepo := &mocks.NetRepo{}
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo.On("GetDefault").Return(
			nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.GetDefault(context.TODO(), &pb.GetDefaultRequest{})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_Get(t *testing.T) {
	t.Run("Network found", func(t *testing.T) {
		netName := TestNetworkName1
		netId := uuid.NewV4()

		networks := pq.StringArray{TestNetworkProvider}
		countries := pq.StringArray{TestCountry}

		netRepo := &mocks.NetRepo{}
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo.On("Get", netId).Return(
			&db.Network{Id: netId,
				Name:             netName,
				AllowedCountries: countries,
				AllowedNetworks:  networks,
				Deactivated:      false,
			}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
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

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.Get(context.TODO(), &pb.GetRequest{
			NetworkId: netId.String()})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Invalid UUID format", func(t *testing.T) {
		msgcRepo := &cmocks.MsgBusServiceClient{}
		netRepo := &mocks.NetRepo{}

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.Get(context.TODO(), &pb.GetRequest{
			NetworkId: "invalid-uuid"})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		netId := uuid.NewV4()
		msgcRepo := &cmocks.MsgBusServiceClient{}
		netRepo := &mocks.NetRepo{}

		netRepo.On("Get", netId).Return(nil, gorm.ErrInvalidDB).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.Get(context.TODO(), &pb.GetRequest{
			NetworkId: netId.String()})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetByName(t *testing.T) {
	t.Run("Network found", func(t *testing.T) {
		netId := uuid.NewV4()
		netName := TestNetworkName1
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", netName).Return(
			&db.Network{Id: netId,
				Name:        netName,
				Deactivated: false,
			}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{
			Name: netName})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netId.String(), netResp.GetNetwork().GetId())
		assert.Equal(t, netName, netResp.Network.Name)
		netRepo.AssertExpectations(t)
	})

	t.Run("Network not found", func(t *testing.T) {
		netName := TestNetworkName1
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetByName", netName).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.GetByName(context.TODO(), &pb.GetByNameRequest{
			Name: netName})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_GetAll(t *testing.T) {
	t.Run("networks found", func(t *testing.T) {
		netId := uuid.NewV4()
		netName := TestNetworkName1
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetAll").Return(
			[]db.Network{
				{Id: netId,
					Name:        netName,
					Deactivated: false,
				}}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.GetAll(context.TODO(),
			&pb.GetNetworksRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Equal(t, netId.String(), netResp.GetNetworks()[0].GetId())
		netRepo.AssertExpectations(t)
	})

	t.Run("no networks found", func(t *testing.T) {
		msgcRepo := &cmocks.MsgBusServiceClient{}
		netRepo := &mocks.NetRepo{}

		netRepo.On("GetAll").Return([]db.Network{}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.GetAll(context.TODO(), &pb.GetNetworksRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Empty(t, netResp.GetNetworks())
		netRepo.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		msgcRepo := &cmocks.MsgBusServiceClient{}
		netRepo := &mocks.NetRepo{}

		netRepo.On("GetAll").Return(nil, gorm.ErrInvalidDB).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.GetAll(context.TODO(), &pb.GetNetworksRequest{})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})

	t.Run("multiple networks found", func(t *testing.T) {
		netId1 := uuid.NewV4()
		netId2 := uuid.NewV4()
		netName1 := TestNetworkName1
		netName2 := TestNetworkName2
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("GetAll").Return(
			[]db.Network{
				{Id: netId1, Name: netName1, Deactivated: false},
				{Id: netId2, Name: netName2, Deactivated: true},
			}, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", "")
		netResp, err := s.GetAll(context.TODO(), &pb.GetNetworksRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, netResp)
		assert.Len(t, netResp.GetNetworks(), 2)
		assert.Equal(t, netId1.String(), netResp.GetNetworks()[0].GetId())
		assert.Equal(t, netId2.String(), netResp.GetNetworks()[1].GetId())
		assert.Equal(t, netName1, netResp.GetNetworks()[0].GetName())
		assert.Equal(t, netName2, netResp.GetNetworks()[1].GetName())
		assert.False(t, netResp.GetNetworks()[0].GetIsDeactivated())
		assert.True(t, netResp.GetNetworks()[1].GetIsDeactivated())
		netRepo.AssertExpectations(t)
	})
}

func TestNetworkServer_Delete(t *testing.T) {
	t.Run("Network exist", func(t *testing.T) {
		orgId := uuid.NewV4()
		netId := uuid.NewV4()

		msgclientRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}
		netRepo.On("Delete", netId).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, &epb.EventNetworkDelete{
			Id:    netId.String(),
			OrgId: orgId.String(),
		}).Return(nil).Once()
		netRepo.On("GetNetworkCount").Return(TestNetworkCount2, nil).Once()
		s := NewNetworkServer(OrgName, netRepo, nil, msgclientRepo, "", "", "", "", orgId.String())
		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			NetworkId: netId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Network does not exist", func(t *testing.T) {
		orgId := uuid.NewV4()

		netId := uuid.NewV4()
		msgcRepo := &cmocks.MsgBusServiceClient{}

		netRepo := &mocks.NetRepo{}

		netRepo.On("Delete", netId).Return(gorm.ErrRecordNotFound).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", orgId.String())
		netResp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			NetworkId: netId.String()})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Invalid UUID format", func(t *testing.T) {
		orgId := uuid.NewV4()
		msgcRepo := &cmocks.MsgBusServiceClient{}
		netRepo := &mocks.NetRepo{}

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", orgId.String())
		netResp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			NetworkId: "invalid-uuid"})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Database error during delete", func(t *testing.T) {
		orgId := uuid.NewV4()
		netId := uuid.NewV4()
		msgcRepo := &cmocks.MsgBusServiceClient{}
		netRepo := &mocks.NetRepo{}

		netRepo.On("Delete", netId).Return(gorm.ErrInvalidDB).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgcRepo, "", "", "", "", orgId.String())
		netResp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			NetworkId: netId.String()})

		assert.Error(t, err)
		assert.Nil(t, netResp)
		netRepo.AssertExpectations(t)
	})

	t.Run("Delete with message bus failure", func(t *testing.T) {
		orgId := uuid.NewV4()
		netId := uuid.NewV4()

		msgclientRepo := &cmocks.MsgBusServiceClient{}
		netRepo := &mocks.NetRepo{}

		netRepo.On("Delete", netId).Return(nil).Once()
		msgclientRepo.On("PublishRequest", mock.Anything, &epb.EventNetworkDelete{
			Id:    netId.String(),
			OrgId: orgId.String(),
		}).Return(gorm.ErrInvalidDB).Once()
		netRepo.On("GetNetworkCount").Return(TestNetworkCount2, nil).Once()

		s := NewNetworkServer(OrgName, netRepo, nil, msgclientRepo, "", "", "", "", orgId.String())
		resp, err := s.Delete(context.TODO(), &pb.DeleteRequest{
			NetworkId: netId.String()})

		// Delete should still succeed even if message bus fails
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		netRepo.AssertExpectations(t)
	})
}
