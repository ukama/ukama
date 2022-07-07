package queue

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ukama/ukama/services/cloud/network/pkg"

	"github.com/streadway/amqp"

	"google.golang.org/grpc/credentials/insecure"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	log "github.com/sirupsen/logrus"
	pb "github.com/ukama/ukama/services/cloud/network/pb/gen"
	"github.com/ukama/ukama/services/common/msgbus"
	"github.com/ukama/ukama/services/common/ukama"
	"google.golang.org/grpc"
)

type QueueListener struct {
	msgBusConn    msgbus.Consumer
	networkClient pb.NetworkServiceClient
	grpcTimeout   time.Duration
	serviceId     string
	grpcConn      *grpc.ClientConn
}

func NewQueueListener(networkGrpcHost string, connectionString string, grpcTimeout time.Duration, serviceId string) (*QueueListener, error) {
	client, err := msgbus.NewConsumerClient(connectionString)
	if err != nil {
		return nil, err
	}

	log.Info("Connecting to network service")
	ctx, cancel := context.WithTimeout(context.Background(), grpcTimeout)
	defer cancel()

	networkConn, err := grpc.DialContext(ctx, networkGrpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect to network service: %v", err)
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
	log.Info("Starting queue listener")
	err = q.msgBusConn.SubscribeToServiceQueue(pkg.ServiceName+"-listener", msgbus.DefaultExchange,
		[]msgbus.RoutingKey{msgbus.NodeUpdatedRoutingKey, msgbus.OrgCreatedRoutingKey}, q.serviceId, q.incomingMessageHandler)
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

	ctx, cancel := context.WithTimeout(context.Background(), q.grpcTimeout)
	defer cancel()

	switch delivery.RoutingKey {

	case string(msgbus.NodeUpdatedRoutingKey):
		q.processNodeUpdatedMsg(ctx, delivery)

	case string(msgbus.OrgCreatedRoutingKey):
		q.processOrgCreatedMsg(ctx, delivery)

	default:
		log.Warning("No handler for routing key ", delivery.RoutingKey)
	}

	done <- true
}

func (q *QueueListener) processOrgCreatedMsg(ctx context.Context, delivery amqp.Delivery) {
	org := &msgbus.OrgCreatedBody{}
	err := json.Unmarshal(delivery.Body, org)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		return
	}
	_, err = q.networkClient.Add(ctx, &pb.AddRequest{
		Name:    "default",
		OrgName: org.Name,
	}, grpc_retry.WithMax(3))

	if err != nil {
		log.Errorf("Failed to add organization '%s'. Error: %v", org.Name, err)
	} else {
		log.Infof("Organization %s added successefully", org.Name)
	}
}

func (q *QueueListener) processNodeUpdatedMsg(ctx context.Context, delivery amqp.Delivery) {
	body := &msgbus.NodeUpdateBody{}
	err := json.Unmarshal(delivery.Body, body)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		return
	}

	log.Debugf("Updating node %+v", delivery)
	nodeId, err := ukama.ValidateNodeId(body.NodeId)
	if err != nil {
		log.Errorf("Invalid Node ID format. Error %v", err)
		return
	}

	nd := &pb.Node{}
	if body.State != "" {
		if ist, ok := pb.NodeState_value[body.State]; ok {
			nd.State = pb.NodeState(ist)
		} else {
			log.Errorf("Invalid node state %s", body.State)
			nd.State = pb.NodeState_UNDEFINED
			return
		}
	}
	if len(body.Name) != 0 {
		nd.Name = body.Name
	}

	_, err = q.networkClient.UpdateNode(ctx,
		&pb.UpdateNodeRequest{
			NodeId: nodeId.String(),
			Node:   nd,
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
