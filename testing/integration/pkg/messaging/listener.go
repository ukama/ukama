package messaging

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/wagslane/go-rabbitmq"
)

const POD_NAME_ENV_VAR = "POD_NAME"

type listener struct {
	conn      *rabbitmq.Conn
	cons      *rabbitmq.Consumer
	serviceId string
	store     map[string]interface{}
	stop      chan bool
	uri       string
}

type ListenerConfig struct {
	config.BaseConfig `mapstructure:",squash"`
	Queue             config.Queue
}

type Listener interface {
	StartListener()
	StopListener()
	GetEvent(key string) (interface{}, bool)
}

func NewListenerConfig() *ListenerConfig {
	return &ListenerConfig{
		Queue: config.Queue{
			Uri: "amqp://guest:guest@localhost:5672/",
		},
	}
}

func NewListener(config *ListenerConfig) Listener {
	conn, err := rabbitmq.NewConn(
		config.Queue.Uri,
		rabbitmq.WithConnectionOptionsLogging,
	)

	if err != nil {
		log.Fatalf("error creating connection. Error: %s", err.Error())
	}

	l := &listener{
		conn:      conn,
		serviceId: os.Getenv(POD_NAME_ENV_VAR),
		store:     make(map[string]interface{}),
		stop:      make(chan bool, 1),
		uri:       config.Queue.Uri,
	}
	log.Debugf("Listener created: %+v.", l)

	return l
}
func (l *listener) StartListener() {

	consumer, err := rabbitmq.NewConsumer(l.conn, l.incomingMessageHandler, l.uri,
		rabbitmq.WithConsumerOptionsRoutingKey("#"),
		rabbitmq.WithConsumerOptionsExchangeName(msgbus.DefaultExchange),
		rabbitmq.WithConsumerOptionsConsumerName(l.serviceId),
		rabbitmq.WithConsumerOptionsExchangeKind(amqp.ExchangeTopic))
	if err != nil {
		log.Fatalf("error creating queue consumer. Error: %s", err.Error())
	}

	l.cons = consumer
	log.Infof("Creating listener for Queue: %s. lsitner: %+v",
		l.uri[strings.LastIndex(l.uri, "@"):], l)

	defer l.conn.Close()

	defer consumer.Close()

	log.Info("Listening for messages...")
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	go func() {

		select {
		case <-sigs:
			sig := <-sigs
			log.Println()
			log.Info(sig)
			done <- true

		case <-l.stop:
			log.Info("Stopping")
			done <- true
		}

	}()

	log.Info("awaiting signal")
	<-done
	log.Info("stopping consumer")
}

func (l *listener) StopListener() {
	l.stop <- true
}

func (l *listener) incomingMessageHandler(delivery rabbitmq.Delivery) rabbitmq.Action {
	log.Infof("Raw message: %+v", delivery)

	l.store[delivery.RoutingKey] = delivery.Body
	log.Infof("Added message %s with body %v", delivery.RoutingKey, delivery.Body)

	return rabbitmq.Ack
}

func (l *listener) GetEvent(key string) (interface{}, bool) {
	if m, ok := l.store[key]; ok {
		return m, ok
	}
	return nil, false
}
