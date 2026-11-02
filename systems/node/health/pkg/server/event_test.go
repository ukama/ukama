/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/node/health/mocks"
	"google.golang.org/protobuf/types/known/anypb"
)

func newNodeDeleteEvent(t *testing.T, nodeId string) *epb.Event {
	msg, err := anypb.New(&epb.EventRegistryNodeDelete{NodeId: nodeId})
	require.NoError(t, err)

	return &epb.Event{RoutingKey: "event.cloud.local.testorg.registry.node.node.delete", Msg: msg}
}

func TestHealthEventServer_NodeDelete(t *testing.T) {
	t.Run("DeletesNodeHealthData", func(t *testing.T) {
		hRepo := &mocks.HealthRepo{}
		hRepo.On("DeleteNode", testNode.StringLowercase()).Return(nil).Once()
		es := NewHealthEventServer(testOrgName, newTestHealthServer(hRepo))

		_, err := es.EventNotification(context.Background(), newNodeDeleteEvent(t, testNode.String()))

		assert.NoError(t, err)
		hRepo.AssertExpectations(t)
	})

	t.Run("RepoError", func(t *testing.T) {
		hRepo := &mocks.HealthRepo{}
		hRepo.On("DeleteNode", testNode.StringLowercase()).Return(errors.New("db down")).Once()
		es := NewHealthEventServer(testOrgName, newTestHealthServer(hRepo))

		_, err := es.EventNotification(context.Background(), newNodeDeleteEvent(t, testNode.String()))

		assert.Error(t, err)
	})

	t.Run("EmptyNodeId", func(t *testing.T) {
		hRepo := &mocks.HealthRepo{}
		es := NewHealthEventServer(testOrgName, newTestHealthServer(hRepo))

		_, err := es.EventNotification(context.Background(), newNodeDeleteEvent(t, ""))

		assert.Error(t, err)
		hRepo.AssertNotCalled(t, "DeleteNode", "")
	})

	t.Run("UnknownRouteIsIgnored", func(t *testing.T) {
		hRepo := &mocks.HealthRepo{}
		es := NewHealthEventServer(testOrgName, newTestHealthServer(hRepo))

		_, err := es.EventNotification(context.Background(), &epb.Event{RoutingKey: "event.cloud.local.testorg.registry.node.node.create"})

		assert.NoError(t, err)
	})
}
