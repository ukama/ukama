package msgbus

import (
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Defines our interface for connecting and consuming messages.
type IMsgBus interface {
	ConnectToBroker(connectionString string)
	GenerateToken() uint64
	Publish(body []byte, queueName string, exchangeName string, routingKey RoutingKeyType, exchangeType string) error
	PublishRPC(body []byte, queueName string, exchangeName string, routingKey RoutingKeyType, exchangeType string, done chan RPCResponse)
	PublishRPCRequest(body []byte, queueName string, exchangeName string, routingKey RoutingKeyType, exchangeType string, done chan bool, respHandleCB func(amqp.Delivery, RoutingKeyType))
	PublishRPCResponse(body []byte, correlationId string, queueName string, exchangeName string, routingKey RoutingKeyType, exchangeType string) error
	PublishOnQueue(msg []byte, queueName string) error
	Subscribe(queueName string, exchangeName string, exchangeType string, routingKeys []RoutingKeyType, consumerName string, handlerFunc func(amqp.Delivery, chan bool)) error
	SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery, chan bool)) error
	Close()
}

type RPCResponse struct {
	Status     bool
	Resp       *amqp.Delivery
	RoutingKey RoutingKeyType
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

// Real implementation, encapsulates a pointer to an amqp.Connection
type MsgClient struct {
	conn *amqp.Connection
}

// Generate token
func (m *MsgClient) GenerateToken() uint64 {
	return rand.Uint64()
}

//Connect to Broker(RabbitMq server)
func (m *MsgClient) ConnectToBroker(connectionString string) {
	if connectionString == "" {
		log.Errorln("Cannot initialize connection to broker, connectionString not set.")
	}

	var err error
	conn := false
	for !conn {

		m.conn, err = amqp.Dial(fmt.Sprintf("%s/", connectionString))
		if err != nil {

			// Error:: Service will not be reachable to any other serveice if not able to connect to MsgBus.
			log.Errorf("MsgBus:: Trying to connect to AMQP compatible broker at: " + connectionString)
			time.Sleep(5 * time.Second)

		} else {
			log.Infof("MsgBus:: Connected to AMQP compatible broker at: " + connectionString)
			conn = true
		}

	}

}

//Publish to queue through exchange
func (m *MsgClient) Publish(body []byte, queueName string, exchangeName string, routingKey RoutingKeyType, exchangeType string) error {
	if m.conn == nil {
		log.Errorln("Tried to send message before connection was initialized.")
	}

	// Get a channel from the connection
	ch, err := m.conn.Channel()
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to connect to channel.", err)
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to register an Exchange.", err)
		return err
	}

	// Declare a queue that will be created if not exists with some args
	queue, err := ch.QueueDeclare(
		queueName, // our queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to declare queue.", err)
		return err
	}

	//Binding queue with exchange
	err = ch.QueueBind(
		queue.Name,         // name of the queue
		string(routingKey), // bindingKey/routingkey
		exchangeName,       // sourceExchange
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to bind queue.", err)
		return err
	}

	// Publishes a message onto the queue.
	err = ch.Publish(
		exchangeName,       // exchange
		string(routingKey), // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			Body: body, // Our JSON body as []byte
		})

	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to publish message.", err)
	} else {
		log.Debugf("MsgBus::Message was sent on Exchange %s Queue %s Routing Key %s ", exchangeName, queue.Name, string(routingKey))
	}
	return err
}

// Publish to Queue.
func (m *MsgClient) PublishOnQueue(body []byte, queueName string) error {
	if m.conn == nil {
		log.Errorln("Tried to send message before connection was initialized.")
	}

	// Get a channel from the connection
	ch, err := m.conn.Channel()
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to connect to channel.", err)
		return err
	}
	defer ch.Close()

	// Declare a queue that will be created if not exists with some args
	queue, err := ch.QueueDeclare(
		queueName, // our queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to declare queue.", err)
		return err
	}

	// Publishes a message onto the queue.
	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body, // Our JSON body as []byte
		})
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to publish message.", err)
	} else {
		log.Debugf("MsgBus::Message was sent on Queue %s ", queue.Name)
	}
	return err
}

