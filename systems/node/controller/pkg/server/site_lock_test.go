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
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/rest/client"
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/node/controller/mocks"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	"github.com/ukama/ukama/systems/node/controller/pkg"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"
	opmonpb "github.com/ukama/ukama/systems/node/operation-monitor/pb/gen"
	opmgrpb "github.com/ukama/ukama/systems/operation/manager/pb/gen"
)

const (
	siteTnode = "uk-sa2645-tnode-v0-aaaa"
	siteAnode = "uk-sa2645-anode-v0-aaaa"
	siteCnode = "uk-sa2645-cnode-v0-aaaa"
	siteID    = "site-1"
)

type siteFixture struct {
	nodes *mbmocks.NodeClient
	mgr   *mbmocks.ManagerClient
	mon   *mocks.OperationMonitor
	bus   *mbmocks.MsgBusServiceClient
	repo  *mocks.NodeLogRepo
}

func newSiteFixture(siteOf map[string]string) *siteFixture {
	f := &siteFixture{
		nodes: &mbmocks.NodeClient{},
		mgr:   &mbmocks.ManagerClient{},
		mon:   &mocks.OperationMonitor{},
		bus:   &mbmocks.MsgBusServiceClient{},
		repo:  &mocks.NodeLogRepo{},
	}
	for id, site := range siteOf {
		f.nodes.On("Get", id).Return(&creg.NodeInfo{
			Id:     id,
			Status: creg.NodeStatusInfo{Connectivity: "Online"},
			Site:   creg.NodeSiteInfo{SiteId: site},
		}, nil).Maybe()
	}
	f.nodes.On("GetNodesBySite", siteID).Return(&creg.NodesBySite{
		SiteId: siteID,
		Nodes: []creg.NodeInfo{
			{Id: strings.ToUpper(siteTnode)}, {Id: siteAnode}, {Id: siteCnode},
		},
	}, nil).Maybe()
	f.repo.On("Get", mock.Anything).Return(&db.NodeLog{}, nil).Maybe()
	return f
}

func (f *siteFixture) server() *ControllerServer {
	return NewControllerServer(testOrgName, f.repo, f.bus, nil, nil, f.nodes, f.mgr,
		f.mon, 180, 120, pkg.IsDebugMode)
}

func (f *siteFixture) startOK(nodeID, opID string, token uint64) {
	f.mgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		return req.ResourceKey == "node:"+nodeID
	})).Return(&copr.StartResponse{Operation: &copr.OperationInfo{
		Id: opID, FencingToken: token, ResourceKey: "node:" + nodeID,
	}}, nil).Once()
}

func (f *siteFixture) idle(nodeIDs ...string) {
	for _, id := range nodeIDs {
		f.mgr.On("GetByResource", "node:"+id).Return(nil, nil).Once()
	}
}

func busySite() map[string]string {
	return map[string]string{siteTnode: siteID, siteAnode: siteID, siteCnode: siteID}
}

func TestLockNodes_BusySiblingRejectsDirectAction(t *testing.T) {
	f := newSiteFixture(busySite())
	f.startOK(siteTnode, "op-restart", 5)
	f.idle(siteAnode)
	f.mgr.On("GetByResource", "node:"+siteCnode).Return(&copr.OperationInfo{
		Id: "op-update", Type: "UpdateSoftware", Status: copr.StatusRunning,
	}, nil).Once()
	f.mgr.On("Fail", "op-restart", uint64(5), pkg.ServiceName,
		mock.MatchedBy(func(reason string) bool {
			return strings.HasPrefix(reason, "not started: ") &&
				strings.Contains(reason, "UpdateSoftware")
		})).Return(&copr.OperationInfo{}, nil).Once()

	_, err := f.server().RestartNode(context.TODO(),
		&pb.RestartNodeRequest{NodeId: siteTnode})

	assert.Equal(t, codes.AlreadyExists, status.Code(err))
	assert.Contains(t, status.Convert(err).Message(), "op-update")
	f.mgr.AssertExpectations(t)
	f.mon.AssertNotCalled(t, "Register", mock.Anything)
	f.bus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
}

func TestLockNodes_LockedTargetIsAlreadyExists(t *testing.T) {
	f := newSiteFixture(busySite())
	f.mgr.On("Start", mock.Anything).Return(nil, &client.ErrorStatus{
		StatusCode: http.StatusConflict, RawError: errors.New("locked"),
	}).Once()

	_, err := f.server().ToggleService(context.TODO(),
		&pb.ToggleServiceRequest{NodeId: siteTnode, State: "on"})

	assert.Equal(t, codes.AlreadyExists, status.Code(err))
	f.mgr.AssertNotCalled(t, "GetByResource", mock.Anything)
}

