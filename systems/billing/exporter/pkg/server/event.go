package server

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ukama/telna/cdr/pb/gen"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	lago "github.com/getlago/lago-go-client"
	log "github.com/sirupsen/logrus"
)

const (
	simUsageRoutingKey   = "event.cloud.cdr.sim.usage"
	handlerTimeoutFactor = 3
)

type BillingExporterEventServer struct {
	c *lago.Client
	epb.UnimplementedEventNotificationServiceServer
}

func NewBillingExporterEventServer(lagoHost, lagoAPIKey string, lagoPort uint) *BillingExporterEventServer {
	log.Warnf("API KEY: %s", lagoAPIKey)

	lagoBaseURL := fmt.Sprintf("http://%s:%d", lagoHost, lagoPort)

	return &BillingExporterEventServer{
		c: lago.New().SetBaseURL(lagoBaseURL).SetApiKey(lagoAPIKey).SetDebug(true),
	}
}

func (b *BillingExporterEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case simUsageRoutingKey:
		msg, err := unmarshalAgentIncommingUsage(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudCdrSimUsage(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalAgentIncommingUsage(msg *anypb.Any) (*pb.SimUsage, error) {
	p := &pb.SimUsage{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func handleEventCloudCdrSimUsage(key string, msg *pb.SimUsage, b *BillingExporterEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	eventInput := &lago.EventInput{
		TransactionID:      fmt.Sprintf("%s%d", msg.Id, time.Now().Unix()),
		ExternalCustomerID: msg.SubscriberID,
		// ExternalSubscriptionID: "1dbe81ce-b092-401c-a00b-314292e17a98",
		Code:      "data_usage",
		Timestamp: time.Now().Unix(),
		Properties: map[string]string{
			"bytes_used": fmt.Sprint(msg.BytesUsed),
		},
	}

	log.Infof("Sending usage event %v to lago billing", eventInput)

	err := b.c.Event().Create(ctx, eventInput)
	if err != nil {
		log.Errorf("Error while sending data usage to lago instance: %v, %v, %v", err.Err, err.HTTPStatusCode, err.Msg)

		return fmt.Errorf("Failed to send data usage to lago instance: %v", err)
	}

	return nil
}
