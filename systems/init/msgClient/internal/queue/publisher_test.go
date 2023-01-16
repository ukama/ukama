package queue

import (
	"testing"

	mocks "github.com/ukama/ukama/systems/common/mocks"
	mb "github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
	pb "github.com/ukama/ukama/systems/init/msgClient/pb/gen"

	"github.com/stretchr/testify/assert"
)

var route1 = db.Route{
	Key: "event.cloud.lookup.organization.create",
}

var ServiceUuid = "1ce2fa2f-2997-422c-83bf-92cf2e7334dd"
var service1 = db.Service{
	Name:        "test",
	InstanceId:  "1",
	MsgBusUri:   "amqp://guest:guest@localhost:5672",
	ListQueue:   "",
	PublQueue:   "",
	Exchange:    "amq.topic",
	ServiceUri:  "localhost:9090",
	GrpcTimeout: 5,
}

func NewTestQueuePublisher() *QueuePublisher {
	pub := &mocks.QPub{}

	qp := &QueuePublisher{
		pub:            pub,
		baseRoutingKey: mb.NewRoutingKeyBuilder().SetCloudSource().SetContainer("test"),
	}

	return qp
}

func TestQueuePublisher_Publish(t *testing.T) {
	pub := &mocks.QPub{}
	qp := &QueuePublisher{
		pub:            pub,
		baseRoutingKey: mb.NewRoutingKeyBuilder().SetCloudSource().SetContainer("test"),
	}

	msg := pb.PublishMsgRequest{
		ServiceUuid: ServiceUuid,
	}

	pub.On("PublishProto", &msg, route1.Key).Return(nil).Once()

	err := qp.Publish(route1.Key, &msg)

	assert.NoError(t, err)
	pub.AssertExpectations(t)

}

func TestQueuePublisher_Close(t *testing.T) {
	pub := &mocks.QPub{}
	qp := &QueuePublisher{
		pub:            pub,
		baseRoutingKey: mb.NewRoutingKeyBuilder().SetCloudSource().SetContainer("test"),
	}

	pub.On("Close").Return(nil).Once()

	err := qp.Close()

	assert.NoError(t, err)
	pub.AssertExpectations(t)

}
