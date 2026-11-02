/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	cmocks "github.com/ukama/ukama/systems/common/mocks"
	pb "github.com/ukama/ukama/systems/common/pb/gen/msgclient"
)

const route = "event.cloud.local.ukama.registry.node.node.create"

type fakeRegistry struct {
	registered []string
}

func (f *fakeRegistry) Register(name, system, instance, uri string, routes []string) (string, error) {
	f.registered = append(f.registered, name)
	return "uuid-" + name, nil
}
func (f *fakeRegistry) Resume(id string) error     { return nil }
func (f *fakeRegistry) Pause(id string) error      { return nil }
func (f *fakeRegistry) Unregister(id string) error { return nil }
func (f *fakeRegistry) Lookup(id string) (string, string, bool) {
	if id == "uuid-node" {
		return "node", "1", true
	}
	return "", "", false
}

type fakePublisher struct {
	err     error
	headers rabbitmq.Table
}

func (f *fakePublisher) Publish(key string, body []byte, headers rabbitmq.Table) error {
	f.headers = headers
	return f.err
}

func newServer(pub *fakePublisher) (*MsgClientServer, *fakeRegistry) {
	r := &fakeRegistry{}
	return NewMsgClientServer(r, pub, &cmocks.MsgBusShovelProvider{}, "registry", 2*time.Second), r
}

func TestRegisterServiceReturnsUuidAndRefreshInterval(t *testing.T) {
	s, r := newServer(&fakePublisher{})

	resp, err := s.RegisterService(context.Background(), &pb.RegisterServiceReq{
		SystemName: "registry", ServiceName: "node", InstanceId: "1", ServiceURI: "node:9090", Routes: []string{route},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.State != pb.REGISTRAION_STATUS_REGISTERED || resp.ServiceUuid != "uuid-node" || resp.RefreshInterval != 2 {
		t.Fatalf("unexpected response %+v", resp)
	}
	if len(r.registered) != 1 {
		t.Fatalf("expected one registration")
	}
}

func TestRegisterServiceRejectsOtherSystem(t *testing.T) {
	s, r := newServer(&fakePublisher{})

	_, err := s.RegisterService(context.Background(), &pb.RegisterServiceReq{SystemName: "billing", ServiceName: "node"})
	if status.Code(err) != codes.InvalidArgument || len(r.registered) != 0 {
		t.Fatalf("expected InvalidArgument without registration, got %v", err)
	}
}

func TestPublishMsgSetsSourceHeaders(t *testing.T) {
	pub := &fakePublisher{}
	s, _ := newServer(pub)
	msg, _ := anypb.New(&emptypb.Empty{})

	if _, err := s.PublishMsg(context.Background(), &pb.PublishMsgRequest{ServiceUuid: "uuid-node", RoutingKey: route, Msg: msg}); err != nil {
		t.Fatal(err)
	}
	if pub.headers["source-service"] != "node" || pub.headers["instance-id"] != "1" {
		t.Fatalf("unexpected headers %v", pub.headers)
	}

	if _, err := s.PublishMsg(context.Background(), &pb.PublishMsgRequest{ServiceUuid: "other", RoutingKey: route, Msg: msg}); err != nil {
		t.Fatal(err)
	}
	if pub.headers["source-service"] != "unknown" {
		t.Fatalf("expected unknown source, got %v", pub.headers)
	}
}

func TestPublishMsgFailsWhenNotAccepted(t *testing.T) {
	s, _ := newServer(&fakePublisher{err: errors.New("no confirm")})
	msg, _ := anypb.New(&emptypb.Empty{})

	_, err := s.PublishMsg(context.Background(), &pb.PublishMsgRequest{ServiceUuid: "uuid-node", RoutingKey: route, Msg: msg})
	if status.Code(err) != codes.Unavailable {
		t.Fatalf("expected Unavailable, got %v", err)
	}
}
