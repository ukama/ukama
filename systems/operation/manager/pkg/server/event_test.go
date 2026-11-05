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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"

	evt "github.com/ukama/ukama/systems/common/events"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/store"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/operation/manager/mocks"
	"github.com/ukama/ukama/systems/operation/manager/pkg/db"
)

func newOperationEvent(t *testing.T, key evt.EventId, msg proto.Message) *epb.Event {
	t.Helper()

	body, err := anypb.New(msg)
	require.NoError(t, err)

	return &epb.Event{
		RoutingKey: msgbus.PrepareRoute(orgName, evt.EventRoutingKey[key]),
		Msg:        body,
	}
}

func TestEventNotification_terminalEvents(t *testing.T) {
	opId := uuid.NewV4()
	precondition := fmt.Errorf("already timeout: %w", store.ErrPrecondition)

	events := []struct {
		name   string
		status db.OperationStatus
		event  func(t *testing.T) *epb.Event
	}{
		{"Completed", db.OperationSuccess, func(t *testing.T) *epb.Event {
			return newOperationEvent(t, evt.EventOperationCompleted,
				&epb.OperationCompletedEvent{OperationId: opId.String(), FencingToken: 5})
		}},
		{"Failed", db.OperationFailed, func(t *testing.T) *epb.Event {
			return newOperationEvent(t, evt.EventOperationFailed,
				&epb.OperationFailedEvent{OperationId: opId.String(), FencingToken: 5,
					Reason: "deadline exceeded"})
		}},
	}

	for _, e := range events {
		t.Run(e.name+"Applied", func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			repo.On("Terminate", mock.Anything, opId, uint64(5), e.status,
				mock.Anything, mock.Anything).
				Return(&db.Operation{}, nil).Once()

			resp, err := NewEventServer(orgName, repo).
				EventNotification(context.Background(), e.event(t))

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			repo.AssertExpectations(t)
		})

		t.Run(e.name+"LateIsAcknowledged", func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			repo.On("Terminate", mock.Anything, opId, uint64(5), e.status,
				mock.Anything, mock.Anything).
				Return(nil, precondition).Once()

			resp, err := NewEventServer(orgName, repo).
				EventNotification(context.Background(), e.event(t))

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			repo.AssertExpectations(t)
		})

		t.Run(e.name+"UnknownOperationIsAcknowledged", func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			repo.On("Terminate", mock.Anything, opId, uint64(5), e.status,
				mock.Anything, mock.Anything).
				Return(nil, gorm.ErrRecordNotFound).Once()

			resp, err := NewEventServer(orgName, repo).
				EventNotification(context.Background(), e.event(t))

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			repo.AssertExpectations(t)
		})

		t.Run(e.name+"RepoFailureIsUnavailable", func(t *testing.T) {
			repo := &mocks.OperationRepo{}
			repo.On("Terminate", mock.Anything, opId, uint64(5), e.status,
				mock.Anything, mock.Anything).
				Return(nil, assert.AnError).Once()

			resp, err := NewEventServer(orgName, repo).
				EventNotification(context.Background(), e.event(t))

			assert.Equal(t, codes.Unavailable, status.Code(err))
			assert.Nil(t, resp)
			repo.AssertExpectations(t)
		})
	}
}
