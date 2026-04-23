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
	"strings"
	"testing"

	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	testOrgName      = "org"
	broadcastRouteT  = "event.cloud.local.{{ .Org}}.ukamaagent.asr.policies.publish"
	unknownRouteKey  = "event.cloud.local.org.ukamaagent.asr.publish.unknown"
	errWrapAnyEvent  = "failed to wrap broadcaster event in any: %v"
	errWrapAnyGeneric = "failed to wrap generic message in any: %v"
)

func TestEventNotificationUnknownRouteReturnsError(t *testing.T) {
	eventServer := NewBroadcasterEventServer(
		testOrgName,
		NewBroadcasterServer(testOrgName, &mockMsgBusClient{}, &mockNodeClient{}, false),
	)

	resp, err := eventServer.EventNotification(context.Background(), &epb.Event{
		RoutingKey: unknownRouteKey,
		Msg:        &anypb.Any{},
	})
	if err == nil || !strings.Contains(err.Error(), "no handler routing key") {
		t.Fatalf("expected unknown routing key error, got %v", err)
	}
	if resp != nil {
		t.Fatalf("expected nil response for unknown route")
	}
}

func TestEventNotificationInvalidBroadcasterPayloadReturnsError(t *testing.T) {
	invalidAny, err := anypb.New(&epb.EventResponse{})
	if err != nil {
		t.Fatalf(errWrapAnyGeneric, err)
	}

	eventServer := NewBroadcasterEventServer(
		testOrgName,
		NewBroadcasterServer(testOrgName, &mockMsgBusClient{}, &mockNodeClient{}, false),
	)

	resp, err := eventServer.EventNotification(context.Background(), &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, broadcastRouteT),
		Msg:        invalidAny,
	})
	if err == nil {
		t.Fatalf("expected unmarshal error, got nil")
	}
	if resp != nil {
		t.Fatalf("expected nil response on unmarshal error")
	}
}

func TestEventNotificationUnsupportedBroadcastTypeReturnsError(t *testing.T) {
	msg := &epb.BroadcasterEvent{
		Type:       epb.BroadcastType_UNKNOWN_BROADCAST,
		RoutingKey: testRoutingKey,
	}
	anyMsg, err := anypb.New(msg)
	if err != nil {
		t.Fatalf(errWrapAnyEvent, err)
	}

	eventServer := NewBroadcasterEventServer(
		testOrgName,
		NewBroadcasterServer(testOrgName, &mockMsgBusClient{}, &mockNodeClient{}, false),
	)

	resp, err := eventServer.EventNotification(context.Background(), &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, broadcastRouteT),
		Msg:        anyMsg,
	})
	if err == nil || !strings.Contains(err.Error(), "no handler broadcast type") {
		t.Fatalf("expected unsupported broadcast type error, got %v", err)
	}
	if resp != nil {
		t.Fatalf("expected nil response for unsupported broadcast type")
	}
}

func TestEventNotificationNodeBroadcastSuccess(t *testing.T) {
	nf := &epb.NodeFeederMessage{
		Target:     testTarget,
		HttpMethod: "POST",
		Path:       testPath,
		Msg:        []byte("payload"),
	}
	nfBytes, err := proto.Marshal(nf)
	if err != nil {
		t.Fatalf("failed to marshal node feeder message: %v", err)
	}

	msg := &epb.BroadcasterEvent{
		Msg:        nfBytes,
		Type:       epb.BroadcastType_NODE_BROADCAST,
		Scope:      epb.BroadcastScope_UNKNOWN_SCOPE,
		RoutingKey: testRoutingKey,
	}
	anyMsg, err := anypb.New(msg)
	if err != nil {
		t.Fatalf(errWrapAnyEvent, err)
	}

	msgBus := &mockMsgBusClient{}
	eventServer := NewBroadcasterEventServer(
		testOrgName,
		NewBroadcasterServer(testOrgName, msgBus, &mockNodeClient{}, false),
	)

	resp, err := eventServer.EventNotification(context.Background(), &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, broadcastRouteT),
		Msg:        anyMsg,
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp == nil {
		t.Fatalf("expected non-nil response")
	}
	if msgBus.publishCalls != 1 {
		t.Fatalf("expected one publish call, got %d", msgBus.publishCalls)
	}
}

func TestEventNotificationNodeBroadcastReturnsUnderlyingError(t *testing.T) {
	msg := &epb.BroadcasterEvent{
		Msg:        []byte("invalid-node-feeder-msg"),
		Type:       epb.BroadcastType_NODE_BROADCAST,
		Scope:      epb.BroadcastScope_UNKNOWN_SCOPE,
		RoutingKey: testRoutingKey,
	}
	anyMsg, err := anypb.New(msg)
	if err != nil {
		t.Fatalf(errWrapAnyEvent, err)
	}

	eventServer := NewBroadcasterEventServer(
		testOrgName,
		NewBroadcasterServer(testOrgName, &mockMsgBusClient{}, &mockNodeClient{}, false),
	)

	resp, err := eventServer.EventNotification(context.Background(), &epb.Event{
		RoutingKey: msgbus.PrepareRoute(testOrgName, broadcastRouteT),
		Msg:        anyMsg,
	})
	if err == nil {
		t.Fatalf("expected downstream broadcaster error, got nil")
	}
	if resp != nil {
		t.Fatalf("expected nil response when broadcaster fails")
	}
}
