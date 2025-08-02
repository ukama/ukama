/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/common/uuid"

	crest "github.com/ukama/ukama/systems/common/rest"
	cclient "github.com/ukama/ukama/systems/common/rest/client"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
)

func TestCient_GetNode(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	nodeId := "uk-sa2341-hnode-v0-a1a0"
	nodeName := "node-1"

	n := client.NewNodeClientSet(nodeClient)

	t.Run("NodeFound", func(t *testing.T) {
		nodeClient.On("Get", nodeId).
			Return(&creg.NodeInfo{
				Id:   nodeId,
				Name: nodeName,
			}, nil).Once()

		nodeInfo, err := n.GetNode(nodeId)

		assert.NoError(t, err)

		assert.NotNil(t, nodeInfo)
		assert.Equal(t, nodeInfo.Id, nodeId)
		assert.Equal(t, nodeInfo.Name, nodeName)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeClient.On("Get", nodeId).
			Return(nil, fmt.Errorf("GetNode failure: %w",
				&cclient.ErrorStatus{
					StatusCode: 404,
					RawError:   crest.ErrorResponse{Err: "not found"},
				})).Once()

		nodeInfo, err := n.GetNode(nodeId)

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "404")

		assert.Nil(t, nodeInfo)
	})

	t.Run("NodeGetError", func(t *testing.T) {
		nodeClient.On("Get", nodeId).
			Return(nil,
				fmt.Errorf("Some unexpected error")).Once()

		nodeInfo, err := n.GetNode(nodeId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")

		assert.Nil(t, nodeInfo)
	})
}

func TestCient_RegisterNode(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	nodeId := "uk-sa2341-hnode-v0-a1a0"
	nodeName := "node-1"
	orgId := uuid.NewV4()
	state := "pending"

	n := client.NewNodeClientSet(nodeClient)

	t.Run("NodeRegistered", func(t *testing.T) {
		nodeClient.On("Add", creg.AddNodeRequest{
			NodeId: nodeId,
			Name:   nodeName,
			OrgId:  orgId.String(),
			State:  state,
		}).Return(&creg.NodeInfo{
			Id:   nodeId,
			Name: nodeName,
		}, nil).Once()

		nodeInfo, err := n.RegisterNode(nodeId, nodeName, orgId.String(), state)

		assert.NoError(t, err)

		assert.Equal(t, nodeInfo.Id, nodeId)
		assert.Equal(t, nodeInfo.Name, nodeName)
	})

	t.Run("NodeNotRegistered", func(t *testing.T) {
		nodeClient.On("Add", creg.AddNodeRequest{
			NodeId: nodeId,
			Name:   nodeName,
			OrgId:  orgId.String(),
			State:  state,
		}).Return(nil, errors.New("some error")).Once()

		nodeInfo, err := n.RegisterNode(nodeId, nodeName, orgId.String(), state)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
		assert.Nil(t, nodeInfo)
	})
}

func TestCient_AttachNode(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	ampNodeL := "uk-sa2341-anode-v0-a1a0"
	ampNodeR := "uk-sa2341-anode-v0-a1a1"

	n := client.NewNodeClientSet(nodeClient)

	t.Run("NodeAttached", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		nodeClient.On("Attach", nodeId, creg.AttachNodesRequest{
			AmpNodeL: ampNodeL,
			AmpNodeR: ampNodeR,
		}).Return(nil).Once()

		err := n.AttachNode(nodeId, ampNodeL, ampNodeR)

		assert.NoError(t, err)
	})

	t.Run("NodeNotAttached", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a1"

		nodeClient.On("Attach", nodeId, creg.AttachNodesRequest{
			AmpNodeL: ampNodeL,
			AmpNodeR: ampNodeR,
		}).Return(errors.New("some error")).Once()

		err := n.AttachNode(nodeId, ampNodeL, ampNodeR)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
	})
}

