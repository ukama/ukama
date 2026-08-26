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
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/server"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	mocks "github.com/ukama/ukama/systems/ukama-agent/asr/mocks"
	pm "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"
)

func TestUkamaAgentEventServer_HandleSimAllocationEvent(t *testing.T) {
	factory := &cmocks.SimFactoryClient{}
	network := &cmocks.NetworkClient{}
	pc := &mocks.Controller{}
	msgbusClient := &cmocks.MsgBusServiceClient{}

	routingKey := msgbus.PrepareRoute(server.Org,
		"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate")

	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	t.Run("ASRActivationSuccess", func(t *testing.T) {
		asrRepo := &mocks.AsrRecordRepo{}

		startDate := timestamppb.New(time.Unix(1700000000, 0))
		endDate := timestamppb.New(time.Unix(1700100000, 0))

		evt := &epb.EventSimAllocation{
			Iccid:            server.Iccid,
			Imsi:             server.Imsi,
			DataPlanId:       uuid.NewV4().String(),
			NetworkId:        uuid.NewV4().String(),
			Type:             ukama.SimTypeUkamaData.String(),
			PackageId:        uuid.NewV4().String(),
			PackageTotalData: 1024000000,
			PackageDlbr:      15000,
			PackageUlbr:      2000,
			PackageStartDate: startDate,
			PackageEndDate:   endDate,
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		network.On("Get", evt.NetworkId).Return(&registry.NetworkInfo{}, nil).Once()
		factory.On("ReadSimCardInfo", evt.Iccid).Return(&server.Sim, nil).Once()
		pc.On("NewPolicy", pm.PolicyInput{
			TotalData: evt.PackageTotalData,
			Dlbr:      evt.PackageDlbr,
			Ulbr:      evt.PackageUlbr,
			StartTime: uint64(startDate.AsTime().Unix()),
			EndTime:   uint64(endDate.AsTime().Unix()),
		}).Return(&server.Policy, nil).Once()
		asrRepo.On("Add", mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == evt.Iccid
		})).Return(nil).Once()
		pc.On("RunPolicyControl", evt.Imsi, false).Return(nil, false).Once()
		pc.On("SyncProfile", mock.Anything, mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == evt.Iccid
		}), msgbus.ACTION_CRUD_CREATE, "activesubscriber", true).Return(nil, false).Once()

		s := server.NewAsrEventServer(asrRepo, nil, nil, factory, network, pc, msgbusClient, server.Atos, server.Org)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("ASRActivationSucceedsDespitePcrfSyncFailure", func(t *testing.T) {
		asrRepo := &mocks.AsrRecordRepo{}
		noRollbackMsgbus := &cmocks.MsgBusServiceClient{}

		startDate := timestamppb.New(time.Unix(1700000000, 0))
		endDate := timestamppb.New(time.Unix(1700100000, 0))

		evt := &epb.EventSimAllocation{
			Iccid:            server.Iccid,
			Imsi:             server.Imsi,
			DataPlanId:       uuid.NewV4().String(),
			NetworkId:        uuid.NewV4().String(),
			Type:             ukama.SimTypeUkamaData.String(),
			PackageId:        uuid.NewV4().String(),
			PackageTotalData: 1024000000,
			PackageDlbr:      15000,
			PackageUlbr:      2000,
			PackageStartDate: startDate,
			PackageEndDate:   endDate,
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		network.On("Get", evt.NetworkId).Return(&registry.NetworkInfo{}, nil).Once()
		factory.On("ReadSimCardInfo", evt.Iccid).Return(&server.Sim, nil).Once()
		pc.On("NewPolicy", pm.PolicyInput{
			TotalData: evt.PackageTotalData,
			Dlbr:      evt.PackageDlbr,
			Ulbr:      evt.PackageUlbr,
			StartTime: uint64(startDate.AsTime().Unix()),
			EndTime:   uint64(endDate.AsTime().Unix()),
		}).Return(&server.Policy, nil).Once()
		asrRepo.On("Add", mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == evt.Iccid
		})).Return(nil).Once()
		pc.On("RunPolicyControl", evt.Imsi, false).Return(nil, false).Once()
		pc.On("SyncProfile", mock.Anything, mock.MatchedBy(func(a1 *db.Asr) bool {
			return a1.Iccid == evt.Iccid
		}), msgbus.ACTION_CRUD_CREATE, "activesubscriber", true).Return(errors.New("pcrf unreachable")).Once()

		s := server.NewAsrEventServer(asrRepo, nil, nil, factory, network, pc, noRollbackMsgbus, server.Atos, server.Org)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
		noRollbackMsgbus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
	})

	t.Run("ASRActivationError", func(t *testing.T) {
		repo := &mocks.AsrRecordRepo{}
		evt := &epb.EventSimAllocation{
			Type:      ukama.SimTypeUkamaData.String(),
			PackageId: "lol",
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewAsrEventServer(repo, nil, nil, factory, network, pc, msgbusClient, server.Atos, server.Org)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("InvalidSimType", func(t *testing.T) {
		repo := &mocks.AsrRecordRepo{}
		evt := &epb.EventSimAllocation{
			Type: ukama.SimTypeTest.String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewAsrEventServer(repo, nil, nil, factory, network, pc, msgbusClient, server.Atos, server.Org)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("InvalidEventTypeSent", func(t *testing.T) {
		repo := &mocks.AsrRecordRepo{}
		evt := &epb.EventAddSite{
			SiteId: uuid.NewV4().String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewAsrEventServer(repo, nil, nil, factory, network, pc, msgbusClient, server.Atos, server.Org)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestUkamaAgentEventServer_HandleSimServiceOnEvent(t *testing.T) {
	pc := &mocks.Controller{}
	msgbusClient := &cmocks.MsgBusServiceClient{}

	routingKey := msgbus.PrepareRoute(server.Org,
		"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.serviceon")

	t.Run("ServiceOnSuccess", func(t *testing.T) {
		asrRepo := &mocks.AsrRecordRepo{}

		evt := &epb.EventSimServiceOn{
			Iccid: server.Iccid,
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		asrRecord := &db.Asr{Iccid: server.Iccid, Imsi: server.Imsi}

		asrRepo.On("GetByIccid", server.Iccid).Return(asrRecord, nil)
		asrRepo.On("Update", server.Imsi, mock.MatchedBy(func(a *db.Asr) bool {
			return a.ServiceStatus == ukama.SimStatusServiceOn
		})).Return(nil).Once()
		pc.On("SyncProfile", mock.Anything, mock.Anything, msgbus.ACTION_CRUD_UPDATE, "activesubscriber", true).
			Return(nil).Once()

		s := server.NewAsrEventServer(asrRepo, nil, nil, nil, nil, pc, msgbusClient, server.Atos, server.Org)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
		asrRepo.AssertExpectations(t)
	})

	t.Run("NoAsrProfileIsNoOp", func(t *testing.T) {
		asrRepo := &mocks.AsrRecordRepo{}

		evt := &epb.EventSimServiceOn{
			Iccid: server.Iccid,
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		asrRepo.On("GetByIccid", server.Iccid).Return(nil, gorm.ErrRecordNotFound).Once()

		s := server.NewAsrEventServer(asrRepo, nil, nil, nil, nil, pc, msgbusClient, server.Atos, server.Org)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
		asrRepo.AssertExpectations(t)
	})

	t.Run("InvalidEventTypeSent", func(t *testing.T) {
		repo := &mocks.AsrRecordRepo{}
		evt := &epb.EventAddSite{
			SiteId: uuid.NewV4().String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewAsrEventServer(repo, nil, nil, nil, nil, pc, msgbusClient, server.Atos, server.Org)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestUkamaAgentEventServer_HandleSimServiceOffEvent(t *testing.T) {
	pc := &mocks.Controller{}
	msgbusClient := &cmocks.MsgBusServiceClient{}

	routingKey := msgbus.PrepareRoute(server.Org,
		"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.serviceoff")

	t.Run("ServiceOffSuccess", func(t *testing.T) {
		asrRepo := &mocks.AsrRecordRepo{}

		evt := &epb.EventSimServiceOff{
			Iccid: server.Iccid,
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		asrRecord := &db.Asr{Iccid: server.Iccid, Imsi: server.Imsi}

		asrRepo.On("GetByIccid", server.Iccid).Return(asrRecord, nil)
		asrRepo.On("Update", server.Imsi, mock.MatchedBy(func(a *db.Asr) bool {
			return a.ServiceStatus == ukama.SimStatusServiceOff
		})).Return(nil).Once()
		pc.On("SyncProfile", mock.Anything, mock.Anything, msgbus.ACTION_CRUD_UPDATE, "activesubscriber", true).
			Return(nil).Once()

		s := server.NewAsrEventServer(asrRepo, nil, nil, nil, nil, pc, msgbusClient, server.Atos, server.Org)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
		asrRepo.AssertExpectations(t)
	})

	t.Run("SyncFailurePropagatesError", func(t *testing.T) {
		asrRepo := &mocks.AsrRecordRepo{}
		localPc := &mocks.Controller{}

		evt := &epb.EventSimServiceOff{
			Iccid: server.Iccid,
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		asrRecord := &db.Asr{Iccid: server.Iccid, Imsi: server.Imsi}

		asrRepo.On("GetByIccid", server.Iccid).Return(asrRecord, nil)
		asrRepo.On("Update", server.Imsi, mock.Anything).Return(nil).Once()
		localPc.On("SyncProfile", mock.Anything, mock.Anything, msgbus.ACTION_CRUD_UPDATE, "activesubscriber", true).
			Return(errors.New("pcrf unreachable")).Once()

		s := server.NewAsrEventServer(asrRepo, nil, nil, nil, nil, localPc, msgbusClient, server.Atos, server.Org)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
		asrRepo.AssertExpectations(t)
	})
}
