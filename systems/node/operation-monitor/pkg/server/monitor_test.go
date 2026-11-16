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

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/node/operation-monitor/mocks"
	pb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
	"github.com/ukama/ukama/systems/node/operation-monitor/pkg/db"
)

func newMonitorServer(repo db.IntentRepo) *MonitorServer {
	return NewMonitorServer(testOrg, "org-id", repo, &mbmocks.MsgBusServiceClient{})
}

func registerRequest() *pb.RegisterIntentRequest {
	return &pb.RegisterIntentRequest{
		OperationId:  uuid.NewV4().String(),
		ResourceKey:  "node:" + testNodeId,
		ActionType:   "RestartNode",
		FencingToken: 1,
	}
}

func TestRegisterIntent_Duplicate(t *testing.T) {
	repo := &mocks.IntentRepo{}
	dup := &pgconn.PgError{Code: sql.PGERROR_CODE_UNIQUE_VIOLATION}
	repo.On("Add", mock.Anything).Return(dup).Once()

	_, err := newMonitorServer(repo).RegisterIntent(context.TODO(), registerRequest())

	assert.Equal(t, codes.AlreadyExists, status.Code(err))
	repo.AssertExpectations(t)
}

func TestRegisterIntent_StoreFailure(t *testing.T) {
	repo := &mocks.IntentRepo{}
	repo.On("Add", mock.Anything).Return(errors.New("database unavailable")).Once()

	_, err := newMonitorServer(repo).RegisterIntent(context.TODO(), registerRequest())

	assert.Equal(t, codes.Internal, status.Code(err))
	repo.AssertExpectations(t)
}

func TestRegisterIntent_Success(t *testing.T) {
	repo := &mocks.IntentRepo{}
	repo.On("Add", mock.MatchedBy(func(i *db.MonitoredIntent) bool {
		return i.CompletionRule == "substate=on" && i.Status == db.IntentWatching
	})).Return(nil).Once()

	resp, err := newMonitorServer(repo).RegisterIntent(context.TODO(), registerRequest())

	assert.NoError(t, err)
	assert.Equal(t, "substate=on", resp.Intent.CompletionRule)
	repo.AssertExpectations(t)
}

func TestCancelIntent_Unknown(t *testing.T) {
	repo := &mocks.IntentRepo{}
	opId := uuid.NewV4()
	repo.On("MarkTerminal", opId, db.IntentCancelled).
		Return(nil, gorm.ErrRecordNotFound).Once()

	_, err := newMonitorServer(repo).CancelIntent(context.TODO(),
		&pb.CancelIntentRequest{OperationId: opId.String()})

	assert.Equal(t, codes.NotFound, status.Code(err))
	repo.AssertExpectations(t)
}

func TestCancelIntent_StoreFailure(t *testing.T) {
	repo := &mocks.IntentRepo{}
	opId := uuid.NewV4()
	repo.On("MarkTerminal", opId, db.IntentCancelled).
		Return(nil, errors.New("database unavailable")).Once()

	_, err := newMonitorServer(repo).CancelIntent(context.TODO(),
		&pb.CancelIntentRequest{OperationId: opId.String()})

	assert.Equal(t, codes.Internal, status.Code(err))
	repo.AssertExpectations(t)
}

func TestCancelIntent_Success(t *testing.T) {
	repo := &mocks.IntentRepo{}
	opId := uuid.NewV4()
	repo.On("MarkTerminal", opId, db.IntentCancelled).
		Return(&db.MonitoredIntent{OperationId: opId}, nil).Once()

	_, err := newMonitorServer(repo).CancelIntent(context.TODO(),
		&pb.CancelIntentRequest{OperationId: opId.String()})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
