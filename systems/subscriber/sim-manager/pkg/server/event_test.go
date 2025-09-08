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

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/mocks"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	sims "github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/db"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/server"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const (
	OrgName = "testOrg"
	orgId   = "592f7a8e-f318-4d3a-aab8-8d4187cde7f9"
	simId   = "e044081b-fbbe-45e9-8f78-0f9c0f112977"
)

func TestSimManagerEventServer_HandleSimManagerSimAllocateEvent(t *testing.T) {
	msgbusClient := &mbmocks.MsgBusServiceClient{}

	routingKey := msgbus.PrepareRoute(OrgName,
		"event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate")

	msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Once()

	t.Run("AllocatedSimNotFound", func(t *testing.T) {
		repo := mocks.SimRepo{}
		repo.On("Get", mock.Anything).Return(nil, errors.New("sim not found"))

		allocatedSim := epb.EventSimAllocation{
			Id: simId,
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ActiveSimFound", func(t *testing.T) {
		repo := mocks.SimRepo{}
		simId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&db.Sim{
			Id:     simId,
			Status: ukama.SimStatusActive,
		}, nil)

		allocatedSim := epb.EventSimAllocation{
			Id:     simId.String(),
			Status: ukama.SimStatusActive.String(),
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ActiveSimUpdateFailure", func(t *testing.T) {
		repo := mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&db.Sim{
			Id:     simId,
			Status: ukama.SimStatusInactive,
			Type:   ukama.SimTypeUkamaData,
		}, nil)

		repo.On("Update", mock.Anything, mock.Anything).
			Return(errors.New("failed to update sim"))

		allocatedSim := epb.EventSimAllocation{
			Id:     simId.String(),
			Status: ukama.SimStatusInactive.String(),
			Type:   ukama.SimTypeUkamaData.String(),
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, agentFactory, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ActiveSimAgentFailure", func(t *testing.T) {
		repo := mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&db.Sim{
			Id:     simId,
			Status: ukama.SimStatusInactive,
		}, nil)

		repo.On("Update", mock.Anything, mock.Anything).Return(nil)

		agentFactory.On("GetAgentAdapter", mock.Anything).
			Return(nil, false).Once()

		allocatedSim := epb.EventSimAllocation{
			Id:     simId.String(),
			Status: ukama.SimStatusInactive.String(),
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, agentFactory, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ActiveSimAgentActivateSimFailure", func(t *testing.T) {
		repo := mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&db.Sim{
			Id:     simId,
			Status: ukama.SimStatusInactive,
			Type:   ukama.SimTypeUkamaData,
		}, nil)

		repo.On("Update", mock.Anything, mock.Anything).Return(nil)

		agentAdapter := agentFactory.On("GetAgentAdapter", ukama.SimTypeUkamaData).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("ActivateSim", mock.Anything, mock.Anything).
			Return(errors.New("fail to activate sim on remote agent")).Once()

		allocatedSim := epb.EventSimAllocation{
			Id:     simId.String(),
			Status: ukama.SimStatusInactive.String(),
			Type:   ukama.SimTypeUkamaData.String(),
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, agentFactory, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("ActiveSimAgentActivateSimSuccess", func(t *testing.T) {
		repo := mocks.SimRepo{}
		agentFactory := &mocks.AgentFactory{}

		simId := uuid.NewV4()

		repo.On("Get", mock.Anything).Return(&db.Sim{
			Id:     simId,
			Status: ukama.SimStatusInactive,
			Type:   ukama.SimTypeUkamaData,
		}, nil)

		repo.On("Update", mock.Anything, mock.Anything).Return(nil)

		agentAdapter := agentFactory.On("GetAgentAdapter", ukama.SimTypeUkamaData).
			Return(&mocks.AgentAdapter{}, true).
			Once().
			ReturnArguments.Get(0).(*mocks.AgentAdapter)

		agentAdapter.On("ActivateSim", mock.Anything, mock.Anything).
			Return(nil).Once()

		repo.On("List", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]sims.Sim{}, nil)

		allocatedSim := epb.EventSimAllocation{
			Id:     simId.String(),
			Status: ukama.SimStatusInactive.String(),
			Type:   ukama.SimTypeUkamaData.String(),
		}

		anyE, err := anypb.New(&allocatedSim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, agentFactory, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("InvalidEventTypeSent", func(t *testing.T) {
		repo := mocks.SimRepo{}
		evt := &epb.Payment{
			Id: uuid.NewV4().String(),
		}

		anyE, err := anypb.New(evt)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewSimManagerEventServer(OrgName, orgId, &repo, nil, msgbusClient, "", nil)
		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}
