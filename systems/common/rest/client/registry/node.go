/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package registry

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const NodeEndpoint = "/v1/nodes"

type NodeSiteInfo struct {
	NodeId    string    `json:"node_id,omitempty"`
	SiteId    string    `json:"site_id,omitempty"`
	NetworkId string    `json:"network_id,omitempty"`
	AddedAt   time.Time `json:"added_at,omitempty"`
}

type NodeStatusInfo struct {
	Connectivity string `json:"connectivity,omitempty"`
	State        string `json:"state,omitempty"`
}

type NodeInfo struct {
	Id        string         `json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	Type      string         `json:"type,omitempty"`
	Status    NodeStatusInfo `json:"status,omitempty"`
	Site      NodeSiteInfo   `json:"site,omitempty"`
	Latitude  string        `json:"latitude,omitempty"`
	Longitude string        `json:"longitude,omitempty"`
	Attahced  []NodeInfo     `json:"attached_nodes,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
}

type NodesBySite struct {
	Nodes  []NodeInfo `json:"nodes"`
	SiteId string     `json:"site_id"`
}

type Nodes struct {
	Nodes []*NodeInfo `json:"nodes"`
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

type ListNodesRequest struct {
	NetworkId    string `json:"network_id"`
	SiteId       string `json:"site_id"`
	Connectivity string `json:"connectivity,omitempty"`
	State        string `json:"state,omitempty"`
	NodeId       string `json:"node_id"`
	Type         string `json:"type"`
}

type ListNodesResponse struct {
	Nodes []*NodeInfo `json:"nodes"`
}

type NodeClient interface {
	Get(string) (*NodeInfo, error)
	GetAll() (*Nodes, error)
	GetNodesBySite(string) (*NodesBySite, error)
	List(ListNodesRequest) (*ListNodesResponse, error)
	Add(AddNodeRequest) (*NodeInfo, error)
	Attach(string, AttachNodesRequest) error
	Detach(string) error
	AddToSite(string, AddToSiteRequest) error
	RemoveFromSite(string) error
	Delete(string) error
}

type nodeClient struct {
	u *url.URL
	R *client.Resty
}

func NewNodeClient(h string, options ...client.Option) *nodeClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &nodeClient{
		u: u,
		R: client.NewResty(options...),
	}
}

// TODO check upstream add returned payload
func (n *nodeClient) Add(req AddNodeRequest) (*NodeInfo, error) {
	log.Debugf("Adding node: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %w", err)
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

		return nil, fmt.Errorf("node info deserialization failure: %w", err)
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

		return nil, fmt.Errorf("node info deserialization failure: %w", err)
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

func (n *nodeClient) GetAll() (*Nodes, error) {
	log.Debugf("Getting all nodes.")

	nodes := Nodes{}

	resp, err := n.R.Get(n.u.String() + NodeEndpoint)
	if err != nil {
		log.Errorf("GetNode failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetNode failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &nodes)
	if err != nil {
		log.Tracef("Failed to deserialize node info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("node info deserialization failure: %w", err)
	}

	return &nodes, nil
}

func (n *nodeClient) GetNodesBySite(id string) (*NodesBySite, error) {
	log.Debugf("Getting all nodes by site.")

	nodes := NodesBySite{}

	resp, err := n.R.Get(n.u.String() + NodeEndpoint + "/sites/" + id)
	if err != nil {
		log.Errorf("GetNodeBySite failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetNodeBySite failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &nodes)
	if err != nil {
		log.Tracef("Failed to deserialize node info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("node info deserialization failure: %w", err)
	}

	return &nodes, nil
}

func (n *nodeClient) List(req ListNodesRequest) (*ListNodesResponse, error) {
	log.Debugf("Listing nodes: %v", req)

	nodes := ListNodesResponse{}

	queryParams := url.Values{}

	if req.NetworkId != "" {
		queryParams.Add("network_id", req.NetworkId)
	}
	if req.SiteId != "" {
		queryParams.Add("site_id", req.SiteId)
	}
	if req.Connectivity != "" {
		queryParams.Add("connectivity", req.Connectivity)
	}

	if req.State != "" {
		queryParams.Add("state", req.State)
	}
	if req.NodeId != "" {
		queryParams.Add("node_id", req.NodeId)
	}
	if req.Type != "" {
		queryParams.Add("type", req.Type)
	}

	resp, err := n.R.Get(n.u.String() + NodeEndpoint + "/list?" + queryParams.Encode())
	if err != nil {
		log.Errorf("ListNodes failure. error: %s", err.Error())
		return nil, fmt.Errorf("ListNodes failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &nodes)
	if err != nil {
		log.Tracef("Failed to deserialize node info. Error message is: %s", err.Error())
		return nil, fmt.Errorf("node info deserialization failure: %w", err)
	}

	return &nodes, nil
}