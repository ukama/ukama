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

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	testRoutingKey = "event.cloud.local.org.messaging.broadcast.nodefeeder"
	testTarget     = "node-feeder"
	testPath       = "/v1/sync"
	testNode1      = "node-1"
	testNode2      = "node-2"
	testNetworkID  = "network-123"
	errMarshalNF   = "failed to marshal node feeder message: %v"
)

type mockMsgBusClient struct {
	publishErr   error
	publishCalls int
	routes       []string
	messages     []*epb.NodeFeederMessage
}

func (m *mockMsgBusClient) Register() error { return nil }

func (m *mockMsgBusClient) Start() error { return nil }

func (m *mockMsgBusClient) Stop() error { return nil }

func (m *mockMsgBusClient) PublishRequest(route string, msg protoreflect.ProtoMessage) error {
	m.publishCalls++
	m.routes = append(m.routes, route)
	if nf, ok := msg.(*epb.NodeFeederMessage); ok {
		m.messages = append(m.messages, proto.Clone(nf).(*epb.NodeFeederMessage))
	}
	return m.publishErr
}

type mockNodeClient struct {
	listResp  *creg.ListNodesResponse
	listErr   error
	lastReq   creg.ListNodesRequest
	listCalls int
}

func (m *mockNodeClient) Get(string) (*creg.NodeInfo, error) { return nil, nil }

func (m *mockNodeClient) GetAll() (*creg.Nodes, error) { return nil, nil }

func (m *mockNodeClient) GetNodesBySite(string) (*creg.NodesBySite, error) { return nil, nil }

func (m *mockNodeClient) List(req creg.ListNodesRequest) (*creg.ListNodesResponse, error) {
	m.listCalls++
	m.lastReq = req
	return m.listResp, m.listErr
}

func (m *mockNodeClient) Add(creg.AddNodeRequest) (*creg.NodeInfo, error) { return nil, nil }

func (m *mockNodeClient) Attach(string, creg.AttachNodesRequest) error { return nil }

func (m *mockNodeClient) Detach(string) error { return nil }

func (m *mockNodeClient) AddToSite(string, creg.AddToSiteRequest) error { return nil }

func (m *mockNodeClient) RemoveFromSite(string) error { return nil }

func (m *mockNodeClient) Delete(string) error { return nil }

func TestNodeFeederBroadcastUnmarshalError(t *testing.T) {
	msgBus := &mockMsgBusClient{}
	nodeClient := &mockNodeClient{}
	s := NewBroadcasterServer("org", msgBus, nodeClient, false)

	err := s.NodeFeederBroadcast(context.Background(), &epb.BroadcasterEvent{
		Msg:        []byte("not-a-valid-proto-payload"),
		RoutingKey: testRoutingKey,
	})
	if err == nil {
		t.Fatalf("expected unmarshal error, got nil")
	}
	if nodeClient.listCalls != 0 {
		t.Fatalf("expected no node list calls, got %d", nodeClient.listCalls)
	}
	if msgBus.publishCalls != 0 {
		t.Fatalf("expected no publish calls, got %d", msgBus.publishCalls)
	}
}

