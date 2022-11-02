package queue

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	validations "github.com/ukama/ukama/systems/data-plan/base-rate/pkg/validations"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/data-plan/base-rate/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QueueListener struct {
	msgBusConn  msgbus.Consumer
	baseRateClient  pb.BaseRatesServiceClient
	grpcTimeout time.Duration
	serviceId   string
	serviceName string
	// keep it to be able to close it
	grpcConn *grpc.ClientConn
}
type BaseRateRegisteredBody struct {
	FileURL   string `json:"file_url"`
	EffectiveAt string `json:"effective_at"`
	SimType string `json:"simType"`
}
type QueueListenerConfig struct {
	config.BaseConfig `mapstructure:",squash"`
	Service           struct {
		Host    string        `default:"localhost:7070"`
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
		baseRateClient:  pb.NewBaseRatesServiceClient(conn),
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
	
	log.Debugf("Received message: %s", delivery.Body)
	rate := &BaseRateRegisteredBody{}
	err := json.Unmarshal(delivery.Body, rate)
	if err != nil {
		log.Errorf("Error unmarshaling message. Error %v", err)
		return
	}

	log.Debugf("Uploading rate %v", rate.SimType)
	_, err = q.baseRateClient.UploadBaseRates(ctx, &pb.UploadBaseRatesRequest{
		FileURL:rate.FileURL,
		EffectiveAt:rate.EffectiveAt,
		SimType:validations.ReqStrTopb(rate.SimType),
	}, grpc_retry.WithMax(3))

	if err != nil {
		log.Errorf("Failed to add upload rate '%s'. Error: %v", rate.SimType, err)
	} else {
		log.Infof("Rate of simType %s uploaded successefully", rate.SimType)
	}
	
}

func (q *QueueListener) Close() {
	q.msgBusConn.Close()
	q.grpcConn.Close()
}