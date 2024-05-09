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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/api"
	"github.com/ukama/ukama/systems/common/rest"
)

const NodeEndpoint = "/noded/v1/nodeinfo"

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

type NodeInfo struct {
	Id            string  `json:"uuid"`
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

type NodedProvider interface {
	GetNodeInfo() (*api.Spr, error)
}

type nodedClient struct {
	u     *url.URL
	R     *rest.RestClient
	debug bool
}

func NewNodedClient(h string, debug bool) (*nodedClient, error) {
	u, err := url.Parse(h)

	if err != nil {
		log.Errorf("Can't parse  %s url. Error: %s", h, err.Error())
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

	return ni.Id, nil
}
func (r *nodedClient) GetNodeInfo() (*NodeInfo, error) {
	log.Debugf("Geeting NodeInfo from Noded.")

	node := &Node{}

	url := r.u.String() + NodeEndpoint

	resp, err := r.R.C.R().
		Get(url)
	if err != nil {
		log.Errorf("Get NodeInfo failure. error: %s", err.Error())
		return nil, fmt.Errorf("get NodeInfo failure: %s", err.Error())
	}

	err = json.Unmarshal(resp.Body(), node)
	if err != nil {
		log.Errorf("Failed to deserialize node info. Error message is: %s", err.Error())
		return nil, fmt.Errorf("node info deserailization failure: %w", err)
	}

	log.Infof("NodeInfo Info: %+v", node)

	return &node.NodeInfo, nil
}