func TestLockNodes_SiteReadFailureIsUnavailable(t *testing.T) {
	f := newSiteFixture(nil)
	f.nodes.On("Get", siteTnode).Return(&creg.NodeInfo{
		Status: creg.NodeStatusInfo{Connectivity: "Online"},
		Site:   creg.NodeSiteInfo{SiteId: "site-x"},
	}, nil).Once()
	f.nodes.On("GetNodesBySite", "site-x").Return(nil, errors.New("down")).Once()

	_, err := f.server().ToggleService(context.TODO(),
		&pb.ToggleServiceRequest{NodeId: siteTnode, State: "on"})

	assert.Equal(t, codes.Unavailable, status.Code(err))
	f.mgr.AssertNotCalled(t, "Start", mock.Anything)
}

func TestLockNodes_NodesFromDifferentSitesAreRejected(t *testing.T) {
	f := newSiteFixture(map[string]string{siteTnode: siteID, siteAnode: "site-2"})

	_, err := f.server().RestartNodes(context.TODO(),
		&pb.RestartNodesRequest{NodeIds: []string{siteTnode, siteAnode}})

	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	f.mgr.AssertNotCalled(t, "Start", mock.Anything)
}

func TestLockNodes_RegisterFailureReleasesEveryLock(t *testing.T) {
	f := newSiteFixture(busySite())
	f.startOK(siteTnode, "op-t", 1)
	f.startOK(siteAnode, "op-a", 2)
	f.idle(siteCnode)
	f.mon.On("Register", mock.MatchedBy(func(r *opmonpb.RegisterIntentRequest) bool {
		return r.OperationId == "op-t"
	})).Return(&opmonpb.RegisterIntentResponse{}, nil).Once()
	f.mon.On("Register", mock.MatchedBy(func(r *opmonpb.RegisterIntentRequest) bool {
		return r.OperationId == "op-a"
	})).Return(nil, errors.New("monitor down")).Once()
	f.mgr.On("Fail", "op-t", uint64(1), pkg.ServiceName, mock.Anything).
		Return(&copr.OperationInfo{}, nil).Once()
	f.mgr.On("Fail", "op-a", uint64(2), pkg.ServiceName, mock.Anything).
		Return(&copr.OperationInfo{}, nil).Once()
	f.mon.On("Cancel", &opmonpb.CancelIntentRequest{OperationId: "op-t"}).
		Return(&opmonpb.CancelIntentResponse{}, nil).Once()

	_, err := f.server().RestartNodes(context.TODO(),
		&pb.RestartNodesRequest{NodeIds: []string{siteTnode, siteAnode}})

	assert.Equal(t, codes.Internal, status.Code(err))
	f.mgr.AssertExpectations(t)
	f.mon.AssertExpectations(t)
	f.mon.AssertNotCalled(t, "Cancel", &opmonpb.CancelIntentRequest{OperationId: "op-a"})
}

