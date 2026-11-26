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
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/node/controller/mocks"
	"github.com/ukama/ukama/systems/node/controller/pkg"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	opmonpb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
)

const (
	resultRoute   = "event.cloud.local.{{ .Org}}.messaging.nodefeeder.request.result"
	resultOpId    = "op-result"
	resultToken   = uint64(4)
	servicePath   = "/device/v1/service"
	rebootPath    = "/device/v1/reboot"
	retriesReason = "retries exhausted"
)

type resultFixture struct {
	mgr *mbmocks.ManagerClient
	mon *mocks.OperationMonitor
	srv *ControllerEventServer
}

func newResultFixture() *resultFixture {
	f := &resultFixture{mgr: &mbmocks.ManagerClient{}, mon: &mocks.OperationMonitor{}}
	s := NewControllerServer(testOrgName, nil, nil, nil, nil, nil, f.mgr, f.mon, 30, 60,
		false)
	f.srv = NewControllerEventServer(testOrgName, s)
	return f
}

func (f *resultFixture) send(t *testing.T, msg proto.Message) error {
	body, err := anypb.New(msg)
	require.NoError(t, err)

	_, err = f.srv.EventNotification(context.TODO(), &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, resultRoute),
		Msg:        body,
	})
	return err
}

func result(path string, outcome epb.NodeFeederRequestResult_Outcome, code int32,
	reason string) *epb.NodeFeederRequestResult {
	return &epb.NodeFeederRequestResult{OperationId: resultOpId, FencingToken: resultToken,
		NodeId: "uk-983794-tnode-78-7830", HttpMethod: "POST", Path: path,
		HttpStatus: code, Outcome: outcome, Reason: reason}
}

func (f *resultFixture) expectCancel() {
	f.mon.On("Cancel", &opmonpb.CancelIntentRequest{OperationId: resultOpId}).
		Return(&opmonpb.CancelIntentResponse{}, nil).Once()
}

func (f *resultFixture) assertExpectations(t *testing.T) {
	f.mgr.AssertExpectations(t)
	f.mon.AssertExpectations(t)
}

func TestNodeFeederResult_ToggleDeliveredCompletes(t *testing.T) {
	f := newResultFixture()
	f.mgr.On("Complete", resultOpId, resultToken, pkg.ServiceName, reasonNodeAccepted).
		Return(&copr.OperationInfo{}, nil).Once()
	f.expectCancel()

	assert.NoError(t, f.send(t,
		result(servicePath, epb.NodeFeederRequestResult_DELIVERED, 202, "")))
	f.assertExpectations(t)
}

func TestNodeFeederResult_RestartDeliveredMarksDispatched(t *testing.T) {
	for name, dispatchErr := range map[string]error{
		"Marked":        nil,
		"MonitorFailed": assert.AnError,
	} {
		t.Run(name, func(t *testing.T) {
			f := newResultFixture()
			f.mon.On("Dispatch", &opmonpb.DispatchIntentRequest{OperationId: resultOpId}).
				Return(&opmonpb.DispatchIntentResponse{}, dispatchErr).Once()

			assert.NoError(t, f.send(t,
				result(rebootPath, epb.NodeFeederRequestResult_DELIVERED, 202, "")))
			f.assertExpectations(t)
			f.mgr.AssertNotCalled(t, "Complete", mock.Anything, mock.Anything,
				mock.Anything, mock.Anything)
			f.mon.AssertNotCalled(t, "Cancel", mock.Anything)
		})
	}
}

func TestDispatchIntent_WithoutMonitorIsNoop(t *testing.T) {
	s := NewControllerServer(testOrgName, nil, nil, nil, nil, nil, nil, nil, 30, 60, false)

	assert.NotPanics(t, func() {
		s.dispatchIntent(&copr.OperationInfo{Id: resultOpId}, rebootPath)
	})
}

func TestNodeFeederResult_FailureFailsOperation(t *testing.T) {
	for name, tc := range map[string]struct {
		path    string
		outcome epb.NodeFeederRequestResult_Outcome
		code    int32
		reason  string
		failure string
	}{
		"ToggleRejected": {servicePath, epb.NodeFeederRequestResult_REJECTED, 409,
			"server error: 409", fmt.Sprintf(reasonNodeRejectedFmt, 409, "server error: 409")},
		"RestartRejected": {rebootPath, epb.NodeFeederRequestResult_REJECTED, 409,
			"server error: 409", fmt.Sprintf(reasonNodeRejectedFmt, 409, "server error: 409")},
		"ToggleDropped": {servicePath, epb.NodeFeederRequestResult_DROPPED, 0,
			retriesReason, fmt.Sprintf(reasonNodeDroppedFmt, retriesReason)},
		"RestartDropped": {rebootPath, epb.NodeFeederRequestResult_DROPPED, 0,
			retriesReason, fmt.Sprintf(reasonNodeDroppedFmt, retriesReason)},
	} {
		t.Run(name, func(t *testing.T) {
			f := newResultFixture()
			f.mgr.On("Fail", resultOpId, resultToken, pkg.ServiceName, tc.failure).
				Return(&copr.OperationInfo{}, nil).Once()
			f.expectCancel()

			assert.NoError(t, f.send(t, result(tc.path, tc.outcome, tc.code, tc.reason)))
			f.assertExpectations(t)
		})
	}
}

func TestNodeFeederResult_LateResultIsOnlyLogged(t *testing.T) {
	f := newResultFixture()
	f.mgr.On("Fail", resultOpId, resultToken, pkg.ServiceName, mock.Anything).
		Return(nil, copr.ErrNotRunnable).Once()

	assert.NoError(t, f.send(t,
		result(servicePath, epb.NodeFeederRequestResult_REJECTED, 409, "late")))
	f.assertExpectations(t)
	f.mon.AssertNotCalled(t, "Cancel", mock.Anything)
}

func TestNodeFeederResult_Ignored(t *testing.T) {
	noOp := result(servicePath, epb.NodeFeederRequestResult_DELIVERED, 200, "")
	noOp.OperationId = ""

	for name, msg := range map[string]*epb.NodeFeederRequestResult{
		"WithoutOperation": noOp,
		"UnknownOutcome": result(servicePath, epb.NodeFeederRequestResult_UNKNOWN, 0,
			""),
	} {
		t.Run(name, func(t *testing.T) {
			f := newResultFixture()

			assert.NoError(t, f.send(t, msg))
			f.mgr.AssertNotCalled(t, "Complete", mock.Anything, mock.Anything,
				mock.Anything, mock.Anything)
			f.mgr.AssertNotCalled(t, "Fail", mock.Anything, mock.Anything, mock.Anything,
				mock.Anything)
		})
	}
}

func TestNodeFeederResult_BadPayload(t *testing.T) {
	f := newResultFixture()

	assert.Error(t, f.send(t, &epb.EventRegistryNodeDelete{NodeId: "x"}))
}
