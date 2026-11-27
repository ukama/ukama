/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/ukama/ukama/systems/messaging/nns/pb/gen"
	"github.com/ukama/ukama/systems/messaging/nns/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const testValidNodeID = "UK-SA2156-CNODE-A1-XXXX"
const testInvalidNodeID = "UK-SA2156-MNODE-A1-XXXX"

var errFakeNotImplemented = errors.New("not implemented")

type fakeNnsStore struct {
	get        func(ctx context.Context, nodeId string) (*pkg.NodeMeshMap, error)
	getAll     func(ctx context.Context) ([]pkg.NodeMeshMap, error)
	add        func(ctx context.Context, obj pkg.NodeMeshMap) error
	delete     func(ctx context.Context, nodeId string) error
	updateMesh func(ctx context.Context, nodeId string, ip string, port int32) error
	updateNode func(ctx context.Context, nodeId string, nodeIp string, nodePort int32) error
}

func (f *fakeNnsStore) Get(ctx context.Context, nodeId string) (*pkg.NodeMeshMap, error) {
	if f.get != nil {
		return f.get(ctx, nodeId)
	}
	return nil, errFakeNotImplemented
}

func (f *fakeNnsStore) GetAll(ctx context.Context) ([]pkg.NodeMeshMap, error) {
	if f.getAll != nil {
		return f.getAll(ctx)
	}
	return nil, errFakeNotImplemented
}

func (f *fakeNnsStore) Add(ctx context.Context, obj pkg.NodeMeshMap) error {
	if f.add != nil {
		return f.add(ctx, obj)
	}
	return errFakeNotImplemented
}

func (f *fakeNnsStore) Delete(ctx context.Context, nodeId string) error {
	if f.delete != nil {
		return f.delete(ctx, nodeId)
	}
	return errFakeNotImplemented
}

func (f *fakeNnsStore) UpdateNodeMesh(ctx context.Context, nodeId string, ip string, port int32) error {
	if f.updateMesh != nil {
		return f.updateMesh(ctx, nodeId, ip, port)
	}
	return errFakeNotImplemented
}

func (f *fakeNnsStore) UpdateNode(ctx context.Context, nodeId string, nodeIp string, nodePort int32) error {
	if f.updateNode != nil {
		return f.updateNode(ctx, nodeId, nodeIp, nodePort)
	}
	return errFakeNotImplemented
}

func testConfig() *pkg.Config {
	return &pkg.Config{OrgName: "test-org"}
}

func testDns() *pkg.DnsConfig {
	return &pkg.DnsConfig{NodeDomain: "mesh.test"}
}

func TestNewNnsServer(t *testing.T) {
	fake := &fakeNnsStore{}
	cfg := testConfig()
	dns := testDns()

	t.Run("returnsNonNil", func(t *testing.T) {
		srv := NewNnsServer(fake, cfg, dns)
		require.NotNil(t, srv)
	})

	t.Run("preservesDependencies", func(t *testing.T) {
		srv := NewNnsServer(fake, cfg, dns)
		assert.Same(t, cfg, srv.config)
		assert.Same(t, dns, srv.dnsConfig)
		assert.Same(t, fake, srv.nns)
	})
}

