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
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/mocks"

	log "github.com/sirupsen/logrus"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	cnotif "github.com/ukama/ukama/systems/common/rest/client/notification"
	cnuc "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	subspb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	subsmocks "github.com/ukama/ukama/systems/subscriber/registry/pb/gen/mocks"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/adapters"
	db "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	splpb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	splmocks "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen/mocks"
)

const (
	testIccid = "890000-this-is-a-test-iccid"
	OrgName   = "testorg"
	// testIccid = "890000000000000001234"
	simTypeOperator = "operator_data"
	simTypeTest     = "test"
	cdrType         = "data"
	from            = "2022-12-01T00:00:00Z"
	to              = "2023-12-01T00:00:00Z"
	bytesUsed       = 28901234567
	cost            = 100.99
)

func TestSimManagerServer_GetSim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var simID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		sim := simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("GetSim", mock.Anything,
			sim.Iccid).Return(nil, nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimId: simID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		assert.Equal(t, simID.String(), simResp.GetSim().GetId())
		assert.Equal(t, false, simResp.GetSim().IsPhysical)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		var simID = uuid.Nil

		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimId: simID.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDInvalid", func(t *testing.T) {
		var simID = "1"

		simRepo := &mocks.SimRepo{}

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimId: simID})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_GetUsages(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var simID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		sim := simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("GetUsages", mock.Anything,
			sim.Iccid, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(map[string]any{}, map[string]any{}, nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		usagesResp, err := s.GetUsages(context.TODO(), &pb.UsageRequest{
			SimId: simID.String(),
			Type:  ukama.CdrTypeData.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, usagesResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		var simID = uuid.Nil

		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).Return(nil, gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		usagesResp, err := s.GetUsages(context.TODO(), &pb.UsageRequest{
			SimId: simID.String(),
			Type:  ukama.CdrTypeData.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, usagesResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDInvalid", func(t *testing.T) {
		var simID = "1"

		simRepo := &mocks.SimRepo{}

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		usagesResp, err := s.GetUsages(context.TODO(), &pb.UsageRequest{
			SimId: simID})

		assert.Error(t, err)
		assert.Nil(t, usagesResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimTypeSupported", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		sType := ukama.ParseSimType(simTypeOperator)

		agentAdapter := agentFactory.On("GetAgentAdapter", sType).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("GetUsages", mock.Anything,
			"", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(map[string]any{}, map[string]any{}, nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		usagesResp, err := s.GetUsages(context.TODO(), &pb.UsageRequest{
			SimType: simTypeOperator,
			Type:    ukama.CdrTypeData.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, usagesResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimTypeNotSupported", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		s := NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		usagesResp, err := s.GetUsages(context.TODO(), &pb.UsageRequest{
			SimType: "lol",
		})

		assert.Error(t, err)
		assert.Nil(t, usagesResp)
		simRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_GetSimsBySubscriber(t *testing.T) {
	t.Run("SubscriberFound", func(t *testing.T) {
		var simID = uuid.NewV4()
		var subscriberID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetBySubscriber", subscriberID).Return(
			[]db.Sim{
				{Id: simID,
					SubscriberId: subscriberID,
					IsPhysical:   false,
				}}, nil).Once()

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSimsBySubscriber(context.TODO(),
			&pb.GetSimsBySubscriberRequest{SubscriberId: subscriberID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		assert.Equal(t, simID.String(), simResp.GetSims()[0].GetId())
		assert.Equal(t, subscriberID.String(), simResp.SubscriberId)
		simRepo.AssertExpectations(t)
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		var subscriberID = uuid.Nil

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetBySubscriber", subscriberID).Return(
			nil, errors.New("some unexpected error has occured")).Once()

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSimsBySubscriber(context.TODO(), &pb.GetSimsBySubscriberRequest{
			SubscriberId: subscriberID.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SubscriberUUIDInvalid", func(t *testing.T) {
		var subscriberID = "1"

		simRepo := &mocks.SimRepo{}

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		simResp, err := s.GetSimsBySubscriber(context.TODO(), &pb.GetSimsBySubscriberRequest{
			SubscriberId: subscriberID})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_GetSimsByNetwork(t *testing.T) {
	t.Run("NetworkFound", func(t *testing.T) {
		var simID = uuid.NewV4()
		var networkID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetByNetwork", networkID).Return(
			[]db.Sim{
				{Id: simID,
					NetworkId:  networkID,
					IsPhysical: false,
				}}, nil).Once()

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSimsByNetwork(context.TODO(),
			&pb.GetSimsByNetworkRequest{NetworkId: networkID.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		assert.Equal(t, simID.String(), simResp.GetSims()[0].GetId())
		assert.Equal(t, networkID.String(), simResp.NetworkId)
		simRepo.AssertExpectations(t)
	})

	t.Run("SomeUnexpectedErrorAsOccured", func(t *testing.T) {
		var networkID = uuid.Nil

		simRepo := &mocks.SimRepo{}

		simRepo.On("GetByNetwork", networkID).Return(
			nil, errors.New("some unexpected error has occured")).Once()

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSimsByNetwork(context.TODO(), &pb.GetSimsByNetworkRequest{
			NetworkId: networkID.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("NetworkUUIDInvalid", func(t *testing.T) {
		var networkID = "1"

		simRepo := &mocks.SimRepo{}

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSimsByNetwork(context.TODO(), &pb.GetSimsByNetworkRequest{
			NetworkId: networkID})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_GetPackagesForSim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var simId = uuid.NewV4()
		var packageId = uuid.NewV4()

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("GetBySim", simId).Return(
			[]db.Package{
				{Id: packageId,
					SimId:    simId,
					IsActive: false,
				}}, nil).Once()

		s := NewSimManagerServer(OrgName, nil, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.GetPackagesForSim(context.TODO(),
			&pb.GetPackagesForSimRequest{SimId: simId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, packageId.String(), resp.GetPackages()[0].GetId())
		assert.Equal(t, simId.String(), resp.SimId)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SomeUnexpectedErrorAsOccured", func(t *testing.T) {
		var simId = uuid.Nil

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("GetBySim", simId).Return(
			nil, errors.New("some unexpected error has occured")).Once()

		s := NewSimManagerServer(OrgName, nil, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.GetPackagesForSim(context.TODO(), &pb.GetPackagesForSimRequest{
			SimId: simId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDInvalid", func(t *testing.T) {
		var simID = "1"

		packageRepo := &mocks.PackageRepo{}

		s := NewSimManagerServer(OrgName, nil, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.GetPackagesForSim(context.TODO(), &pb.GetPackagesForSimRequest{
			SimId: simID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_AllocateSim(t *testing.T) {
	t.Run("SimAndPackageFound", func(t *testing.T) {
		var subscriberID = uuid.NewV4()
		var networkID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var OrgId = uuid.NewV4()

		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}

		msgbusClient := &cmocks.MsgBusServiceClient{}

		simPoolService := &mocks.SimPoolClientProvider{}
		subscriberService := &mocks.SubscriberRegistryClientProvider{}

		packageClient := &cmocks.PackageClient{}
		netClient := &cmocks.NetworkClient{}
		orgClient := &cmocks.OrgClient{}
		userClient := &cmocks.UserClient{}
		mailerClient := &cmocks.MailerClient{}

		agentFactory := &mocks.AgentFactory{}

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberID.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberID.String(),
					NetworkId:    networkID.String(),
				},
			}, nil).Once()

		packageClient.On("Get", packageID.String()).
			Return(&cdplan.PackageInfo{
				IsActive: true,
				Duration: 3600,
				SimType:  simTypeTest,
			}, nil).Once()

		netClient.On("Get", networkID.String()).
			Return(&creg.NetworkInfo{
				IsDeactivated: false,
				TrafficPolicy: 0,
			}, nil).Once()

		simPoolClient := simPoolService.On("GetClient").
			Return(&splmocks.SimServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*splmocks.SimServiceClient)

		simPoolResp := simPoolClient.On("Get", mock.Anything,
			&splpb.GetRequest{IsPhysicalSim: false,
				SimType: simTypeTest,
			}).
			Return(&splpb.GetResponse{
				Sim: &splpb.Sim{
					Iccid:      testIccid,
					IsPhysical: false,
					SimType:    simTypeTest,
				},
			}, nil).Once().
			ReturnArguments.Get(0).(*splpb.GetResponse)

		sim := &db.Sim{
			SubscriberId: subscriberID,
			NetworkId:    networkID,
			Iccid:        testIccid,
			Type:         ukama.SimTypeTest,
			Status:       ukama.SimStatusInactive,
			IsPhysical:   simPoolResp.Sim.IsPhysical,
			SyncStatus:   ukama.StatusTypePending,
		}

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("BindSim", mock.Anything,
			sim.Iccid).Return(nil, nil).Once()

		simRepo.On("Add", sim,
			mock.Anything).Return(nil).Once()

		pkg := &db.Package{
			SimId:     sim.Id,
			PackageId: packageID,
			IsActive:  true,
		}

		packageRepo.On("Add", pkg,
			mock.Anything).Return(nil).Once()

		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		simRepo.On("GetSimMetrics").Return(int64(0), int64(0), int64(0), int64(0), nil).Once()
		mailerClient.On("SendEmail", mock.MatchedBy(func(req cnotif.SendEmailReq) bool {
			return req.To[0] == "test@example.com"
		})).Return(nil).Once()

		orgClient.On("Get", OrgName).
			Return(&cnuc.OrgInfo{
				Name:  OrgName,
				Owner: OrgId.String(),
			}, nil).Once()

		userClient.On("GetById", OrgId.String()).
			Return(&cnuc.UserInfo{
				Name: "test-user",
			},
				nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, packageRepo, agentFactory,
			packageClient, subscriberService, simPoolService, "", msgbusClient, OrgId.String(), "", mailerClient, netClient, orgClient, userClient)

		log.Info("SimManagerServer: ", s)
		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberID.String(),
			NetworkId:    networkID.String(),
			PackageId:    packageID.String(),
			SimType:      simTypeTest,
			SimToken:     "",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		simRepo.AssertExpectations(t)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)

		simPoolService.AssertExpectations(t)
		simPoolClient.AssertExpectations(t)

		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)

		msgbusClient.AssertExpectations(t)
	})

	t.Run("SubscriberNotRegisteredOnProvidedNetwork", func(t *testing.T) {
		var subscriberID = uuid.NewV4()
		var networkID = uuid.NewV4()
		var packageID = uuid.NewV4()

		subscriberService := &mocks.SubscriberRegistryClientProvider{}

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberID.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberID.String(),
					NetworkId:    uuid.NewV4().String(),
				},
			}, nil).Once()

		s := NewSimManagerServer(OrgName, nil, nil, nil, nil, subscriberService,
			nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberID.String(),
			NetworkId:    networkID.String(),
			PackageId:    packageID.String(),
			SimType:      simTypeTest,
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)
	})

	t.Run("OrgPackageNoMoreActive", func(t *testing.T) {
		var subscriberID = uuid.NewV4()
		var networkID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()

		subscriberService := &mocks.SubscriberRegistryClientProvider{}
		packageClient := &cmocks.PackageClient{}

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberID.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberID.String(),
					NetworkId:    networkID.String(),
				},
			}, nil).Once()

		packageClient.On("Get", packageID.String()).
			Return(&cdplan.PackageInfo{
				OrgId:    orgID.String(),
				IsActive: false,
				Duration: 3600,
				SimType:  simTypeTest,
			}, nil).Once()

		s := NewSimManagerServer(OrgName, nil, nil, nil,
			packageClient, subscriberService, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberID.String(),
			NetworkId:    networkID.String(),
			PackageId:    packageID.String(),
			SimType:      simTypeTest,
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)

		packageClient.AssertExpectations(t)
	})

	t.Run("PackageSimtypeAndProvidedSimtypeMismatch", func(t *testing.T) {
		var subscriberID = uuid.NewV4()
		var networkID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()

		subscriberService := &mocks.SubscriberRegistryClientProvider{}
		packageClient := &cmocks.PackageClient{}

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberID.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberID.String(),
					NetworkId:    networkID.String(),
				},
			}, nil).Once()

		packageClient.On("Get", packageID.String()).
			Return(
				&cdplan.PackageInfo{
					OrgId:    orgID.String(),
					IsActive: true,
					Duration: 3600,
					SimType:  ukama.SimTypeUnknown.String(),
				}, nil).Once()

		s := NewSimManagerServer(OrgName, nil, nil, nil,
			packageClient, subscriberService, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberID.String(),
			NetworkId:    networkID.String(),
			PackageId:    packageID.String(),
			SimType:      simTypeTest,
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)

		packageClient.AssertExpectations(t)
	})

}

func TestSimManagerServer_SetActivePackageForSim(t *testing.T) {
	t.Run("SimAndPackageFound", func(t *testing.T) {
		var packageID = uuid.NewV4()
		var simID = uuid.NewV4()

		msgbusClient := &cmocks.MsgBusServiceClient{}
		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}
		simd := &db.Sim{Id: simID,
			IsPhysical: false,
			Status:     ukama.SimStatusActive,
			Type:       ukama.SimTypeTest,
		}
		simRepo.On("Get", simID).
			Return(simd, nil).
			Once()

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    simID,
				EndDate:  time.Now().UTC().AddDate(0, 1, 0), // next month
				IsActive: false,
			}, nil).Once()

		packageRepo.On("Update",
			&db.Package{
				Id:       packageID,
				IsActive: true,
			},
			mock.Anything).Return(nil).Once()

		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		agentAdapter := agentFactory.On("GetAgentAdapter", simd.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("UpdatePackage", mock.Anything,
			mock.MatchedBy(func(a adapters.ReqData) bool {
				return a.Iccid == simd.Iccid
			})).Return(nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, packageRepo,
			agentFactory, nil, nil, nil, "", msgbusClient, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
		agentAdapter.AssertExpectations(t)
	})

	t.Run("SimStatusInvalid", func(t *testing.T) {
		var packageID = uuid.NewV4()
		var simID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				IsPhysical: false,
				Status:     ukama.SimStatusUnknown,
			}, nil).
			Once()

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimIdAndPackageSimIdMismatch", func(t *testing.T) {
		var packageID = uuid.NewV4()
		var simID = uuid.NewV4()

		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    uuid.NewV4(),
				EndDate:  time.Now().UTC().AddDate(0, 1, 0), // next month
				IsActive: false,
			}, nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
	})

	t.Run("PackageAlreadyExpired", func(t *testing.T) {
		var packageID = uuid.NewV4()
		var simID = uuid.NewV4()

		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    simID,
				EndDate:  time.Now().UTC().AddDate(0, -1, 0), // one month ago
				IsActive: false,
			}, nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_RemovePackageForSim(t *testing.T) {
	t.Run("PackageFound", func(t *testing.T) {
		var packageID = uuid.NewV4()
		var simID = uuid.NewV4()
		simRepo := &mocks.SimRepo{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		packageRepo := &mocks.PackageRepo{}

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    simID,
				IsActive: false,
			}, nil).Once()

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Delete", packageID,
			mock.Anything).Return(nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", msgbusClient, "", "", nil, nil, nil, nil)
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageID.String(),
			SimId:     simID.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		packageRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
		simRepo.AssertExpectations(t)
	})

	t.Run("PackageDeleteError", func(t *testing.T) {
		var packageID = uuid.Nil
		var simID = uuid.NewV4()

		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    simID,
				IsActive: false,
			}, nil).Once()

		packageRepo.On("Delete", packageID,
			mock.Anything).Return(gorm.ErrRecordNotFound).Once()

		s := NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageID.String(),
			SimId:     simID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("PackageUUIDInvalid", func(t *testing.T) {
		var packageID = "1"

		packageRepo := &mocks.PackageRepo{}

		s := NewSimManagerServer(OrgName, nil, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageID})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDMismatch", func(t *testing.T) {
		var packageID = uuid.Nil
		var simID = uuid.NewV4()

		packageRepo := &mocks.PackageRepo{}

		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageID).Return(
			&db.Package{Id: packageID,
				SimId:    simID,
				IsActive: false,
			}, nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageID.String(),
			SimId:     uuid.NewV4().String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_AddPackageForSim(t *testing.T) {
	t.Run("SimAndPackageFound", func(t *testing.T) {
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()
		msgbusClient := &cmocks.MsgBusServiceClient{}

		var orgID = uuid.NewV4()
		// startDate := time.Now().UTC().Add(5 * time.Minute).Truncate(time.Second)
		// startDate := time.Now().UTC())
		startDate := time.Now().UTC()

		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		packageClient := &cmocks.PackageClient{}

		sim := simRepo.On("Get", simID, mock.Anything).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		packageClient.On("Get", packageID.String(), mock.Anything).
			Return(&cdplan.PackageInfo{
				Id:       packageID.String(),
				OrgId:    orgID.String(),
				IsActive: true,
				Duration: 1,
				SimType:  simTypeTest,
			}, nil).
			Once()

		pkg := &db.Package{
			SimId:     sim.Id,
			StartDate: startDate.Truncate(time.Second),
			EndDate:   startDate.Add(time.Hour * 24).Truncate(time.Second),
			PackageId: packageID,
			IsActive:  true,
		}

		packageRepo.On("List", sim.Id.String(), mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, uint32(0), true).Return([]db.Package{}, nil).Once()

		packageRepo.On("GetOverlap", pkg).Return([]db.Package{}, nil).Once()
		packageRepo.On("Add", pkg, mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, packageRepo, nil, packageClient,
			nil, nil, "", msgbusClient, orgID.String(), "", nil, nil, nil, nil)

		resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
			StartDate: startDate.Format(time.RFC3339),
		})

		fmt.Printf("Response: %v\n, %v", resp, pkg.EndDate)

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})

	t.Run("PackageStartDateNotValid", func(t *testing.T) {
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()

		s := NewSimManagerServer(OrgName, nil, nil, nil, nil,
			nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
			StartDate: "xxxx/12/xx",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("OrgPackageNoMoreActive", func(t *testing.T) {
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()
		startDate := time.Now().UTC().AddDate(0, 0, 1) // tomorrow

		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		packageClient := &cmocks.PackageClient{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once()

		packageClient.On("Get", packageID.String()).
			Return(&cdplan.PackageInfo{
				Id:       packageID.String(),
				OrgId:    orgID.String(),
				IsActive: false,
				Duration: 3600,
				SimType:  simTypeTest,
			}, nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, packageClient, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
			StartDate: startDate.Format(time.RFC3339),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})

	t.Run("SimSimtypeAndPackageSimtypeMismatch", func(t *testing.T) {
		var simID = uuid.NewV4()
		var packageID = uuid.NewV4()
		var orgID = uuid.NewV4()
		startDate := time.Now().UTC().AddDate(0, 0, 1) // tomorrow

		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		packageClient := &cmocks.PackageClient{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeOperatorData,
				IsPhysical: false,
			}, nil).
			Once()

		packageClient.On("Get", packageID.String()).
			Return(&cdplan.PackageInfo{
				Id:       packageID.String(),
				OrgId:    orgID.String(),
				IsActive: true,
				Duration: 3600,
				SimType:  simTypeTest,
			}, nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, packageRepo, nil, packageClient,
			nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
			SimId:     simID.String(),
			PackageId: packageID.String(),
			StartDate: startDate.Format(time.RFC3339),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})

	// t.Run("PackageValidityPeriodOverlapsWithExistingPackages", func(t *testing.T) {
	// 	var simID = uuid.NewV4()
	// 	var packageID = uuid.NewV4()
	// 	var orgID = uuid.NewV4()
	// 	startDate := time.Now().UTC().AddDate(0, 0, 1) // tomorrow

	// 	simRepo := &mocks.SimRepo{}
	// 	packageRepo := &mocks.PackageRepo{}
	// 	packageClient := &cmocks.PackageClient{}

	// 	sim := simRepo.On("Get", simID).
	// 		Return(&db.Sim{Id: simID,
	// 			Iccid:      testIccid,
	// 			Status:     ukama.SimStatusInactive,
	// 			Type:       ukama.SimTypeTest,
	// 			IsPhysical: false,
	// 		}, nil).
	// 		Once().
	// 		ReturnArguments.Get(0).(*db.Sim)

	// 	pkgInfo := packageClient.On("Get", packageID.String()).
	// 		Return(&cdplan.PackageInfo{
	// 			Id:       packageID.String(),
	// 			OrgId:    orgID.String(),
	// 			IsActive: true,
	// 			Duration: 3600,
	// 			SimType:  simTypeTest,
	// 		}, nil).Once().
	// 		ReturnArguments.Get(0).(*cdplan.PackageInfo)

	// 	pkg := &db.Package{
	// 		SimId:     sim.Id,
	// 		StartDate: startDate,
	// 		EndDate:   startDate.Add(time.Duration(pkgInfo.Duration)),
	// 		PackageId: packageID,
	// 		IsActive:  false,
	// 	}

	// 	packageRepo.On("GetOverlap", pkg).
	// 		Return([]db.Package{
	// 			{}, {},
	// 		}, nil).Once()

	// 	s := NewSimManagerServer(OrgName, simRepo, packageRepo, nil, packageClient,
	// 		nil, nil, "", nil, orgID.String(), "", nil, nil)

	// 	resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
	// 		SimId:     simID.String(),
	// 		PackageId: packageID.String(),
	// 		StartDate: startDate.Format(time.RFC3339),
	// 	})

	// 	assert.Error(t, err)
	// 	assert.Nil(t, resp)

	// 	simRepo.AssertExpectations(t)
	// 	packageRepo.AssertExpectations(t)
	// 	packageClient.AssertExpectations(t)
	// })
}

func TestSimManagerServer_DeleteSim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		var simID = uuid.NewV4()
		msgbusClient := &cmocks.MsgBusServiceClient{}

		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		sim := simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("TerminateSim", mock.Anything,
			sim.Iccid).Return(nil).Once()

		simRepo.On("Update",
			&db.Sim{
				Id:     sim.Id,
				Status: ukama.SimStatusTerminated,
			},
			mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()
		simRepo.On("GetSimMetrics").Return(int64(0), int64(0), int64(0), int64(0), nil).Once()

		s := NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", msgbusClient, "", "", nil, nil, nil, nil)

		resp, err := s.DeleteSim(context.TODO(), &pb.DeleteSimRequest{
			SimId: simID.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		simRepo.AssertExpectations(t)
		agentFactory.AssertExpectations(t)
	})

	t.Run("SimStatusInvalid", func(t *testing.T) {
		var simID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}

		simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     ukama.SimStatusActive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once()

		s := NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.DeleteSim(context.TODO(), &pb.DeleteSimRequest{
			SimId: simID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
	})

	t.Run("SimTypeNotSupported", func(t *testing.T) {
		var simID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		sim := simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       100,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, false).
			Once()

		s := NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.DeleteSim(context.TODO(), &pb.DeleteSimRequest{
			SimId: simID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
	})

	t.Run("SimAgentFailToTerminate", func(t *testing.T) {
		var simID = uuid.NewV4()

		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		sim := simRepo.On("Get", simID).
			Return(&db.Sim{Id: simID,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*db.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("TerminateSim", mock.Anything,
			sim.Iccid).Return(errors.New("anyError")).Once()

		s := NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.DeleteSim(context.TODO(), &pb.DeleteSimRequest{
			SimId: simID.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
	})
}
