package listener

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/msgbus"
	commonpb "github.com/ukama/ukama/systems/common/pb/gen/ukamaos/mesh"
	pb "github.com/ukama/ukama/systems/messaging/net/pb/gen"
	regpb "github.com/ukama/ukama/systems/registry/pb/gen"
	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

const POD_NAME_ENV_VAR = "POD_NAME"

type listener struct {
	msgBusConn  rabbitmq.Consumer
	nnsClient   pb.NnsClient
	grpcTimeout int
	serviceId   string
	registry    regpb.RegistryServiceClient
}

type ListenerConfig struct {
	config.BaseConfig `mapstructure:",squash"`
	Queue             config.Queue
	GrpcTimeout       int
	NnsGrpcHost       string
	RegistryHost      string
}

func NewLiseterConfig() *ListenerConfig {
	return &ListenerConfig{
		Queue: config.Queue{
			Uri: "amqp://guest:guest@localhost:5672/",
		},
		NnsGrpcHost:  "localhost:9090",
		RegistryHost: "localhost:9090",
		GrpcTimeout:  3,
	}
}

func StartListener(config *ListenerConfig) {
	client, err := rabbitmq.NewConsumer(config.Queue.Uri, amqp.Config{},
		rabbitmq.WithConsumerOptionsLogger(logrus.WithField("service", "rabbitmq")))
	if err != nil {
		logrus.Fatalf("error creating queue consumer. Error: %s", err.Error())
	}

	nnsConn := newGrpcConnection(config.NnsGrpcHost, config.GrpcTimeout)
	regConn := newGrpcConnection(config.RegistryHost, config.GrpcTimeout)

	logrus.Infof("Creating listener. Queue: %s. Nns: %s",
		config.Queue.Uri[strings.LastIndex(config.Queue.Uri, "@"):], config.NnsGrpcHost)
	l := listener{
		nnsClient:   pb.NewNnsClient(nnsConn),
		registry:    regpb.NewRegistryServiceClient(regConn),
		msgBusConn:  client,
		grpcTimeout: config.GrpcTimeout,
		serviceId:   os.Getenv(POD_NAME_ENV_VAR),
	}
	l.startQueueListening()
}

func newGrpcConnection(nnsGrpcHost string, grpcTimeout int) *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(grpcTimeout)*time.Second)
	defer cancel()
	nnsConn, err := grpc.DialContext(ctx, nnsGrpcHost, grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Could not connect: %v", err)
	}

	return nnsConn
}

func (l listener) startQueueListening() {
	err := l.msgBusConn.StartConsuming(l.incomingMessageHandler, "net-listener", []string{string(msgbus.DeviceConnectedRoutingKey)},
		rabbitmq.WithConsumeOptionsBindingExchangeName(msgbus.DefaultExchange),
		rabbitmq.WithConsumeOptionsConsumerName(l.serviceId),
		rabbitmq.WithConsumeOptionsQueueDurable,
		rabbitmq.WithConsumeOptionsBindingExchangeDurable,
		rabbitmq.WithConsumeOptionsBindingExchangeKind(amqp.ExchangeTopic),
	)

	if err != nil {
		logrus.Fatalf("Error subscribing for a queue messages. Error: %+v", err)
	}

	logrus.Info("Listening for messages...")
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}

func (l *listener) incomingMessageHandler(delivery rabbitmq.Delivery) rabbitmq.Action {
	logrus.Infof("Received message: %s", delivery.Body)

	var link commonpb.Link
	err := proto.Unmarshal(delivery.Body, &link)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message. Error: %+v", err)
		return rabbitmq.NackDiscard
	}

	logrus.Infof("Getting org and network for %s", link.GetNodeId())
	orgName, network, err := l.getOrgAndNetwork(*link.NodeId)
	if err != nil {
		logrus.Errorf("Failed to get org and network. Error: %+v", err)
		logrus.Warningf("Node id %s won't have org and network info", link.GetNodeId())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.grpcTimeout)*time.Second)
	defer cancel()
	_, err = l.nnsClient.Set(ctx, &pb.SetRequest{
		NodeId:  link.GetNodeId(),
		Ip:      link.GetIp(),
		OrgName: orgName,
		Network: network,
	}, grpc_retry.WithMax(3))

	if err != nil {
		logrus.Errorf("Failed to set node IP. Error: %+v", err)
		return rabbitmq.NackRequeue
	}
	logrus.Infof("Node %s IP set to %s", link.GetNodeId(), link.GetIp())
	return rabbitmq.Ack
}

func (l listener) getOrgAndNetwork(nodeId string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.grpcTimeout)*time.Second)
	defer cancel()
	r, err := l.registry.GetNode(ctx, &regpb.GetNodeRequest{
		NodeId: nodeId,
	})

	if err != nil {
		logrus.Errorf("Failed to get node from registry. Error: %+v", err)
		return "", "", errors.Wrap(err, "error getting node")
	}
	return r.Org.Name, r.Network.Name, nil
}
