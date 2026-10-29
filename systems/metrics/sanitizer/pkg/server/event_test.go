/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/metrics/sanitizer/pkg/server"

	evt "github.com/ukama/ukama/systems/common/events"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const testOrgName = "testorg"

func newEventServer(t *testing.T) (*server.SanitizerEventServer, *atomic.Int32) {
	syncs := &atomic.Int32{}

	registrySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		syncs.Add(1)

		_ = json.NewEncoder(w).Encode(registry.Nodes{Nodes: []*registry.NodeInfo{}})
	}))
	t.Cleanup(registrySrv.Close)

	ss, err := server.NewSanitizerServer(registrySrv.URL, "http://localhost:0", testOrgName,
		"testorg-id", nil)
	require.NoError(t, err)

	syncs.Store(0)

	return server.NewSanitizerEventServer(testOrgName, ss, nil), syncs
}

func newEvent(t *testing.T, key evt.EventId, msg proto.Message) *epb.Event {
	anyMsg, err := anypb.New(msg)
	require.NoError(t, err)

	return &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, evt.EventRoutingKey[key]),
		Msg:        anyMsg,
	}
}

func TestSanitizerEventServer_EventNotification(t *testing.T) {
	tests := []struct {
		name string
		key  evt.EventId
		msg  proto.Message
	}{
		{
			name: "NodeAssign",
			key:  evt.EventNodeAssign,
			msg: &epb.EventRegistryNodeAssign{NodeId: "uk-sa2341-tnode-a1-1234",
				Site: "site-1", Network: "net-1"},
		},
		{
			name: "NodeRelease",
			key:  evt.EventNodeRelease,
			msg:  &epb.NodeReleasedEvent{NodeId: "uk-sa2341-tnode-a1-1234"},
		},
		{
			name: "NodeDelete",
			key:  evt.EventNodeDelete,
			msg:  &epb.EventRegistryNodeDelete{NodeId: "uk-sa2341-tnode-a1-1234"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			es, syncs := newEventServer(t)

			resp, err := es.EventNotification(context.TODO(), newEvent(t, tc.key, tc.msg))

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, int32(1), syncs.Load())
		})
	}

	t.Run("MismatchedMessageType", func(t *testing.T) {
		es, syncs := newEventServer(t)

		_, err := es.EventNotification(context.TODO(), newEvent(t, evt.EventNodeDelete,
			&epb.NodeReleasedEvent{NodeId: "uk-sa2341-tnode-a1-1234"}))

		assert.Error(t, err)
		assert.Equal(t, int32(0), syncs.Load())
	})

	t.Run("UnknownRoutingKey", func(t *testing.T) {
		es, syncs := newEventServer(t)

		_, err := es.EventNotification(context.TODO(), &epb.Event{
			RoutingKey: "event.cloud.local.testorg.registry.node.node.update",
		})

		assert.Error(t, err)
		assert.Equal(t, int32(0), syncs.Load())
	})
}
