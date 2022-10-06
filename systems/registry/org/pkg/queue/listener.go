package queue

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"github.com/ukama/ukama/systems/registry/org/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QueueListener struct {
	msgBusConn  msgbus.Consumer
	orgClient   pb.OrgServiceClient
	grpcTimeout time.Duration
	serviceId   string
	// keep it here to be able to close it in Close()
	grpcConn *grpc.ClientConn
}

type QueueListenerConfig struct {
	OrgService struct {
		Host    string        `default:"localhost:9090"`
		Timeout time.Duration `default:"2s"`
	}
	Queue   *config.Queue   `default:"{}"`
	Metrics *config.Metrics `default:"{}"`
}

type UserRegisteredBody struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func NewQueueListener(conf QueueListenerConfig, serviceId string) (*QueueListener, error) {
	log.Info("Starting queue listener")
	client, err := msgbus.NewConsumerClient(conf.Queue.Uri)
	if err != nil {
		return nil, err
	}

	log.Info("Connecting to org service")
	ctx, cancel := context.WithTimeout(context.Background(), conf.OrgService.Timeout)
	defer cancel()
	networkConn, err := grpc.DialContext(ctx, conf.OrgService.Host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	return &QueueListener{
		orgClient:   pb.NewOrgServiceClient(networkConn),
		msgBusConn:  client,
		grpcTimeout: conf.OrgService.Timeout,
		serviceId:   serviceId,
		grpcConn:    networkConn,
	}, nil
}

func (q *QueueListener) StartQueueListening() (err error) {

	err = q.msgBusConn.SubscribeToServiceQueue(pkg.ServiceName+"-listener", msgbus.DefaultExchange,
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

	ctx, cancel := context.WithTimeout(context.Background(), q.grpcTimeout)
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
	log.Debugf("Received message: %s", delivery.Body)
	user := &UserRegisteredBody{}
	err := json.Unmarshal(delivery.Body, user)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		return
	}

	log.Debugf("Adding org %v", user.Id)
	_, err = q.orgClient.Add(ctx, &pb.AddRequest{
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

func (q *QueueListener) Close() {
	q.msgBusConn.Close()
	q.grpcConn.Close()
}
