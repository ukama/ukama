/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package messaging

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest/client"
)

const NnsEndpoint = "/v1/nns"

type MeshInfo struct {
	MeshIp string `json:"meshIp"`
	MeshPort int `json:"meshPort"`
}

type NodeInfo struct {
	NodeId string `json:"nodeId"`
	NodeIp string `json:"nodeIp"`
	NodePort int `json:"nodePort"`
}

type NodeMeshInfo struct {
	NodeId string `json:"nodeId"`
	NodeIp string `json:"nodeIp"`
	NodePort int `json:"nodePort"`
	MeshIp string `json:"meshIp"`
	MeshPort int `json:"meshPort"`
	Org string `json:"org"`
	Network string `json:"network"`
	Site string `json:"site"`
	MeshHostName string `json:"meshHostName"`
}

type ListNodeMeshInfo struct {
	List []NodeMeshInfo `json:"list"`
}

type NnsClient interface {
	GetMesh(nodeId string) (*MeshInfo, error)
	GetNode(nodeId string) (*NodeInfo, error)
	List() (*ListNodeMeshInfo, error)
}

type nnsClient struct {
	u *url.URL
	R *client.Resty
}

func NewNnsClient(h string, options ...client.Option) *nnsClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &nnsClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (c *nnsClient) GetMesh(nodeId string) (*MeshInfo, error) {
	log.Debugf("Getting mesh for node: %v", nodeId)

	mesh := MeshInfo{}

	resp, err := c.R.Get(c.u.String() + NnsEndpoint + "/mesh/" + nodeId)
	if err != nil {
		log.Errorf("GetMesh failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetMesh failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &mesh)
	if err != nil {
		log.Tracef("Failed to deserialize mesh info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("mesh info deserialization failure: %w", err)
	}

	log.Infof("Mesh Info: %+v", mesh)

	return &mesh, nil
}

func (c *nnsClient) GetNode(nodeId string) (*NodeInfo, error) {
	log.Debugf("Getting node: %v", nodeId)

	node := NodeInfo{}

	resp, err := c.R.Get(c.u.String() + NnsEndpoint + "/node/" + nodeId)
	if err != nil {
		log.Errorf("GetNode failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetNode failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &node)
	if err != nil {
		log.Tracef("Failed to deserialize node info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("node info deserialization failure: %w", err)
	}

	log.Infof("Node Info: %+v", node)

	return &node, nil
}

func (c *nnsClient) List() (*ListNodeMeshInfo, error) {
	log.Debugf("Listing node mesh info")

	list := ListNodeMeshInfo{}

	resp, err := c.R.Get(c.u.String() + NnsEndpoint + "/list")
	if err != nil {
		log.Errorf("List failure. error: %s", err.Error())

		return nil, fmt.Errorf("List failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &list)
	if err != nil {
		log.Tracef("Failed to deserialize list node mesh info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("list node mesh info deserialization failure: %w", err)
	}

	log.Infof("List Node Mesh Info: %+v", list)

	return &list, nil
}