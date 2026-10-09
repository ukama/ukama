/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package delivery

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

type fakeService struct {
	code  codes.Code
	calls atomic.Int32
}

func (f *fakeService) EventNotification(ctx context.Context, in *pb.Event, opts ...grpc.CallOption) (*pb.EventResponse, error) {
	f.calls.Add(1)
	if f.code == codes.OK {
		return &pb.EventResponse{}, nil
	}
	return nil, status.Error(f.code, "failed")
}

func newListener(code codes.Code) (*Listener, *fakeService, *atomic.Int32) {
	svc := &fakeService{code: code}
	paused := &atomic.Int32{}
	l := &Listener{
		service:       "node",
		timeout:       time.Second,
		client:        svc,
		onUnreachable: func() { paused.Add(1) },
	}
	return l, svc, paused
}

func event(t *testing.T) rabbitmq.Delivery {
	t.Helper()
	msg, err := anypb.New(&emptypb.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	body, err := proto.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	return rabbitmq.Delivery{Delivery: amqp.Delivery{RoutingKey: "event.cloud.local.ukama.registry.node.node.create", Body: body}}
}

func TestDeliveredEventsAreAcked(t *testing.T) {
	for _, code := range []codes.Code{codes.OK, codes.Internal, codes.Unimplemented, codes.InvalidArgument} {
		l, _, paused := newListener(code)
		if got := l.handle(event(t)); got != rabbitmq.Ack {
			t.Errorf("%s: expected Ack, got %v", code, got)
		}
		if paused.Load() != 0 {
			t.Errorf("%s: delivery paused", code)
		}
	}
}

func TestUndeliveredEventsAreRequeuedAndPauseDelivery(t *testing.T) {
	for _, code := range []codes.Code{codes.Unavailable, codes.DeadlineExceeded} {
		l, svc, paused := newListener(code)
		if got := l.handle(event(t)); got != rabbitmq.NackRequeue {
			t.Errorf("%s: expected NackRequeue, got %v", code, got)
		}

		if got := l.handle(event(t)); got != rabbitmq.NackRequeue || svc.calls.Load() != 1 {
			t.Errorf("%s: expected later events requeued without delivery", code)
		}

		deadline := time.Now().Add(time.Second)
		for paused.Load() == 0 && time.Now().Before(deadline) {
			time.Sleep(time.Millisecond)
		}
		if paused.Load() != 1 {
			t.Errorf("%s: expected one pause, got %d", code, paused.Load())
		}
	}
}

func TestUnparseableEventIsAcked(t *testing.T) {
	l, svc, _ := newListener(codes.OK)
	d := rabbitmq.Delivery{Delivery: amqp.Delivery{RoutingKey: "x", Body: []byte{0xff, 0xff}}}

	if got := l.handle(d); got != rabbitmq.Ack || svc.calls.Load() != 0 {
		t.Fatalf("expected Ack without delivery, got %v", got)
	}
}
