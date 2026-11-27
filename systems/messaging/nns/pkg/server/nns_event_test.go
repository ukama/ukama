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
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/messaging/nns/pkg"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
)

const (
	testEventOrg  = "test-org"
	assignRoute   = "event.cloud.local.{{ .Org}}.registry.node.node.assign"
	releaseRoute  = "event.cloud.local.{{ .Org}}.registry.node.node.release"
	assignedRoute = "event.cloud.local.{{ .Org}}.registry.node.node.assigned"
	onlineRoute   = "event.cloud.local.{{ .Org}}.messaging.mesh.node.online"
	offlineRoute  = "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline"
)

func registryError(code int) error {
	return fmt.Errorf("GetNode failure: %w", fmt.Errorf("rest api GET failure: %w",
		&client.ErrorStatus{StatusCode: code, RawError: errors.New("registry error")}))
}

func onlineEvent() *epb.NodeOnlineEvent {
	return &epb.NodeOnlineEvent{
		NodeId:       testValidNodeID,
		NodeIp:       "10.0.0.2",
		NodePort:     3001,
		MeshIp:       "10.1.0.2",
		MeshPort:     4001,
		MeshHostName: "mesh-2",
	}
}

func testEvent(t *testing.T, route string, msg proto.Message) *epb.Event {
	anyMsg, err := anypb.New(msg)
	require.NoError(t, err)

	return &epb.Event{RoutingKey: msgbus.PrepareRoute(testEventOrg, route), Msg: anyMsg}
}

func storedNode() *pkg.NodeMeshMap {
	return &pkg.NodeMeshMap{
		NodeId:       testValidNodeID,
		NodeIp:       "10.0.0.1",
		NodePort:     3000,
		MeshIp:       "10.1.0.1",
		MeshPort:     4000,
		MeshHostName: "mesh-1",
		Site:         "old-site",
		Network:      "old-network",
	}
}

func TestNnsEventServerNodeSiteEvents(t *testing.T) {
	t.Run("assignSetsSiteAndNetwork", func(t *testing.T) {
		var added []pkg.NodeMeshMap
		store := &fakeNnsStore{
			get: func(_ context.Context, _ string) (*pkg.NodeMeshMap, error) {
				return storedNode(), nil
			},
			add: func(_ context.Context, obj pkg.NodeMeshMap) error {
				added = append(added, obj)
				return nil
			},
		}
		srv := NewNnsEventServer(testEventOrg, nil, NewNnsServer(store, nil, nil),
			testEventOrg)

		_, err := srv.EventNotification(context.Background(),
			testEvent(t, assignRoute, &epb.EventRegistryNodeAssign{
				NodeId: testValidNodeID, Site: "site-1", Network: "network-1",
			}))

		require.NoError(t, err)
		require.Len(t, added, 1)
		assert.Equal(t, "site-1", added[0].Site)
		assert.Equal(t, "network-1", added[0].Network)
		assert.Equal(t, testEventOrg, added[0].Org)
		assert.Equal(t, "10.0.0.1", added[0].NodeIp)
		assert.Equal(t, "10.1.0.1", added[0].MeshIp)
	})

	t.Run("releaseClearsSiteAndNetwork", func(t *testing.T) {
		var added []pkg.NodeMeshMap
		store := &fakeNnsStore{
			get: func(_ context.Context, _ string) (*pkg.NodeMeshMap, error) {
				return storedNode(), nil
			},
			add: func(_ context.Context, obj pkg.NodeMeshMap) error {
				added = append(added, obj)
				return nil
			},
		}
		srv := NewNnsEventServer(testEventOrg, nil, NewNnsServer(store, nil, nil),
			testEventOrg)

		_, err := srv.EventNotification(context.Background(),
			testEvent(t, releaseRoute, &epb.NodeReleasedEvent{NodeId: testValidNodeID}))

		require.NoError(t, err)
		require.Len(t, added, 1)
		assert.Empty(t, added[0].Site)
		assert.Empty(t, added[0].Network)
		assert.Equal(t, "10.0.0.1", added[0].NodeIp)
	})

	noStoreCases := []struct {
		name    string
		event   func(t *testing.T) *epb.Event
		wantErr bool
	}{
		{"legacyAssignedKeyIsIgnored", func(t *testing.T) *epb.Event {
			return testEvent(t, assignedRoute,
				&epb.EventRegistryNodeAssign{NodeId: testValidNodeID})
		}, false},
		{"wrongPayloadOnAssignFails", func(t *testing.T) *epb.Event {
			return testEvent(t, assignRoute,
				&epb.NodeReleasedEvent{NodeId: testValidNodeID})
		}, true},
		{"wrongPayloadOnReleaseFails", func(t *testing.T) *epb.Event {
			return testEvent(t, releaseRoute,
				&epb.EventRegistryNodeAssign{NodeId: testValidNodeID})
		}, true},
		{"wrongPayloadOnOnlineFails", func(t *testing.T) *epb.Event {
			return testEvent(t, onlineRoute,
				&epb.NodeOfflineEvent{NodeId: testValidNodeID})
		}, true},
		{"wrongPayloadOnOfflineFails", func(t *testing.T) *epb.Event {
			return testEvent(t, offlineRoute,
				&epb.NodeOnlineEvent{NodeId: testValidNodeID})
		}, true},
		{"offlineIsAcknowledged", func(t *testing.T) *epb.Event {
			return testEvent(t, offlineRoute,
				&epb.NodeOfflineEvent{NodeId: testValidNodeID})
		}, false},
	}

	for _, tc := range noStoreCases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			store := &fakeNnsStore{
				get: func(_ context.Context, _ string) (*pkg.NodeMeshMap, error) {
					calls++
					return storedNode(), nil
				},
				add: func(_ context.Context, _ pkg.NodeMeshMap) error {
					calls++
					return nil
				},
			}
			srv := NewNnsEventServer(testEventOrg, nil, NewNnsServer(store, nil, nil),
				testEventOrg)

			_, err := srv.EventNotification(context.Background(), tc.event(t))

			assert.Equal(t, tc.wantErr, err != nil, err)
			assert.Zero(t, calls)
		})
	}

	storeFailureCases := []struct {
		name   string
		event  func(t *testing.T) *epb.Event
		getErr error
		addErr error
		added  int
	}{
		{"assignForUnknownNodeFails", func(t *testing.T) *epb.Event {
			return testEvent(t, assignRoute,
				&epb.EventRegistryNodeAssign{NodeId: testValidNodeID})
		}, errFakeNotImplemented, nil, 0},
		{"releaseForUnknownNodeFails", func(t *testing.T) *epb.Event {
			return testEvent(t, releaseRoute,
				&epb.NodeReleasedEvent{NodeId: testValidNodeID})
		}, errFakeNotImplemented, nil, 0},
		{"assignStoreWriteFails", func(t *testing.T) *epb.Event {
			return testEvent(t, assignRoute,
				&epb.EventRegistryNodeAssign{NodeId: testValidNodeID})
		}, nil, errFakeNotImplemented, 1},
		{"releaseStoreWriteFails", func(t *testing.T) *epb.Event {
			return testEvent(t, releaseRoute,
				&epb.NodeReleasedEvent{NodeId: testValidNodeID})
		}, nil, errFakeNotImplemented, 1},
	}

	for _, tc := range storeFailureCases {
		t.Run(tc.name, func(t *testing.T) {
			added := 0
			store := &fakeNnsStore{
				get: func(_ context.Context, _ string) (*pkg.NodeMeshMap, error) {
					if tc.getErr != nil {
						return nil, tc.getErr
					}
					return storedNode(), nil
				},
				add: func(_ context.Context, _ pkg.NodeMeshMap) error {
					added++
					return tc.addErr
				},
			}
			srv := NewNnsEventServer(testEventOrg, nil, NewNnsServer(store, nil, nil),
				testEventOrg)

			_, err := srv.EventNotification(context.Background(), tc.event(t))

			assert.Error(t, err)
			assert.Equal(t, tc.added, added)
		})
	}
}

