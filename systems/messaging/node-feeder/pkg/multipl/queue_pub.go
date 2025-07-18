/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package multipl

import (
	"encoding/json"

	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/protobuf/proto"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/messaging/node-feeder/pkg/global"

	log "github.com/sirupsen/logrus"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

type queuePublisher struct {
	conn        *rabbitmq.Conn
	publisher   *rabbitmq.Publisher
	serviceName string
	instanceId  string
}

type QueuePublisher interface {
	Publish(msg *cpb.NodeFeederMessage) error
	PublishProto(payload proto.Message, routingKey string) error
	PublishToQueue(queueName string, payload any) error
	Close() error
}

func NewQPub(queueUri string, serviceName string, exchange string, instanceId string) (*queuePublisher, error) {
	conn, err := rabbitmq.NewConn(
		queueUri,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Infof("Error creating publisher %s.", err.Error())
		return nil, err
	}

	publisher, err := rabbitmq.NewPublisher(conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(exchange),
	)
	if err != nil {
		return nil, err
	}

	return &queuePublisher{
		conn:        conn,
		publisher:   publisher,
		serviceName: serviceName,
		instanceId:  instanceId,
	}, nil
}

func (q *queuePublisher) Publish(msg *cpb.NodeFeederMessage) error {

	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = q.publisher.Publish(b, []string{string(msgbus.NodeFeederRequestRoutingKey)},
		rabbitmq.WithPublishOptionsHeaders(map[string]interface{}{
			global.OptionalTargetHeaderName: msg.Target,
		}),
		rabbitmq.WithPublishOptionsExchange(msgbus.DefaultExchange))

	if err != nil {
		return err
	}

	return nil
}
func (q *queuePublisher) PublishToQueue(queueName string, payload any) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = q.publisher.Publish(b, []string{queueName},
		rabbitmq.WithPublishOptionsHeaders(map[string]interface{}{
			"source-service": q.serviceName,
			"instance-id":    q.instanceId,
		}))

	if err != nil {
		return err
	}

	return nil
}
func (q *queuePublisher) PublishProto(payload proto.Message, routingKey string) error {

	b, err := proto.Marshal(payload)
	if err != nil {
		return err
	}

	err = q.publisher.Publish(b, []string{routingKey},
		rabbitmq.WithPublishOptionsHeaders(map[string]interface{}{
			"source-service": q.serviceName,
			"instance-id":    q.instanceId,
		}),
		rabbitmq.WithPublishOptionsExchange(msgbus.DefaultExchange))

	if err != nil {
		return err
	}

	return nil
}

func (q *queuePublisher) Close() error {
	return q.conn.Close()
}
