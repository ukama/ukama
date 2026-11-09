/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"errors"
	"testing"
	"time"

	amqp091 "github.com/rabbitmq/amqp091-go"
	amqp "github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/messaging/node-feeder/mocks"
)

const testTarget = "ukama.org.node.uk-sa2645-anode-v0-a90e"

var errUnavailable = errors.New("server error: 503")

type listenerFixture struct {
	listener *QueueListener
	exec     *mocks.RequestExecutor
	mult     *mocks.RequestMultiplier
	pub      *mocks.RetryPublisher
	now      time.Time
}

func newFixture(t *testing.T) *listenerFixture {
	f := &listenerFixture{
		exec: mocks.NewRequestExecutor(t),
		mult: mocks.NewRequestMultiplier(t),
		pub:  mocks.NewRetryPublisher(t),
		now:  time.UnixMilli(1_794_207_918_000),
	}
	f.listener = &QueueListener{
		requestExec:   f.exec,
		requestMult:   f.mult,
		maxRetryCount: 3,
		retryPub:      f.pub,
		retryBackoff:  backoffDurations([]int{2, 5, 15, 30}),
		maxRequestAge: 60 * time.Second,
		latest:        newLatestRequests(60 * time.Second),
		now:           func() time.Time { return f.now },
	}
	return f
}

func (f *listenerFixture) handle(t *testing.T, d amqp.Delivery) bool {
	done := make(chan bool, 1)
	f.listener.incomingMessageHandler(d, done)
	select {
	case ack := <-done:
		return ack
	default:
		t.Fatal("handler did not report a result")
		return false
	}
}

func newDelivery(t *testing.T, target, method, path string, headers amqp.Table) amqp.Delivery {
	msg, err := anypb.New(&epb.NodeFeederMessage{Target: target, HttpMethod: method, Path: path, Msg: []byte(`{}`)})
	require.NoError(t, err)
	body, err := proto.Marshal(msg)
	require.NoError(t, err)
	return amqp.Delivery{Body: body, Headers: headers}
}

func retryOf(attempt int, firstSeen time.Time) amqp.Table {
	return amqp.Table{retryAttemptHeaderName: int64(attempt), retryFirstSeenHeaderName: firstSeen.UnixMilli()}
}

// encodable reports whether the publisher can send these headers.
func encodable(h amqp.Table) bool {
	return amqp091.Table(h).Validate() == nil
}

