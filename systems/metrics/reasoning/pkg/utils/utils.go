/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/metric"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/store"
)

const startEndKey = "start_end"

type Nodes struct {
	TNode string
	ANode string
}

func startEndStoreKey(nodeID string) string {
	return nodeID + "/" + startEndKey
}

// SortNodeIds validates a tower node ID and returns the tower + amp node pair.
func SortNodeIds(nodeID string) (Nodes, error) {
	nid, err := ukama.ValidateNodeId(nodeID)
	if err != nil {
		return Nodes{}, fmt.Errorf("validate node id: %w", err)
	}

	nodeType := nid.GetNodeType()
	if nodeType != ukama.NODE_ID_TYPE_TOWERNODE {
		return Nodes{}, fmt.Errorf("expected tower node, got %s", nid.String())
	}

	tNode := nid.String()
	aNode := strings.Replace(tNode, nodeType, ukama.NODE_ID_TYPE_AMPNODE, 1)
	return Nodes{TNode: tNode, ANode: aNode}, nil
}

// GetStartEndFromStore returns the next rolling window (start, end) for a node.
// Stores previous end as new start to avoid overlapping queries.
func GetStartEndFromStore(s *store.Store, nodeID string, interval int) (start, end string, err error) {
	key := startEndStoreKey(nodeID)
	value, err := s.Get(key)
	if err != nil {
		now := time.Now().Unix()
		start = strconv.FormatInt(now-int64(interval), 10)
		end = strconv.FormatInt(now, 10)
		_ = s.Put(key, start+":"+end)
		log.Warnf("No start/end in store for node %s, using current window", nodeID)
		return start, end, nil
	}

	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid stored value %q, expected start:end", value)
	}

	prevEnd, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", "", fmt.Errorf("invalid end timestamp %q: %w", parts[1], err)
	}

	start = parts[1]
	end = strconv.FormatInt(prevEnd+int64(interval), 10)
	_ = s.Put(key, start+":"+end)
	return start, end, nil
}

func StoreMetricResults(s *store.Store, nodeID string, metricName string, results []metric.FilteredPrometheusResult) {
	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Errorf("Failed to marshal metric results: %v", err)
		return
	}
	s.Put(nodeID + "/" + metricName, string(jsonData))
}

func GetMetricResults(s *store.Store, nodeID string, metricName string) ([]metric.FilteredPrometheusResult, error) {
	jsonData, err := s.Get(nodeID + "/" + metricName)
	if err != nil {
		log.Errorf("Failed to get metric results: %v", err)
		return nil, err
	}
	var results []metric.FilteredPrometheusResult
	dec := json.NewDecoder(bytes.NewReader([]byte(jsonData)))
	dec.UseNumber() 
	if err := dec.Decode(&results); err != nil {
		log.Errorf("Failed to unmarshal metric results: %v", err)
		return nil, err
	}
	return results, nil
}