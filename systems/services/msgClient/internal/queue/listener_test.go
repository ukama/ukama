/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package queue

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocks "github.com/ukama/ukama/systems/common/mocks"
	mb "github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/services/msgClient/internal/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var route = []mb.RoutingKey{mb.RoutingKey("event.cloud.local.ukama.init.lookup.organization.create")}

var service = db.Service{
	Name:        "test",
	ServiceUuid: "1ce2fa2f-2997-422c-83bf-92cf2e7334dd",
	InstanceId:  "1",
	MsgBusUri:   "amqp://guest:guest@localhost:5672",
	ListQueue:   "",
	PublQueue:   "",
	Exchange:    "amq.topic",
	ServiceUri:  "localhost:9090",
	GrpcTimeout: 5,
	Routes:      []db.Route{{Key: "event.cloud.local.ukama.init.lookup.organization.create"}},
}

func NewTestQueueListener(s db.Service) *QueueListener {

	routes := make([]string, len(s.Routes))
	for idx, r := range s.Routes {
		/*  Create a queue listner for each service */
		routes[idx] = r.Key
	}

	ch := make(chan bool, 1)

	return &QueueListener{
		serviceUuid: s.ServiceUuid,
		serviceName: s.Name,
		serviceHost: s.ServiceUri,
		c:           ch,
		routes:      routes,
		queue:       s.ListQueue,
		exchange:    s.Exchange,
		mConn:       &mocks.Consumer{},
	}
}

func TestQueuePublisher_startstopQueueListening(t *testing.T) {
	client := &mocks.Consumer{}
	qp := NewTestQueueListener(service)
	qp.mConn = client

	client.On("SubscribeToServiceQueue", qp.serviceName, qp.exchange, route, qp.serviceUuid, mock.AnythingOfType("func(amqp.Delivery, chan<- bool)")).Return(nil).Once()
	client.On("Close").Return(nil).Once()

	go qp.startQueueListening()

	time.Sleep(2 * time.Second)

	qp.stopQueueListening()

	time.Sleep(2 * time.Second)

	client.AssertExpectations(t)

}

func TestQueueListenerRetriesFailedSubscribe(t *testing.T) {
	previous := listenerRetryMin
	listenerRetryMin = 10 * time.Millisecond
	defer func() { listenerRetryMin = previous }()

	client := &mocks.Consumer{}
	qp := NewTestQueueListener(service)
	qp.mConn = client

	handler := mock.AnythingOfType("func(amqp.Delivery, chan<- bool)")
	client.On("SubscribeToServiceQueue", qp.serviceName, qp.exchange, route, qp.serviceUuid, handler).
		Return(errors.New("dial tcp: connection refused")).Twice()
	client.On("SubscribeToServiceQueue", qp.serviceName, qp.exchange, route, qp.serviceUuid, handler).Return(nil).Once()
	client.On("Close").Return(nil).Once()

	qp.startQueueListening()
	require.Eventually(t, func() bool { return qp.state.Load() }, 2*time.Second, 5*time.Millisecond)

	qp.stopQueueListening()
	require.Eventually(t, func() bool { return !qp.state.Load() }, 2*time.Second, 5*time.Millisecond)
	client.AssertExpectations(t)
}

func TestQueueListenerStopsWhileRetrying(t *testing.T) {
	previous := listenerRetryMin
	listenerRetryMin = time.Hour
	defer func() { listenerRetryMin = previous }()

	client := &mocks.Consumer{}
	qp := NewTestQueueListener(service)
	qp.mConn = client

	client.On("SubscribeToServiceQueue", qp.serviceName, qp.exchange, route, qp.serviceUuid, mock.AnythingOfType("func(amqp.Delivery, chan<- bool)")).
		Return(errors.New("dial tcp: connection refused")).Once()
	client.On("Close").Return(nil).Once()

	qp.startQueueListening()
	require.Eventually(t, func() bool { return qp.retrying.Load() }, 2*time.Second, 5*time.Millisecond)

	qp.stopQueueListening()
	require.Eventually(t, func() bool { return !qp.retrying.Load() }, 2*time.Second, 5*time.Millisecond)
	require.False(t, qp.state.Load())
	client.AssertExpectations(t)
}

type fakeEventClient struct {
	errs  []error
	calls int
}

func (f *fakeEventClient) EventNotification(ctx context.Context, in *epb.Event, opts ...grpc.CallOption) (*epb.EventResponse, error) {
	f.calls++
	if f.calls <= len(f.errs) {
		return nil, f.errs[f.calls-1]
	}
	return &epb.EventResponse{}, nil
}

func newDeliveryListener(t *testing.T, errs ...error) (*QueueListener, *fakeEventClient) {
	prevMin, prevMax, prevHold := listenerRetryMin, listenerRetryMax, maxHold
	listenerRetryMin, listenerRetryMax = time.Millisecond, 5*time.Millisecond
	t.Cleanup(func() { listenerRetryMin, listenerRetryMax, maxHold = prevMin, prevMax, prevHold })

	conn, err := grpc.NewClient("localhost:0", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	f := &fakeEventClient{errs: errs}
	q := NewTestQueueListener(service)
	q.gConn = conn
	q.gClient = f
	q.grpcTimeout = time.Second
	q.stopped = make(chan struct{})
	return q, f
}

func eventDelivery(t *testing.T) amqp.Delivery {
	msg, err := anypb.New(&epb.EventResponse{})
	require.NoError(t, err)
	body, err := proto.Marshal(msg)
	require.NoError(t, err)
	return amqp.Delivery{RoutingKey: string(route[0]), Body: body}
}

func unavailable(n int) []error {
	errs := make([]error, n)
	for i := range errs {
		errs[i] = status.Error(codes.Unavailable, "connection refused")
	}
	return errs
}

func TestProcessEventMsgAcksHandlerError(t *testing.T) {
	q, f := newDeliveryListener(t, status.Error(codes.Unknown, "failed to unmarshal"))

	require.True(t, q.processEventMsg(eventDelivery(t)))
	require.Equal(t, 1, f.calls)
}

func TestProcessEventMsgRetriesUntilDelivered(t *testing.T) {
	q, f := newDeliveryListener(t,
		status.Error(codes.Unavailable, "connection refused"),
		status.Error(codes.DeadlineExceeded, "context deadline exceeded"))

	require.True(t, q.processEventMsg(eventDelivery(t)))
	require.Equal(t, 3, f.calls)
}

func TestProcessEventMsgRequeuesOnStop(t *testing.T) {
	q, _ := newDeliveryListener(t, unavailable(100000)...)

	result := make(chan bool, 1)
	go func() { result <- q.processEventMsg(eventDelivery(t)) }()

	time.Sleep(20 * time.Millisecond)
	close(q.stopped)

	select {
	case acked := <-result:
		require.False(t, acked)
	case <-time.After(2 * time.Second):
		t.Fatal("processEventMsg did not return after stop")
	}
}

func TestProcessEventMsgRequeuesAfterMaxHold(t *testing.T) {
	q, f := newDeliveryListener(t, unavailable(100000)...)
	maxHold = 20 * time.Millisecond

	require.False(t, q.processEventMsg(eventDelivery(t)))
	require.Greater(t, f.calls, 1)
}

func TestProcessEventMsgAcksUnparseableBody(t *testing.T) {
	q, f := newDeliveryListener(t)

	require.True(t, q.processEventMsg(amqp.Delivery{RoutingKey: string(route[0]), Body: []byte{0x0f}}))
	require.Equal(t, 0, f.calls)
}
