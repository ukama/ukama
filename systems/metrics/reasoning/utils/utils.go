/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/store"
)

// MetricResultLogData holds the data needed to build metric result log output.
type MetricResultLogData struct {
	Value           float64
	State           string
	Trend           string
	Confidence      float64
	ProjectionType  string
	ProjectionEtaSec float64
}

// MetricResultLogFields returns log.Fields for structured metric result logging.
func MetricResultLogFields(d MetricResultLogData) log.Fields {
	fields := log.Fields{
		"value":      d.Value,
		"state":      d.State,
		"trend":      d.Trend,
		"confidence": d.Confidence,
	}
	if d.ProjectionType != "" {
		fields["projection"] = d.ProjectionType + " in " + FormatSec(d.ProjectionEtaSec)
	}
	return fields
}

const startEndKey = "start_end"

// RoundToDecimalPoints rounds value to the specified number of decimal places.
// Preserves NaN and Inf unchanged. Uses half-up rounding.
func RoundToDecimalPoints(value float64, decimalPoints int) float64 {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return value
	}
	if decimalPoints < 0 {
		return value
	}
	shift := math.Pow(10, float64(decimalPoints))
	return math.Round(value*shift) / shift
}

type Nodes struct {
	TNode string
	ANode string
}

func startEndStoreKey(nodeID string) string {
	return nodeID + "/" + startEndKey
}

func GetAlgoStatsStoreKey(nodeID string, metricKey string) string {
	return nodeID + "/" + metricKey + "/" + "algo_stats"
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

// FormatSec formats seconds for display (e.g. "45s" or "3m").
func FormatSec(sec float64) string {
	if sec < 60 {
		return strconv.FormatInt(int64(sec), 10) + "s"
	}
	return strconv.FormatInt(int64(sec/60), 10) + "m"
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