func TestCient_DetachNode(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	nodeId := "uk-sa2341-hnode-v0-a1a0"

	n := client.NewNodeClientSet(nodeClient)

	t.Run("NodeDetached", func(t *testing.T) {
		nodeClient.On("Detach", nodeId).
			Return(nil).Once()

		err := n.DetachNode(nodeId)

		assert.NoError(t, err)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeClient.On("Detach", nodeId).
			Return(fmt.Errorf("DetachNode failure: %w",
				&cclient.ErrorStatus{
					StatusCode: 404,
					RawError:   crest.ErrorResponse{Err: "not found"},
				})).Once()

		err := n.DetachNode(nodeId)

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("NodeDetachError", func(t *testing.T) {
		nodeClient.On("Detach", nodeId).
			Return(fmt.Errorf("Some unexpected error")).Once()

		err := n.DetachNode(nodeId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
	})
}

func TestCient_AddToSite(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	networkId := uuid.NewV4().String()
	siteId := uuid.NewV4().String()

	n := client.NewNodeClientSet(nodeClient)

	t.Run("NodeAdded", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a0"

		nodeClient.On("AddToSite", nodeId, creg.AddToSiteRequest{
			NetworkId: networkId,
			SiteId:    siteId,
		}).Return(nil).Once()

		err := n.AddNodeToSite(nodeId, networkId, siteId)

		assert.NoError(t, err)
	})

	t.Run("NodeNotAdded", func(t *testing.T) {
		nodeId := "uk-sa2341-tnode-v0-a1a1"

		nodeClient.On("AddToSite", nodeId, creg.AddToSiteRequest{
			NetworkId: networkId,
			SiteId:    siteId,
		}).Return(errors.New("some error")).Once()

		err := n.AddNodeToSite(nodeId, networkId, siteId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
	})
}

func TestCient_RemoveNodeFromSite(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	nodeId := "uk-sa2341-hnode-v0-a1a0"

	n := client.NewNodeClientSet(nodeClient)

	t.Run("NodeRemoved", func(t *testing.T) {
		nodeClient.On("RemoveFromSite", nodeId).
			Return(nil).Once()

		err := n.RemoveNodeFromSite(nodeId)

		assert.NoError(t, err)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeClient.On("RemoveFromSite", nodeId).
			Return(fmt.Errorf("DetachNode failure: %w",
				&cclient.ErrorStatus{
					StatusCode: 404,
					RawError:   crest.ErrorResponse{Err: "not found"},
				})).Once()

		err := n.RemoveNodeFromSite(nodeId)

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("NodeRemoveError", func(t *testing.T) {
		nodeClient.On("RemoveFromSite", nodeId).
			Return(fmt.Errorf("Some unexpected error")).Once()

		err := n.RemoveNodeFromSite(nodeId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
	})
}

func TestCient_DeleteNode(t *testing.T) {
	nodeClient := &mocks.NodeClient{}

	nodeId := "uk-sa2341-hnode-v0-a1a0"

	n := client.NewNodeClientSet(nodeClient)

	t.Run("NodeDeleted", func(t *testing.T) {
		nodeClient.On("Delete", nodeId).
			Return(nil).Once()

		err := n.DeleteNode(nodeId)

		assert.NoError(t, err)
	})

	t.Run("NodeNotFound", func(t *testing.T) {
		nodeClient.On("Delete", nodeId).
			Return(fmt.Errorf("DeleteNode failure: %w",
				&cclient.ErrorStatus{
					StatusCode: 404,
					RawError:   crest.ErrorResponse{Err: "not found"},
				})).Once()

		err := n.DeleteNode(nodeId)

		assert.Error(t, err)
		assert.IsType(t, err, crest.HttpError{})
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("NodeDeleteError", func(t *testing.T) {
		nodeClient.On("Delete", nodeId).
			Return(fmt.Errorf("Some unexpected error")).Once()

		err := n.DeleteNode(nodeId)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error")
	})
}
