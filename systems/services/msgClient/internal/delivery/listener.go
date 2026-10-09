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
	"fmt"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	pb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

// Listener delivers a service's queue to it one event at a time.
type Listener struct {
	service       string
	timeout       time.Duration
	conn          *grpc.ClientConn
	client        pb.EventNotificationServiceClient
	consumer      *rabbitmq.Consumer
	unreachable   atomic.Bool
	onUnreachable func()
}

func Start(conn *rabbitmq.Conn, service, uri string, timeout time.Duration, onUnreachable func()) (*Listener, error) {
	gc, err := grpc.NewClient(uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc client for %s at %s: %w", service, uri, err)
	}

	consumer, err := rabbitmq.NewConsumer(conn, service,
		rabbitmq.WithConsumerOptionsQueueNoDeclare,
		rabbitmq.WithConsumerOptionsQOSPrefetch(1),
		rabbitmq.WithConsumerOptionsLogging,
	)
	if err != nil {
		_ = gc.Close()
		return nil, fmt.Errorf("consumer for %s: %w", service, err)
	}

	l := &Listener{
		service:       service,
		timeout:       timeout,
		conn:          gc,
		client:        pb.NewEventNotificationServiceClient(gc),
		consumer:      consumer,
		onUnreachable: onUnreachable,
	}

	go func() {
		if err := consumer.Run(l.handle); err != nil {
			log.Errorf("[%s] Listener stopped. Error %s", service, err.Error())
		}
	}()

	log.Infof("[%s] Listener started", service)
	return l, nil
}

// Stop waits for the event in flight, if any, then stops consuming.
func (l *Listener) Stop() {
	l.consumer.Close()
	_ = l.conn.Close()
	log.Infof("[%s] Listener stopped", l.service)
}

func (l *Listener) handle(d rabbitmq.Delivery) rabbitmq.Action {
	if l.unreachable.Load() {
		return rabbitmq.NackRequeue
	}

	msg := new(anypb.Any)
	if err := proto.Unmarshal(d.Body, msg); err != nil {
		log.Errorf("[%s] Dropping unparseable event %s. Error %s", l.service, d.RoutingKey, err.Error())
		return rabbitmq.Ack
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	_, err := l.client.EventNotification(ctx, &pb.Event{RoutingKey: d.RoutingKey, Msg: msg})

	switch status.Code(err) {
	case codes.OK:
		return rabbitmq.Ack
	case codes.Unavailable, codes.DeadlineExceeded:
		log.Warnf("[%s] Event %s not delivered, pausing delivery. Error %s", l.service, d.RoutingKey, err.Error())
		if l.unreachable.CompareAndSwap(false, true) {
			go l.onUnreachable()
		}
		return rabbitmq.NackRequeue
	default:
		log.Warnf("[%s] Event %s delivered, listener returned error %s", l.service, d.RoutingKey, err.Error())
		return rabbitmq.Ack
	}
}
