/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package msgbus

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ukama/ukama/systems/common/errors"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const CONNECTION_NOT_INIT_ERR_MSG = "Connection is not initialized"

// Backoff bounds used when a consumer channel/connection dies and we
// try to re-establish the subscription.
const (
	consumerMinBackoff = 1 * time.Second
	consumerMaxBackoff = 30 * time.Second
)

// defaultHeartbeat is the AMQP heartbeat interval negotiated with the broker.
// A dead TCP socket (broker restart, LB/NAT idle timeout, network partition)
// is detected after ~2 missed heartbeats, which then triggers NotifyClose and
// the consumer reconnect logic. 10s matches RabbitMQ's own default, so this is
// a safeguard (e.g. against a broker proposing 0/disabled) rather than a
// behaviour change.
const defaultHeartbeat = 10 * time.Second

// ConsumerPrefetchCount, when > 0, sets a per-consumer channel QoS prefetch so
// the broker will not deliver more than this many unacked messages at once
// (backpressure). It defaults to 0 (unlimited) to preserve existing throughput
// behaviour across all services; set it at process start-up to opt in.
var ConsumerPrefetchCount = 0

// Defines our interface for connecting and consuming messages.
// Consider using github.com/wagslane/go-rabbitmq instead. It provides similar functionality.
type IMsgBus interface {
	ConnectToBroker(connectionString string)
	Publish(body []byte, queueName string, exchangeName string, routingKey RoutingKey, exchangeType string) error
	PublishOnQueue(body []byte, queueName string, initQueue bool) error
	Subscribe(queueName string, exchangeName string, exchangeType string, routingKeys []RoutingKey, consumerName string, handlerFunc func(amqp.Delivery, chan<- bool)) error
	SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery, chan<- bool)) error
	Close()
}

type Publisher interface {
	Publish(body []byte, queueName string, exchangeName string, routingKey RoutingKey, exchangeType string) error
	PublishOnQueue(msg []byte, queueName string, initQueue bool) error
	PublishOnExchange(exchange string, routingKey string, body interface{}) error
	DeclareQueue(queueName string, durable bool) (*amqp.Queue, error)
	IsClosed() bool
	Close()
}

type Consumer interface {
	Subscribe(queueName string, exchangeName string, exchangeType string, routingKeys []RoutingKey, consumerName string, handlerFunc func(amqp.Delivery, chan<- bool)) error
	SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery, chan<- bool)) error
	SubscribeToServiceQueue(serviceName string, exchangeName string, routingKeys []RoutingKey, consumerId string, handlerFunc func(amqp.Delivery, chan<- bool)) error
	SubscribeWithArgs(queueName string, exchangeName string, exchangeType string,
		routingKeys []RoutingKey, consumerName string, queueArgs map[string]interface{}, handlerFunc func(amqp.Delivery, chan<- bool)) error
	IsClosed() bool
	Close()
}

// Real implementation, encapsulates a pointer to an amqp.Connection.
// The consumer path is reconnect-aware: if the underlying channel or
// connection dies, active subscriptions are automatically re-established
// (see subscribe/superviseConsumer). The publisher path (m.channel) is
// unchanged.
type MsgClient struct {
	conn             *amqp.Connection
	log              *logrus.Entry
	channel          *amqp.Channel
	connectionString string
	connMu           sync.Mutex // guards conn re-dial + consumer channel creation
	mu               sync.Mutex // guards closed
	closed           bool       // set once Close() is called; stops reconnection
}

// subscription captures everything required to (re)establish a consumer so
// it can be replayed after a connection/channel drop.
type subscription struct {
	queueName       string
	exchangeName    string
	exchangeType    string
	declareExchange bool
	durableQueue    bool
	autoAck         bool
	routingKeys     []RoutingKey
	consumerName    string
	queueArgs       map[string]interface{}
	handlerFunc     func(amqp.Delivery, chan<- bool)
}

// Servcie Config
type Config struct {
}

// Queue Config
type MsgBusQConfig struct {
	Exchange         string
	Queue            string
	ExchangeType     string
	ReqRountingKeys  []RoutingKey
	RespRountingKeys []RoutingKey
}

type RPCResponse struct {
	Status     bool
	Resp       *amqp.Delivery
	RoutingKey RoutingKey
}

// creates a message consumer and initializes connection

func NewConsumerClient(connectionString string) (Consumer, error) {
	return createClient(connectionString)
}

// NewPublisherClient creates a publisher and opens connection and channel
// Use one publisher per thread as it's common practice to use one channel per thread
func NewPublisherClient(connectionString string) (Publisher, error) {
	return createClient(connectionString)
}

