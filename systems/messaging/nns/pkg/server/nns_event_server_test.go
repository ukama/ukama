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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/messaging/nns/pkg"
)

type fakeNodeClient struct {
	get func(nodeId string) (*creg.NodeInfo, error)
}

func (f *fakeNodeClient) Get(nodeId string) (*creg.NodeInfo, error) {
	if f.get != nil {
		return f.get(nodeId)
	}

	return nil, errFakeNotImplemented
}

func (f *fakeNodeClient) GetAll() (*creg.Nodes, error) { return nil, errFakeNotImplemented }

func (f *fakeNodeClient) GetNodesBySite(string) (*creg.NodesBySite, error) {
	return nil, errFakeNotImplemented
}

func (f *fakeNodeClient) List(creg.ListNodesRequest) (*creg.ListNodesResponse, error) {
	return nil, errFakeNotImplemented
}

func (f *fakeNodeClient) Add(creg.AddNodeRequest) (*creg.NodeInfo, error) {
	return nil, errFakeNotImplemented
}

func (f *fakeNodeClient) Attach(string, creg.AttachNodesRequest) error {
	return errFakeNotImplemented
}
func (f *fakeNodeClient) Detach(string) error                           { return errFakeNotImplemented }
func (f *fakeNodeClient) AddToSite(string, creg.AddToSiteRequest) error { return errFakeNotImplemented }
func (f *fakeNodeClient) RemoveFromSite(string) error                   { return errFakeNotImplemented }
func (f *fakeNodeClient) Delete(string) error                           { return errFakeNotImplemented }

func nodeInfoWithSite(networkId, siteId string) *creg.NodeInfo {
	return &creg.NodeInfo{
		Id:   testValidNodeID,
		Site: creg.NodeSiteInfo{NodeId: testValidNodeID, NetworkId: networkId, SiteId: siteId},
	}
}

func newEventServer(store NnsStore, nodeClient creg.NodeClient) *NnsEventServer {
	return NewNnsEventServer("test-org", nodeClient, NewNnsServer(store, testConfig(), testDns()), "org-id")
}

func onlineEvent() *epb.NodeOnlineEvent {
	return &epb.NodeOnlineEvent{
		NodeId:       testValidNodeID,
		NodeIp:       "10.0.0.1",
		NodePort:     100,
		MeshIp:       "10.0.0.2",
		MeshPort:     200,
		MeshHostName: "mesh.host",
	}
}

func TestHandleNodeOnlineEvent(t *testing.T) {
	t.Run("usesRegistryLineage", func(t *testing.T) {
		var got pkg.NodeMeshMap
		store := &fakeNnsStore{
			add: func(_ context.Context, obj pkg.NodeMeshMap) error {
				got = obj

				return nil
			},
		}
		srv := newEventServer(store, &fakeNodeClient{
			get: func(string) (*creg.NodeInfo, error) {
				return nodeInfoWithSite("net-1", "site-1"), nil
			},
		})

		require.NoError(t, srv.handleNodeOnlineEvent("k", onlineEvent()))
		assert.Equal(t, "net-1", got.Network)
		assert.Equal(t, "site-1", got.Site)
		assert.Equal(t, "10.0.0.2", got.MeshIp)
	})

	t.Run("unattachedNodeStoresEmptyLineage", func(t *testing.T) {
		var got pkg.NodeMeshMap
		store := &fakeNnsStore{
			get: func(context.Context, string) (*pkg.NodeMeshMap, error) {
				return &pkg.NodeMeshMap{Network: "stale-net", Site: "stale-site"}, nil
			},
			add: func(_ context.Context, obj pkg.NodeMeshMap) error {
				got = obj

				return nil
			},
		}
		srv := newEventServer(store, &fakeNodeClient{
			get: func(string) (*creg.NodeInfo, error) {
				return nodeInfoWithSite("", ""), nil
			},
		})

		require.NoError(t, srv.handleNodeOnlineEvent("k", onlineEvent()))
		assert.Empty(t, got.Network)
		assert.Empty(t, got.Site)
	})

	t.Run("registryFailureKeepsStoredLineage", func(t *testing.T) {
		var got pkg.NodeMeshMap
		store := &fakeNnsStore{
			get: func(context.Context, string) (*pkg.NodeMeshMap, error) {
				return &pkg.NodeMeshMap{Network: "net-1", Site: "site-1"}, nil
			},
			add: func(_ context.Context, obj pkg.NodeMeshMap) error {
				got = obj

				return nil
			},
		}
		srv := newEventServer(store, &fakeNodeClient{
			get: func(string) (*creg.NodeInfo, error) {
				return nil, errors.New("registry unreachable")
			},
		})

		require.NoError(t, srv.handleNodeOnlineEvent("k", onlineEvent()))
		assert.Equal(t, "net-1", got.Network)
		assert.Equal(t, "site-1", got.Site)
	})

	t.Run("registryNilInfoKeepsStoredLineage", func(t *testing.T) {
		var got pkg.NodeMeshMap
		store := &fakeNnsStore{
			get: func(context.Context, string) (*pkg.NodeMeshMap, error) {
				return &pkg.NodeMeshMap{Network: "net-1", Site: "site-1"}, nil
			},
			add: func(_ context.Context, obj pkg.NodeMeshMap) error {
				got = obj

				return nil
			},
		}
		srv := newEventServer(store, &fakeNodeClient{
			get: func(string) (*creg.NodeInfo, error) {
				return nil, nil
			},
		})

		require.NoError(t, srv.handleNodeOnlineEvent("k", onlineEvent()))
		assert.Equal(t, "net-1", got.Network)
		assert.Equal(t, "site-1", got.Site)
	})

	t.Run("registryFailureOnFirstContactStoresEmptyLineage", func(t *testing.T) {
		var got pkg.NodeMeshMap
		store := &fakeNnsStore{
			get: func(context.Context, string) (*pkg.NodeMeshMap, error) {
				return nil, errors.New("node not found")
			},
			add: func(_ context.Context, obj pkg.NodeMeshMap) error {
				got = obj

				return nil
			},
		}
		srv := newEventServer(store, &fakeNodeClient{
			get: func(string) (*creg.NodeInfo, error) {
				return nil, errors.New("registry unreachable")
			},
		})

		require.NoError(t, srv.handleNodeOnlineEvent("k", onlineEvent()))
		assert.Empty(t, got.Network)
		assert.Empty(t, got.Site)
		assert.Equal(t, testValidNodeID, got.NodeId)
	})
}