// Subscribe to exchange with option to listen to particular typr of message
func (m *MsgClient) Subscribe(queueName string, exchangeName string, exchangeType string, routingKeys []RoutingKeyType, consumerName string, handlerFunc func(amqp.Delivery, chan bool)) error {
	ch, err := m.conn.Channel()
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	failOnError(err, "Failed to register an Exchange")

	log.Printf("declared Exchange, declaring Queue (%s)", "")
	queue, err := ch.QueueDeclare(
		"",    // name of the queue
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to register an Queue")

	log.Debugf("MsgBus::declared Queue (%d messages, %d consumers), binding to Exchange (key '%s')",
		queue.Messages, queue.Consumers, exchangeName)

	//Binding queue with exchange
	for _, routingKey := range routingKeys {
		err = ch.QueueBind(
			queue.Name,         // name of the queue
			string(routingKey), // routing key /bindingKey
			exchangeName,       // sourceExchange
			false,              // noWait
			nil,                // arguments
		)
		if err != nil {
			log.Errorf("MsgBus::Queue Bind: %s", err)
			failOnError(err, "Failed to bind queue")
		}
	}

	msgs, err := ch.Consume(
		queue.Name,   // queue
		consumerName, // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	go consumeLoop(msgs, handlerFunc)
	return nil
}

//Subscribe directly to queue
func (m *MsgClient) SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery, chan bool)) error {
	ch, err := m.conn.Channel()
	failOnError(err, "Failed to open a channel")

	log.Debugf("MsgBus::Declaring Queue (%s)", queueName)
	queue, err := ch.QueueDeclare(
		queueName, // name of the queue
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	failOnError(err, "Failed to register an Queue")

	msgs, err := ch.Consume(
		queue.Name,   // queue
		consumerName, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	go consumeLoop(msgs, handlerFunc)
	return nil
}

// Close connection
func (m *MsgClient) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}

// Read messages from Queue.
func consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc func(d amqp.Delivery, ch chan bool)) {
	for d := range deliveries {

		// Invoke the handlerFunc func we passed as parameter.
		go handleTransit(d, handlerFunc)

	}
}

// This Go-Routine is transit between message consumer and handler.
func handleTransit(msg amqp.Delivery, handlerFunc func(d amqp.Delivery, ch chan bool)) {

	//channel to sync
	done := make(chan bool, 1)

	// handler for incoming messages.
	handlerFunc(msg, done)

	// Ack response
	select {

	// Request processed but it may be success or failure
	case <-done:
		sendAck(msg)
		// TODO: Check if it fails what to do.
		// if resp {
		// 	log.Debugf("MsgBus::Responding with Ack to the Request msg %+v", msg)
		// 	sendAck(msg)
		// } else {
		// 	log.Errorf("MsgBus::Internal server error to the Request msg %+v", msg)
		// 	sendNack(msg)
		// }

		//Time out TODO: Read the timeout from config.
	case <-time.After(1 * time.Second):
		log.Errorf("MsgBus::Timeout while respodning to request.")
		sendNack(msg)
	}

}

// Fatal Logging error info
func failOnError(err error, msg string) {
	if err != nil {
		log.Errorf("MsgBus::%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

//Logging error info
func LogOnError(err error, msg string) {
	if err != nil {
		log.Errorf("MsgBus::%s: %s", msg, err)
	}
}

// Ack to send message
func sendAck(msg amqp.Delivery) {
	if err := msg.Ack(false); err != nil {
		log.Errorf("MsgBus::Error acknowledging message [%+v]:: %s", msg, err)
	}
}

// Nack to handle negative messages
func sendNack(msg amqp.Delivery) {
	if err := msg.Nack(true, true); err != nil {
		log.Errorf("MsgBus::Error acknowledging message [%+v]:: %s", msg, err)
	} else {
		log.Debugf("MsgBus::Acknowledged message [%+v]", msg)
	}
}

// Publish RPC request. After sending a request meesage on topic exchange wait on the amq.rabbitmq.reply-to for response.
func (m *MsgClient) PublishRPCRequest(body []byte, queueName string, exchangeName string, routingKey RoutingKeyType, exchangeType string, done chan bool, respHandleCB func(amqp.Delivery, RoutingKeyType)) {
	log.Debugf("MsgBus::Publishing RPC messages Queue Name %s Exchange Name %s Routing Key %s Exchange Type %s ", queueName, exchangeName, routingKey, exchangeType)
	if m.conn == nil {
		log.Errorln("Tried to send message before connection was initialized.")
	}

	// Get a channel from the connection
	ch, err := m.conn.Channel()
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to connect to channel.", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to register an Exchange.", err)
		return
	}

	// Declare a queue that will be created if not exists with some args
	queue, err := ch.QueueDeclare(
		queueName, // our queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to declare queue.", err)
		return
	}

	err = ch.QueueBind(
		queue.Name,         // name of the queue
		string(routingKey), // bindingKey/routingkey
		exchangeName,       // sourceExchange
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to bind queue.", err)
		return
	}

	resp, err := ch.Consume(
		"amq.rabbitmq.reply-to", // queue
		"ReplyToRPCConsumer",    // consumer
		true,                    // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to start receiver for RPC message.", err)
		return
	}

	corrId := randomString(32)

	err = ch.Publish( // Publishes a message onto the queue.
		exchangeName,       // exchange
		string(routingKey), // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       "amq.rabbitmq.reply-to",
			Body:          body, // Our JSON body as []byte
		})

	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to publish message.", err)
		return
	} else {
		log.Debugf("MsgBus::RPC Message was sent on Exchange %s Queue %s Routing Key %s", exchangeName, queue.Name, string(routingKey))
	}

	// Check amq.rabbitmq.reply-to queue for messages
	for d := range resp {
		if corrId == d.CorrelationId {
			log.Debugf("MsgBus::MSGBUS:: RPC response message is recieved for correlationID: %s and routing key %s ", d.CorrelationId, d.RoutingKey)
			respHandleCB(d, routingKey)
			break
		}
	}

	// Send status to Cleint
	done <- true

}

