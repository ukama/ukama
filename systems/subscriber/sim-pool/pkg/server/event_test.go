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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	spmock "github.com/ukama/ukama/systems/subscriber/sim-pool/mocks"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	testOrgName       = "test-org"
	testIccid         = "test-iccid-123"
	testIccidError    = "test-iccid-789"
	testId            = "test-id"
	unknownRoutingKey = "unknown.routing.key"
)

func TestSimAllocationEvent(t *testing.T) {
	mockRepo := &spmock.SimRepo{}
	server := NewSimPoolEventServer(testOrgName, mockRepo)
	msgbusClient := &mbmocks.MsgBusServiceClient{}
	routingKey := msgbus.PrepareRoute(testOrgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate")

	t.Run("Success", func(t *testing.T) {
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		simAllocation := &epb.EventSimAllocation{
			Iccid: testIccid,
		}
		anyMsg, err := anypb.New(simAllocation)
		assert.NoError(t, err)

		event := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyMsg,
		}

		mockRepo.On("UpdateStatus", testIccid, true, false).Return(nil)

		// Act
		response, err := server.EventNotification(context.Background(), event)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateStatusError", func(t *testing.T) {
		msg := &epb.EventSimAllocation{
			Iccid: testIccidError,
		}

		expectedError := errors.New("database update failed")
		mockRepo.On("UpdateStatus", testIccidError, true, false).Return(expectedError)

		// Act
		err := handleEventCloudSimManagerSimAllocate(routingKey, msg, server)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UnmarshalError", func(t *testing.T) {
		msgbusClient.On("PublishRequest", mock.Anything, mock.Anything).
			Return(nil).Once()

		invalidMsg := &epb.EventSimUsage{
			Id: testId,
		}
		anyMsg, err := anypb.New(invalidMsg)
		assert.NoError(t, err)

		event := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyMsg,
		}

		// Act
		response, err := server.EventNotification(context.Background(), event)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, response)
	})

	t.Run("UnknownRoutingKey", func(t *testing.T) {
		simAllocation := &epb.EventSimAllocation{
			Iccid: testIccid,
		}
		anyMsg, err := anypb.New(simAllocation)
		assert.NoError(t, err)

		event := &epb.Event{
			RoutingKey: unknownRoutingKey,
			Msg:        anyMsg,
		}

		// Act
		response, err := server.EventNotification(context.Background(), event)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, response)
	})
}
