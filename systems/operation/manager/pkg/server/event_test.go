/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/protobuf/types/known/anypb"

	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/operation/manager/mocks"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

func completedRoute() string {
	return msgbus.PrepareRoute(orgName, evt.EventRoutingKey[evt.EventOperationCompleted])
}

func failedRoute() string {
	return msgbus.PrepareRoute(orgName, evt.EventRoutingKey[evt.EventOperationFailed])
}

func TestEventNotification(t *testing.T) {
	t.Run("NilEvent", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		e := NewEventServer(orgName, repo)

		resp, err := e.EventNotification(context.Background(), nil)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("UnhandledRoutingKey", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		e := NewEventServer(orgName, repo)

		resp, err := e.EventNotification(context.Background(), &epb.Event{
			RoutingKey: "event.cloud.global.ukama.something.else",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		repo.AssertNotCalled(t, "Terminate", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("CompletedReleasesLock", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		e := NewEventServer(orgName, repo)

		opId := uuid.NewV4()
		repo.On("Terminate", opId, uint64(7), db.OperationSuccess,
			mock.AnythingOfType("db.OperationAudit"), "").
			Return(&db.Operation{Id: opId, Status: db.OperationSuccess}, nil)

		msg, err := anypb.New(&epb.OperationCompletedEvent{
			OperationId:  opId.String(),
			FencingToken: 7,
			ResourceKey:  "node:abc",
		})
		assert.NoError(t, err)

		resp, err := e.EventNotification(context.Background(), &epb.Event{
			RoutingKey: completedRoute(),
			Msg:        msg,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("FailedReleasesLock", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		e := NewEventServer(orgName, repo)

		opId := uuid.NewV4()
		repo.On("Terminate", opId, uint64(4), db.OperationFailed,
			mock.AnythingOfType("db.OperationAudit"), "boom").
			Return(&db.Operation{Id: opId, Status: db.OperationFailed}, nil)

		msg, err := anypb.New(&epb.OperationFailedEvent{
			OperationId:  opId.String(),
			FencingToken: 4,
			ResourceKey:  "node:abc",
			Reason:       "boom",
		})
		assert.NoError(t, err)

		resp, err := e.EventNotification(context.Background(), &epb.Event{
			RoutingKey: failedRoute(),
			Msg:        msg,
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		repo.AssertExpectations(t)
	})

	t.Run("CompletedInvalidOperationId", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		e := NewEventServer(orgName, repo)

		msg, err := anypb.New(&epb.OperationCompletedEvent{
			OperationId:  "not-a-uuid",
			FencingToken: 1,
		})
		assert.NoError(t, err)

		resp, err := e.EventNotification(context.Background(), &epb.Event{
			RoutingKey: completedRoute(),
			Msg:        msg,
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("CompletedTerminateError", func(t *testing.T) {
		repo := &mocks.OperationRepo{}
		e := NewEventServer(orgName, repo)

		opId := uuid.NewV4()
		repo.On("Terminate", opId, uint64(2), db.OperationSuccess,
			mock.AnythingOfType("db.OperationAudit"), "").
			Return(nil, errors.New("db down"))

		msg, err := anypb.New(&epb.OperationCompletedEvent{
			OperationId:  opId.String(),
			FencingToken: 2,
		})
		assert.NoError(t, err)

		resp, err := e.EventNotification(context.Background(), &epb.Event{
			RoutingKey: completedRoute(),
			Msg:        msg,
		})

		assert.Error(t, err)
		assert.Nil(t, resp)
		repo.AssertExpectations(t)
	})
}
