package queue

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	mocks "github.com/ukama/ukama/systems/common/mocks"
	mb "github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
)

var route = []mb.RoutingKey{mb.RoutingKey("event.cloud.lookup.organization.create")}

var service = db.Service{
	Name:        "test",
	ServiceUuid: "1ce2fa2f-2997-422c-83bf-92cf2e7334dd",
	InstanceId:  "1",
	MsgBusUri:   "amqp://guest:guest@localhost:5672",
	ListQueue:   "",
	PublQueue:   "",
	Exchange:    "amq.topic",
	ServiceUri:  "localhost:9090",
	GrpcTimeout: 5,
	Routes:      []db.Route{db.Route{Key: "event.cloud.lookup.organization.create"}},
}

func NewTestQueueListener(s db.Service) *QueueListener {

	routes := make([]string, len(s.Routes))
	for idx, r := range s.Routes {
		/*  Create a queue listner for each service */
		routes[idx] = r.Key
	}

	ch := make(chan bool, 1)

	return &QueueListener{
		serviceUuid: s.ServiceUuid,
		serviceName: s.Name,
		serviceHost: s.ServiceUri,
		c:           ch,
		state:       false,
		routes:      routes,
		queue:       s.ListQueue,
		exchange:    s.Exchange,
	}
}

func TestQueuePublisher_startstopQueueListening(t *testing.T) {
	client := &mocks.Consumer{}
	qp := NewTestQueueListener(service)
	qp.mConn = client

	client.On("SubscribeToServiceQueue", qp.serviceName, qp.exchange, route, qp.serviceUuid, mock.AnythingOfType("func(amqp.Delivery, chan<- bool)")).Return(nil).Once()

	go qp.startQueueListening()

	time.Sleep(2 * time.Second)

	qp.stopQueueListening()

	client.AssertExpectations(t)

}
