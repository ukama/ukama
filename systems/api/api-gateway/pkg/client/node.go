/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client/rest"
)

type Node interface {
	GetNode(string) (*rest.NodeInfo, error)
	RegisterNode(string, string, string, string) (*rest.NodeInfo, error)
	AttachNode(string, string, string) error
	DetachNode(string) error
	AddNodeToSite(string, string, string) error
	RemoveNodeFromSite(string) error
	DeleteNode(string) error
}

type node struct {
	nc rest.NodeClient
}

func NewNodeClientSet(nd rest.NodeClient) Node {
	n := &node{
		nc: nd,
	}

	return n
}

func (n *node) GetNode(id string) (*rest.NodeInfo, error) {
	node, err := n.nc.Get(id)
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	return node, nil
}

func (n *node) RegisterNode(nodeId, nodeName, orgId, state string) (*rest.NodeInfo, error) {
	node, err := n.nc.Add(rest.AddNodeRequest{
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
	err := n.nc.Attach(id, rest.AttachNodesRequest{
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
	err := n.nc.AddToSite(id, rest.AddToSiteRequest{
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