func createClient(connectionString string) (*MsgClient, error) {
	conn, err := connectClient(connectionString)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	client := &MsgClient{
		conn:             conn,
		channel:          channel,
		log:              logrus.WithField("prefix", ""),
		connectionString: connectionString,
	}

	return client, nil
}

func RemovePassFromConnection(connectioStr string) string {
	return connectioStr[strings.LastIndex(connectioStr, "@"):]
}

func connectClient(connectionString string) (*amqp.Connection, error) {
	conn, err := amqp.DialConfig(fmt.Sprintf("%s/", connectionString), amqp.Config{
		Heartbeat: defaultHeartbeat,
		Locale:    "en_US",
	})
	if err != nil {
		logrus.Errorf("Trying to connect to AMQP compatible broker at: %s", RemovePassFromConnection(connectionString))

		return nil, err
	}

	return conn, nil
}

// Connect to Broker(RabbitMq server)
func (m *MsgClient) ConnectToBroker(connectionString string) {
	if connectionString == "" {
		panic("Cannot initialize connection to broker, connectionString not set.")
	}

	m.connectionString = connectionString

	conn := false
	for !conn {
		c, err := connectClient(connectionString)
		if err != nil {
			m.log.Infof("could not establish connection. Waiting for 5 seconds to re-connect")
			time.Sleep(5 * time.Second)
		} else {
			m.conn = c
			conn = true
		}
	}
}

// Publish to queue through exchange
func (m *MsgClient) Publish(body []byte, queueName string, exchangeName string, routingKey RoutingKey, exchangeType string) error {

	err := m.declareExchange(m.channel, exchangeName, exchangeType)
	if err != nil {
		return err
	}

	queue, err := m.declareQueue(m.channel, queueName, false, nil)
	if err != nil {
		return err
	}

	err = m.bindQueue(m.channel, queue.Name, routingKey, exchangeName)
	if err != nil {
		return err
	}

	// Publishes a message onto the queue.
	err = m.channel.Publish(
		exchangeName,       // exchange
		string(routingKey), // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			Body: body, // Our JSON body as []byte
		})

	if err != nil {
		m.log.Errorf("Err: %s .Failed to publish message to exchange.", err)
	} else {
		m.log.Debugf("Message was sent on Exchange %s Queue %s Routing Key %s ", exchangeName, queue.Name, string(routingKey))
	}
	return err
}

// Publish to Queue.
func (m *MsgClient) PublishOnQueue(body []byte, queueName string, initQueue bool) error {
	if initQueue {
		_, err := m.declareQueue(m.channel, queueName, false, nil)
		if err != nil {
			return errors.Wrap(err, "error declaring queue")
		}
	}

	// Publishes a message onto the queue.
	err := m.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body, // Our JSON body as []byte
		})
	if err != nil {
		m.log.Errorf("Err: %s Failed to publish message to queue.", err)
	} else {
		m.log.Debugf("Message was sent on Queue %s ", queueName)
	}
	return err
}

// PublishOnExchange publishes event to an exchange
// body - an object that is marshalled to json
func (m *MsgClient) PublishOnExchange(exchange string, routingKey string, body interface{}) error {
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "error marshalling the body")
	}

	// Publishes a message onto the queue.
	err = m.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyJson, // Our JSON body as []byte
		})
	if err != nil {
		return errors.Wrap(err, "failed to publish message to queue")
	}

	m.log.Debugf("Message was sent to exchange %s ", exchange)
	return nil
}

// Subscribe to exchange with option to listen to particular type of message
func (m *MsgClient) Subscribe(queueName string, exchangeName string, exchangeType string, routingKeys []RoutingKey, consumerName string, handlerFunc func(amqp.Delivery, chan<- bool)) error {
	return m.SubscribeWithArgs(queueName, exchangeName, exchangeType, routingKeys, consumerName, nil, handlerFunc)
}

func (m *MsgClient) SubscribeWithArgs(queueName string, exchangeName string, exchangeType string,
	routingKeys []RoutingKey, consumerName string, queueArgs map[string]interface{}, handlerFunc func(amqp.Delivery, chan<- bool)) error {
	return m.subscribe(&subscription{
		queueName:       queueName,
		exchangeName:    exchangeName,
		exchangeType:    exchangeType,
		declareExchange: true,
		durableQueue:    false,
		autoAck:         false,
		routingKeys:     routingKeys,
		consumerName:    consumerName,
		queueArgs:       queueArgs,
		handlerFunc:     handlerFunc,
	})
}

