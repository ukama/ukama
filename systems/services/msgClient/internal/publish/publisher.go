/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package publish

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"

	mb "github.com/ukama/ukama/systems/common/msgbus"
)

var unroutable = promauto.NewCounter(prometheus.CounterOpts{
	Name: "msgclient_unroutable_events_total",
	Help: "Published events that matched no queue.",
})

type Publisher struct {
	pub     *rabbitmq.Publisher
	timeout time.Duration
}

func New(conn *rabbitmq.Conn, timeout time.Duration) (*Publisher, error) {
	pub, err := rabbitmq.NewPublisher(conn,
		rabbitmq.WithPublisherOptionsConfirm,
		rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		return nil, err
	}

	pub.NotifyReturn(func(r rabbitmq.Return) {
		unroutable.Inc()
		log.Warnf("Event %s from %v matched no queue", r.RoutingKey, r.Headers["source-service"])
	})

	return &Publisher{pub: pub, timeout: timeout}, nil
}

// Publish returns nil only once RabbitMQ has confirmed it stored the event.
func (p *Publisher) Publish(key string, body []byte, headers rabbitmq.Table) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	confirms, err := p.pub.PublishWithDeferredConfirmWithContext(ctx, body, []string{key},
		rabbitmq.WithPublishOptionsExchange(mb.DefaultExchange),
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsHeaders(headers),
	)
	if err != nil {
		return fmt.Errorf("publish %s: %w", key, err)
	}

	if len(confirms) == 0 || confirms[0] == nil {
		return errors.New("publish " + key + ": no confirmation from rabbitmq")
	}

	ok, err := confirms[0].WaitContext(ctx)
	if err != nil {
		return fmt.Errorf("publish %s: %w", key, err)
	}

	if !ok {
		return errors.New("publish " + key + ": rejected by rabbitmq")
	}

	return nil
}

func (p *Publisher) Close() {
	p.pub.Close()
}
