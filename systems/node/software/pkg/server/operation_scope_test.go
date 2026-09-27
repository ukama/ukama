/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	copr "github.com/ukama/ukama/systems/common/rest/client/operation"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/node/software/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSoftwareUpdateAdmissionChecksAllSiteNodes(t *testing.T) {
	nodes := &mbmocks.NodeClient{}
	nodes.On("Get", "cnode").Return(&creg.NodeInfo{Site: creg.NodeSiteInfo{SiteId: "site"}}, nil).Once()
	nodes.On("GetNodesBySite", "site").Return(&creg.NodesBySite{Nodes: []creg.NodeInfo{
		{Id: "tnode"}, {Id: "cnode"}, {Id: "anode"},
	}}, nil).Once()
	manager := &mbmocks.ManagerClient{}
	manager.On("Start", mock.MatchedBy(func(req copr.StartRequest) bool {
		assert.ElementsMatch(t, []string{"node:tnode", "node:cnode", "node:anode"}, req.ConflictResourceKeys)
		return req.ResourceKey == "node:cnode" && req.Type == "UpdateSoftware"
	})).Return(nil, status.Error(codes.AlreadyExists, "tower restart active")).Once()
	monitor := &mocks.OperationMonitor{}
	s := &SoftwareServer{nodeClient: nodes, opManager: manager, opMonitor: monitor}
	_, err := s.acquireAndRegister("UpdateSoftware", "node:cnode")
	require.Equal(t, codes.AlreadyExists, status.Code(err))
	assert.Empty(t, monitor.Calls)
	manager.AssertExpectations(t)
	nodes.AssertExpectations(t)
}