// subscribe establishes the consumer once synchronously (so the caller still
// receives immediate setup errors, preserving previous behaviour) and then
// hands the subscription to a supervisor goroutine that re-establishes it if
// the channel or connection is ever lost.
func (m *MsgClient) subscribe(sub *subscription) error {
	closeCh, err := m.establishConsumer(sub)
	if err != nil {
		return err
	}

	go m.superviseConsumer(sub, closeCh)
	return nil
}

// newConsumerChannel returns a fresh channel, re-dialling the connection first
// if it has been lost. It is safe for concurrent use across subscriptions.
func (m *MsgClient) newConsumerChannel() (*amqp.Channel, error) {
	if m.isIntentionallyClosed() {
		return nil, fmt.Errorf("client is closed; not reconnecting consumer")
	}

	m.connMu.Lock()
	defer m.connMu.Unlock()

	if m.conn == nil || m.conn.IsClosed() {
		if m.connectionString == "" {
			return nil, fmt.Errorf("connection string not set; cannot (re)connect consumer")
		}

		c, err := connectClient(m.connectionString)
		if err != nil {
			return nil, err
		}
		m.conn = c
	}

	return m.conn.Channel()
}

// establishConsumer (re)declares the queue, binds routes, starts consuming and
// returns the channel's close-notification so the supervisor can react to drops.
func (m *MsgClient) establishConsumer(sub *subscription) (<-chan *amqp.Error, error) {
	ch, err := m.newConsumerChannel()
	if err != nil {
		return nil, err
	}

	if sub.declareExchange {
		if err := m.declareExchange(ch, sub.exchangeName, sub.exchangeType); err != nil {
			_ = ch.Close()
			return nil, err
		}
	}

	queue, err := m.declareQueue(ch, sub.queueName, sub.durableQueue, sub.queueArgs)
	if err != nil {
		_ = ch.Close()
		return nil, err
	}

	m.log.Debugf("declared Queue (%d messages, %d consumers), binding to Exchange (key '%s')",
		queue.Messages, queue.Consumers, sub.exchangeName)

	for _, routingKey := range sub.routingKeys {
		if err := m.bindQueue(ch, queue.Name, routingKey, sub.exchangeName); err != nil {
			_ = ch.Close()
			return nil, err
		}
	}

	if ConsumerPrefetchCount > 0 {
		if err := ch.Qos(ConsumerPrefetchCount, 0, false); err != nil {
			_ = ch.Close()
			return nil, err
		}
	}

	msgs, err := m.consume(ch, queue.Name, sub.consumerName, sub.autoAck)
	if err != nil {
		_ = ch.Close()
		return nil, err
	}

	closeCh := ch.NotifyClose(make(chan *amqp.Error, 1))
	go m.consumeLoop(msgs, sub.handlerFunc)
	return closeCh, nil
}

// superviseConsumer blocks until the consumer channel/connection closes, then
// re-establishes the subscription with exponential backoff. It exits only when
// the client is intentionally closed.
func (m *MsgClient) superviseConsumer(sub *subscription, closeCh <-chan *amqp.Error) {
	backoff := consumerMinBackoff

	for {
		amqpErr := <-closeCh
		if m.isIntentionallyClosed() {
			return
		}

		m.log.Warnf("[msgbus] consumer %q on queue %q lost its channel (err: %v). Re-establishing.",
			sub.consumerName, sub.queueName, amqpErr)

		for {
			if m.isIntentionallyClosed() {
				return
			}

			time.Sleep(backoff)

			newCloseCh, err := m.establishConsumer(sub)
			if err != nil {
				backoff = nextBackoff(backoff, consumerMaxBackoff)
				m.log.Errorf("[msgbus] failed to re-subscribe consumer %q on queue %q: %s. Retrying in %s.",
					sub.consumerName, sub.queueName, err, backoff)
				continue
			}

			m.log.Infof("[msgbus] consumer %q re-subscribed on queue %q routes %v.",
				sub.consumerName, sub.queueName, sub.routingKeys)
			backoff = consumerMinBackoff
			closeCh = newCloseCh
			break
		}
	}
}

func (m *MsgClient) isIntentionallyClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}

func nextBackoff(cur, max time.Duration) time.Duration {
	next := cur * 2
	if next > max {
		return max
	}
	return next
}

func (m *MsgClient) createChannel() (*amqp.Channel, error) {
	if m.conn == nil {
		m.log.Errorln(CONNECTION_NOT_INIT_ERR_MSG)
		return nil, fmt.Errorf("connection not initialized")
	}

	// Get a channel from the connection
	ch, err := m.conn.Channel()
	if err != nil {
		m.log.Errorf("Err: %s Failed to connect to channel.", err.Error())
		return nil, err
	}
	return ch, nil
}