// Publish RPC After sending a request meesage on topic exchange wait on the amq.rabbitmq.reply-to for response.
func (m *MsgClient) PublishRPC(body []byte, queueName string, exchangeName string, routingKey RoutingKeyType, exchangeType string, done chan RPCResponse) {
	log.Debugf("MsgBus::Publishing RPC messages Queue Name %s Exchange Name %s Routing Key %s Exchange Type %s ", queueName, exchangeName, routingKey, exchangeType)
	if m.conn == nil {
		log.Errorln("Tried to send message before connection was initialized.")
	}

	// Get a channel from the connection
	ch, err := m.conn.Channel()
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to connect to channel.", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to register an Exchange.", err)
		return
	}

	// Declare a queue that will be created if not exists with some args
	queue, err := ch.QueueDeclare(
		queueName, // our queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to declare queue.", err)
		return
	}

	err = ch.QueueBind(
		queue.Name,         // name of the queue
		string(routingKey), // bindingKey/routingkey
		exchangeName,       // sourceExchange
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to bind queue.", err)
		return
	}

	resp, err := ch.Consume(
		"amq.rabbitmq.reply-to", // queue
		"ReplyToRPCConsumer",    // consumer
		true,                    // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to start receiver for RPC message.", err)
		return
	}

	corrId := randomString(32)

	err = ch.Publish( // Publishes a message onto the queue.
		exchangeName,       // exchange
		string(routingKey), // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       "amq.rabbitmq.reply-to",
			Body:          body, // Our JSON body as []byte
		})
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to publish message.", err)
		return
	} else {
		log.Debugf("MsgBus::RPC Message was sent on Exchange %s Queue %s Routing Key %s CorrelationID %s", exchangeName, queue.Name, string(routingKey), corrId)
	}

	var rpcResp RPCResponse

	// Check amq.rabbitmq.reply-to queue for messages
	for d := range resp {
		if corrId == d.CorrelationId {
			rpcResp.Status = true
			rpcResp.Resp = &d
			rpcResp.RoutingKey = routingKey
			log.Debugf("MsgBus::MSGBUS:: RPC response message is recieved for correlationID: %s and routing key %s ", d.CorrelationId, d.RoutingKey)
			break
		}
	}

	// Send status to Cleint
	done <- rpcResp

}

//Publish RPC response.
func (m *MsgClient) PublishRPCResponse(body []byte, correlationId string, queueName string, exchangeName string, routingKey RoutingKeyType, exchangeType string) error {
	if m.conn == nil {
		log.Errorln("Tried to send message before connection was initialized.")
	}

	// Get a channel from the connection
	ch, err := m.conn.Channel()
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to connect to channel.", err)
		return err
	}
	defer ch.Close()

	// Publishes a message onto the amq.rabbitmq.reply-to queue.
	err = ch.Publish(
		"",                 // exchange
		string(routingKey), // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: correlationId,
			Body:          body, // Our JSON body as []byte
		})
	if err != nil {
		log.Errorf("MsgBus::Err: %s Failed to publish message.", err)
	} else {
		log.Debugf("MsgBus::RPC Message response was sent with Routing Key %s and CorrelationId %s", string(routingKey), correlationId)
	}
	return err
}
