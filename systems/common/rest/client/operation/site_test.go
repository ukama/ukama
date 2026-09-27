/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package operation_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mbmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/operation"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
)

func TestNodeConflictResources(t *testing.T) {
	for _, scenario := range []string{"unassigned", "site", "registry error", "missing member"} {
		t.Run(scenario, func(t *testing.T) {
			nodes := &mbmocks.NodeClient{}
			node := &creg.NodeInfo{Id: "c"}
			if scenario != "unassigned" {
				node.Site.SiteId = "site"
			}
			nodes.On("Get", "c").Return(node, nil).Once()
			if scenario == "registry error" {
				nodes.On("GetNodesBySite", "site").Return(nil, assert.AnError).Once()
			} else if scenario != "unassigned" {
				members := []creg.NodeInfo{{Id: "t"}, {Id: "a"}, {Id: "a"}}
				if scenario == "site" { members = append(members, creg.NodeInfo{Id: "c"}) }
				nodes.On("GetNodesBySite", "site").Return(&creg.NodesBySite{Nodes: members}, nil).Once()
			}
			keys, err := operation.NodeConflictResources(nodes, "c")
			if scenario == "registry error" || scenario == "missing member" {
				require.Error(t, err)
				assert.Nil(t, keys)
			} else {
				require.NoError(t, err)
				want := []string{"node:c"}
				if scenario == "site" { want = []string{"node:a", "node:c", "node:t"} }
				assert.Equal(t, want, keys)
			}
			nodes.AssertExpectations(t)
		})
	}
}

func TestScopedStartRequiresBackendAcknowledgement(t *testing.T) {
	for _, checked := range []bool{false, true} {
		manager := operation.NewManagerClient("http://operation")
		manager.R.C.SetTransport(client.RoundTripFunc(func(req *http.Request) *http.Response {
			var body operation.StartRequest
			require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
			assert.Equal(t, []string{"node:c", "node:t", "node:a"}, body.ConflictResourceKeys)
			assert.Equal(t, []string{"node:a"}, body.AdditionalResourceKeys)
			data, err := json.Marshal(operation.StartResponse{ConflictsChecked: checked})
			require.NoError(t, err)
			return &http.Response{StatusCode: 201, Header: make(http.Header), Body: io.NopCloser(bytes.NewReader(data))}
		}))
		_, err := manager.Start(operation.StartRequest{Type: "RestartNode", System: "node",
			ResourceKey: "node:t", AdditionalResourceKeys: []string{"node:a"},
			ConflictResourceKeys: []string{"node:c", "node:t", "node:a"}})
		if checked { require.NoError(t, err) } else { require.Error(t, err) }
	}
}

func TestForceUnlockSendsActorAndReason(t *testing.T) {
	manager := operation.NewManagerClient("http://operation")
	manager.R.C.SetTransport(client.RoundTripFunc(func(req *http.Request) *http.Response {
		assert.Equal(t, "/v1/operations/reservation/force-unlock", req.URL.Path)
		assert.Equal(t, http.MethodPost, req.Method)
		var body operation.ForceUnlockRequest
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		assert.Equal(t, "site-controller", body.Actor)
		assert.Equal(t, "not dispatched", body.Reason)
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(bytes.NewBufferString(`{"operation":{"status":"CANCELLED"}}`))}
	}))
	_, err := manager.ForceUnlock("reservation", "site-controller", "not dispatched")
	require.NoError(t, err)
}