func (m *MsgClient) declareExchange(ch *amqp.Channel, exchangeName string, exchangeType string) error {
	m.log.Debugf("Channel %+v exchange name %s exchange type %s", ch, exchangeName, exchangeType)
	err := ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		m.log.Errorf("%s: %s", "Error creating an exchange", err.Error())
		return err
	}
	return nil
}

// SubscribeToServiceQueue creates a durable queue with a serviceName name and routes messages from an exchange
// If queue does not exist then it will be created with the `serviceName`
func (m *MsgClient) SubscribeToServiceQueue(serviceName string, exchangeName string, routingKeys []RoutingKey, consumerId string, handlerFunc func(amqp.Delivery, chan<- bool)) error {
	return m.subscribe(&subscription{
		queueName:       serviceName,
		exchangeName:    exchangeName,
		declareExchange: false,
		durableQueue:    true,
		autoAck:         false,
		routingKeys:     routingKeys,
		consumerName:    consumerId,
		handlerFunc:     handlerFunc,
	})
}

func (m *MsgClient) consume(ch *amqp.Channel, queueName string, consumerId string, autoAck bool) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(
		queueName,  // queue
		consumerId, // consumer
		autoAck,    // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		m.log.Errorf("%s: %s", "Failed to register a consumer", err)
		return nil, err
	}
	return msgs, nil
}

// Subscribe directly to queue
func (m *MsgClient) SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery, chan<- bool)) error {
	return m.subscribe(&subscription{
		queueName:       queueName,
		declareExchange: false,
		durableQueue:    false,
		autoAck:         true,
		routingKeys:     nil,
		consumerName:    consumerName,
		handlerFunc:     handlerFunc,
	})
}

// Close connection. Also signals any consumer supervisor goroutines to stop
// reconnecting.
func (m *MsgClient) Close() {
	m.mu.Lock()
	m.closed = true
	m.mu.Unlock()

	m.connMu.Lock()
	defer m.connMu.Unlock()
	if m.conn != nil && !m.conn.IsClosed() {
		m.conn.Close()
	}
}

func (m *MsgClient) IsClosed() bool {
	return m.conn.IsClosed()
}

// Read messages from Queue.
func (m *MsgClient) consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc func(d amqp.Delivery, ch chan<- bool)) {
	for d := range deliveries {

		// Invoke the handlerFunc func we passed as parameter.
		go m.handleTransit(d, handlerFunc)
	}
}

// This Go-Routine is transit between message consumer and handler.
func (m *MsgClient) handleTransit(msg amqp.Delivery, handlerFunc func(d amqp.Delivery, ch chan<- bool)) {

	//channel to sync
	done := make(chan bool, 1)

	// handler for incoming messages.
	handlerFunc(msg, done)

	// Ack response
	select {

	// Request processed but it may be success or failure
	case res := <-done:
		logrus.Debugf("Message %s acknowledged with result %v", msg.MessageId, res)
		m.sendAck(msg)

	case <-time.After(1 * time.Second):
		logrus.Errorf("Timeout while responding to request.")
		m.sendNack(msg)
	}

}

// Ack to send message
func (m *MsgClient) sendAck(msg amqp.Delivery) {
	if err := msg.Ack(false); err != nil {
		m.log.Errorf("Error acknowledging message [%+v]:: %s", msg, err)
	}
}

// Nack to handle negative messages
func (m *MsgClient) sendNack(msg amqp.Delivery) {
	if err := msg.Nack(true, true); err != nil {
		m.log.Errorf("Error acknowledging message [%+v]:: %s", msg, err)
	} else {
		m.log.Debugf("Acknowledged message [%+v]", msg)
	}
}

func (m *MsgClient) DeclareQueue(queueName string, durable bool) (*amqp.Queue, error) {
	queue, err := m.channel.QueueDeclare(
		queueName, // our queue name
		durable,   // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		m.log.Errorf("Err: %s Failed to declare queue.", err)
		return nil, err
	}
	return &queue, nil
}

func (m *MsgClient) declareQueue(ch *amqp.Channel, queueName string, durable bool, args map[string]interface{}) (*amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		queueName, // our queue name
		durable,   // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	if err != nil {
		m.log.Errorf("Err: %s Failed to declare queue.", err)
		return nil, err
	}
	return &queue, nil
}

func (m *MsgClient) bindQueue(ch *amqp.Channel, queueName string, routingKey RoutingKey, exchangeName string) error {
	err := ch.QueueBind(
		queueName,          // name of the queue
		string(routingKey), // bindingKey/routingkey
		exchangeName,       // sourceExchange
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		m.log.Errorf("Err: %s Failed to bind queue.", err)
		return err
	}
	return nil
}
