/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/node/operation-monitor/mocks"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg/db"
)

func TestSweeper_PublishFailureCanRetry(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(true)
	intent.Deadline = time.Now().Add(-time.Minute)
	s := NewSweeper(NewMonitorServer(testOrg, "org-id", repo, mb))

	repo.On("FindExpired", mock.Anything, 100).
		Return([]db.MonitoredIntent{intent}, nil).Twice()
	mb.On("PublishRequest", mock.Anything, mock.Anything).
		Return(errors.New("message bus unavailable")).Once()

	s.sweepOnce()
	repo.AssertNotCalled(t, "MarkTerminal", mock.Anything, mock.Anything)

	mb.On("PublishRequest", mock.Anything, mock.MatchedBy(func(event *epb.OperationFailedEvent) bool {
		return event.OperationId == intent.OperationId.String() &&
			event.FencingToken == intent.FencingToken &&
			event.Reason == "deadline exceeded"
	})).Return(nil).Once()
	repo.On("MarkTerminal", intent.OperationId, db.IntentExpired).
		Return(&intent, nil).Once()

	s.sweepOnce()
	repo.AssertExpectations(t)
	mb.AssertExpectations(t)
}

func TestSweeper_TerminalWriteFailureCanRetry(t *testing.T) {
	repo := &mocks.IntentRepo{}
	mb := &mbmocks.MsgBusServiceClient{}
	intent := watchingIntent(true)
	intent.Deadline = time.Now().Add(-time.Minute)
	s := NewSweeper(NewMonitorServer(testOrg, "org-id", repo, mb))

	repo.On("FindExpired", mock.Anything, 100).
		Return([]db.MonitoredIntent{intent}, nil).Twice()
	mb.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Twice()
	repo.On("MarkTerminal", intent.OperationId, db.IntentExpired).
		Run(func(_ mock.Arguments) {
			mb.AssertNumberOfCalls(t, "PublishRequest", 1)
		}).Return(nil, errors.New("database unavailable")).Once()

	s.sweepOnce()
	mb.AssertNumberOfCalls(t, "PublishRequest", 1)

	repo.On("MarkTerminal", intent.OperationId, db.IntentExpired).
		Return(&intent, nil).Once()
	s.sweepOnce()

	repo.AssertExpectations(t)
	mb.AssertExpectations(t)
}
