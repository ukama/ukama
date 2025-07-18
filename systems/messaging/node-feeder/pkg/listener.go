/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/messaging/node-feeder/pkg/global"
	"github.com/ukama/ukama/systems/messaging/node-feeder/pkg/metrics"

	log "github.com/sirupsen/logrus"
	amqp "github.com/streadway/amqp"
	mb "github.com/ukama/ukama/systems/common/msgbus"
	cpb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

const (
	deadLetterExchangeName         = "device-feeder.dead-letter"
	deadLetterExchangeHeaderName   = "x-dead-letter-exchange"
	errorCreatingWaitingQueueErr   = "error declaring waiting queue"
	deadLetterRoutingKeyHeaderName = "x-dead-letter-routing-key"
)

type QueueListener struct {
	service        string
	consumer       mb.Consumer
	serviceId      string
	requestMult    RequestMultiplier
	requestExec    RequestExecutor
	maxRetryCount  int64
	retryPeriodSec int
	listenerConfig ListenerConfig
}

type RequestMultiplier interface {
	Process(body *cpb.NodeFeederMessage) error
}

func NewQueueListener(service string, queueUri string, serviceId string, requestMult RequestMultiplier, requestExec RequestExecutor, conf ListenerConfig) (*QueueListener, error) {

	client, err := mb.NewConsumerClient(queueUri)
	if err != nil {
		return nil, err
	}
	q := &QueueListener{
		service:        service,
		consumer:       client,
		serviceId:      serviceId,
		requestMult:    requestMult,
		requestExec:    requestExec,
		listenerConfig: conf,
	}

	err = q.declareQueueTopology(queueUri)
	if err != nil {
		return nil, errors.Wrap(err, "error declaring queue topology")
	}

	return q, nil
}

func (q *QueueListener) declareQueueTopology(queueUri string) error {
	conn, err := amqp.Dial(queueUri)
	if err != nil {
		return errors.Wrap(err, "error connecting to queue")
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Warnf("failed to close connection: %v", err)
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, "error creating amqp channel")
	}

	defer func() {
		err := ch.Close()
		if err != nil {
			log.Warnf("failed to close connection: %v", err)
		}
	}()

	err = ch.ExchangeDeclare(deadLetterExchangeName, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "error declaring dead letter exchange")
	}

	waitingQueue, err := q.createWaitingQueue(ch)
	if err != nil {
		if err.(*amqp.Error).Code == amqp.PreconditionFailed {
			log.Warnf("Did you change waiting queue message TTL? If so then delete queue manually")
		}

		return errors.Wrap(err, errorCreatingWaitingQueueErr)
	}

	err = ch.QueueBind(waitingQueue.Name, "*", deadLetterExchangeName, false, nil)
	if err != nil {
		return errors.Wrap(err, "error binding waiting-queue")
	}

	// data feeder queue
	dataFeederQueue, err := ch.QueueDeclare(
		"node-feeder", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		map[string]interface{}{
			deadLetterExchangeHeaderName:   deadLetterExchangeName,
			deadLetterRoutingKeyHeaderName: string(mb.NodeFeederRequestRoutingKey),
		}, // arguments
	)
	if err != nil {
		return errors.Wrap(err, errorCreatingWaitingQueueErr)
	}
	err = ch.QueueBind(dataFeederQueue.Name, string(mb.NodeFeederRequestRoutingKey), mb.DefaultExchange, false, nil)

	if err != nil {
		return errors.Wrap(err, errorCreatingWaitingQueueErr)
	}

	return nil
}

func (q *QueueListener) createWaitingQueue(ch *amqp.Channel) (amqp.Queue, error) {

	// TODO: set TTL via policy
	waitingQueue, err := ch.QueueDeclare(
		"node-feeder.waiting-queue", // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		map[string]interface{}{
			"x-message-ttl":                q.retryPeriodSec * 1000,
			deadLetterExchangeHeaderName:   mb.DefaultExchange,
			deadLetterRoutingKeyHeaderName: string(mb.NodeFeederRequestRoutingKey),
		},
	)
	return waitingQueue, err
}

func (q *QueueListener) StartQueueListening() (err error) {

	err = q.consumer.SubscribeToServiceQueue("nodefeederService", q.listenerConfig.Exchange, q.listenerConfig.Routes, q.serviceId, q.incomingMessageHandler)
	if err != nil {
		log.Errorf("Error subscribing for queue messages. Error: %+v", err)
		return err
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	return nil
}

func (q *QueueListener) incomingMessageHandler(delivery amqp.Delivery, done chan<- bool) {
	if q.isRetryLimitReached(delivery) {
		metrics.RecordFailedRequestMetric()
		done <- true
		return
	}

	err := q.processRequest(delivery)
	if err != nil {
		log.Errorf("Error processing request. Error: %+v", err)
		metrics.RecordFailedRequestMetric()
		done <- false
		return
	}
	metrics.RecordSuccessfulRequestMetric()

	done <- true
}

func (q *QueueListener) isRetryLimitReached(delivery amqp.Delivery) bool {
	const deathHeader = "x-death"
	if delivery.Headers[deathHeader] == nil {
		return false
	}
	death := delivery.Headers[deathHeader].([]interface{})
	for _, d := range death {
		vals := d.(amqp.Table)
		if vals == nil {
			log.Errorf("Unexpected format of death header")
			return false
		}

		if vals["exchange"] == deadLetterExchangeName {
			count := vals["count"].(int64)
			if count > q.maxRetryCount {
				log.Infof("Retry limit reached for message: %v, target: %v", delivery.MessageId, delivery.Headers[global.OptionalTargetHeaderName])
				return true
			} else {
				log.Infof("Retry count: %v, target: %v", count, delivery.Headers[global.OptionalTargetHeaderName])
				return false
			}
		}
	}

	log.Warning("Cannot get retry count from message headers")
	return false
}

// process node msg also
func (q *QueueListener) processRequest(delivery amqp.Delivery) error {

	log.Infof("Raw message %v", delivery.Body)

	evtAny := new(anypb.Any)
	err := proto.Unmarshal(delivery.Body, evtAny)
	if err != nil {
		log.Errorf("Failed to parse message with key %s. Error %s", delivery.RoutingKey, err.Error())
		return nil
	}

	request := &cpb.NodeFeederMessage{}
	err = evtAny.UnmarshalTo(request)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		return nil
	}

	log.Infof("Received request: %+v", request)
	if strings.HasSuffix(request.Target, "*") {
		log.Infof("Wildcarded target: %s", request.Target)
		err = q.requestMult.Process(request)
		return err
	} else {
		log.Infof("Direct node target: %s", request.Target)
		err = q.requestExec.Execute(request)

		if err != nil {
			if dErr, ok := err.(Device4xxServerError); ok {
				log.Warningf("Request failed but won't be retried. Error %v", dErr)
				return nil
			}
		}

		return err
	}

}
