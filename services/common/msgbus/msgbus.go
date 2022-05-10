package msgbus

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const CONNECTION_NOT_INIT_ERR_MSG = "Connection is not initialized"

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

// Real implementation, encapsulates a pointer to an amqp.Connection
// Does not reconnect if connection is lost
type MsgClient struct {
	conn    *amqp.Connection
	log     *logrus.Entry
	channel *amqp.Channel
}

//Servcie Config
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

//Random integer generation
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// CorrelationID for the RPC message
func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
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
		conn:    conn,
		channel: channel,
		log:     logrus.WithField("prefix", ""),
	}

	return client, nil
}

func RemovePassFromConnection(connectioStr string) string {
	return connectioStr[strings.LastIndex(connectioStr, "@"):]
}

func connectClient(connectionString string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("%s/", connectionString))
	if err != nil {
		logrus.Errorf("Trying to connect to AMQP compatible broker at: " + RemovePassFromConnection(connectionString))
		return nil, err
	}

	return conn, nil
}

//Connect to Broker(RabbitMq server)
func (m *MsgClient) ConnectToBroker(connectionString string) {
	if connectionString == "" {
		panic("Cannot initialize connection to broker, connectionString not set.")
	}

	conn := false
	for !conn {
		c, err := connectClient(connectionString)
		if err != nil {
			m.log.Infof("could not establis connection. Waiting for 5 seconds to re-connect")
			time.Sleep(5 * time.Second)
		} else {
			m.conn = c
			conn = true
		}
	}
}

//Publish to queue through exchange
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
	ch, err := m.createChannel()
	if err != nil {
		return err
	}

	err = m.declareExchange(ch, exchangeName, exchangeType)
	if err != nil {
		return err
	}

	log.Printf("declared Exchange, declaring Queue (%s)", "")
	queue, err := m.declareQueue(ch, queueName, false, queueArgs)
	if err != nil {
		return err
	}

	m.log.Debugf("declared Queue (%d messages, %d consumers), binding to Exchange (key '%s')",
		queue.Messages, queue.Consumers, exchangeName)

	//Binding queue with exchange
	for _, routingKey := range routingKeys {
		err = m.bindQueue(ch, queue.Name, routingKey, exchangeName)
		if err != nil {
			return err
		}
	}

	msgs, err := m.consume(ch, queue.Name, consumerName, false)
	if err != nil {
		return err
	}

	go m.consumeLoop(msgs, handlerFunc)
	return nil
}

func (m *MsgClient) createChannel() (*amqp.Channel, error) {
	if m.conn == nil {
		m.log.Errorln(CONNECTION_NOT_INIT_ERR_MSG)
		return nil, fmt.Errorf("connection not initialized")
	}

	// Get a channel from the connection
	ch, err := m.conn.Channel()
	if err != nil {
		m.log.Errorf("Err: %s Failed to connect to channel.", err)
		return nil, err
	}
	return ch, nil
}

func (m *MsgClient) declareExchange(ch *amqp.Channel, exchangeName string, exchangeType string) error {
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
		m.log.Errorf("%s: %s", "Error creating an exchange", err)
		return err
	}
	return nil
}

// SubscribeToServiceQueue creates a durable queue with a serviceName name and routes messages from an exchange
// If queue does not exist then it will be created with the `serviceName`
func (m *MsgClient) SubscribeToServiceQueue(serviceName string, exchangeName string, routingKeys []RoutingKey, consumerId string, handlerFunc func(amqp.Delivery, chan<- bool)) error {
	ch, err := m.createChannel()
	if err != nil {
		return err
	}

	queue, err := m.declareQueue(ch, serviceName, true, nil)
	if err != nil {
		return err
	}

	m.log.Debugf("declared Queue (%d messages, %d consumers), binding to Exchange (key '%s')",
		queue.Messages, queue.Consumers, exchangeName)

	//Binding queue with exchange
	for _, routingKey := range routingKeys {
		err = m.bindQueue(ch, queue.Name, routingKey, exchangeName)
		if err != nil {
			return err
		}
	}

	msgs, err := m.consume(ch, queue.Name, consumerId, false)
	if err != nil {
		return err
	}

	go m.consumeLoop(msgs, handlerFunc)
	return nil
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

//Subscribe directly to queue
func (m *MsgClient) SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery, chan<- bool)) error {
	ch, err := m.createChannel()
	if err != nil {
		return err
	}

	queue, err := m.declareQueue(ch, queueName, false, nil)
	if err != nil {
		return err
	}

	msgs, err := m.consume(ch, queue.Name, consumerName, true)
	if err != nil {
		return err
	}

	go m.consumeLoop(msgs, handlerFunc)
	return nil
}

// Close connection
func (m *MsgClient) Close() {
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
		logrus.Debugf("Message %s acknowladged with result %v", msg.MessageId, res)
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

func (m *MsgClient) prepareQueueForPublishing(body []byte, queueName string, exchangeName string, routingKey RoutingKey, exchangeType string) (<-chan amqp.Delivery, string, error) {
	m.log.Debugf("Publishing RPC messages Queue Name %s Exchange Name %s Routing Key %s Exchange Type %s ", queueName, exchangeName, routingKey, exchangeType)
	ch, err := m.createChannel()
	if err != nil {
		return nil, "", err
	}
	defer ch.Close()

	err = m.declareExchange(ch, exchangeName, exchangeType)
	if err != nil {
		return nil, "", err
	}

	queue, err := m.declareQueue(ch, queueName, false, nil)
	if err != nil {
		return nil, "", err
	}

	err = m.bindQueue(ch, queue.Name, routingKey, exchangeName)
	if err != nil {
		return nil, "", err
	}

	const replyQueueName = "amq.rabbitmq.reply-to"
	resp, err := m.consume(ch, replyQueueName, "ReplyToRPCConsumer", true)
	if err != nil {
		return nil, "", err
	}

	corrId, err := m.publishMessage(body, ch, exchangeName, routingKey, queue, replyQueueName)
	if err != nil {
		return nil, "", err
	}
	return resp, corrId, nil
}

func (m *MsgClient) publishMessage(body []byte, ch *amqp.Channel, exchangeName string,
	routingKey RoutingKey, queue *amqp.Queue, replyTo string) (string, error) {
	corrId := randomString(32)

	err := ch.Publish( // Publishes a message onto the queue.
		exchangeName,       // exchange
		string(routingKey), // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       replyTo,
			Body:          body, // Our JSON body as []byte
		})

	if err != nil {
		m.log.Errorf("Err: %s Failed to publish message.", err)
		return "", err
	} else {
		m.log.Debugf("RPC Message was sent on Exchange %s Queue %s Routing Key %s", exchangeName, queue.Name, string(routingKey))
	}
	return corrId, nil
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