func TestRestartNodes(t *testing.T) {
	t.Run("LocksTheSiteOnceAndRestartsEveryNode", func(t *testing.T) {
		f := newSiteFixture(busySite())
		f.startOK(siteTnode, "op-t", 1)
		f.startOK(siteAnode, "op-a", 2)
		f.idle(siteCnode)
		f.mon.On("Register", mock.Anything).
			Return(&opmonpb.RegisterIntentResponse{}, nil).Twice()
		f.mgr.On("MarkRunning", "op-t", uint64(1), pkg.ServiceName).
			Return(&copr.OperationInfo{}, nil).Once()
		f.mgr.On("MarkRunning", "op-a", uint64(2), pkg.ServiceName).
			Return(&copr.OperationInfo{}, nil).Once()
		f.bus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Twice()

		resp, err := f.server().RestartNodes(context.TODO(),
			&pb.RestartNodesRequest{NodeIds: []string{siteTnode, siteAnode, siteTnode}})

		assert.NoError(t, err)
		assert.Len(t, resp.Operations, 2)
		for _, op := range resp.Operations {
			assert.Equal(t, opmgrpb.OperationStatus_RUNNING.String(), op.Status)
		}
		f.mgr.AssertExpectations(t)
		f.mgr.AssertNotCalled(t, "GetByResource", "node:"+siteTnode)
		f.mgr.AssertNotCalled(t, "GetByResource", "node:"+siteAnode)
	})

	t.Run("PublishFailureFailsTheRemainingNodes", func(t *testing.T) {
		f := newSiteFixture(busySite())
		f.startOK(siteTnode, "op-t", 1)
		f.startOK(siteAnode, "op-a", 2)
		f.idle(siteCnode)
		f.mon.On("Register", mock.Anything).
			Return(&opmonpb.RegisterIntentResponse{}, nil).Twice()
		f.mgr.On("MarkRunning", mock.Anything, mock.Anything, pkg.ServiceName).
			Return(&copr.OperationInfo{}, nil).Twice()
		f.bus.On("PublishRequest", mock.Anything,
			mock.MatchedBy(func(m *epb.NodeFeederMessage) bool {
				return m.NodeId == siteTnode
			})).Return(nil).Once()
		f.bus.On("PublishRequest", mock.Anything, mock.Anything).
			Return(errors.New("bus down")).Once()
		f.mgr.On("Fail", "op-a", uint64(2), pkg.ServiceName,
			mock.MatchedBy(func(r string) bool {
				return strings.HasPrefix(r, "publish failed")
			})).Return(&copr.OperationInfo{}, nil).Once()
		f.mon.On("Cancel", &opmonpb.CancelIntentRequest{OperationId: "op-a"}).
			Return(&opmonpb.CancelIntentResponse{}, nil).Once()

		_, err := f.server().RestartNodes(context.TODO(),
			&pb.RestartNodesRequest{NodeIds: []string{siteTnode, siteAnode}})

		assert.Equal(t, codes.Internal, status.Code(err))
		f.mgr.AssertExpectations(t)
		f.mgr.AssertNotCalled(t, "Fail", "op-t", mock.Anything, mock.Anything,
			mock.Anything)
	})

	t.Run("MarkRunningFailureFailsEveryNode", func(t *testing.T) {
		f := newSiteFixture(busySite())
		f.startOK(siteTnode, "op-t", 1)
		f.startOK(siteAnode, "op-a", 2)
		f.idle(siteCnode)
		f.mon.On("Register", mock.Anything).
			Return(&opmonpb.RegisterIntentResponse{}, nil).Twice()
		f.mgr.On("MarkRunning", "op-t", uint64(1), pkg.ServiceName).
			Return(nil, errors.New("manager down")).Once()
		f.mgr.On("Fail", mock.Anything, mock.Anything, pkg.ServiceName,
			mock.MatchedBy(func(r string) bool {
				return strings.HasPrefix(r, "mark running failed")
			})).Return(&copr.OperationInfo{}, nil).Twice()
		f.mon.On("Cancel", mock.Anything).
			Return(&opmonpb.CancelIntentResponse{}, nil).Twice()

		_, err := f.server().RestartNodes(context.TODO(),
			&pb.RestartNodesRequest{NodeIds: []string{siteTnode, siteAnode}})

		assert.Equal(t, codes.Internal, status.Code(err))
		f.mgr.AssertExpectations(t)
		f.bus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
	})

	t.Run("UnregisteredNodeIsRejected", func(t *testing.T) {
		f := newSiteFixture(nil)
		repo := &mocks.NodeLogRepo{}
		repo.On("Get", siteTnode).Return(nil, errors.New("not found")).Once()
		f.repo = repo

		_, err := f.server().RestartNodes(context.TODO(),
			&pb.RestartNodesRequest{NodeIds: []string{siteTnode}})

		assert.Equal(t, codes.InvalidArgument, status.Code(err))
		f.mgr.AssertNotCalled(t, "Start", mock.Anything)
	})
}

func TestToggleRadios(t *testing.T) {
	f := newSiteFixture(busySite())
	f.startOK(siteTnode, "op-t", 1)
	f.startOK(siteAnode, "op-a", 2)
	f.idle(siteCnode)
	f.mon.On("Register", mock.Anything).
		Return(&opmonpb.RegisterIntentResponse{}, nil).Twice()
	f.mgr.On("MarkRunning", mock.Anything, mock.Anything, pkg.ServiceName).
		Return(&copr.OperationInfo{}, nil).Twice()
	f.bus.On("PublishRequest", mock.Anything, mock.Anything).Return(nil).Twice()
	f.mgr.On("Complete", "op-t", uint64(1), pkg.ServiceName, "dispatched").
		Return(&copr.OperationInfo{}, nil).Once()
	f.mgr.On("Complete", "op-a", uint64(2), pkg.ServiceName, "dispatched").
		Return(&copr.OperationInfo{}, nil).Once()
	f.mon.On("Cancel", mock.Anything).Return(&opmonpb.CancelIntentResponse{}, nil).Twice()

	resp, err := f.server().ToggleRadios(context.TODO(), &pb.ToggleRadiosRequest{
		NodeIds: []string{siteTnode, siteAnode}, State: "off",
	})

	assert.NoError(t, err)
	assert.Len(t, resp.Operations, 2)
	for _, op := range resp.Operations {
		assert.Equal(t, opmgrpb.OperationStatus_SUCCESS.String(), op.Status)
	}
	f.mgr.AssertExpectations(t)
	f.mon.AssertExpectations(t)
}

func TestValidNodeIds(t *testing.T) {
	ids, err := validNodeIds([]string{strings.ToUpper(siteTnode), siteTnode, siteAnode})
	assert.NoError(t, err)
	assert.Equal(t, []string{siteTnode, siteAnode}, ids)

	_, err = validNodeIds(nil)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))

	_, err = validNodeIds([]string{"bad"})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestLockNodes_WithoutManagerIsUnavailable(t *testing.T) {
	c := &ControllerServer{nodeClient: onlineNodeClient()}
	_, err := c.lockNodes("RestartNode", siteTnode)
	assert.Equal(t, codes.Unavailable, status.Code(err))
}
