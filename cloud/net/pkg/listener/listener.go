package listener

import (
	"context"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	pb "github.com/ukama/ukamaX/cloud/net/pb/gen"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/common/msgbus"
	commonpb "github.com/ukama/ukamaX/common/pb/gen/ukamaos/mesh"
	"github.com/wagslane/go-rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const POD_NAME_ENV_VAR = "POD_NAME"

type listener struct {
	msgBusConn  rabbitmq.Consumer
	nnsClient   pb.NnsClient
	grpcTimeout int
	serviceId   string
}

type ListenerConfig struct {
	config.BaseConfig `mapstructure:",squash"`
	Queue             config.Queue
	GrpcTimeout       int
	NnsGrpcHost       string
}

func NewLiseterConfig() *ListenerConfig {
	return &ListenerConfig{
		Queue: config.Queue{
			Uri: "amqp://guest:guest@localhost:5672/",
		},
		NnsGrpcHost: "localhost:9090",
		GrpcTimeout: 3,
	}
}

func StartListener(config *ListenerConfig) {
	client, err := rabbitmq.NewConsumer(config.Queue.Uri, amqp.Config{},
		rabbitmq.WithConsumerOptionsLogger(logrus.WithField("service", "rabbitmq")))
	if err != nil {
		logrus.Fatalf("error creating queue consumer. Error: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GrpcTimeout)*time.Second)
	defer cancel()
	nnsConn, err := grpc.DialContext(ctx, config.NnsGrpcHost, grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Could not connect: %v", err)
	}

	logrus.Infof("Creating listener. Queue: %s. Nns: %s",
		config.Queue.Uri[strings.LastIndex(config.Queue.Uri, "@"):], config.NnsGrpcHost)

	l := listener{
		nnsClient:   pb.NewNnsClient(nnsConn),
		msgBusConn:  client,
		grpcTimeout: config.GrpcTimeout,
		serviceId:   os.Getenv(POD_NAME_ENV_VAR),
	}
	l.startQueueListening()
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.grpcTimeout)*time.Second)
	defer cancel()
	_, err = l.nnsClient.Set(ctx, &pb.SetRequest{
		NodeId: link.GetUuid(),
		Ip:     link.GetIp(),
	}, grpc_retry.WithMax(2))

	if err != nil {
		logrus.Errorf("Failed to set node IP. Error: %+v", err)
		return rabbitmq.NackRequeue
	}
	logrus.Infof("Node %s IP set to %s", link.GetUuid(), link.GetIp())
	return rabbitmq.Ack
}