func TestNnsEventServerNodeOnlineEvent(t *testing.T) {
	registered := &creg.NodeInfo{
		Id:   testValidNodeID,
		Site: creg.NodeSiteInfo{SiteId: "site-1", NetworkId: "network-1"},
	}

	cases := []struct {
		name        string
		node        *creg.NodeInfo
		registryErr error
		storeErr    error
		addErr      error
		wantSite    string
		wantNetwork string
		wantErr     bool
	}{
		{name: "registeredNodeTakesRegistrySite", node: registered,
			wantSite: "site-1", wantNetwork: "network-1"},
		{name: "unregisteredNodeStoresAddressOnly",
			registryErr: registryError(http.StatusNotFound)},
		{name: "registryFailureKeepsStoredSite",
			registryErr: registryError(http.StatusServiceUnavailable),
			wantSite:    "old-site", wantNetwork: "old-network"},
		{name: "registryUnreachableKeepsStoredSite",
			registryErr: errors.New("connection refused"),
			wantSite:    "old-site", wantNetwork: "old-network"},
		{name: "registryFailureWithoutStoredEntry",
			registryErr: registryError(http.StatusServiceUnavailable),
			storeErr:    errFakeNotImplemented},
		{name: "storeWriteFailureIsReturned", node: registered,
			addErr: errFakeNotImplemented, wantSite: "site-1", wantNetwork: "network-1",
			wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nodeClient := &mocks.NodeClient{}
			nodeClient.On("Get", testValidNodeID).Return(tc.node, tc.registryErr).Once()

			var added []pkg.NodeMeshMap
			deleted := 0
			store := &fakeNnsStore{
				get: func(_ context.Context, _ string) (*pkg.NodeMeshMap, error) {
					if tc.storeErr != nil {
						return nil, tc.storeErr
					}
					return storedNode(), nil
				},
				add: func(_ context.Context, obj pkg.NodeMeshMap) error {
					added = append(added, obj)
					return tc.addErr
				},
				delete: func(_ context.Context, _ string) error {
					deleted++
					return nil
				},
			}
			nnsSrv := NewNnsServer(store, &pkg.Config{OrgName: testEventOrg}, nil)
			srv := NewNnsEventServer(testEventOrg, nodeClient, nnsSrv, testEventOrg)

			_, err := srv.EventNotification(context.Background(),
				testEvent(t, onlineRoute, onlineEvent()))

			assert.Equal(t, tc.wantErr, err != nil, err)
			assert.Zero(t, deleted)
			require.Len(t, added, 1)
			assert.Equal(t, tc.wantSite, added[0].Site)
			assert.Equal(t, tc.wantNetwork, added[0].Network)
			assert.Equal(t, testEventOrg, added[0].Org)
			assert.Equal(t, "10.0.0.2", added[0].NodeIp)
			assert.Equal(t, int32(3001), added[0].NodePort)
			assert.Equal(t, "10.1.0.2", added[0].MeshIp)
			assert.Equal(t, int32(4001), added[0].MeshPort)
			assert.Equal(t, "mesh-2", added[0].MeshHostName)
			nodeClient.AssertExpectations(t)
		})
	}
}
