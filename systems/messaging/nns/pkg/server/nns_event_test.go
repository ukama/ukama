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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/messaging/nns/pkg"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const (
	testEventOrg  = "test-org"
	assignRoute   = "event.cloud.local.{{ .Org}}.registry.node.node.assign"
	releaseRoute  = "event.cloud.local.{{ .Org}}.registry.node.node.release"
	assignedRoute = "event.cloud.local.{{ .Org}}.registry.node.node.assigned"
)

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
			return testEvent(t, assignRoute, &epb.NodeReleasedEvent{NodeId: testValidNodeID})
		}, true},
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

	t.Run("assignForUnknownNodeFails", func(t *testing.T) {
		added := 0
		store := &fakeNnsStore{
			get: func(_ context.Context, _ string) (*pkg.NodeMeshMap, error) {
				return nil, errFakeNotImplemented
			},
			add: func(_ context.Context, _ pkg.NodeMeshMap) error {
				added++
				return nil
			},
		}
		srv := NewNnsEventServer(testEventOrg, nil, NewNnsServer(store, nil, nil),
			testEventOrg)

		_, err := srv.EventNotification(context.Background(),
			testEvent(t, assignRoute, &epb.EventRegistryNodeAssign{NodeId: testValidNodeID}))

		assert.Error(t, err)
		assert.Zero(t, added)
	})
}
