/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package msgbus

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type QPub interface {
	Publish(payload any, routingKey string) error
	PublishProto(ctx context.Context, payload proto.Message, routingKey string) error
	PublishToQueue(queueName string, payload any) error
	Close() error
}

// QPub is a simplified AMQP client that publishes messages to a default exchange
// Client reconnect in case of connection loss.
type qPub struct {
	conn        *rabbitmq.Conn
	publisher   *rabbitmq.Publisher
	serviceName string
	instanceId  string
}

func NewQPub(queueUri string, serviceName string, exchange string, instanceId string) (*qPub, error) {
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

	return &qPub{
		conn:        conn,
		publisher:   publisher,
		serviceName: serviceName,
		instanceId:  instanceId,
	}, nil
}

// Publish publishes a message in json format to the default topic exchange with a routing key specified
func (q *qPub) Publish(payload any, routingKey string) error {

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = q.publisher.Publish(b, []string{routingKey},
		rabbitmq.WithPublishOptionsHeaders(map[string]interface{}{
			"source-service": q.serviceName,
			"instance-id":    q.instanceId,
		}),
		rabbitmq.WithPublishOptionsExchange(DefaultExchange))

	if err != nil {
		return err
	}

	return nil
}

// PublishProto publishes a proto message on the default exchange. The
// caller's trace context is carried in the message headers (traceparent), so
// consumers continue the same trace.
func (q *qPub) PublishProto(ctx context.Context, payload proto.Message, routingKey string) error {
	ctx, span := otel.Tracer("ukama/msgbus").Start(ctx, "publish "+DefaultExchange,
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "rabbitmq"),
			attribute.String("messaging.destination.name", DefaultExchange),
			attribute.String("messaging.rabbitmq.destination.routing_key", routingKey),
		))
	defer span.End()

	b, err := proto.Marshal(payload)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	headers := HeaderCarrier{
		"source-service": q.serviceName,
		"instance-id":    q.instanceId,
	}
	otel.GetTextMapPropagator().Inject(ctx, headers)

	err = q.publisher.Publish(b, []string{routingKey},
		rabbitmq.WithPublishOptionsHeaders(map[string]interface{}(headers)),
		rabbitmq.WithPublishOptionsExchange(DefaultExchange))

	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (q *qPub) Close() error {
	q.conn.Close()
	return nil
}

func (q *qPub) PublishToQueue(queueName string, payload any) error {
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
