package queue

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/ukama/ukama/systems/common/config"
	uconf "github.com/ukama/ukama/systems/common/config"
	mb "github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
)

type QueueListener struct {
	s           db.ServiceRepo
	r           db.RoutingKeyRepo
	m           mb.Consumer
	grpcTimeout time.Duration
	serviceId   string
	serviceName string
	state       bool
	c           chan bool
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
		m:           client,
		serviceId:   serviceId,
		serviceName: serviceName,
		c:           ch,
		state:       false,
		s:           sRepo,
		r:           kRepo,
	}, nil
}

func (q *QueueListener) StartQueueListening() (err error) {

	/* Read the possible list of the Routes from db if empty */
	strRoutes, err := q.r.ReadAllRoutes()
	if err != nil {
		log.Errorf("Error reading routes. Error %s", err.Error())
		return err
	}

	if len(strRoutes) <= 0 {
		log.Errorf("No routes available.")
		return fmt.Errorf("no routes to monitor")
	}

	routes, err := mb.ParseRouteList(strRoutes)
	if err != nil {
		return err
	}

	/* Subscribe to exchange for the routes */
	err = q.m.SubscribeToServiceQueue(q.serviceName, mb.DeviceQ.Exchange,
		routes, q.serviceId, q.incomingMessageHandler)
	if err != nil {
		log.Errorf("Error subscribing for a queue messages. Error: %+v", err)
		return err
	}

	q.state = true

	/* Waiting for stop */
	<-q.c

	log.Info("Shutting down...")
	q.m.Close()
	q.state = false

	return nil
}

func (q *QueueListener) StopQueueListening() {
	q.c <- true
	time.Sleep(1 * time.Second) // TODO: Update this
}

func (q *QueueListener) RetstartListening() (err error) {
	if q.state {
		q.StopQueueListening()
	}
	return q.StartQueueListening()
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
