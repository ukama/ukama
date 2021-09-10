package queue

import (
	"context"
	"github.com/ukama/ukamaX/cloud/registry/pb/gen/external"
	"google.golang.org/protobuf/proto"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/generated"
	"github.com/ukama/ukamaX/common/msgbus"
	"github.com/ukama/ukamaX/common/ukama"
	"google.golang.org/grpc"
)

const DEVICE_CONNECTED_ROUTING_KEY = "event.device.mesh.link.connect"

type QueueListener struct {
	msgBusConn     msgbus.Consumer
	registryClient pb.RegistryServiceClient
	grpcTimeout    int
	serviceId      string
}

func NewQueueListener(registryGrpcHost string, connectionString string, grpcTimeout int, serviceId string) (*QueueListener, error) {
	client, err := msgbus.NewConsumerClient(connectionString)
	if err != nil {
		return nil, err
	}

	registryConn, err := grpc.Dial(registryGrpcHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	return &QueueListener{
		registryClient: pb.NewRegistryServiceClient(registryConn),
		msgBusConn:     client,
		grpcTimeout:    grpcTimeout,
		serviceId:      serviceId,
	}, nil
}

func (q *QueueListener) StartQueueListening() (err error) {
	routingKeys := msgbus.RoutingKeyType("event.device.mesh.link.connect")

	err = q.msgBusConn.SubscribeToServiceQueue("registry-listener", msgbus.DeviceQ.Exchange,
		[]msgbus.RoutingKeyType{routingKeys}, q.serviceId, q.incomingMessageHandler)
	if err != nil {
		log.Errorf("Error subscribing for a queue messages. Error: %+v", err)
		return err
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	return nil
}

func (q *QueueListener) incomingMessageHandler(delivery amqp.Delivery, done chan<- bool) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(q.grpcTimeout))
	defer cancel()

	link := &external.Link{}
	err := proto.Unmarshal(delivery.Body, link)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		done <- true
		return
	}

	log.Debugf("Updating node %+v", delivery)
	nodeId, err := ukama.ValidateNodeId(link.GetUuid())
	if err != nil {
		log.Errorf("Invalid Node ID format. Error %v", err)
		done <- true
		return
	}

	_, err = q.registryClient.UpdateNode(ctx,
		&pb.UpdateNodeRequest{
			NodeId: nodeId.String(),
			State:  pb.NodeState_ONBOARDED})

	if err != nil {
		log.Errorf("Failed to update node %s status. Error: %v", nodeId.String(), err)
	} else {
		log.Infof("Node %s updated successefully", nodeId.String())
	}
	done <- true
}
