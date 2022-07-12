package queue

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	pb "github.com/ukama/ukama/services/cloud/node/pb/gen"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/msgbus"
	pbmesh "github.com/ukama/ukama/services/common/pb/gen/ukamaos/mesh"

	"github.com/ukama/ukama/services/common/ukama"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type QueueListener struct {
	msgBusConn  msgbus.Consumer
	nodeClient  pb.NodeServiceClient
	grpcTimeout time.Duration
	serviceId   string
	serviceName string
	// keep it to be able to close it
	grpcConn *grpc.ClientConn
}

type QueueListenerConfig struct {
	config.BaseConfig `mapstructure:",squash"`
	Service           struct {
		Host    string        `default:"localhost:9090"`
		Timeout time.Duration `default:"3s"`
	}
	Queue   config.Queue
	Metrics config.Metrics
}

func NewQueueListener(conf QueueListenerConfig, serviceName string, serviceId string) (*QueueListener, error) {
	ctx, cancel := context.WithTimeout(context.Background(), conf.Service.Timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, conf.Service.Host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	client, err := msgbus.NewConsumerClient(conf.Queue.Uri)
	if err != nil {
		return nil, err
	}

	return &QueueListener{
		nodeClient:  pb.NewNodeServiceClient(conn),
		msgBusConn:  client,
		grpcTimeout: conf.Service.Timeout,
		serviceId:   serviceId,
		serviceName: serviceName,
		grpcConn:    conn,
	}, nil
}

func (q *QueueListener) StartQueueListening() (err error) {
	err = q.msgBusConn.SubscribeToServiceQueue(q.serviceName, msgbus.DeviceQ.Exchange,
		[]msgbus.RoutingKey{msgbus.DeviceConnectedRoutingKey}, q.serviceId, q.incomingMessageHandler)
	if err != nil {
		log.Errorf("Error subscribing for a queue messages. Error: %+v", err)
		return err
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	log.Info("Shutting down...")
	q.msgBusConn.Close()

	return nil
}

func (q *QueueListener) incomingMessageHandler(delivery amqp.Delivery, done chan<- bool) {
	ctx, cancel := context.WithTimeout(context.Background(), q.grpcTimeout)
	defer cancel()

	switch delivery.RoutingKey {
	case string(msgbus.DeviceConnectedRoutingKey):
		q.processDeviceConnectedMsg(ctx, delivery)

	default:
		log.Warning("No handler for routing key ", delivery.RoutingKey)
	}

	done <- true
}

func (q *QueueListener) processDeviceConnectedMsg(ctx context.Context, delivery amqp.Delivery) {
	link := &pbmesh.Link{}
	err := proto.Unmarshal(delivery.Body, link)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		MessageProcessFailedMetric()
		return
	}

	log.Debugf("Updating node %+v", delivery)
	nodeId, err := ukama.ValidateNodeId(link.GetNodeId())
	if err != nil {
		log.Errorf("Invalid Node ID format. Error %v", err)
		MessageProcessFailedMetric()
		return
	}

	_, err = q.nodeClient.UpdateNodeState(ctx,
		&pb.UpdateNodeStateRequest{
			NodeId: nodeId.String(),
			State:  pb.NodeState_ONBOARDED},
		grpc_retry.WithMax(3))

	if err != nil {
		log.Errorf("Failed to update node %s status. Error: %v", nodeId.String(), err)
		MessageProcessFailedMetric()
	} else {
		MessageProcessedMetric()
		log.Infof("Node %s updated successefully", nodeId.String())
	}
}

func (q *QueueListener) Close() {
	q.msgBusConn.Close()
	q.grpcConn.Close()
}
