/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package store

import (
	"context"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

// MemKV is an in-memory implementation of kvOps for testing.
type MemKV struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewMemKV creates a new in-memory key-value store.
func NewMemKV() *MemKV {
	return &MemKV{data: make(map[string]string)}
}

func (m *MemKV) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val
	return &clientv3.PutResponse{}, nil
}

func (m *MemKV) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if key == "" {
		kvs := make([]*mvccpb.KeyValue, 0, len(m.data))
		for k, v := range m.data {
			kvs = append(kvs, &mvccpb.KeyValue{Key: []byte(k), Value: []byte(v)})
		}
		return &clientv3.GetResponse{Kvs: kvs}, nil
	}
	v, ok := m.data[key]
	if !ok {
		return &clientv3.GetResponse{Kvs: nil}, nil
	}
	return &clientv3.GetResponse{Kvs: []*mvccpb.KeyValue{{Key: []byte(key), Value: []byte(v)}}}, nil
}

func (m *MemKV) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return &clientv3.DeleteResponse{}, nil
}

// NewInMemoryStore creates a Store backed by in-memory storage. Used for testing.
func NewInMemoryStore() *Store {
	return NewStoreWithKV(NewMemKV(), 5*time.Second)
}
