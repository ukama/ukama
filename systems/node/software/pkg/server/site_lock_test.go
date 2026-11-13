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
	"github.com/ukama/ukama/systems/common/rest/client"
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/node/software/mocks"
	pb "github.com/ukama/ukama/systems/node/software/pb/gen"
	"github.com/ukama/ukama/systems/node/software/pkg"
	"github.com/ukama/ukama/systems/node/software/pkg/db"
)

const (
	updateSiteID = "site-1"
	siblingTower = "uk-sa2156-tnode-a1-xxxx"
)

type updateFixture struct {
	nodes *mbmocks.NodeClient
	mgr   *mbmocks.ManagerClient
	mon   *mocks.OperationMonitor
	bus   *mbmocks.MsgBusServiceClient
	repo  *mocks.SoftwareRepo
}

func newUpdateFixture(t *testing.T) *updateFixture {
	f := &updateFixture{
		nodes: &mbmocks.NodeClient{},
		mgr:   &mbmocks.ManagerClient{},
		mon:   &mocks.OperationMonitor{},
		bus:   &mbmocks.MsgBusServiceClient{},
		repo:  mocks.NewSoftwareRepo(t),
	}
	f.repo.On("List", testNodeIdNormalized, ukama.UpdateAvailable, testAppNameForUpdate).
		Return([]*db.Software{dbSoftwareFixture()}, nil).Once()
	f.nodes.On("Get", testNodeIdNormalized).Return(&creg.NodeInfo{
		Status: creg.NodeStatusInfo{Connectivity: "Online"},
		Site:   creg.NodeSiteInfo{SiteId: updateSiteID},
	}, nil).Once()
	return f
}

func (f *updateFixture) siteNodes(err error) {
	if err != nil {
		f.nodes.On("GetNodesBySite", updateSiteID).Return(nil, err).Once()
		return
	}
	f.nodes.On("GetNodesBySite", updateSiteID).Return(&creg.NodesBySite{
		Nodes: []creg.NodeInfo{{Id: testNodeIdNormalized}, {Id: strings.ToUpper(siblingTower)}},
	}, nil).Once()
}

func (f *updateFixture) update(t *testing.T) error {
	s := NewSoftwareServer(testOrgName, f.repo, mocks.NewAppRepo(t), mocks.NewNodeRepo(t),
		nil, nil, nil, f.bus, false, []string{testNodeGwIP}, f.mgr, f.mon, 360, 300,
		f.nodes)
	_, err := s.UpdateSoftware(context.Background(), &pb.UpdateSoftwareRequest{
		NodeId: testNodeId, Name: testAppNameForUpdate, Tag: testTagVersion})
	return err
}

func TestUpdateSoftware_SiteLock(t *testing.T) {
	t.Run("BusySiblingIsAlreadyExists", func(t *testing.T) {
		f := newUpdateFixture(t)
		f.siteNodes(nil)
		op := &copr.OperationInfo{Id: "op-update", FencingToken: 9,
			ResourceKey: "node:" + testNodeIdNormalized}
		f.mgr.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
			return req.ResourceKey == "node:"+testNodeIdNormalized &&
				req.LeaseSeconds == 360
		})).Return(&copr.StartResponse{Operation: op}, nil).Once()
		f.mgr.On("GetByResource", "node:"+siblingTower).Return(&copr.OperationInfo{
			Id: "op-restart", Type: "RestartNode", Status: copr.StatusRunning,
		}, nil).Once()
		f.mgr.On("Fail", "op-update", uint64(9), pkg.ServiceName,
			mock.MatchedBy(func(r string) bool {
				return strings.HasPrefix(r, "not started: ") &&
					strings.Contains(r, "RestartNode")
			})).Return(op, nil).Once()

		err := f.update(t)

		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		assert.Contains(t, status.Convert(err).Message(), "op-restart")
		f.mgr.AssertExpectations(t)
		f.mgr.AssertNotCalled(t, "GetByResource", "node:"+testNodeIdNormalized)
		f.mon.AssertNotCalled(t, "Register", mock.Anything)
		f.bus.AssertNotCalled(t, "PublishRequest", mock.Anything, mock.Anything)
	})

	t.Run("LockedNodeIsAlreadyExists", func(t *testing.T) {
		f := newUpdateFixture(t)
		f.siteNodes(nil)
		f.mgr.On("Start", mock.Anything).Return(nil, &client.ErrorStatus{
			StatusCode: http.StatusConflict, RawError: errors.New("locked"),
		}).Once()

		err := f.update(t)

		assert.Equal(t, codes.AlreadyExists, status.Code(err))
		f.mgr.AssertNotCalled(t, "GetByResource", mock.Anything)
		f.mgr.AssertNotCalled(t, "Fail", mock.Anything, mock.Anything, mock.Anything,
			mock.Anything)
	})

	t.Run("ManagerFailureIsUnavailable", func(t *testing.T) {
		f := newUpdateFixture(t)
		f.siteNodes(nil)
		f.mgr.On("Start", mock.Anything).Return(nil, &client.ErrorStatus{
			StatusCode: http.StatusServiceUnavailable, RawError: errors.New("down"),
		}).Once()

		err := f.update(t)

		assert.Equal(t, codes.Unavailable, status.Code(err))
	})

	t.Run("SiteReadFailureIsUnavailable", func(t *testing.T) {
		f := newUpdateFixture(t)
		f.siteNodes(errors.New("registry down"))

		err := f.update(t)

		assert.Equal(t, codes.Unavailable, status.Code(err))
		f.mgr.AssertNotCalled(t, "Start", mock.Anything)
	})
}