func TestNodeFeederBroadcastUnknownScopePublishesAsIs(t *testing.T) {
	nf := &epb.NodeFeederMessage{
		Target:     testTarget,
		HttpMethod: "POST",
		Path:       testPath,
		Msg:        []byte("payload"),
	}
	encoded, err := proto.Marshal(nf)
	if err != nil {
		t.Fatalf(errMarshalNF, err)
	}

	msgBus := &mockMsgBusClient{}
	nodeClient := &mockNodeClient{}
	s := NewBroadcasterServer("org", msgBus, nodeClient, false)

	err = s.NodeFeederBroadcast(context.Background(), &epb.BroadcasterEvent{
		Msg:        encoded,
		Scope:      epb.BroadcastScope_UNKNOWN_SCOPE,
		RoutingKey: testRoutingKey,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if nodeClient.listCalls != 0 {
		t.Fatalf("expected no node list calls, got %d", nodeClient.listCalls)
	}
	if msgBus.publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", msgBus.publishCalls)
	}
	if len(msgBus.messages) != 1 || msgBus.messages[0].NodeId != "" {
		t.Fatalf("expected published message with empty node_id")
	}
}

func TestNodeFeederBroadcastNetworkScopePublishesForEachNode(t *testing.T) {
	nf := &epb.NodeFeederMessage{
		Target:     testTarget,
		HttpMethod: "POST",
		Path:       testPath,
		Msg:        []byte("payload"),
	}
	encoded, err := proto.Marshal(nf)
	if err != nil {
		t.Fatalf(errMarshalNF, err)
	}

	msgBus := &mockMsgBusClient{}
	nodeClient := &mockNodeClient{
		listResp: &creg.ListNodesResponse{
			Nodes: []*creg.NodeInfo{
				{Id: testNode1},
				{Id: testNode2},
			},
		},
	}
	s := NewBroadcasterServer("org", msgBus, nodeClient, false)

	err = s.NodeFeederBroadcast(context.Background(), &epb.BroadcasterEvent{
		Msg:        encoded,
		Scope:      epb.BroadcastScope_NETWORK_SCOPE,
		TargetId:   testNetworkID,
		RoutingKey: testRoutingKey,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if nodeClient.listCalls != 1 {
		t.Fatalf("expected one node list call, got %d", nodeClient.listCalls)
	}
	if nodeClient.lastReq.NetworkId != testNetworkID || nodeClient.lastReq.SiteId != "" {
		t.Fatalf("unexpected list request: %+v", nodeClient.lastReq)
	}
	if msgBus.publishCalls != 2 {
		t.Fatalf("expected two publish calls, got %d", msgBus.publishCalls)
	}
	if len(msgBus.messages) != 2 || msgBus.messages[0].NodeId != testNode1 || msgBus.messages[1].NodeId != testNode2 {
		t.Fatalf("unexpected published node ids: %+v", msgBus.messages)
	}
}

func TestNodeFeederBroadcastSiteScopeListFailure(t *testing.T) {
	nf := &epb.NodeFeederMessage{
		Target:     testTarget,
		HttpMethod: "POST",
		Path:       testPath,
		Msg:        []byte("payload"),
	}
	encoded, err := proto.Marshal(nf)
	if err != nil {
		t.Fatalf(errMarshalNF, err)
	}

	listErr := errors.New("registry unavailable")
	msgBus := &mockMsgBusClient{}
	nodeClient := &mockNodeClient{
		listErr: listErr,
	}
	s := NewBroadcasterServer("org", msgBus, nodeClient, false)

	err = s.NodeFeederBroadcast(context.Background(), &epb.BroadcasterEvent{
		Msg:        encoded,
		Scope:      epb.BroadcastScope_SITE_SCOPE,
		TargetId:   "site-123",
		RoutingKey: testRoutingKey,
	})
	if !errors.Is(err, listErr) {
		t.Fatalf("expected list error, got %v", err)
	}
	if nodeClient.lastReq.SiteId != "site-123" || nodeClient.lastReq.NetworkId != "" {
		t.Fatalf("unexpected list request: %+v", nodeClient.lastReq)
	}
	if msgBus.publishCalls != 0 {
		t.Fatalf("expected no publish calls, got %d", msgBus.publishCalls)
	}
}

func TestNodeFeederBroadcastPublishFailureStopsIteration(t *testing.T) {
	nf := &epb.NodeFeederMessage{
		Target:     testTarget,
		HttpMethod: "POST",
		Path:       testPath,
		Msg:        []byte("payload"),
	}
	encoded, err := proto.Marshal(nf)
	if err != nil {
		t.Fatalf(errMarshalNF, err)
	}

	publishErr := errors.New("publish failed")
	msgBus := &mockMsgBusClient{
		publishErr: publishErr,
	}
	nodeClient := &mockNodeClient{
		listResp: &creg.ListNodesResponse{
			Nodes: []*creg.NodeInfo{
				{Id: testNode1},
				{Id: testNode2},
			},
		},
	}
	s := NewBroadcasterServer("org", msgBus, nodeClient, false)

	err = s.NodeFeederBroadcast(context.Background(), &epb.BroadcasterEvent{
		Msg:        encoded,
		Scope:      epb.BroadcastScope_NETWORK_SCOPE,
		TargetId:   testNetworkID,
		RoutingKey: testRoutingKey,
	})
	if !errors.Is(err, publishErr) {
		t.Fatalf("expected publish error, got %v", err)
	}
	if msgBus.publishCalls != 1 {
		t.Fatalf("expected publish loop to stop after first error, got %d calls", msgBus.publishCalls)
	}
}