func TestUpdateNodeLineage(t *testing.T) {
	t.Run("assignUpdatesLineageAndKeepsAddressing", func(t *testing.T) {
		var got pkg.NodeMeshMap
		store := &fakeNnsStore{
			get: func(context.Context, string) (*pkg.NodeMeshMap, error) {
				return &pkg.NodeMeshMap{
					NodeId:       testValidNodeID,
					NodeIp:       "10.0.0.1",
					NodePort:     100,
					MeshIp:       "10.0.0.2",
					MeshPort:     200,
					MeshHostName: "mesh.host",
					Network:      "net-1",
					Site:         "site-1",
				}, nil
			},
			add: func(_ context.Context, obj pkg.NodeMeshMap) error {
				got = obj

				return nil
			},
		}
		srv := newEventServer(store, &fakeNodeClient{})

		require.NoError(t, srv.handleNodeAssignedEvent("k", &epb.EventRegistryNodeAssign{
			NodeId:  testValidNodeID,
			Network: "net-2",
			Site:    "site-2",
		}))
		assert.Equal(t, "net-2", got.Network)
		assert.Equal(t, "site-2", got.Site)
		assert.Equal(t, "10.0.0.2", got.MeshIp)
		assert.Equal(t, int32(200), got.MeshPort)
		assert.Equal(t, "mesh.host", got.MeshHostName)
	})

	t.Run("releaseClearsLineage", func(t *testing.T) {
		var got pkg.NodeMeshMap
		store := &fakeNnsStore{
			get: func(context.Context, string) (*pkg.NodeMeshMap, error) {
				return &pkg.NodeMeshMap{
					NodeId:  testValidNodeID,
					MeshIp:  "10.0.0.2",
					Network: "net-1",
					Site:    "site-1",
				}, nil
			},
			add: func(_ context.Context, obj pkg.NodeMeshMap) error {
				got = obj

				return nil
			},
		}
		srv := newEventServer(store, &fakeNodeClient{})

		require.NoError(t, srv.handleNodeReleaseEvent("k", &epb.NodeReleasedEvent{
			NodeId:  testValidNodeID,
			Network: "net-1",
			Site:    "site-1",
		}))
		assert.Empty(t, got.Network)
		assert.Empty(t, got.Site)
		assert.Equal(t, "10.0.0.2", got.MeshIp)
	})

	t.Run("unknownNodeIsNotAnError", func(t *testing.T) {
		store := &fakeNnsStore{
			get: func(context.Context, string) (*pkg.NodeMeshMap, error) {
				return nil, errors.New("node not found")
			},
			add: func(context.Context, pkg.NodeMeshMap) error {
				t.Fatal("Add must not be called for a node with no mapping")

				return nil
			},
		}
		srv := newEventServer(store, &fakeNodeClient{})

		require.NoError(t, srv.handleNodeAssignedEvent("k", &epb.EventRegistryNodeAssign{
			NodeId:  testValidNodeID,
			Network: "net-2",
			Site:    "site-2",
		}))
	})

	t.Run("unchangedLineageSkipsWrite", func(t *testing.T) {
		store := &fakeNnsStore{
			get: func(context.Context, string) (*pkg.NodeMeshMap, error) {
				return &pkg.NodeMeshMap{Network: "net-1", Site: "site-1"}, nil
			},
			add: func(context.Context, pkg.NodeMeshMap) error {
				t.Fatal("Add must not be called when lineage is unchanged")

				return nil
			},
		}
		srv := newEventServer(store, &fakeNodeClient{})

		require.NoError(t, srv.handleNodeAssignedEvent("k", &epb.EventRegistryNodeAssign{
			NodeId:  testValidNodeID,
			Network: "net-1",
			Site:    "site-1",
		}))
	})
}
