/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
)

const NodeEndpoint = "/v1/nodeinfo"

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

type NodeInfo struct {
	Id            string  `json:"UUID"`
	Name          string  `json:"name"`
	NodeType      int     `json:"type"`
	PartNumber    string  `json:"partNumber"`
	Skew          string  `json:"skew"`
	Mac           string  `json:"mac"`
	ProdSwVersion Version `json:"prodSwVersion"`
	SwVersion     Version `json:"swVersion"`
	AssemblyDate  string  `json:"assemblyDate"`
	OemName       string  `json:"oemName"`
	ModuleCount   int     `json:"moduleCount"`
}

type Node struct {
	NodeInfo NodeInfo `json:"nodeInfo"`
}

type nodedClient struct {
	u     *url.URL
	R     *rest.RestClient
	debug bool
}

func NewNodedClient(h string, debug bool) (*nodedClient, error) {
	u, err := url.Parse(strings.TrimRight(h, "/"))
	if err != nil {
		log.Errorf("Can't parse %s url. Error: %s", h, err.Error())
		return nil, err
	}

	return &nodedClient{
		u:     u,
		R:     rest.NewRestyClient(u, debug),
		debug: debug,
	}, nil
}

func (r *nodedClient) GetNodeId() (string, error) {
	ni, err := r.GetNodeInfo()
	if err != nil {
		return "", err
	}

	if ni.Id == "" {
		return "", fmt.Errorf("node info missing UUID")
	}

	return ni.Id, nil
}

func (r *nodedClient) GetNodeInfo() (*NodeInfo, error) {
	var err error
	var respBody string
	node := &Node{}

	log.Debugf("Getting NodeInfo from Noded.")

	resp, err := r.R.C.R().Get(r.u.String() + NodeEndpoint)
	if err != nil {
		log.Errorf("Get NodeInfo failure. error: %s", err.Error())
		return nil, fmt.Errorf("get NodeInfo failure: %w", err)
	}

	if resp.StatusCode() != 200 {
		respBody = strings.TrimSpace(string(resp.Body()))
		return nil, fmt.Errorf("node info returned http %d: %s",
			resp.StatusCode(), respBody)
	}

	err = json.Unmarshal(resp.Body(), node)
	if err != nil {
		log.Errorf("Failed to deserialize node info. Error message is: %s",
			err.Error())
		return nil, fmt.Errorf("node info deserialization failure: %w", err)
	}

	log.Infof("NodeInfo Info: %+v", node)

	return &node.NodeInfo, nil
}

type NodedClientAlias nodedClient

func (r *NodedClientAlias) GetNodeId() (string, error) {
	return (*nodedClient)(r).GetNodeId()
}