func TestNnsServerGetNode(t *testing.T) {
	ctx := context.Background()

	t.Run("invalidNodeId", func(t *testing.T) {
		srv := NewNnsServer(&fakeNnsStore{}, testConfig(), testDns())
		_, err := srv.GetNode(ctx, &pb.GetNodeRequest{NodeId: "short"})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("success", func(t *testing.T) {
		fake := &fakeNnsStore{
			get: func(_ context.Context, nodeId string) (*pkg.NodeMeshMap, error) {
				return &pkg.NodeMeshMap{
					NodeId:   nodeId,
					NodeIp:   "1.1.1.1",
					NodePort: 1000,
				}, nil
			},
		}
		srv := NewNnsServer(fake, testConfig(), testDns())
		resp, err := srv.GetNode(ctx, &pb.GetNodeRequest{NodeId: testValidNodeID})
		require.NoError(t, err)
		assert.Equal(t, testValidNodeID, resp.NodeId)
		assert.Equal(t, "1.1.1.1", resp.NodeIp)
		assert.Equal(t, int32(1000), resp.NodePort)
	})
}

func TestNnsServerGetMesh(t *testing.T) {
	ctx := context.Background()

	t.Run("invalidNodeId", func(t *testing.T) {
		srv := NewNnsServer(&fakeNnsStore{}, testConfig(), testDns())
		_, err := srv.GetMesh(ctx, &pb.GetMeshRequest{NodeId: testInvalidNodeID})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("success", func(t *testing.T) {
		fake := &fakeNnsStore{
			get: func(_ context.Context, _ string) (*pkg.NodeMeshMap, error) {
				return &pkg.NodeMeshMap{MeshIp: "10.0.0.5", MeshPort: 5000}, nil
			},
		}
		srv := NewNnsServer(fake, testConfig(), testDns())
		resp, err := srv.GetMesh(ctx, &pb.GetMeshRequest{NodeId: testValidNodeID})
		require.NoError(t, err)
		assert.Equal(t, "10.0.0.5", resp.MeshIp)
		assert.Equal(t, int32(5000), resp.MeshPort)
	})
}

func TestNnsServerSet(t *testing.T) {
	ctx := context.Background()

	t.Run("invalidNodeId", func(t *testing.T) {
		srv := NewNnsServer(&fakeNnsStore{}, testConfig(), testDns())
		_, err := srv.Set(ctx, &pb.SetRequest{NodeId: "x"})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("success", func(t *testing.T) {
		var got pkg.NodeMeshMap
		fake := &fakeNnsStore{
			add: func(_ context.Context, obj pkg.NodeMeshMap) error {
				got = obj
				return nil
			},
		}
		cfg := testConfig()
		srv := NewNnsServer(fake, cfg, testDns())
		_, err := srv.Set(ctx, &pb.SetRequest{
			NodeId:       testValidNodeID,
			NodeIp:       "1.1.1.1",
			NodePort:     11,
			MeshPort:     22,
			Network:      "net-a",
			Site:         "site-b",
			MeshHostName: "mesh.host",
			MeshIp:       "2.2.2.2",
		})
		require.NoError(t, err)
		assert.Equal(t, testValidNodeID, got.NodeId)
		assert.Equal(t, "1.1.1.1", got.NodeIp)
		assert.Equal(t, int32(11), got.NodePort)
		assert.Equal(t, int32(22), got.MeshPort)
		assert.Equal(t, cfg.OrgName, got.Org)
		assert.Equal(t, "net-a", got.Network)
		assert.Equal(t, "site-b", got.Site)
		assert.Equal(t, "mesh.host", got.MeshHostName)
		assert.Equal(t, "2.2.2.2", got.MeshIp)
	})
}

func TestNnsServerUpdateMesh(t *testing.T) {
	ctx := context.Background()

	t.Run("backendError", func(t *testing.T) {
		want := errors.New("etcd down")
		fake := &fakeNnsStore{
			updateMesh: func(context.Context, string, string, int32) error {
				return want
			},
		}
		srv := NewNnsServer(fake, testConfig(), testDns())
		_, err := srv.UpdateMesh(ctx, &pb.UpdateMeshRequest{
			NodeId:    testValidNodeID,
			MeshIp:    "9.9.9.9",
			MeshPort:  3333,
		})
		require.ErrorIs(t, err, want)
	})

	t.Run("success", func(t *testing.T) {
		var gotID string
		var gotIP string
		var gotPort int32
		fake := &fakeNnsStore{
			updateMesh: func(_ context.Context, nodeId string, ip string, port int32) error {
				gotID, gotIP, gotPort = nodeId, ip, port
				return nil
			},
		}
		srv := NewNnsServer(fake, testConfig(), testDns())
		_, err := srv.UpdateMesh(ctx, &pb.UpdateMeshRequest{
			NodeId:   testValidNodeID,
			MeshIp:   "8.8.8.8",
			MeshPort: 4444,
		})
		require.NoError(t, err)
		assert.Equal(t, testValidNodeID, gotID)
		assert.Equal(t, "8.8.8.8", gotIP)
		assert.Equal(t, int32(4444), gotPort)
	})
}

func TestNnsServerUpdateNode(t *testing.T) {
	ctx := context.Background()

	t.Run("invalidNodeId", func(t *testing.T) {
		srv := NewNnsServer(&fakeNnsStore{}, testConfig(), testDns())
		_, err := srv.UpdateNode(ctx, &pb.UpdateNodeRequest{NodeId: "nope"})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("success", func(t *testing.T) {
		fake := &fakeNnsStore{
			updateNode: func(_ context.Context, nodeId string, nodeIp string, nodePort int32) error {
				assert.Equal(t, testValidNodeID, nodeId)
				assert.Equal(t, "3.3.3.3", nodeIp)
				assert.Equal(t, int32(777), nodePort)
				return nil
			},
		}
		srv := NewNnsServer(fake, testConfig(), testDns())
		_, err := srv.UpdateNode(ctx, &pb.UpdateNodeRequest{
			NodeId:   testValidNodeID,
			NodeIp:   "3.3.3.3",
			NodePort: 777,
		})
		require.NoError(t, err)
	})
}

func TestNnsServerDelete(t *testing.T) {
	ctx := context.Background()

	t.Run("invalidNodeId", func(t *testing.T) {
		srv := NewNnsServer(&fakeNnsStore{}, testConfig(), testDns())
		_, err := srv.Delete(ctx, &pb.DeleteRequest{NodeId: ""})
		require.Error(t, err)
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("success", func(t *testing.T) {
		var deleted string
		fake := &fakeNnsStore{
			delete: func(_ context.Context, nodeId string) error {
				deleted = nodeId
				return nil
			},
		}
		srv := NewNnsServer(fake, testConfig(), testDns())
		_, err := srv.Delete(ctx, &pb.DeleteRequest{NodeId: testValidNodeID})
		require.NoError(t, err)
		assert.Equal(t, testValidNodeID, deleted)
	})
}

func TestNnsServerList(t *testing.T) {
	ctx := context.Background()

	t.Run("backendError", func(t *testing.T) {
		want := errors.New("scan failed")
		fake := &fakeNnsStore{
			getAll: func(context.Context) ([]pkg.NodeMeshMap, error) {
				return nil, want
			},
		}
		srv := NewNnsServer(fake, testConfig(), testDns())
		_, err := srv.List(ctx, &pb.ListRequest{})
		require.ErrorIs(t, err, want)
	})

	t.Run("success", func(t *testing.T) {
		fake := &fakeNnsStore{
			getAll: func(context.Context) ([]pkg.NodeMeshMap, error) {
				return []pkg.NodeMeshMap{
					{
						NodeId: "UK-SA2156-TNODE-A1-XXXX",
						NodeIp: "10.0.0.1",
						Org:    "o",
					},
				}, nil
			},
		}
		srv := NewNnsServer(fake, testConfig(), testDns())
		resp, err := srv.List(ctx, &pb.ListRequest{})
		require.NoError(t, err)
		require.Len(t, resp.List, 1)
		assert.Equal(t, "UK-SA2156-TNODE-A1-XXXX", resp.List[0].NodeId)
		assert.Equal(t, "10.0.0.1", resp.List[0].NodeIp)
		assert.Equal(t, "o", resp.List[0].Org)
	})
}
