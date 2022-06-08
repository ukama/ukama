package queue

import (
	"fmt"
	"testing"

	uuid "github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/ukama/ukama/services/cloud/network/pb/gen"
	pbmocks "github.com/ukama/ukama/services/cloud/network/pb/gen/mocks"
	"github.com/ukama/ukama/services/common/msgbus"
	"github.com/ukama/ukama/services/common/pb/gen/ukamaos/mesh"
	"github.com/ukama/ukama/services/common/ukama"
	"google.golang.org/protobuf/proto"
)

func TestDeviceIncomingMessageHandler(t *testing.T) {
	// Arrange
	reg := &pbmocks.NetworkServiceClient{}
	nodeId := string(ukama.NewVirtualNodeId("homenode"))

	reg.On("UpdateNode", mock.Anything, mock.MatchedBy(func(r *pb.UpdateNodeRequest) bool {
		return r.NodeId == nodeId && r.GetNode().State == pb.NodeState_ONBOARDED
	}), mock.Anything).Return(nil, nil)

	message, err := proto.Marshal(&mesh.Link{NodeId: &nodeId, Ip: proto.String("192.168.0.1")})
	assert.NoError(t, err)
	delivery := amqp.Delivery{Body: message, RoutingKey: string(msgbus.DeviceConnectedRoutingKey)}

	q := &QueueListener{
		networkClient: reg,
	}
	done := make(chan bool)

	// Act
	go func() { q.incomingMessageHandler(delivery, done) }()
	ret := <-done
	// Assert
	reg.AssertExpectations(t)
	assert.Equal(t, true, ret)
}

func TestUserRegisteredMessageHandler(t *testing.T) {
	// Arrange
	reg := &pbmocks.NetworkServiceClient{}
	userId := uuid.NewString()

	reg.On("AddOrg", mock.Anything, mock.MatchedBy(func(r *pb.AddOrgRequest) bool {
		return r.Name == userId && r.Owner == userId
	}), mock.Anything).Return(nil, nil)

	message := fmt.Sprintf(`{
 "email": "dev+a19996db-417a-410e-a7b5-d1623f232697@dev.ukama.com",
 "id": "%s"
}`, userId)

	delivery := amqp.Delivery{Body: []byte(message), RoutingKey: string(msgbus.UserRegisteredRoutingKey)}

	q := &QueueListener{
		networkClient: reg,
	}
	done := make(chan bool)

	// Act
	go func() { q.incomingMessageHandler(delivery, done) }()
	ret := <-done
	// Assert
	reg.AssertExpectations(t)
	assert.Equal(t, true, ret)
}

func TestIncomingMessageHandler_MessageFormatErrors(t *testing.T) {
	nodeId := "random node id"
	badUuidMessage, _ := proto.Marshal(&mesh.Link{NodeId: &nodeId})

	tests := []struct {
		name       string
		message    []byte
		routingKey msgbus.RoutingKey
	}{
		{name: "DeviceRegistered", message: []byte("random message"), routingKey: msgbus.DeviceConnectedRoutingKey},
		{name: "DeviceRegistered_BadUuid", message: badUuidMessage, routingKey: msgbus.DeviceConnectedRoutingKey},
		{name: "UserRegisteredMessage", message: []byte("random message"), routingKey: msgbus.UserRegisteredRoutingKey},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &pbmocks.NetworkServiceClient{}

			delivery := amqp.Delivery{Body: tt.message, RoutingKey: string(tt.routingKey)}

			q := &QueueListener{
				networkClient: reg,
			}
			done := make(chan bool)

			// Act
			go func() { q.incomingMessageHandler(delivery, done) }()
			ret := <-done

			// Assert
			assert.Equal(t, true, ret)
			// make sure we don't call update node
			reg.AssertExpectations(t)
		})
	}
}
