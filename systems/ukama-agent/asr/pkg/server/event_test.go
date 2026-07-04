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
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/server"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	mocks "github.com/ukama/ukama/systems/ukama-agent/asr/mocks"
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

		evt := &epb.EventSimAllocation{
			Iccid:      server.Iccid,
			Imsi:       server.Imsi,
			DataPlanId: uuid.NewV4().String(),
			NetworkId:  uuid.NewV4().String(),
			Type:       ukama.SimTypeUkamaData.String(),
			PackageId:  uuid.NewV4().String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		network.On("Get", evt.NetworkId).Return(&registry.NetworkInfo{}, nil).Once()
		factory.On("ReadSimCardInfo", evt.Iccid).Return(&server.Sim, nil).Once()
		pc.On("NewPolicy", mock.Anything).Return(&server.Policy, nil).Once()
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
