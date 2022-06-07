package queue

import (
	"context"
	"encoding/json"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	pb "github.com/ukama/ukama/services/cloud/org/pb/gen"
	"github.com/ukama/ukama/services/cloud/org/pkg"
	"github.com/ukama/ukama/services/common/config"
	"github.com/ukama/ukama/services/common/msgbus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type QueueListener struct {
	msgBusConn  msgbus.Consumer
	orgClient   pb.OrgServiceClient
	grpcTimeout time.Duration
	serviceId   string
}

type QueueListenerConfig struct {
	Registry struct {
		Host    string `default:"localhost:9090"`
		Timeout time.Duration
	}
	Queue   config.Queue
	Metrics config.Metrics
}

type UserRegisteredBody struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func NewQueueListener(conf QueueListenerConfig, serviceName string, serviceId string) (*QueueListener, error) {
	client, err := msgbus.NewConsumerClient(conf.Queue.Uri)
	if err != nil {
		return nil, err
	}

	registryConn, err := grpc.Dial(conf.Registry.Host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	return &QueueListener{
		orgClient:   pb.NewOrgServiceClient(registryConn),
		msgBusConn:  client,
		grpcTimeout: conf.Registry.Timeout,
		serviceId:   serviceId,
	}, nil
}

func (q *QueueListener) StartQueueListening() (err error) {

	err = q.msgBusConn.SubscribeToServiceQueue(pkg.ServiceName+"-listener", msgbus.DeviceQ.Exchange,
		[]msgbus.RoutingKey{msgbus.UserRegisteredRoutingKey}, q.serviceId, q.incomingMessageHandler)
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
	_, err = q.orgClient.AddOrg(ctx, &pb.AddOrgRequest{
		Org: &pb.Organization{
			Name:  user.Id,
			Owner: user.Id,
		},
	}, grpc_retry.WithMax(3))

	if err != nil {
		log.Errorf("Failed to add organization '%s'. Error: %v", user.Id, err)
	} else {
		log.Infof("Organization %s added successefully", user.Id)
	}
}
