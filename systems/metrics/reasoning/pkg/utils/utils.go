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
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	ukama "github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/metrics/reasoning/pkg/store"
)

type Nodes struct {
	TNode string
	ANode string
}

func SortNodeIds(nodeId string) (Nodes, error) {
	nodes := Nodes{}
	nid, err := ukama.ValidateNodeId(nodeId)
	if err != nil {
		log.Errorf("Failed to validate node id: %v", err)
		return nodes, fmt.Errorf("failed to validate node id: %v", err)
	}

	ntype := ukama.GetNodeType(nid.String()) 
	if ntype == nil {
		log.Errorf("Failed to get node type: %v", err)
		return nodes, fmt.Errorf("failed to get node type: %v", err)
	}

	if *ntype != ukama.NODE_ID_TYPE_TOWERNODE {
		log.Errorf("Node type is not a tower node: %v", nid.String())
		return nodes, fmt.Errorf("node type is not a tower node: %v", nid.String())
	}

	nodes.TNode = nid.String()
	nodes.ANode = strings.Replace(nid.String(), *ntype, ukama.NODE_ID_TYPE_AMPNODE, 1)
	log.Infof("Sorted nodes: %v", nodes)

	return nodes, nil
}

func GetToNFromStore(store *store.Store, nodeId string, interval int) (toN, fromN string, err error) {
	value, err := store.Get(nodeId + "/to_n_from")
	if err != nil {
		return "", "", err
	}

	if value == "" {
		now := time.Now()
		log.Errorf("no To and From value found for node: %s", nodeId)
		return strconv.FormatInt(now.Unix(), 10), strconv.FormatInt(now.Unix() - int64(interval), 10), nil
	}

	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid To and From value: %s", value)
	}
	toN = parts[0]
	fromN = parts[1]

	log.Infof("To and From value: %s, %s", toN, fromN)

	return toN, fromN, nil
}