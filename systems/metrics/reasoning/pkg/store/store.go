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
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const defaultRequestTimeout = 5 * time.Second

type Store struct {
	etcd            *clientv3.Client
	requestTimeout  time.Duration
}

func NewStore(config *pkg.Config) *Store {
	etcd, err := clientv3.New(clientv3.Config{
		DialTimeout: config.DialTimeoutSecond,
		Endpoints:   []string{config.EtcdHost},
	})
	if err != nil {
		log.Fatalf("Cannot connect to etcd: %v", err)
	}
	timeout := config.Timeout
	if timeout == 0 {
		timeout = defaultRequestTimeout
	}
	return &Store{
		etcd:           etcd,
		requestTimeout: timeout,
	}
}

func (s *Store) opCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.requestTimeout)
}	

func (s *Store) Put(key string, value string) error {
	ctx, cancel := s.opCtx()
	defer cancel()
	_, err := s.etcd.Put(ctx, key, value)
	if err != nil {
		return fmt.Errorf("failed to add record %s with value %s. Error: %v", key, value, err)
	}
	log.Infof("Added node %s with value %s to etcd", key, value)
	return nil
}

func (s *Store) PutJson(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value %v. Error: %v", value, err)
	}
	if err := s.Put(key, string(jsonData)); err != nil {
		return fmt.Errorf("failed to store json data: %v", err)
	}
	return nil
}

func (s *Store) GetJson(key string) ([]byte, error) {
	jsonData, err := s.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get json data: %v", err)
	}
	return []byte(jsonData), nil
}

func (s *Store) Get(key string) (string, error) {
	ctx, cancel := s.opCtx()
	defer cancel()
	response, err := s.etcd.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get record %s. Error: %v", key, err)
	}
	if len(response.Kvs) == 0 {
		return "", fmt.Errorf("record %s not found", key)
	}
	return string(response.Kvs[0].Value), nil
}

func (s *Store) GetAll() ([]string, error) {
	ctx, cancel := s.opCtx()
	defer cancel()
	response, err := s.etcd.Get(ctx, "", clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get all records. Error: %v", err)
	}
	keys := make([]string, len(response.Kvs))
	for i, kv := range response.Kvs {
		keys[i] = string(kv.Key)
	}
	return keys, nil
}

func (s *Store) Delete(key string) error {
	ctx, cancel := s.opCtx()
	defer cancel()
	_, err := s.etcd.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete record %s. Error: %v", key, err)
	}
	return nil
}