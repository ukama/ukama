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

// Header keys used by the managed dead-letter retry (see RetryPolicy).
const (
	retryCountHeader = "x-retry-count"
	origRKHeader     = "x-original-routing-key"
)

// RetryPolicy configures a managed dead-letter retry for a consumer client
// (see NewConsumerClientWithRetry). When set, a handler failure does not drop
// the message: it is republished to a per-queue delay queue ("<queue>.retry",
// with a TTL) up to MaxAttempts times, then parked in "<queue>.parking" for
// inspection. This needs no change to the main queue, so existing durable
// queues do not have to be recreated.
//
// Consumers created without a policy (NewConsumerClient) keep the legacy
// behaviour (ack on completion, requeue only on timeout) and manage their own
// retry/dead-lettering — e.g. node-feeder.
type RetryPolicy struct {
	MaxAttempts int           // total delivery attempts before parking (e.g. 3)
	Delay       time.Duration // delay between attempts (retry-queue TTL, e.g. 30s)
}

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

	// retry, when non-nil, enables the managed dead-letter retry for this
	// client's subscriptions. wPub is the publisher used to route failed
	// messages to the retry/parking queues.
	retry *RetryPolicy
	wPub  *rabbitmq.Publisher
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

// NewConsumerClientWithRetry creates a consumer whose subscriptions use a
// managed dead-letter retry (see RetryPolicy): failed messages are retried with
// a delay up to policy.MaxAttempts times and then parked, instead of being
// dropped. Intended for the msgClient sidecar, which has no dead-letter queue
// of its own.
func NewConsumerClientWithRetry(connectionString string, policy RetryPolicy) (Consumer, error) {
	p := policy
	return &MsgClient{
		connectionString: connectionString,
		log:              logrus.WithField("prefix", ""),
		retry:            &p,
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

	// Managed retry: declare the per-queue delay + parking topology and a
	// publisher used to route failed messages there.
	if m.retry != nil {
		if err := m.declareRetryTopology(sub.queueName); err != nil {
			return err
		}
		if err := m.ensurePublisher(); err != nil {
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

// wagslaneHandler bridges a wagslane delivery to the existing handlerFunc
// signature (streadway amqp.Delivery + done channel) and maps the outcome to a
// wagslane Action. When autoAck is set the returned Action is ignored.
//
// Behaviour on failure depends on whether this client has a RetryPolicy:
//   - no policy (raw, e.g. node-feeder): preserve the legacy semantics — ack on
//     handler completion (success or failure; the caller manages its own
//     dead-lettering), requeue only on timeout.
//   - policy set (msgClient): route the failed message through the managed
//     delay/parking topology (see retryOrPark), then ack the original.
func (m *MsgClient) wagslaneHandler(sub *subscription) rabbitmq.Handler {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		done := make(chan bool, 1)
		sub.handlerFunc(toStreadwayDelivery(d.Delivery), done)

		if m.retry == nil {
			// Legacy/raw path: ack on completion, requeue on timeout.
			select {
			case <-done:
				return rabbitmq.Ack
			case <-time.After(handlerAckTimeout):
				return rabbitmq.NackRequeue
			}
		}

		// Managed-retry path.
		select {
		case ok := <-done:
			if ok {
				return rabbitmq.Ack
			}
			return m.retryOrPark(sub, d, "handler reported failure")
		case <-time.After(handlerAckTimeout):
			return m.retryOrPark(sub, d, "handler timed out")
		}
	}
}

// retryOrPark republishes a failed message to the queue's delay queue for a
// later retry, or to its parking queue once MaxAttempts is reached, then acks
// the original so it is not left unacked. The original routing key is preserved
// in a header because the delay queue dead-letters via the default exchange.
func (m *MsgClient) retryOrPark(sub *subscription, d rabbitmq.Delivery, reason string) rabbitmq.Action {
	attempt := headerInt(d.Headers, retryCountHeader)

	origRK := d.RoutingKey
	if v, ok := d.Headers[origRKHeader].(string); ok && v != "" {
		origRK = v
	}

	headers := amqp091.Table{
		retryCountHeader: int32(attempt + 1),
		origRKHeader:     origRK,
	}

	if attempt+1 >= m.retry.MaxAttempts {
		parking := sub.queueName + ".parking"
		m.log.Errorf("[msgbus] %s for queue %q key %q; exhausted %d attempts, parking in %q",
			reason, sub.queueName, origRK, m.retry.MaxAttempts, parking)
		if err := m.republish(parking, d.Body, headers); err != nil {
			// Could not park: requeue so the message is not lost.
			m.log.Errorf("[msgbus] failed to park message for queue %q: %s; requeueing", sub.queueName, err)
			return rabbitmq.NackRequeue
		}
		return rabbitmq.Ack
	}

	retryQ := sub.queueName + ".retry"
	m.log.Warnf("[msgbus] %s for queue %q key %q; scheduling retry %d/%d via %q",
		reason, sub.queueName, origRK, attempt+1, m.retry.MaxAttempts, retryQ)
	if err := m.republish(retryQ, d.Body, headers); err != nil {
		m.log.Errorf("[msgbus] failed to schedule retry for queue %q: %s; requeueing", sub.queueName, err)
		return rabbitmq.NackRequeue
	}
	return rabbitmq.Ack
}

// republish sends body to a queue via the default exchange (routing key == queue
// name), carrying the given headers.
func (m *MsgClient) republish(queue string, body []byte, headers amqp091.Table) error {
	return m.wPub.Publish(
		body,
		[]string{queue},
		rabbitmq.WithPublishOptionsExchange(""),
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsHeaders(rabbitmq.Table(headers)),
	)
}

// declareRetryTopology declares the delay and parking queues for a managed-retry
// main queue. The delay queue holds messages for RetryPolicy.Delay, then
// dead-letters them (via the default exchange) straight back to the main queue.
// The main queue itself is not touched, so existing durable queues are unchanged.
func (m *MsgClient) declareRetryTopology(mainQueue string) error {
	conn, err := amqp091.Dial(m.connectionString)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func() { _ = ch.Close() }()

	// Delay queue: TTL then dead-letter back to the main queue.
	_, err = ch.QueueDeclare(mainQueue+".retry", true, false, false, false, amqp091.Table{
		"x-message-ttl":             m.retry.Delay.Milliseconds(),
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": mainQueue,
	})
	if err != nil {
		return errors.Wrap(err, "failed to declare retry queue")
	}

	// Parking queue: terminal failures land here for inspection.
	_, err = ch.QueueDeclare(mainQueue+".parking", true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "failed to declare parking queue")
	}

	return nil
}

// ensurePublisher lazily creates the wagslane publisher used to route failed
// messages to the retry/parking queues.
func (m *MsgClient) ensurePublisher() error {
	m.connMu.Lock()
	defer m.connMu.Unlock()

	if m.wPub != nil {
		return nil
	}

	pub, err := rabbitmq.NewPublisher(m.wConn, rabbitmq.WithPublisherOptionsLogging)
	if err != nil {
		return err
	}
	m.wPub = pub
	return nil
}

func headerInt(h amqp091.Table, key string) int {
	if h == nil {
		return 0
	}
	switch v := h[key].(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	default:
		return 0
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
	if m.wPub != nil {
		m.wPub.Close()
		m.wPub = nil
	}
	if m.wConn != nil {
		_ = m.wConn.Close()
		m.wConn = nil
	}
	// streadway connection (only present for publisher clients).
	if m.conn != nil && !m.conn.IsClosed() {
		_ = m.conn.Close()
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
