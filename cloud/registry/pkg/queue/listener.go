package queue

import (
	"context"
	"encoding/json"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ukama/ukamaX/cloud/registry/pb/gen/external"
	"google.golang.org/protobuf/proto"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	pb "github.com/ukama/ukamaX/cloud/registry/pb/gen"
	"github.com/ukama/ukamaX/common/msgbus"
	"github.com/ukama/ukamaX/common/ukama"
	"google.golang.org/grpc"
)

type QueueListener struct {
	msgBusConn     msgbus.Consumer
	registryClient pb.RegistryServiceClient
	grpcTimeout    int
	serviceId      string
}

type UserRegisteredBody struct {
	Id    string `json:"id"`
	Email string `json:"email"`
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

	err = q.msgBusConn.SubscribeToServiceQueue("registry-listener", msgbus.DeviceQ.Exchange,
		[]msgbus.RoutingKey{msgbus.DeviceConnectedRoutingKey, msgbus.UserRegisteredRoutingKey}, q.serviceId, q.incomingMessageHandler)
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

	switch delivery.RoutingKey {

	case string(msgbus.DeviceConnectedRoutingKey):
		q.processDeviceConnectedMsg(ctx, delivery)

	case string(msgbus.UserRegisteredRoutingKey):
		q.processUserRegisteredMsg(ctx, delivery)

	default:
		log.Warning("No handler for routing key ", delivery.RoutingKey)
	}

	done <- true
}

func (q *QueueListener) processUserRegisteredMsg(ctx context.Context, delivery amqp.Delivery) {
	user := &UserRegisteredBody{}
	err := json.Unmarshal(delivery.Body, user)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		return
	}
	_, err = q.registryClient.AddOrg(ctx, &pb.AddOrgRequest{
		Name:  user.Id,
		Owner: user.Id,
	}, grpc_retry.WithMax(3))

	if err != nil {
		log.Errorf("Failed to add organization '%s'. Error: %v", user.Id, err)
	} else {
		log.Infof("Organization %s added successefully", user.Id)
	}
}

func (q *QueueListener) processDeviceConnectedMsg(ctx context.Context, delivery amqp.Delivery) {
	link := &external.Link{}
	err := proto.Unmarshal(delivery.Body, link)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		return
	}

	log.Debugf("Updating node %+v", delivery)
	nodeId, err := ukama.ValidateNodeId(link.GetUuid())
	if err != nil {
		log.Errorf("Invalid Node ID format. Error %v", err)
		return
	}

	_, err = q.registryClient.UpdateNode(ctx,
		&pb.UpdateNodeRequest{
			NodeId: nodeId.String(),
			State:  pb.NodeState_ONBOARDED},
		grpc_retry.WithMax(3))

	if err != nil {
		log.Errorf("Failed to update node %s status. Error: %v", nodeId.String(), err)
	} else {
		log.Infof("Node %s updated successefully", nodeId.String())
	}
}
