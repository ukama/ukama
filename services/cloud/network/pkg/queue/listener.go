package queue

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/streadway/amqp"

	"google.golang.org/grpc/credentials/insecure"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	pbmesh "github.com/ukama/ukama/services/common/pb/gen/ukamaos/mesh"
	"google.golang.org/protobuf/proto"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/services/cloud/network/pb/gen"
	"github.com/ukama/ukama/services/common/msgbus"
	"github.com/ukama/ukama/services/common/ukama"
	"google.golang.org/grpc"
)

type QueueListener struct {
	msgBusConn    msgbus.Consumer
	networkClient pb.NetworkServiceClient
	grpcTimeout   int
	serviceId     string
	grpcConn      *grpc.ClientConn
}

type OrgCreatedBody struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

func NewQueueListener(networkGrpcHost string, connectionString string, grpcTimeout int, serviceId string) (*QueueListener, error) {
	client, err := msgbus.NewConsumerClient(connectionString)
	if err != nil {
		return nil, err
	}

	networkConn, err := grpc.Dial(networkGrpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	return &QueueListener{
		networkClient: pb.NewNetworkServiceClient(networkConn),
		msgBusConn:    client,
		grpcTimeout:   grpcTimeout,
		serviceId:     serviceId,
		grpcConn:      networkConn,
	}, nil
}

func (q *QueueListener) StartQueueListening() (err error) {
	err = q.msgBusConn.SubscribeToServiceQueue("network-listener", msgbus.DeviceQ.Exchange,
		[]msgbus.RoutingKey{msgbus.DeviceConnectedRoutingKey, msgbus.OrgCreatedRoutingKey}, q.serviceId, q.incomingMessageHandler)
	if err != nil {
		log.Errorf("Error subscribing for a queue messages. Error: %+v", err)
		return err
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	log.Info("Shutting down...")
	q.Close()
	return nil
}

func (q *QueueListener) incomingMessageHandler(delivery amqp.Delivery, done chan<- bool) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(q.grpcTimeout))
	defer cancel()

	switch delivery.RoutingKey {

	case string(msgbus.DeviceConnectedRoutingKey):
		q.processDeviceConnectedMsg(ctx, delivery)

	case string(msgbus.OrgCreatedRoutingKey):
		q.processOrgCreatedMsg(ctx, delivery)

	default:
		log.Warning("No handler for routing key ", delivery.RoutingKey)
	}

	done <- true
}

func (q *QueueListener) processOrgCreatedMsg(ctx context.Context, delivery amqp.Delivery) {
	org := &OrgCreatedBody{}
	err := json.Unmarshal(delivery.Body, org)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		return
	}
	_, err = q.networkClient.AddNetwork(ctx, &pb.AddNetworkRequest{
		Name:    "default",
		OrgName: org.Name,
	}, grpc_retry.WithMax(3))

	if err != nil {
		log.Errorf("Failed to add organization '%s'. Error: %v", org.Name, err)
	} else {
		log.Infof("Organization %s added successefully", org.Name)
	}
}

func (q *QueueListener) processDeviceConnectedMsg(ctx context.Context, delivery amqp.Delivery) {
	link := &pbmesh.Link{}
	err := proto.Unmarshal(delivery.Body, link)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		return
	}

	log.Debugf("Updating node %+v", delivery)
	nodeId, err := ukama.ValidateNodeId(link.GetNodeId())
	if err != nil {
		log.Errorf("Invalid Node ID format. Error %v", err)
		return
	}

	_, err = q.networkClient.UpdateNode(ctx,
		&pb.UpdateNodeRequest{
			NodeId: nodeId.String(),
			Node: &pb.Node{
				State: pb.NodeState_ONBOARDED,
			},
		},
		grpc_retry.WithMax(3))

	if err != nil {
		log.Errorf("Failed to update node %s status. Error: %v", nodeId.String(), err)
	} else {
		log.Infof("Node %s updated successefully", nodeId.String())
	}
}

func (q *QueueListener) Close() {
	q.msgBusConn.Close()
	q.grpcConn.Close()
}
