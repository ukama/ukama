package queue

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"github.com/ukama/ukamaX/cloud/registry/pb/gen/external"
	pbmocks "github.com/ukama/ukamaX/cloud/registry/pb/gen/mocks"
	"github.com/ukama/ukamaX/common/ukama"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestIncomingMessageHandler(t *testing.T) {
	// Arrange
	reg := &pbmocks.RegistryServiceClient{}
	nodeId := string(ukama.NewVirtualNodeId("homenode"))

	reg.On("UpdateNode", mock.Anything, mock.MatchedBy(func(r *pb.UpdateNodeRequest) bool {
		return r.NodeId == nodeId && r.State == pb.NodeState_ONBOARDED
	})).Return(nil, nil)

	message, err := proto.Marshal(&external.Link{Uuid: &nodeId})
	assert.NoError(t, err)
	delivery := amqp.Delivery{Body: message}

	q := &QueueListener{
		registryClient: reg,
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
	nodeId := "randome node id"
	badUuidMessage, _ := proto.Marshal(&external.Link{Uuid: &nodeId})

	tests := []struct {
		name    string
		message []byte
	}{
		{name: "BadMessageFormat", message: []byte("random message")},
		{name: "BadMessageFormat", message: badUuidMessage},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := &pbmocks.RegistryServiceClient{}

			delivery := amqp.Delivery{Body: tt.message}

			q := &QueueListener{
				registryClient: reg,
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