func TestHandler_SuccessIsAcked(t *testing.T) {
	f := newFixture(t)
	f.exec.On("Execute", mock.Anything).Return(nil).Once()

	assert.True(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", nil)))
	f.pub.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestHandler_FirstFailureRetriesAfterFirstStep(t *testing.T) {
	f := newFixture(t)
	f.exec.On("Execute", mock.Anything).Return(errUnavailable).Once()
	f.pub.On("Publish", "node-feeder.retry.2s", mock.Anything, mock.MatchedBy(func(h amqp.Table) bool {
		return h[retryAttemptHeaderName] == int64(1) && h[retryFirstSeenHeaderName] == f.now.UnixMilli() &&
			h["df-target"] == testTarget && encodable(h)
	})).Return(nil).Once()

	assert.True(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", amqp.Table{"df-target": testTarget})))
}

func TestHandler_LaterFailureUsesNextStepAndKeepsFirstSeen(t *testing.T) {
	f := newFixture(t)
	firstSeen := f.now.Add(-8 * time.Second)
	f.exec.On("Execute", mock.Anything).Return(errUnavailable).Once()
	f.pub.On("Publish", "node-feeder.retry.15s", mock.Anything, mock.MatchedBy(func(h amqp.Table) bool {
		_, hasDeath := h["x-death"]
		return h[retryAttemptHeaderName] == int64(3) && h[retryFirstSeenHeaderName] == firstSeen.UnixMilli() &&
			h["df-target"] == testTarget && !hasDeath && h["x-first-death-queue"] == nil && encodable(h)
	})).Return(nil).Once()

	// A retry comes back from its TTL queue with the broker's dead-letter headers.
	headers := retryOf(2, firstSeen)
	headers["df-target"] = testTarget
	headers["x-death"] = []interface{}{amqp.Table{"count": int64(1), "queue": "node-feeder.retry.5s", "reason": "expired"}}
	headers["x-first-death-queue"] = "node-feeder.retry.2s"

	assert.True(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", headers)))
}

func TestHandler_RetiresWhenRetriesExhausted(t *testing.T) {
	f := newFixture(t)
	f.listener.maxRequestAge = 0
	f.exec.On("Execute", mock.Anything).Return(errUnavailable).Once()

	assert.True(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", retryOf(4, f.now.Add(-50*time.Second)))))
	f.pub.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestHandler_RetiresWhenNextRetryPassesMaxAge(t *testing.T) {
	f := newFixture(t)
	f.exec.On("Execute", mock.Anything).Return(errUnavailable).Once()

	// 40s old and the next step is 30s, so the retry would land at 70s.
	assert.True(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", retryOf(3, f.now.Add(-40*time.Second)))))
	f.pub.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestHandler_DropsRequestOlderThanMaxAge(t *testing.T) {
	f := newFixture(t)

	assert.True(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", retryOf(4, f.now.Add(-61*time.Second)))))
	f.exec.AssertNotCalled(t, "Execute", mock.Anything)
}

func TestHandler_DropsRetryReplacedByNewerRequest(t *testing.T) {
	f := newFixture(t)
	older := f.now
	f.exec.On("Execute", mock.Anything).Return(errUnavailable).Once()
	f.pub.On("Publish", "node-feeder.retry.2s", mock.Anything, mock.Anything).Return(nil).Once()
	assert.True(t, f.handle(t, newDelivery(t, testTarget, "DELETE", "configd/v1/config", nil)))

	f.now = f.now.Add(time.Second)
	f.exec.On("Execute", mock.Anything).Return(nil).Once()
	assert.True(t, f.handle(t, newDelivery(t, testTarget, "DELETE", "/configd/v1/config", nil)))

	f.now = f.now.Add(time.Second)
	assert.True(t, f.handle(t, newDelivery(t, testTarget, "DELETE", "configd/v1/config", retryOf(1, older))))
	f.exec.AssertNumberOfCalls(t, "Execute", 2)
}

func TestHandler_KeepsRetryOfNewestRequest(t *testing.T) {
	f := newFixture(t)
	firstSeen := f.now
	f.exec.On("Execute", mock.Anything).Return(errUnavailable).Once()
	f.pub.On("Publish", "node-feeder.retry.2s", mock.Anything, mock.Anything).Return(nil).Once()
	assert.True(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", nil)))

	f.now = f.now.Add(2 * time.Second)
	f.exec.On("Execute", mock.Anything).Return(nil).Once()
	assert.True(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", retryOf(1, firstSeen))))
	f.exec.AssertNumberOfCalls(t, "Execute", 2)
}

func TestHandler_OtherEndpointDoesNotSupersede(t *testing.T) {
	f := newFixture(t)
	firstSeen := f.now
	f.listener.latest.record(requestKey(&epb.NodeFeederMessage{Target: testTarget, HttpMethod: "PUT", Path: "v1/service"}), f.now.Add(time.Second))
	f.listener.latest.record(requestKey(&epb.NodeFeederMessage{Target: testTarget, HttpMethod: "DELETE", Path: "configd/v1/config"}), f.now.Add(time.Second))

	f.now = f.now.Add(2 * time.Second)
	f.exec.On("Execute", mock.Anything).Return(nil).Once()
	assert.True(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", retryOf(1, firstSeen))))
}

func TestHandler_PublishFailureFallsBackToNack(t *testing.T) {
	f := newFixture(t)
	f.exec.On("Execute", mock.Anything).Return(errUnavailable).Once()
	f.pub.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("channel closed")).Once()

	assert.False(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", nil)))
}

func TestHandler_NoBackoffUsesWaitingQueue(t *testing.T) {
	f := newFixture(t)
	f.listener.retryBackoff = nil
	f.exec.On("Execute", mock.Anything).Return(errUnavailable).Once()

	assert.False(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", nil)))
}

func TestHandler_Device4xxIsNotRetried(t *testing.T) {
	f := newFixture(t)
	f.exec.On("Execute", mock.Anything).Return(Device4xxServerError{errors.New("409")}).Once()

	assert.True(t, f.handle(t, newDelivery(t, testTarget, "POST", "configd/v1/config", nil)))
	f.pub.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)
}

func TestHandler_WildcardFailureIsRetried(t *testing.T) {
	f := newFixture(t)
	f.mult.On("Process", mock.Anything).Return(errUnavailable).Once()
	f.pub.On("Publish", "node-feeder.retry.2s", mock.Anything, mock.Anything).Return(nil).Once()

	assert.True(t, f.handle(t, newDelivery(t, "ukama.org.node.*", "PUT", "v1/service", nil)))
	assert.Empty(t, f.listener.latest.seen)
}

func TestRetryState_ReadsBrokerIntegerTypes(t *testing.T) {
	received := time.UnixMilli(5_000)

	attempt, firstSeen := retryState(amqp.Table{retryAttemptHeaderName: int32(2), retryFirstSeenHeaderName: int64(1_000)}, received)
	assert.Equal(t, 2, attempt)
	assert.Equal(t, time.UnixMilli(1_000), firstSeen)

	attempt, firstSeen = retryState(amqp.Table{"x-death": []interface{}{}}, received)
	assert.Equal(t, 0, attempt)
	assert.Equal(t, received, firstSeen)
}

func TestLatestRequests_PrunesOldEntries(t *testing.T) {
	l := newLatestRequests(time.Minute)
	start := time.UnixMilli(1_000_000)
	l.record("old", start)
	l.record("new", start.Add(2*time.Minute))

	assert.NotContains(t, l.seen, "old")
	assert.Contains(t, l.seen, "new")
}
