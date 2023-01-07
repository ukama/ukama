package queue

import (
	"context"
	"fmt"
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
	serviceId   string
	serviceName string
	serviceHost string
	state       bool
	queue       string
	exchange    string
	c           chan bool
	routes      []string
}

type MsgBusListener struct {
	ql map[string]*QueueListener
	s  db.ServiceRepo
	r  db.RouteRepo
}

func NewMessageBusListener(s db.ServiceRepo, r db.RouteRepo) *MsgBusListener {

	mbl := &MsgBusListener{
		s: s,
		r: r,
	}
	mbl.ql = make(map[string]*QueueListener)

	return mbl

}

func newQueueListener(s db.Service) (*QueueListener, error) {

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
		serviceId:   s.ServiceUuid,
		serviceName: s.Name,
		c:           ch,
		state:       false,
		gConn:       conn,
		routes:      routes,
		queue:       s.QueueName,
		exchange:    s.Exchange,
	}, nil
}

func startQueueListening(q *QueueListener) {

	log.Debugf("[%s]: Starting listener routine.", q.serviceName)
	/* Validate routes */ // TODO: Update ParseRoutesList implementation
	routes, err := mb.ParseRouteList(q.routes)
	if err != nil {
		log.Errorf("[%s] Failed to create listener. Error %s", q.serviceName, err.Error())
	}

	/* Subscribe to exchange for the routes */
	err = q.mConn.SubscribeToServiceQueue(q.serviceName, q.exchange,
		routes, q.serviceId, q.incomingMessageHandler)
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

func stopQueueListening(q *QueueListener) {
	if q.state {
		log.Infof("Stopping queue listener routine for service %s on %v routes", q.serviceName, q.routes)
		q.c <- true
	}
}

func (m *MsgBusListener) CreateQueueListeners() error {

	services, err := m.s.List()
	if err != nil {
		log.Errorf("Error reading services. Error %s", err.Error())
		return err
	}

	if len(services) <= 0 {
		log.Errorf("No services available.")
	}

	for _, s := range services {
		/*  Create a queue listner for each service */
		log.Infof("Starting listener for %s service.", s.ServiceUuid)
		listener, err := newQueueListener(s)
		if err != nil {
			log.Errorf("Failed to create listner for %s. Error %s", s.Name, err.Error())
			return err
		}

		m.ql[s.ServiceUuid] = listener

	}

	m.StartQueueListeners()

	return nil
}

func (m *MsgBusListener) StartQueueListeners() {

	for _, q := range m.ql {
		/*  Create a queue listner for each service */
		log.Infof("Starting new queue listener routine for service %s on %v routes", q.serviceName, q.routes)
		go startQueueListening(q)

	}

}

func (m *MsgBusListener) StopQueueListener() {

	for _, q := range m.ql {
		stopQueueListening(q)
	}
}

func (m *MsgBusListener) RetstartServiceQueueListening(service string) (err error) {
	q, ok := m.ql[service]
	if ok {
		stopQueueListening(q)
		time.Sleep(500 * time.Millisecond)
		if !q.state {
			startQueueListening(q)
		}
	}
	return nil
}

func (m *MsgBusListener) UpdateServiceQueueListening(s *db.Service) (err error) {
	_, ok := m.ql[s.ServiceUuid]
	if ok {
		m.RemoveServiceQueueListening(s.ServiceUuid)
	}

	listener, err := newQueueListener(*s)
	if err != nil {
		log.Errorf("Failed to create listener for %s. Error %s", s.Name, err.Error())
		return err
	}

	m.ql[s.ServiceUuid] = listener

	go startQueueListening(listener)

	time.Sleep(500 * time.Millisecond)

	if !listener.state {
		return fmt.Errorf("failed to start listener for service %s", listener.serviceName)
	}

	return nil
}

func (m *MsgBusListener) StopServiceQueueListening(service string) (err error) {
	q, ok := m.ql[service]
	if ok {
		stopQueueListening(q)
		time.Sleep(500 * time.Millisecond)
		if q.state {
			return fmt.Errorf("failed to stop queue listening service for %s", q.serviceName)
		}
	} else {
		return fmt.Errorf("no service with id %s registered", service)
	}

	return nil
}

func (m *MsgBusListener) RemoveServiceQueueListening(service string) error {
	log.Infof("Removing queue listener for %s service", service)

	err := m.StopServiceQueueListening(service)
	if err != nil {
		return err
	}
	delete(m.ql, service)

	log.Infof("Removed queue listener for %s service", service)

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
