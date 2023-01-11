package queue

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	mb "github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type QueueListener struct {
	mConn       mb.Consumer
	gConn       *grpc.ClientConn
	gClient     pb.EventNotificationServiceClient
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
	var gc pb.EventNotificationServiceClient
	log.Debugf("Listener Config %+v", s)
	routes := make([]string, len(s.Routes))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.GrpcTimeout))
	defer cancel()

	conn, err := grpc.DialContext(ctx, s.ServiceUri, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Errorf("Could not connect to %s. Error %s Will try again at message reception.", s.ServiceUri, err.Error())
	} else {
		gc = pb.NewEventNotificationServiceClient(conn)
	}

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
		gClient:     gc,
		routes:      routes,
		queue:       s.ListQueue,
		exchange:    s.Exchange,
		grpcTimeout: time.Duration(s.GrpcTimeout) * time.Second,
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
		q.processEventMsg(ctx, delivery)
	}

	done <- true
}

func (q *QueueListener) processEventMsg(ctx context.Context, d amqp.Delivery) {
	// Read Db for the key and find the services which we need to post message to.
	log.Debugf("Raw message: %+v", d)

	evtAny := new(anypb.Any)
	err := proto.Unmarshal(d.Body, evtAny)
	if err != nil {
		log.Errorf("Failed to parse message with key %s. Error %s", d.RoutingKey, err.Error())
		return
	}
	e := &pb.Event{
		RoutingKey: d.RoutingKey,
		Msg:        evtAny,
	}

	log.Infof("Received a message: %+v", e)

	if q.gConn == nil {
		conn, err := grpc.DialContext(ctx, q.serviceHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Errorf("Could not connect to %s. Error %s Will try again at message reception.", q.serviceHost, err.Error())
			return
		} else {
			q.gClient = pb.NewEventNotificationServiceClient(conn)
		}

		q.gConn = conn
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = q.gClient.EventNotification(ctx, e)
	if err != nil {
		log.Errorf("Failed to send message to %s with key %s. Error %s", q.serviceHost, d.RoutingKey, err.Error())
	}
}
