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

	amqp091 "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

const CONNECTION_NOT_INIT_ERR_MSG = "Connection is not initialized"

// defaultHeartbeat is the AMQP heartbeat interval negotiated with the broker.
// Heartbeats keep the TCP connection active (so k8s conntrack / LB idle timers
// don't silently drop it) and let the client detect a dead peer. The wagslane
// consumer below auto-reconnects and re-subscribes when that happens.
const defaultHeartbeat = 10 * time.Second

// handlerAckTimeout bounds how long we wait for a handler to signal completion
// on its done channel before requeueing the message.
const handlerAckTimeout = 30 * time.Second

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

// Real implementation. The consumer path uses the wagslane/go-rabbitmq
// transport (see subscribe/ensureConn), which auto-reconnects and re-subscribes
// on channel/connection loss or server-side consumer cancellation. The legacy
// publisher path (m.channel, streadway) is unchanged.
type MsgClient struct {
	conn             *amqp.Connection // streadway connection, used by the (legacy) publisher methods
	log              *logrus.Entry
	channel          *amqp.Channel
	connectionString string
	connMu           sync.Mutex // guards wConn creation
	mu               sync.Mutex // guards closed + consumers
	closed           bool       // set once Close() is called

	// Consumer transport. Uses github.com/wagslane/go-rabbitmq, which manages
	// heartbeats, dead-connection detection and automatic re-subscription of
	// consumers on reconnect (the same library the publisher/qpub uses).
	wConn     *rabbitmq.Conn
	consumers []*rabbitmq.Consumer
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
	// Consumer clients use the wagslane transport, created lazily on the first
	// Subscribe call. No eager streadway dial is needed here.
	return &MsgClient{
		connectionString: connectionString,
		log:              logrus.WithField("prefix", ""),
	}, nil
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

// subscribe creates a wagslane consumer for the subscription and runs it. The
// wagslane Conn transparently reconnects and re-subscribes the consumer on any
// channel/connection loss OR server-side consumer cancellation, so no custom
// reconnect supervisor is needed here.
func (m *MsgClient) subscribe(sub *subscription) error {
	conn, err := m.ensureConn()
	if err != nil {
		return err
	}

	// Declare the queue and its route bindings synchronously so any setup error
	// is returned to the caller. wagslane performs declare/bind inside
	// consumer.Run() on a background goroutine, where a failure is only logged
	// and the consumer silently receives nothing. Doing it here guarantees the
	// queue is bound (or fails loudly) before we start consuming. wagslane still
	// re-declares/re-binds idempotently on reconnect for self-healing.
	if len(sub.routingKeys) > 0 || sub.declareExchange {
		if err := m.declareTopology(sub); err != nil {
			return err
		}
	}

	opts := []func(*rabbitmq.ConsumerOptions){
		rabbitmq.WithConsumerOptionsConsumerName(sub.consumerName),
	}

	if sub.durableQueue {
		opts = append(opts, rabbitmq.WithConsumerOptionsQueueDurable)
	}

	if sub.exchangeName != "" {
		opts = append(opts, rabbitmq.WithConsumerOptionsExchangeName(sub.exchangeName))
	}

	// Only declare the exchange when the caller asked for it. SubscribeToServiceQueue
	// binds to the pre-existing amq.topic and must NOT redeclare it.
	if sub.declareExchange {
		opts = append(opts, rabbitmq.WithConsumerOptionsExchangeDeclare, rabbitmq.WithConsumerOptionsExchangeDurable)
		if sub.exchangeType != "" {
			opts = append(opts, rabbitmq.WithConsumerOptionsExchangeKind(sub.exchangeType))
		}
	}

	for _, rk := range sub.routingKeys {
		opts = append(opts, rabbitmq.WithConsumerOptionsRoutingKey(string(rk)))
	}

	if sub.queueArgs != nil {
		opts = append(opts, rabbitmq.WithConsumerOptionsQueueArgs(rabbitmq.Table(sub.queueArgs)))
	}

	if sub.autoAck {
		opts = append(opts, rabbitmq.WithConsumerOptionsConsumerAutoAck(true))
	}

	if ConsumerPrefetchCount > 0 {
		opts = append(opts, rabbitmq.WithConsumerOptionsQOSPrefetch(ConsumerPrefetchCount))
	}

	consumer, err := rabbitmq.NewConsumer(conn, sub.queueName, opts...)
	if err != nil {
		return err
	}

	m.mu.Lock()
	m.consumers = append(m.consumers, consumer)
	m.mu.Unlock()

	// Run blocks and keeps consuming across reconnects until the consumer is
	// closed, so it must run in its own goroutine.
	go func() {
		if err := consumer.Run(m.wagslaneHandler(sub)); err != nil && !m.isIntentionallyClosed() {
			m.log.Errorf("[msgbus] consumer %q on queue %q stopped: %s", sub.consumerName, sub.queueName, err)
		}
	}()

	m.log.Infof("[msgbus] consumer %q subscribed on queue %q routes %v", sub.consumerName, sub.queueName, sub.routingKeys)
	return nil
}

// ensureConn lazily creates the shared wagslane connection (heartbeat enabled,
// auto-reconnecting).
func (m *MsgClient) ensureConn() (*rabbitmq.Conn, error) {
	m.connMu.Lock()
	defer m.connMu.Unlock()

	if m.wConn != nil {
		return m.wConn, nil
	}

	if m.connectionString == "" {
		return nil, fmt.Errorf("connection string not set; cannot connect consumer")
	}

	conn, err := rabbitmq.NewConn(
		m.connectionString,
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsConfig(rabbitmq.Config{
			Heartbeat: defaultHeartbeat,
			Locale:    "en_US",
		}),
	)
	if err != nil {
		return nil, err
	}

	m.wConn = conn
	return conn, nil
}

// declareTopology synchronously declares the exchange (when requested), the
// queue, and the route bindings for a subscription, returning any error to the
// caller. This mirrors what wagslane does inside consumer.Run(), but done here
// the errors surface at subscribe time instead of being swallowed in the
// background goroutine (which previously left consumers bound to nothing).
func (m *MsgClient) declareTopology(sub *subscription) error {
	tconn, err := amqp091.Dial(m.connectionString)
	if err != nil {
		return fmt.Errorf("topology setup: dial failed: %w", err)
	}
	defer func() { _ = tconn.Close() }()

	ch, err := tconn.Channel()
	if err != nil {
		return fmt.Errorf("topology setup: channel failed: %w", err)
	}
	defer func() { _ = ch.Close() }()

	if sub.declareExchange {
		kind := sub.exchangeType
		if kind == "" {
			kind = "topic"
		}
		if err := ch.ExchangeDeclare(sub.exchangeName, kind, true, false, false, false, nil); err != nil {
			return fmt.Errorf("failed to declare exchange %q: %w", sub.exchangeName, err)
		}
	}

	var args amqp091.Table
	if sub.queueArgs != nil {
		args = amqp091.Table(sub.queueArgs)
	}
	if _, err := ch.QueueDeclare(sub.queueName, sub.durableQueue, false, false, false, args); err != nil {
		return fmt.Errorf("failed to declare queue %q: %w", sub.queueName, err)
	}

	for _, rk := range sub.routingKeys {
		if err := ch.QueueBind(sub.queueName, string(rk), sub.exchangeName, false, nil); err != nil {
			return fmt.Errorf("failed to bind queue %q to exchange %q with key %q: %w",
				sub.queueName, sub.exchangeName, rk, err)
		}
		m.log.Infof("[msgbus] bound queue %q to exchange %q with key %q", sub.queueName, sub.exchangeName, rk)
	}

	return nil
}

// wagslaneHandler bridges a wagslane delivery to the existing handlerFunc
// signature (streadway amqp.Delivery + done channel) and maps the outcome to a
// wagslane Action. When autoAck is set the returned Action is ignored.
func (m *MsgClient) wagslaneHandler(sub *subscription) rabbitmq.Handler {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		done := make(chan bool, 1)
		sub.handlerFunc(toStreadwayDelivery(d.Delivery), done)

		select {
		case ok := <-done:
			if ok {
				return rabbitmq.Ack
			}
			return rabbitmq.NackRequeue
		case <-time.After(handlerAckTimeout):
			m.log.Errorf("[msgbus] handler timed out for queue %q key %q; requeueing", sub.queueName, d.RoutingKey)
			return rabbitmq.NackRequeue
		}
	}
}

