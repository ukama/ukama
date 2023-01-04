package queue

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/ukama/ukama/systems/common/config"
	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/msgbus"
	mb "github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
)

type QueueListener struct {
	serviceRepo    db.ServiceRepo
	routingKeyRepo db.RoutingKeyRepo
	msgBusConn     mb.Consumer
	grpcTimeout    time.Duration
	serviceId      string
	serviceName    string
	state          bool
	channel        chan bool
}

type QueueListenerConfig struct {
	config.BaseConfig `mapstructure:",squash"`
	Service           struct {
		Host    string        `default:"localhost:9090"`
		Timeout time.Duration `default:"3s"`
	}
	Queue   config.Queue
	Metrics config.Metrics
}

func NewQueueListener(conf *uconf.Queue, serviceName string, serviceId string, sRepo db.ServiceRepo, kRepo db.RoutingKeyRepo) (*QueueListener, error) {

	client, err := mb.NewConsumerClient(conf.Uri)
	if err != nil {
		return nil, err
	}

	ch := make(chan bool, 1)

	return &QueueListener{
		msgBusConn:     client,
		serviceId:      serviceId,
		serviceName:    serviceName,
		channel:        ch,
		state:          false,
		serviceRepo:    sRepo,
		routingKeyRepo: kRepo,
	}, nil
}

func (q *QueueListener) StartQueueListening() (err error) {

	/* Read the possible list of the Routes from db if empty */

	/* Subscribe to exchange for the routes */
	err = q.msgBusConn.SubscribeToServiceQueue(q.serviceName, mb.DeviceQ.Exchange,
		[]mb.RoutingKey{msgbus.DeviceConnectedRoutingKey}, q.serviceId, q.incomingMessageHandler)
	if err != nil {
		log.Errorf("Error subscribing for a queue messages. Error: %+v", err)
		return err
	}

	q.state = true

	/* Waiting for stop */
	<-q.channel

	log.Info("Shutting down...")
	q.msgBusConn.Close()
	q.state = false

	return nil
}

func (q *QueueListener) StopQueueListening() (err error) {
	q.channel <- true
	return nil
}

func (q *QueueListener) incomingMessageHandler(delivery amqp.Delivery, done chan<- bool) {
	ctx, cancel := context.WithTimeout(context.Background(), q.grpcTimeout)
	defer cancel()

	switch delivery.RoutingKey {
	case string(mb.DeviceConnectedRoutingKey):
		q.processEventMsg(ctx, delivery)

	default:
		log.Warning("No handler for routing key ", delivery.RoutingKey)
	}

	done <- true
}

func (q *QueueListener) processEventMsg(ctx context.Context, delivery amqp.Delivery) {
	// Read Db for the key and find the services which we need to post message to.

}
