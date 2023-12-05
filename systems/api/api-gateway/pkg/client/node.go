/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import cclient "github.com/ukama/ukama/systems/common/rest/client"

type Node interface {
	GetNode(string) (*cclient.NodeInfo, error)
	RegisterNode(string, string, string, string) (*cclient.NodeInfo, error)
	AttachNode(string, string, string) error
	DetachNode(string) error
	AddNodeToSite(string, string, string) error
	RemoveNodeFromSite(string) error
	DeleteNode(string) error
}

type node struct {
	nc cclient.NodeClient
}

func NewNodeClientSet(nd cclient.NodeClient) Node {
	n := &node{
		nc: nd,
	}

	return n
}

func (n *node) GetNode(id string) (*cclient.NodeInfo, error) {
	node, err := n.nc.Get(id)
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	return node, nil
}

func (n *node) RegisterNode(nodeId, nodeName, orgId, state string) (*cclient.NodeInfo, error) {
	node, err := n.nc.Add(cclient.AddNodeRequest{
		NodeId: nodeId,
		Name:   nodeName,
		OrgId:  orgId,
		State:  state,
	})
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	return node, nil
}

func (n *node) AttachNode(id, left, right string) error {
	err := n.nc.Attach(id, cclient.AttachNodesRequest{
		AmpNodeL: left,
		AmpNodeR: right,
	})
	if err != nil {
		return handleRestErrorStatus(err)
	}

	return nil
}

func (n *node) DetachNode(id string) error {
	err := n.nc.Detach(id)
	if err != nil {
		return handleRestErrorStatus(err)
	}

	return nil
}

func (n *node) AddNodeToSite(id, networkId, siteId string) error {
	err := n.nc.AddToSite(id, cclient.AddToSiteRequest{
		NetworkId: networkId,
		SiteId:    siteId,
	})
	if err != nil {
		return handleRestErrorStatus(err)
	}

	return nil
}

func (n *node) RemoveNodeFromSite(id string) error {
	err := n.nc.RemoveFromSite(id)
	if err != nil {
		return handleRestErrorStatus(err)
	}

	return nil
}

func (n *node) DeleteNode(id string) error {
	err := n.nc.Delete(id)
	if err != nil {
		return handleRestErrorStatus(err)
	}

	return nil
}
