/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ukama/ukama/systems/node/api-gateway/pkg/client"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	cmocks "github.com/ukama/ukama/systems/node/controller/pb/gen/mocks"
)

const nodeID = "uk-sa2645-tnode-v0-aaaa"

func TestController_RestartNode(t *testing.T) {
	t.Run("SendsAOneNodeList", func(t *testing.T) {
		c := &cmocks.ControllerServiceClient{}
		c.On("RestartNodes", mock.Anything,
			&pb.RestartNodesRequest{NodeIds: []string{nodeID}}).
			Return(&pb.RestartNodesResponse{Operations: []*pb.RestartNodeResponse{
				{OperationId: "op-1", Status: "RUNNING"},
			}}, nil).Once()

		resp, err := client.NewControllerFromClient(c).RestartNode(nodeID)

		assert.NoError(t, err)
		assert.Equal(t, "op-1", resp.OperationId)
		c.AssertExpectations(t)
	})

	t.Run("PassesErrorsThrough", func(t *testing.T) {
		c := &cmocks.ControllerServiceClient{}
		c.On("RestartNodes", mock.Anything, mock.Anything).
			Return(nil, errors.New("busy")).Once()

		_, err := client.NewControllerFromClient(c).RestartNode(nodeID)

		assert.EqualError(t, err, "busy")
	})

	t.Run("EmptyResponseIsAnError", func(t *testing.T) {
		c := &cmocks.ControllerServiceClient{}
		c.On("RestartNodes", mock.Anything, mock.Anything).
			Return(&pb.RestartNodesResponse{}, nil).Once()

		_, err := client.NewControllerFromClient(c).RestartNode(nodeID)

		assert.Error(t, err)
	})
}

func TestController_ToggleRadio(t *testing.T) {
	t.Run("SendsAOneNodeList", func(t *testing.T) {
		c := &cmocks.ControllerServiceClient{}
		c.On("ToggleRadios", mock.Anything,
			&pb.ToggleRadiosRequest{NodeIds: []string{nodeID}, State: "off"}).
			Return(&pb.ToggleRadiosResponse{Operations: []*pb.ToggleRadioResponse{
				{OperationId: "op-2", Status: "SUCCESS"},
			}}, nil).Once()

		resp, err := client.NewControllerFromClient(c).ToggleRadio(nodeID, "off")

		assert.NoError(t, err)
		assert.Equal(t, "op-2", resp.OperationId)
		c.AssertExpectations(t)
	})

	t.Run("PassesErrorsThrough", func(t *testing.T) {
		c := &cmocks.ControllerServiceClient{}
		c.On("ToggleRadios", mock.Anything, mock.Anything).
			Return(nil, errors.New("busy")).Once()

		_, err := client.NewControllerFromClient(c).ToggleRadio(nodeID, "on")

		assert.EqualError(t, err, "busy")
	})

	t.Run("EmptyResponseIsAnError", func(t *testing.T) {
		c := &cmocks.ControllerServiceClient{}
		c.On("ToggleRadios", mock.Anything, mock.Anything).
			Return(&pb.ToggleRadiosResponse{}, nil).Once()

		_, err := client.NewControllerFromClient(c).ToggleRadio(nodeID, "on")

		assert.Error(t, err)
	})
}
