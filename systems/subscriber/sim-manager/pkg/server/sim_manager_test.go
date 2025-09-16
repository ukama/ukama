/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/mocks"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/server"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	cnotif "github.com/ukama/ukama/systems/common/rest/client/notification"
	cnuc "github.com/ukama/ukama/systems/common/rest/client/nucleus"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	subspb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	subsmocks "github.com/ukama/ukama/systems/subscriber/registry/pb/gen/mocks"
	pb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	splpb "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen"
	splmocks "github.com/ukama/ukama/systems/subscriber/sim-pool/pb/gen/mocks"
)

const (
	testIccid       = "890000-this-is-a-test-iccid"
	testImsi        = "890000-this-is-a-test-iccid"
	simId           = "e044081b-fbbe-45e9-8f78-0f9c0f112977"
	OrgName         = "testorg"
	orgId           = "592f7a8e-f318-4d3a-aab8-8d4187cde7f9"
	simTypeOperator = "operator_data"
	simTypeTest     = "test"
	cdrType         = "data"
	from            = "2022-12-01T00:00:00Z"
	to              = "2023-12-01T00:00:00Z"
	bytesUsed       = 28901234567
	cost            = 100.99
)

func TestSimManagerServer_ListPackagesForSim(t *testing.T) {
	resp := make([]sims.Package, 1)

	t.Run("SimPackagesFound", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()
		dataplanId := uuid.NewV4()

		pckgResp := sims.Package{
			Id:        packageId,
			SimId:     simId,
			IsActive:  true,
			AsExpired: true,
			PackageId: dataplanId,
		}

		resp[0] = pckgResp

		packageRepo.On("List", simId.String(), dataplanId.String(), from, to,
			from, to, true, true, uint32(0), false).Return(resp, nil)

		s := server.NewSimManagerServer(OrgName, nil, packageRepo, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		list, err := s.ListPackagesForSim(context.TODO(), &pb.ListPackagesForSimRequest{
			SimId:         simId.String(),
			FromStartDate: from,
			ToStartDate:   to,
			FromEndDate:   from,
			ToEndDate:     to,
			IsActive:      true,
			AsExpired:     true,
			DataPlanId:    dataplanId.String(),
			Count:         uint32(0),
			Sort:          false,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
	})

	t.Run("SimIdNotValid", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		dataplanId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, nil, packageRepo, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		list, err := s.ListPackagesForSim(context.TODO(), &pb.ListPackagesForSimRequest{
			SimId:         "lol",
			FromStartDate: from,
			ToStartDate:   to,
			FromEndDate:   from,
			ToEndDate:     to,
			IsActive:      true,
			AsExpired:     true,
			DataPlanId:    dataplanId.String(),
			Count:         uint32(0),
			Sort:          false,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("DataplanIdNotValid", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, nil, packageRepo, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		list, err := s.ListPackagesForSim(context.TODO(), &pb.ListPackagesForSimRequest{
			SimId:         simId.String(),
			FromStartDate: from,
			ToStartDate:   to,
			FromEndDate:   from,
			ToEndDate:     to,
			IsActive:      true,
			AsExpired:     true,
			DataPlanId:    "lol",
			Count:         uint32(0),
			Sort:          false,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("FromStartDateNotValid", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()
		dataplanId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, nil, packageRepo, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		list, err := s.ListPackagesForSim(context.TODO(), &pb.ListPackagesForSimRequest{
			SimId:         simId.String(),
			FromStartDate: "lol",
			ToStartDate:   to,
			FromEndDate:   from,
			ToEndDate:     to,
			IsActive:      true,
			AsExpired:     true,
			DataPlanId:    dataplanId.String(),
			Count:         uint32(0),
			Sort:          false,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("ToStartDateNotValid", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()
		dataplanId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, nil, packageRepo, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		list, err := s.ListPackagesForSim(context.TODO(), &pb.ListPackagesForSimRequest{
			SimId:         simId.String(),
			FromStartDate: from,
			ToStartDate:   "lol",
			FromEndDate:   from,
			ToEndDate:     to,
			IsActive:      true,
			AsExpired:     true,
			DataPlanId:    dataplanId.String(),
			Count:         uint32(0),
			Sort:          false,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("FromEndDateNotValid", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()
		dataplanId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, nil, packageRepo, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		list, err := s.ListPackagesForSim(context.TODO(), &pb.ListPackagesForSimRequest{
			SimId:         simId.String(),
			FromStartDate: from,
			ToStartDate:   to,
			FromEndDate:   "lol",
			ToEndDate:     to,
			IsActive:      true,
			AsExpired:     true,
			DataPlanId:    dataplanId.String(),
			Count:         uint32(0),
			Sort:          false,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("ToDateNotValid", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()
		dataplanId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, nil, packageRepo, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		list, err := s.ListPackagesForSim(context.TODO(), &pb.ListPackagesForSimRequest{
			SimId:         simId.String(),
			FromStartDate: from,
			ToStartDate:   to,
			FromEndDate:   from,
			ToEndDate:     "lol",
			IsActive:      true,
			AsExpired:     true,
			DataPlanId:    dataplanId.String(),
			Count:         uint32(0),
			Sort:          false,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("ToEndDateNotValid", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()
		dataplanId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, nil, packageRepo, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		list, err := s.ListPackagesForSim(context.TODO(), &pb.ListPackagesForSimRequest{
			SimId:         simId.String(),
			FromStartDate: from,
			ToStartDate:   to,
			FromEndDate:   "lol",
			ToEndDate:     to,
			IsActive:      true,
			AsExpired:     true,
			DataPlanId:    dataplanId.String(),
			Count:         uint32(0),
			Sort:          false,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("ListError", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		packageRepo.On("List", simId, "", "", "", "", "", false, false,
			uint32(0), false).Return(nil, errors.New("package list for sim error"))

		s := server.NewSimManagerServer(OrgName, nil, packageRepo, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		list, err := s.ListPackagesForSim(context.TODO(), &pb.ListPackagesForSimRequest{
			SimId: simId,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

}

func TestSimManagerServer_ListSims(t *testing.T) {
	resp := make([]sims.Sim, 1)

	t.Run("SimFound", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()
		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		simResp := sims.Sim{
			Id:           simId,
			SubscriberId: subscriberId,
			NetworkId:    networkId,
			Iccid:        testIccid,
			Status:       ukama.SimStatusInactive,
			Type:         ukama.SimTypeTest,
			IsPhysical:   true,
		}

		resp[0] = simResp

		simRepo.On("List", testIccid, testImsi, subscriberId.String(), networkId.String(),
			ukama.SimTypeTest, ukama.SimStatusInactive, uint32(0), true,
			uint32(0), false).Return(resp, nil)

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		list, err := s.ListSims(context.TODO(), &pb.ListSimsRequest{
			Iccid:         testIccid,
			Imsi:          testImsi,
			SubscriberId:  subscriberId.String(),
			NetworkId:     networkId.String(),
			SimType:       ukama.SimTypeTest.String(),
			SimStatus:     ukama.SimStatusInactive.String(),
			TrafficPolicy: uint32(0),
			IsPhysical:    true,
			Count:         uint32(0),
			Sort:          false,
		})

		assert.NoError(t, err)
		assert.NotNil(t, list)
		assertList(t, list, resp)
	})

	t.Run("SubscriberNotValid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		networkId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		list, err := s.ListSims(context.TODO(), &pb.ListSimsRequest{
			Iccid:         testIccid,
			Imsi:          testImsi,
			SubscriberId:  "lol",
			NetworkId:     networkId.String(),
			SimType:       ukama.SimTypeTest.String(),
			SimStatus:     ukama.SimStatusInactive.String(),
			TrafficPolicy: uint32(0),
			IsPhysical:    true,
			Count:         uint32(0),
			Sort:          false,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("NetworkNotValid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		subscriberId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		list, err := s.ListSims(context.TODO(), &pb.ListSimsRequest{
			Iccid:         testIccid,
			Imsi:          testImsi,
			SubscriberId:  subscriberId.String(),
			NetworkId:     "lol",
			SimType:       ukama.SimTypeTest.String(),
			SimStatus:     ukama.SimStatusInactive.String(),
			TrafficPolicy: uint32(0),
			IsPhysical:    true,
			Count:         uint32(0),
			Sort:          false,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("SimTypeNotValid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		list, err := s.ListSims(context.TODO(), &pb.ListSimsRequest{
			Iccid:         testIccid,
			Imsi:          testImsi,
			SubscriberId:  subscriberId.String(),
			NetworkId:     networkId.String(),
			SimType:       "lol",
			SimStatus:     ukama.SimStatusInactive.String(),
			TrafficPolicy: uint32(0),
			IsPhysical:    true,
			Count:         uint32(0),
			Sort:          false,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("SimStatusNotValid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		list, err := s.ListSims(context.TODO(), &pb.ListSimsRequest{
			Iccid:         testIccid,
			Imsi:          testImsi,
			SubscriberId:  subscriberId.String(),
			NetworkId:     networkId.String(),
			SimType:       ukama.SimTypeTest.String(),
			SimStatus:     "lol",
			TrafficPolicy: uint32(0),
			IsPhysical:    true,
			Count:         uint32(0),
			Sort:          false,
		})

		assert.Error(t, err)
		assert.Nil(t, list)
	})

	t.Run("ListError", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simRepo.On("List", "", "", "", "", ukama.SimTypeUnknown, ukama.SimStatusUnknown, uint32(0), false,
			uint32(0), false).Return(nil, errors.New("sim list error"))

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		list, err := s.ListSims(context.TODO(), &pb.ListSimsRequest{})

		assert.Error(t, err)
		assert.Nil(t, list)
	})
}

func TestSimManagerServer_GetSim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		sim := simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*sims.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("GetSim", mock.Anything,
			sim.Iccid).Return(nil, nil).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimId: simId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		assert.Equal(t, simId.String(), simResp.GetSim().GetId())
		assert.Equal(t, false, simResp.GetSim().IsPhysical)
		simRepo.AssertExpectations(t)
	})

	t.Run("AgentFactoryFailure", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once()

		agentFactory.On("GetAgentAdapter", mock.Anything).
			Return(nil, false)

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimId: simId.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("AgentAdapterFailure", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		sim := simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*sims.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("GetSim", mock.Anything,
			sim.Iccid).Return(nil, errors.New("fail to get sim details from remote agent")).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimId: simId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		simId := uuid.Nil

		simRepo.On("Get", simId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimId: simId.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDInvalid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		simId := "1"

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			SimId: simId})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_GetUsages(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		sim := simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*sims.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("GetUsages", mock.Anything,
			sim.Iccid, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(map[string]any{}, map[string]any{}, nil).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		usagesResp, err := s.GetUsages(context.TODO(), &pb.UsageRequest{
			SimId: simId.String(),
			Type:  ukama.CdrTypeData.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, usagesResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		simId := uuid.Nil

		simRepo.On("Get", simId).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		usagesResp, err := s.GetUsages(context.TODO(), &pb.UsageRequest{
			SimId: simId.String(),
			Type:  ukama.CdrTypeData.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, usagesResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDInvalid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		simId := "1"

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		usagesResp, err := s.GetUsages(context.TODO(), &pb.UsageRequest{
			SimId: simId})

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

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
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

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
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
		simRepo := &mocks.SimRepo{}

		simId := uuid.NewV4()
		subscriberId := uuid.NewV4()

		simRepo.On("GetBySubscriber", subscriberId).Return(
			[]sims.Sim{
				{Id: simId,
					SubscriberId: subscriberId,
					IsPhysical:   false,
				}}, nil).Once()

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSimsBySubscriber(context.TODO(),
			&pb.GetSimsBySubscriberRequest{SubscriberId: subscriberId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		assert.Equal(t, simId.String(), simResp.GetSims()[0].GetId())
		assert.Equal(t, subscriberId.String(), simResp.SubscriberId)
		simRepo.AssertExpectations(t)
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		subscriberId := uuid.Nil

		simRepo.On("GetBySubscriber", subscriberId).Return(
			nil, errors.New("some unexpected error has occurred")).Once()

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSimsBySubscriber(context.TODO(), &pb.GetSimsBySubscriberRequest{
			SubscriberId: subscriberId.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SubscriberUUIDInvalid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		subscriberId := "1"

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		simResp, err := s.GetSimsBySubscriber(context.TODO(), &pb.GetSimsBySubscriberRequest{
			SubscriberId: subscriberId})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_GetSimsByNetwork(t *testing.T) {
	t.Run("NetworkFound", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		simId := uuid.NewV4()
		networkId := uuid.NewV4()

		simRepo.On("GetByNetwork", networkId).Return(
			[]sims.Sim{
				{Id: simId,
					NetworkId:  networkId,
					IsPhysical: false,
				}}, nil).Once()

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSimsByNetwork(context.TODO(),
			&pb.GetSimsByNetworkRequest{NetworkId: networkId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, simResp)
		assert.Equal(t, simId.String(), simResp.GetSims()[0].GetId())
		assert.Equal(t, networkId.String(), simResp.NetworkId)
		simRepo.AssertExpectations(t)
	})

	t.Run("SomeUnexpectedErrorAsOccurred", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		networkId := uuid.Nil

		simRepo.On("GetByNetwork", networkId).Return(
			nil, errors.New("some unexpected error has occurred")).Once()

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSimsByNetwork(context.TODO(), &pb.GetSimsByNetworkRequest{
			NetworkId: networkId.String()})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})

	t.Run("NetworkUUIDInvalid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		networkId := "1"

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)
		simResp, err := s.GetSimsByNetwork(context.TODO(), &pb.GetSimsByNetworkRequest{
			NetworkId: networkId})

		assert.Error(t, err)
		assert.Nil(t, simResp)
		simRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_GetPackagesForSim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		packageRepo.On("GetBySim", simId).Return(
			[]sims.Package{
				{Id: packageId,
					SimId:    simId,
					IsActive: false,
				}}, nil).Once()

		s := server.NewSimManagerServer(OrgName, nil, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.GetPackagesForSim(context.TODO(),
			&pb.GetPackagesForSimRequest{SimId: simId.String()})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, packageId.String(), resp.GetPackages()[0].GetId())
		assert.Equal(t, simId.String(), resp.SimId)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SomeUnexpectedErrorAsOccurred", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}

		simId := uuid.Nil

		packageRepo.On("GetBySim", simId).Return(
			nil, errors.New("some unexpected error has occurred")).Once()

		s := server.NewSimManagerServer(OrgName, nil, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.GetPackagesForSim(context.TODO(), &pb.GetPackagesForSimRequest{
			SimId: simId.String()})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDInvalid", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}

		simId := "1"

		s := server.NewSimManagerServer(OrgName, nil, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.GetPackagesForSim(context.TODO(), &pb.GetPackagesForSimRequest{
			SimId: simId})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_AllocateSim(t *testing.T) {
	t.Run("SimAndPackageFound", func(t *testing.T) {

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

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()
		simPackageIDString := "00000000-0000-0000-0000-000000000000"
		orgId := uuid.NewV4()

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberId.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
					Email:        "test@example.com",
					Name:         "Test User",
				},
			}, nil).Once()

		packageInfo := &cdplan.PackageInfo{
			IsActive:   true,
			Duration:   3600,
			SimType:    simTypeTest,
			DataVolume: 10,
			DataUnit:   "GB",
			Name:       "Test Package",
			Amount:     100,
		}

		packageClient.On("Get", packageId.String()).
			Return(packageInfo, nil).
			Times(1)

		netClient.On("Get", networkId.String()).
			Return(&creg.NetworkInfo{
				IsDeactivated: false,
				TrafficPolicy: 0,
				Name:          "Test Network",
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
					QrCode:     "test-qr-code",
				},
			}, nil).Once().
			ReturnArguments.Get(0).(*splpb.GetResponse)

		sim := &sims.Sim{
			SubscriberId: subscriberId,
			NetworkId:    networkId,
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
			client.AgentRequestData{
				Iccid:        testIccid,
				Imsi:         sim.Imsi,
				SimPackageId: simPackageIDString,
				PackageId:    packageId.String(),
				NetworkId:    sim.NetworkId.String(),
			}).Return(nil, nil).Once()

		simRepo.On("Add", sim,
			mock.Anything).Return(nil).Once()

		pkg := &sims.Package{
			SimId:           sim.Id,
			PackageId:       packageId,
			IsActive:        true,
			DefaultDuration: 3600,
		}

		packageRepo.On("Add", pkg,
			mock.Anything).Return(nil).Once()

		orgClient.On("Get", OrgName).
			Return(&cnuc.OrgInfo{
				Name:  OrgName,
				Owner: orgId.String(),
			}, nil).Once()

		userClient.On("GetById", orgId.String()).
			Return(&cnuc.UserInfo{
				Name: "test-user",
			}, nil).Once()

		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything).Return([]sims.Sim{}, nil).Twice()

		mailerClient.On("SendEmail", mock.MatchedBy(func(req cnotif.SendEmailReq) bool {
			return req.To[0] == "test@example.com"
		})).Return(nil).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo, agentFactory,
			packageClient, subscriberService, simPoolService, "", msgbusClient, orgId.String(), "",
			mailerClient, netClient, orgClient, userClient)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberId.String(), NetworkId: networkId.String(),
			PackageId: packageId.String(), SimType: simTypeTest, SimToken: "",
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
		orgClient.AssertExpectations(t)
		userClient.AssertExpectations(t)
		mailerClient.AssertExpectations(t)
	})

	t.Run("SimTokenNotValid", func(t *testing.T) {
		simPoolService := &mocks.SimPoolClientProvider{}
		subscriberService := &mocks.SubscriberRegistryClientProvider{}
		packageClient := &cmocks.PackageClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()
		OrgId := uuid.NewV4()

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberId.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
					Email:        "test@example.com",
					Name:         "Test User",
				},
			}, nil).Once()

		packageInfo := &cdplan.PackageInfo{
			IsActive:   true,
			Duration:   3600,
			SimType:    simTypeTest,
			DataVolume: 10,
			DataUnit:   "GB",
			Name:       "Test Package",
			Amount:     100,
		}

		packageClient.On("Get", packageId.String()).
			Return(packageInfo, nil).
			Times(1)

		simPoolClient := simPoolService.On("GetClient").
			Return(&splmocks.SimServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*splmocks.SimServiceClient)

		s := server.NewSimManagerServer(OrgName, nil, nil, nil,
			packageClient, subscriberService, simPoolService, "", nil, OrgId.String(), "",
			nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberId.String(), NetworkId: networkId.String(),
			PackageId: packageId.String(), SimType: simTypeTest, SimToken: "lol",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)
		simPoolService.AssertExpectations(t)
		simPoolClient.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})

	t.Run("SimpoolServiceClientNotFound", func(t *testing.T) {
		simPoolService := &mocks.SimPoolClientProvider{}
		subscriberService := &mocks.SubscriberRegistryClientProvider{}
		packageClient := &cmocks.PackageClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()
		OrgId := uuid.NewV4()

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberId.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
					Email:        "test@example.com",
					Name:         "Test User",
				},
			}, nil).Once()

		packageInfo := &cdplan.PackageInfo{
			IsActive:   true,
			Duration:   3600,
			SimType:    simTypeTest,
			DataVolume: 10,
			DataUnit:   "GB",
			Name:       "Test Package",
			Amount:     100,
		}

		packageClient.On("Get", packageId.String()).
			Return(packageInfo, nil).
			Times(1)

		simPoolService.On("GetClient").
			Return(nil, errors.New("failed to get sim pool service client"))

		s := server.NewSimManagerServer(OrgName, nil, nil, nil,
			packageClient, subscriberService, simPoolService, "", nil, OrgId.String(), "",
			nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberId.String(), NetworkId: networkId.String(),
			PackageId: packageId.String(), SimType: simTypeTest, SimToken: "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)
		simPoolService.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})

	t.Run("SubscriberNotRegisteredOnProvidedNetwork", func(t *testing.T) {
		subscriberService := &mocks.SubscriberRegistryClientProvider{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberId.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberId.String(),
					NetworkId:    uuid.NewV4().String(),
				},
			}, nil).Once()

		s := server.NewSimManagerServer(OrgName, nil, nil, nil, nil, subscriberService,
			nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
			SimType:      simTypeTest,
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)
	})

	t.Run("OrgPackageNoMoreActive", func(t *testing.T) {
		subscriberService := &mocks.SubscriberRegistryClientProvider{}
		packageClient := &cmocks.PackageClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()
		orgId := uuid.NewV4()

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberId.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
				},
			}, nil).Once()

		packageClient.On("Get", packageId.String()).
			Return(&cdplan.PackageInfo{
				OrgId:    orgId.String(),
				IsActive: false,
				Duration: 3600,
				SimType:  simTypeTest,
			}, nil).Once()

		s := server.NewSimManagerServer(OrgName, nil, nil, nil,
			packageClient, subscriberService, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
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
		subscriberService := &mocks.SubscriberRegistryClientProvider{}
		packageClient := &cmocks.PackageClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()
		orgId := uuid.NewV4()

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberId.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
				},
			}, nil).Once()

		packageClient.On("Get", packageId.String()).
			Return(
				&cdplan.PackageInfo{
					OrgId:    orgId.String(),
					IsActive: true,
					Duration: 3600,
					SimType:  ukama.SimTypeUnknown.String(),
				}, nil).Once()

		s := server.NewSimManagerServer(OrgName, nil, nil, nil,
			packageClient, subscriberService, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
			SimType:      simTypeTest,
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)

		packageClient.AssertExpectations(t)
	})

	t.Run("PackageNotFound", func(t *testing.T) {
		subscriberService := &mocks.SubscriberRegistryClientProvider{}
		packageClient := &cmocks.PackageClient{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberId.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
				},
			}, nil).Once()

		packageClient.On("Get", packageId.String()).
			Return(nil, errors.New("package not found")).Once()

		s := server.NewSimManagerServer(OrgName, nil, nil, nil,
			packageClient, subscriberService, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
			SimType:      simTypeTest,
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)

		packageClient.AssertExpectations(t)
	})

	t.Run("PackageIdNotValid", func(t *testing.T) {
		subscriberService := &mocks.SubscriberRegistryClientProvider{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberId.String()}).
			Return(&subspb.GetSubscriberResponse{
				Subscriber: &upb.Subscriber{
					SubscriberId: subscriberId.String(),
					NetworkId:    networkId.String(),
				},
			}, nil).Once()

		s := server.NewSimManagerServer(OrgName, nil, nil, nil,
			nil, subscriberService, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    "lol",
			SimType:      simTypeTest,
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
		subscriberClient.AssertExpectations(t)
	})

	t.Run("SubscriberInfoNotFound", func(t *testing.T) {
		subscriberService := &mocks.SubscriberRegistryClientProvider{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()

		subscriberClient := subscriberService.On("GetClient").
			Return(&subsmocks.RegistryServiceClient{}, nil).
			Once().
			ReturnArguments.Get(0).(*subsmocks.RegistryServiceClient)

		subscriberClient.On("Get", mock.Anything,
			&subspb.GetSubscriberRequest{SubscriberId: subscriberId.String()}).
			Return(nil, errors.New("subscriber record not found")).Once()

		s := server.NewSimManagerServer(OrgName, nil, nil, nil,
			nil, subscriberService, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
			SimType:      simTypeTest,
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
	})

	t.Run("SubscriberServiceClientNotFound", func(t *testing.T) {
		subscriberService := &mocks.SubscriberRegistryClientProvider{}

		subscriberId := uuid.NewV4()
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()

		subscriberService.On("GetClient").
			Return(nil, errors.New("failed to get subscriber service client"))

		s := server.NewSimManagerServer(OrgName, nil, nil, nil,
			nil, subscriberService, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: subscriberId.String(),
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
			SimType:      simTypeTest,
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		subscriberService.AssertExpectations(t)
	})

	t.Run("InvalidSubscriberId", func(t *testing.T) {
		networkId := uuid.NewV4()
		packageId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, nil, nil, nil,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.AllocateSim(context.TODO(), &pb.AllocateSimRequest{
			SubscriberId: "lol",
			NetworkId:    networkId.String(),
			PackageId:    packageId.String(),
			SimType:      simTypeTest,
			SimToken:     "",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestSimManagerServer_SetActivePackageForSim(t *testing.T) {
	t.Run("SimAndPackageFound", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		packageId := uuid.NewV4()
		simId := uuid.NewV4()

		simd := &sims.Sim{Id: simId,
			IsPhysical: false,
			Status:     ukama.SimStatusActive,
			Type:       ukama.SimTypeTest,
		}

		simRepo.On("Get", simId).
			Return(simd, nil).
			Once()

		packageRepo.On("Get", packageId).Return(
			&sims.Package{Id: packageId,
				SimId:    simId,
				EndDate:  time.Now().UTC().AddDate(0, 1, 0), // next month
				IsActive: false,
			}, nil).Once()

		packageRepo.On("Update",
			&sims.Package{
				Id:       packageId,
				IsActive: true,
			},
			mock.Anything).Return(nil).Once()

		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		agentAdapter := agentFactory.On("GetAgentAdapter", simd.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("ActivateSim", mock.Anything,
			mock.MatchedBy(func(a client.AgentRequestData) bool {
				return a.Iccid == simd.Iccid
			})).Return(nil).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			agentFactory, nil, nil, nil, "", msgbusClient, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simId.String(),
			PackageId: packageId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
		agentAdapter.AssertExpectations(t)
	})

	t.Run("SimIdNotValid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()

		simRepo.On("Get", mock.Anything).
			Return(nil, errors.New("sim not found")).
			Once()

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simId.String(),
			PackageId: packageId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
	})

	t.Run("SimStatusInvalid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		packageId := uuid.NewV4()
		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				IsPhysical: false,
				Status:     ukama.SimStatusUnknown,
			}, nil).
			Once()

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simId.String(),
			PackageId: packageId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimActivePackageStillExists", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		packageId := uuid.NewV4()
		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{
				Id:         simId,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
				Package: sims.Package{
					Id: packageId,
				},
			}, nil).
			Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simId.String(),
			PackageId: packageId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
	})

	t.Run("PackageIdNotValid", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{
				Id:         simId,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simId.String(),
			PackageId: "lol",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
	})

	t.Run("PackageFetchFailure", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		packageId := uuid.NewV4()
		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{
				Id:         simId,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageId).Return(nil, errors.New("fail to get package Info")).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simId.String(),
			PackageId: packageId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SimIdAndPackageSimIdMismatch", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		packageId := uuid.NewV4()
		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageId).Return(
			&sims.Package{Id: packageId,
				SimId:    uuid.NewV4(),
				EndDate:  time.Now().UTC().AddDate(0, 1, 0), // next month
				IsActive: false,
			}, nil).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simId.String(),
			PackageId: packageId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
	})

	t.Run("PackageAlreadyExpired", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		packageId := uuid.NewV4()
		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageId).Return(
			&sims.Package{Id: packageId,
				SimId:     simId,
				EndDate:   time.Now().UTC().AddDate(0, -1, 0), // one month ago
				IsActive:  false,
				AsExpired: true,
			}, nil).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.SetActivePackageForSim(context.TODO(), &pb.SetActivePackageRequest{
			SimId:     simId.String(),
			PackageId: packageId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_ToggleSimStatus(t *testing.T) {
	t.Run("InvalidSimStatus", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}

		simId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.ToggleSimStatus(context.TODO(), &pb.ToggleSimStatusRequest{
			SimId:  simId.String(),
			Status: "lol",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		packageRepo.AssertExpectations(t)
		simRepo.AssertExpectations(t)
	})

	t.Run("ActiveSimNotFound", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}

		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(nil, errors.New("sim not found")).
			Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.ToggleSimStatus(context.TODO(), &pb.ToggleSimStatusRequest{
			SimId:  simId.String(),
			Status: ukama.SimStatusActive.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		packageRepo.AssertExpectations(t)
		simRepo.AssertExpectations(t)
	})

	t.Run("InactiveSimNotFound", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}

		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(nil, errors.New("sim not found")).
			Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.ToggleSimStatus(context.TODO(), &pb.ToggleSimStatusRequest{
			SimId:  simId.String(),
			Status: ukama.SimStatusInactive.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		packageRepo.AssertExpectations(t)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimAlreadyInactive", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}

		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{
				Id:         simId,
				IsPhysical: false,
				Status:     ukama.SimStatusInactive,
			}, nil).
			Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.ToggleSimStatus(context.TODO(), &pb.ToggleSimStatusRequest{
			SimId:  simId.String(),
			Status: ukama.SimStatusInactive.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		packageRepo.AssertExpectations(t)
		simRepo.AssertExpectations(t)
	})

	t.Run("AgentFactoryFailure", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{
				Id:         simId,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		agentFactory.On("GetAgentAdapter", mock.Anything).
			Return(nil, false).
			Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			agentFactory, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.ToggleSimStatus(context.TODO(), &pb.ToggleSimStatusRequest{
			SimId:  simId.String(),
			Status: ukama.SimStatusInactive.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		packageRepo.AssertExpectations(t)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimStatusUpdateFailure", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{
				Id:         simId,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		simRepo.On("Update",
			&sims.Sim{
				Id:                 simId,
				Status:             ukama.SimStatusInactive,
				DeactivationsCount: uint64(1),
			},
			mock.Anything).Return(errors.New("sim status update failure")).Once()

		agentFactory.On("GetAgentAdapter", mock.Anything).
			Return(&mocks.AgentAdapter{}, true).
			Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			agentFactory, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.ToggleSimStatus(context.TODO(), &pb.ToggleSimStatusRequest{
			SimId:  simId.String(),
			Status: ukama.SimStatusInactive.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		packageRepo.AssertExpectations(t)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimAgentFailure", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		sim := simRepo.On("Get", simId).
			Return(&sims.Sim{
				Id:         simId,
				Iccid:      testIccid,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
				Type:       ukama.SimTypeTest,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*sims.Sim)

		simRepo.On("Update",
			&sims.Sim{
				Id:                 simId,
				Status:             ukama.SimStatusInactive,
				DeactivationsCount: uint64(1),
			},
			mock.Anything).Return(nil).Once()

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("DeactivateSim", mock.Anything,
			mock.Anything).Return(errors.New("fail to deactivate sim on remove agent")).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			agentFactory, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.ToggleSimStatus(context.TODO(), &pb.ToggleSimStatusRequest{
			SimId:  simId.String(),
			Status: ukama.SimStatusInactive.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		packageRepo.AssertExpectations(t)
		simRepo.AssertExpectations(t)
	})

	t.Run("SimDeactivated", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		agentFactory := &mocks.AgentFactory{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		simId := uuid.NewV4()

		sim := simRepo.On("Get", simId).
			Return(&sims.Sim{
				Id:         simId,
				Iccid:      testIccid,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
				Type:       ukama.SimTypeTest,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*sims.Sim)

		simRepo.On("Update",
			&sims.Sim{
				Id:                 simId,
				Status:             ukama.SimStatusInactive,
				DeactivationsCount: uint64(1),
			},
			mock.Anything).Return(nil).Once()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]sims.Sim{}, nil).Twice()

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		agentAdapter.On("DeactivateSim", mock.Anything,
			mock.Anything).Return(nil).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			agentFactory, nil, nil, nil, "", msgbusClient, "", "", nil, nil, nil, nil)

		resp, err := s.ToggleSimStatus(context.TODO(), &pb.ToggleSimStatusRequest{
			SimId:  simId.String(),
			Status: ukama.SimStatusInactive.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		packageRepo.AssertExpectations(t)
		simRepo.AssertExpectations(t)
	})

}

func TestSimManagerServer_RemovePackageForSim(t *testing.T) {
	t.Run("PackageFound", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		msgbusClient := &cmocks.MsgBusServiceClient{}

		packageId := uuid.NewV4()
		simId := uuid.NewV4()

		packageRepo.On("Get", packageId).Return(
			&sims.Package{Id: packageId,
				SimId:    simId,
				IsActive: false,
			}, nil).Once()

		simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Delete", packageId,
			mock.Anything).Return(nil).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", msgbusClient, "", "", nil, nil, nil, nil)
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageId.String(),
			SimId:     simId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		packageRepo.AssertExpectations(t)
		msgbusClient.AssertExpectations(t)
		simRepo.AssertExpectations(t)
	})

	t.Run("PackageDeleteError", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		packageId := uuid.Nil
		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageId).Return(
			&sims.Package{Id: packageId,
				SimId:    simId,
				IsActive: false,
			}, nil).Once()

		packageRepo.On("Delete", packageId,
			mock.Anything).Return(gorm.ErrRecordNotFound).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageId.String(),
			SimId:     simId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("PackageUUIDInvalid", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}

		packageId := "1"

		s := server.NewSimManagerServer(OrgName, nil, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageId})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})

	t.Run("SimUUIDMismatch", func(t *testing.T) {
		packageRepo := &mocks.PackageRepo{}
		simRepo := &mocks.SimRepo{}

		packageId := uuid.Nil
		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				IsPhysical: false,
				Status:     ukama.SimStatusActive,
			}, nil).
			Once()

		packageRepo.On("Get", packageId).Return(
			&sims.Package{Id: packageId,
				SimId:    simId,
				IsActive: false,
			}, nil).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo,
			nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.RemovePackageForSim(context.TODO(), &pb.RemovePackageRequest{
			PackageId: packageId.String(),
			SimId:     uuid.NewV4().String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		packageRepo.AssertExpectations(t)
	})
}

func TestSimManagerServer_AddPackageForSim(t *testing.T) {
	t.Run("PackageStartDateNotValid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		packageClient := &cmocks.PackageClient{}

		simId := uuid.NewV4()
		packageId := uuid.NewV4()
		orgId := uuid.NewV4()
		startDate := "lol"

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo, nil, packageClient,
			nil, nil, "", nil, orgId.String(), "", nil, nil, nil, nil)

		resp, err := s.AddPackageForSim(context.TODO(), &pb.AddPackageRequest{
			SimId:     simId.String(),
			PackageId: packageId.String(),
			StartDate: startDate,
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)

	})
}

func TestSimManagerServer_TerminatePackageForSim(t *testing.T) {
	t.Run("PackageIdNotValid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		packageRepo := &mocks.PackageRepo{}
		packageClient := &cmocks.PackageClient{}

		simId := uuid.NewV4()
		orgId := uuid.NewV4()

		s := server.NewSimManagerServer(OrgName, simRepo, packageRepo, nil, packageClient,
			nil, nil, "", nil, orgId.String(), "", nil, nil, nil, nil)

		resp, err := s.TerminatePackageForSim(context.TODO(), &pb.TerminatePackageRequest{
			SimId:     simId.String(),
			PackageId: "lol",
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
		packageRepo.AssertExpectations(t)
		packageClient.AssertExpectations(t)
	})
}

func TestSimManagerServer_TerminateSim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		msgbusClient := &cmocks.MsgBusServiceClient{}
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		sim := simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*sims.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("TerminateSim", mock.Anything,
			sim.Iccid).Return(nil).Once()

		simRepo.On("Update",
			&sims.Sim{
				Id:     sim.Id,
				Status: ukama.SimStatusTerminated,
			},
			mock.Anything).Return(nil).Once()
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

		simRepo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]sims.Sim{}, nil).Twice()

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", msgbusClient, "", "", nil, nil, nil, nil)

		resp, err := s.TerminateSim(context.TODO(), &pb.TerminateSimRequest{
			SimId: simId.String(),
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)

		simRepo.AssertExpectations(t)
		agentFactory.AssertExpectations(t)
	})

	t.Run("SimStatusInvalid", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}

		simId := uuid.NewV4()

		simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				Iccid:      testIccid,
				Status:     ukama.SimStatusActive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once()

		s := server.NewSimManagerServer(OrgName, simRepo,
			nil, nil, nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.TerminateSim(context.TODO(), &pb.TerminateSimRequest{
			SimId: simId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
	})

	t.Run("SimTypeNotSupported", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		sim := simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       100,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*sims.Sim)

		agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, false).
			Once()

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.TerminateSim(context.TODO(), &pb.TerminateSimRequest{
			SimId: simId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
	})

	t.Run("SimAgentFailToTerminate", func(t *testing.T) {
		simRepo := &mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		sim := simRepo.On("Get", simId).
			Return(&sims.Sim{Id: simId,
				Iccid:      testIccid,
				Status:     ukama.SimStatusInactive,
				Type:       ukama.SimTypeTest,
				IsPhysical: false,
			}, nil).
			Once().
			ReturnArguments.Get(0).(*sims.Sim)

		agentAdapter := agentFactory.On("GetAgentAdapter", sim.Type).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("TerminateSim", mock.Anything,
			sim.Iccid).Return(errors.New("anyError")).Once()

		s := server.NewSimManagerServer(OrgName, simRepo, nil, agentFactory,
			nil, nil, nil, "", nil, "", "", nil, nil, nil, nil)

		resp, err := s.TerminateSim(context.TODO(), &pb.TerminateSimRequest{
			SimId: simId.String(),
		})

		assert.Error(t, err)
		assert.Nil(t, resp)

		simRepo.AssertExpectations(t)
	})
}

func assertList(t *testing.T, list *pb.ListSimsResponse, resp []sims.Sim) {
	for idx, sim := range list.Sims {
		assert.Equal(t, sim.Id, resp[idx].Id.String())
		assert.Equal(t, sim.Iccid, resp[idx].Iccid)
		assert.Equal(t, sim.Type, resp[idx].Type.String())
	}
}
