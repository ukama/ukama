/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package pkg

import (
	"fmt"
	"strings"
	"sync"
	"time"

	amqp "github.com/streadway/amqp"
	"github.com/wagslane/go-rabbitmq"

	mb "github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const (
	retryAttemptHeaderName   = "x-nf-attempt"
	retryFirstSeenHeaderName = "x-nf-first-seen"
)

type RetryPublisher interface {
	Publish(queueName string, body []byte, headers amqp.Table) error
}

type queueRetryPublisher struct {
	publisher *rabbitmq.Publisher
}

// NewRetryPublisher publishes to the default exchange, which routes by queue name.
func NewRetryPublisher(queueUri string) (*queueRetryPublisher, error) {
	conn, err := rabbitmq.NewConn(queueUri, rabbitmq.WithConnectionOptionsLogging)
	if err != nil {
		return nil, err
	}

	publisher, err := rabbitmq.NewPublisher(conn, rabbitmq.WithPublisherOptionsLogging)
	if err != nil {
		return nil, err
	}

	return &queueRetryPublisher{publisher: publisher}, nil
}

func (p *queueRetryPublisher) Publish(queueName string, body []byte, headers amqp.Table) error {
	return p.publisher.Publish(body, []string{queueName},
		rabbitmq.WithPublishOptionsHeaders(rabbitmq.Table(headers)),
		rabbitmq.WithPublishOptionsPersistentDelivery)
}

func retryQueueName(delay time.Duration) string {
	return fmt.Sprintf("node-feeder.retry.%ds", int(delay/time.Second))
}

func retryQueueArgs(delay time.Duration) amqp.Table {
	return amqp.Table{
		"x-message-ttl":                int64(delay / time.Millisecond),
		deadLetterExchangeHeaderName:   mb.DefaultExchange,
		deadLetterRoutingKeyHeaderName: string(mb.NodeFeederRequestRoutingKey),
	}
}

func backoffDurations(seconds []int) []time.Duration {
	delays := make([]time.Duration, 0, len(seconds))
	for _, s := range seconds {
		if s > 0 {
			delays = append(delays, time.Duration(s)*time.Second)
		}
	}
	return delays
}

// retryState returns how many retries a delivery has had and when node-feeder
// first received it. A delivery without our headers is a first attempt.
func retryState(headers amqp.Table, receivedAt time.Time) (int, time.Time) {
	attempt, ok := headerInt(headers, retryAttemptHeaderName)
	if !ok {
		return 0, receivedAt
	}

	firstSeen, ok := headerInt(headers, retryFirstSeenHeaderName)
	if !ok {
		return int(attempt), receivedAt
	}

	return int(attempt), time.UnixMilli(firstSeen)
}

func headerInt(headers amqp.Table, name string) (int64, bool) {
	switch v := headers[name].(type) {
	case int64:
		return v, true
	case int32:
		return int64(v), true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}

func retryHeaders(headers amqp.Table, attempt int, firstSeen time.Time) amqp.Table {
	out := amqp.Table{}
	for k, v := range headers {
		if k == "x-death" || strings.HasPrefix(k, "x-first-death-") || strings.HasPrefix(k, "x-last-death-") {
			continue
		}
		out[k] = v
	}
	out[retryAttemptHeaderName] = int64(attempt)
	out[retryFirstSeenHeaderName] = firstSeen.UnixMilli()
	return out
}

func requestKey(request *epb.NodeFeederMessage) string {
	return request.Target + "|" + strings.ToUpper(request.HttpMethod) + "|" + strings.Trim(request.Path, "/")
}

// latestRequests remembers when the newest request for each node endpoint
// arrived, so a retry of an older request to the same endpoint can be dropped.
type latestRequests struct {
	mu        sync.Mutex
	seen      map[string]int64
	maxAge    time.Duration
	lastPrune time.Time
}

func newLatestRequests(maxAge time.Duration) *latestRequests {
	return &latestRequests{seen: map[string]int64{}, maxAge: maxAge}
}

func (l *latestRequests) record(key string, firstSeen time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ms := firstSeen.UnixMilli(); ms > l.seen[key] {
		l.seen[key] = ms
	}
	l.prune(firstSeen)
}

func (l *latestRequests) superseded(key string, firstSeen time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.seen[key] > firstSeen.UnixMilli()
}

// Entries older than maxAge cannot supersede anything, since retries that old are retired.
func (l *latestRequests) prune(now time.Time) {
	if l.maxAge <= 0 || now.Sub(l.lastPrune) < l.maxAge {
		return
	}
	cutoff := now.Add(-l.maxAge).UnixMilli()
	for key, ms := range l.seen {
		if ms < cutoff {
			delete(l.seen, key)
		}
	}
	l.lastPrune = now
}
