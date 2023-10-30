/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

const NodeEndpoint = "/v1/nodes"

type SiteInfo struct {
	NodeId    string    `json:"node_id,omitempty"`
	SiteId    string    `json:"site_id,omitempty"`
	NetworkId string    `json:"network_id,omitempty"`
	AddedAt   time.Time `json:"added_at,omitempty"`
}

type NodeInfo struct {
	Id        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	OrgId     string     `json:"org_id,omitempty"`
	Type      string     `json:"type,omitempty"`
	State     string     `json:"state,omitempty"`
	Site      SiteInfo   `json:"site,omitempty"`
	Attahced  []NodeInfo `json:"attached_nodes,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
}

type Node struct {
	NodeInfo *NodeInfo `json:"node"`
}

type AddNodeRequest struct {
	NodeId string `json:"node_id,omitempty"`
	Name   string `json:"name,omitempty"`
	OrgId  string `json:"org_id,omitempty"`

	// TODO: open issue to remove state param required on registry api-gateway
	State string `json:"state,omitempty"`
}

type AttachNodesRequest struct {
	AmpNodeL string `json:"anodel"`
	AmpNodeR string `json:"anoder"`
}

type AddToSiteRequest struct {
	// NodeId string `json:"node_id" path:"node_id" validate:"required"`

	SiteId    string `json:"site_id"`
	NetworkId string `json:"net_id"`
}

type NodeClient interface {
	Get(string) (*NodeInfo, error)
	Add(AddNodeRequest) (*NodeInfo, error)
	Attach(string, AttachNodesRequest) error
	Detach(string) error
	AddToSite(string, AddToSiteRequest) error
	RemoveFromSite(string) error
	Delete(string) error
}

type nodeClient struct {
	u *url.URL
	R *Resty
}

func NewNodeClient(h string) *nodeClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &nodeClient{
		u: u,
		R: NewResty(),
	}
}

// TODO check upstream add returns payload
func (n *nodeClient) Add(req AddNodeRequest) (*NodeInfo, error) {
	log.Debugf("Adding node: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	node := Node{}

	resp, err := n.R.Post(n.u.String()+NodeEndpoint, b)
	if err != nil {
		log.Errorf("AddNode failure. error: %s", err.Error())

		return nil, fmt.Errorf("AddNode failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &node)
	if err != nil {
		log.Tracef("Failed to deserialize node info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("node info deserailization failure: %w", err)
	}

	log.Infof("Node Info: %+v", node.NodeInfo)

	return node.NodeInfo, nil
}

func (n *nodeClient) Get(nodeId string) (*NodeInfo, error) {
	log.Debugf("Getting node: %v", nodeId)

	node := Node{}

	resp, err := n.R.Get(n.u.String() + NodeEndpoint + "/" + nodeId)
	if err != nil {
		log.Errorf("GetNode failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetNode failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &node)
	if err != nil {
		log.Tracef("Failed to deserialize node info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("node info deserailization failure: %w", err)
	}

	log.Infof("Node Info: %+v", node.NodeInfo)

	return node.NodeInfo, nil
}

func (n *nodeClient) Attach(parentNodeId string, req AttachNodesRequest) error {
	log.Debugf("Attaching nodes node: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	_, err = n.R.Post(n.u.String()+NodeEndpoint+"/"+parentNodeId+"/attach", b)
	if err != nil {
		log.Errorf("AttachNode failure. error: %s", err.Error())

		return fmt.Errorf("AttachNode failure: %w", err)
	}

	return nil
}

func (n *nodeClient) Detach(nodeId string) error {
	log.Debugf("Detaching node: %v", nodeId)

	_, err := n.R.Delete(n.u.String() + NodeEndpoint + "/" + nodeId + "/attach")
	if err != nil {
		log.Errorf("DetachNode failure. error: %s", err.Error())

		return fmt.Errorf("DetachNode failure: %w", err)
	}

	return nil
}

func (n *nodeClient) AddToSite(nodeId string, req AddToSiteRequest) error {
	log.Debugf("Attaching nodes node: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	_, err = n.R.Post(n.u.String()+NodeEndpoint+"/"+nodeId+"/sites", b)
	if err != nil {
		log.Errorf("AttachNode failure. error: %s", err.Error())

		return fmt.Errorf("AttachNode failure: %w", err)
	}

	return nil
}

func (n *nodeClient) RemoveFromSite(nodeId string) error {
	log.Debugf("Detaching node: %v", nodeId)

	_, err := n.R.Delete(n.u.String() + NodeEndpoint + "/" + nodeId + "/sites")
	if err != nil {
		log.Errorf("DetachNode failure. error: %s", err.Error())

		return fmt.Errorf("DetachNode failure: %w", err)
	}

	return nil
}

func (n *nodeClient) Delete(nodeId string) error {
	log.Debugf("Deleting node: %v", nodeId)

	_, err := n.R.Delete(n.u.String() + NodeEndpoint + "/" + nodeId)
	if err != nil {
		log.Errorf("DeleteNode failure. error: %s", err.Error())

		return fmt.Errorf("DeleteNode failure: %w", err)
	}

	return nil
}
