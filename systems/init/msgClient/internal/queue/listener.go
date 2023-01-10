package queue

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	mb "github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QueueListener struct {
	mConn       mb.Consumer
	gConn       *grpc.ClientConn
	grpcTimeout time.Duration
	serviceUuid string
	serviceName string
	serviceHost string
	state       bool
	queue       string
	exchange    string
	c           chan bool
	routes      []string
}

func NewQueueListener(s db.Service) (*QueueListener, error) {

	log.Debugf("Listener Config %+v", s)
	routes := make([]string, len(s.Routes))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.GrpcTimeout))
	defer cancel()

	conn, err := grpc.DialContext(ctx, s.ServiceUri, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	// if err != nil {
	// 	log.Fatalf("Could not connect: %v", err)
	// }

	client, err := mb.NewConsumerClient(s.MsgBusUri)
	if err != nil {
		return nil, err
	}

	for idx, r := range s.Routes {
		/*  Create a queue listner for each service */
		routes[idx] = r.Key
	}

	ch := make(chan bool, 1)

	return &QueueListener{
		mConn:       client,
		serviceUuid: s.ServiceUuid,
		serviceName: s.Name,
		serviceHost: s.ServiceUri,
		c:           ch,
		state:       false,
		gConn:       conn,
		routes:      routes,
		queue:       s.ListQueue,
		exchange:    s.Exchange,
	}, nil
}

func (q *QueueListener) startQueueListening() {

	log.Debugf("[%s]: Starting listener routine.", q.serviceName)
	/* Validate routes */ // TODO: Update ParseRoutesList implementation
	routes, err := mb.ParseRouteList(q.routes)
	if err != nil {
		log.Errorf("[%s] Failed to create listener. Error %s", q.serviceName, err.Error())
	}

	/* Subscribe to exchange for the routes */
	err = q.mConn.SubscribeToServiceQueue(q.serviceName, q.exchange,
		routes, q.serviceUuid, q.incomingMessageHandler)
	if err != nil {
		log.Errorf("[%s] Failed to create listener. Error %s", q.serviceName, err.Error())
		log.Errorf("[%s] Shutting down listener.", q.serviceName)
		q.mConn.Close()
		q.state = false
		return
	}

	q.state = true
	log.Infof("[%s] Queue listener started on %v routes", q.serviceName, q.routes)
	/* Waiting for stop */
	<-q.c

	log.Infof("[%s] Shutting down queue listener", q.serviceName)
	q.mConn.Close()
	q.state = false

}

func (q *QueueListener) stopQueueListening() {
	if q.state {
		log.Infof("Stopping queue listener routine for service %s on %v routes", q.serviceName, q.routes)
		q.c <- true
	}
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
	log.Infof("Received message %+v", delivery)
}