func (m *MsgClient) isIntentionallyClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}

// toStreadwayDelivery converts an amqp091 delivery to the streadway
// amqp.Delivery type used by existing handlerFunc callers. The returned value
// has no Acknowledger (acking is handled by the returned wagslane Action), so
// callers must not call Ack/Nack on it directly.
func toStreadwayDelivery(d amqp091.Delivery) amqp.Delivery {
	return amqp.Delivery{
		Headers:         convHeadersToStreadway(d.Headers),
		ContentType:     d.ContentType,
		ContentEncoding: d.ContentEncoding,
		DeliveryMode:    d.DeliveryMode,
		Priority:        d.Priority,
		CorrelationId:   d.CorrelationId,
		ReplyTo:         d.ReplyTo,
		Expiration:      d.Expiration,
		MessageId:       d.MessageId,
		Timestamp:       d.Timestamp,
		Type:            d.Type,
		UserId:          d.UserId,
		AppId:           d.AppId,
		ConsumerTag:     d.ConsumerTag,
		MessageCount:    d.MessageCount,
		DeliveryTag:     d.DeliveryTag,
		Redelivered:     d.Redelivered,
		Exchange:        d.Exchange,
		RoutingKey:      d.RoutingKey,
		Body:            d.Body,
	}
}

// convHeadersToStreadway deep-converts amqp091 header tables (including nested
// tables/arrays such as x-death) into streadway amqp.Table so that consumers
// which type-assert to amqp.Table (e.g. node-feeder retry counting) keep working.
func convHeadersToStreadway(in amqp091.Table) amqp.Table {
	if in == nil {
		return nil
	}
	out := make(amqp.Table, len(in))
	for k, v := range in {
		out[k] = convHeaderValue(v)
	}
	return out
}

func convHeaderValue(v interface{}) interface{} {
	switch t := v.(type) {
	case amqp091.Table:
		return convHeadersToStreadway(t)
	case []interface{}:
		s := make([]interface{}, len(t))
		for i, e := range t {
			s[i] = convHeaderValue(e)
		}
		return s
	default:
		return v
	}
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
	consumers := m.consumers
	m.consumers = nil
	m.mu.Unlock()

	// Stop wagslane consumers first, then the shared connection.
	for _, c := range consumers {
		c.Close()
	}

	m.connMu.Lock()
	defer m.connMu.Unlock()
	if m.wConn != nil {
		_ = m.wConn.Close()
		m.wConn = nil
	}
	// streadway connection (only present for publisher clients).
	if m.conn != nil && !m.conn.IsClosed() {
		m.conn.Close()
	}
}

func (m *MsgClient) IsClosed() bool {
	// Consumer clients have no streadway connection; report the intentional
	// close flag instead (the wagslane Conn reconnects on its own otherwise).
	if m.conn == nil {
		return m.isIntentionallyClosed()
	}
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
