/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package queue

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocks "github.com/ukama/ukama/systems/common/mocks"
	mb "github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/services/msgClient/internal/db"
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
