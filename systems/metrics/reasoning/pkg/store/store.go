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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Store struct {
	etcd *clientv3.Client
}

func NewStore(config *pkg.Config) *Store {
	etcd, err := clientv3.New(clientv3.Config{
		DialTimeout: config.DialTimeoutSecond,
		Endpoints:   []string{config.EtcdHost},
	})
	if err != nil {
		log.Fatalf("Cannot connect to etcd: %v", err)
	}
	return &Store{
		etcd: etcd,
	}
}	

func (s *Store) Put(key string, value string) error {
	_, err := s.etcd.Put(context.Background(), key, value)
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

func (s *Store) GetJson(key string) (interface{}, error) {
	jsonData, err := s.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get json data: %v", err)
	}
	var value interface{}
	if err := json.Unmarshal([]byte(jsonData), &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json data: %v", err)
	}
	return value, nil
}

func (s *Store) Get(key string) (string, error) {
	response, err := s.etcd.Get(context.Background(), key)
	if err != nil {
		return "", fmt.Errorf("failed to get record %s. Error: %v", key, err)
	}
	if len(response.Kvs) == 0 {
		return "", fmt.Errorf("record %s not found", key)
	}
	return string(response.Kvs[0].Value), nil
}

func (s *Store) GetAll() ([]string, error) {
	response, err := s.etcd.Get(context.Background(), "", clientv3.WithPrefix())
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
	_, err := s.etcd.Delete(context.Background(), key)
	if err != nil {
		return fmt.Errorf("failed to delete record %s. Error: %v", key, err)
	}
	return nil
}