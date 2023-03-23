package server

import (
	"context"
	"fmt"
	"time"

	lago "github.com/getlago/lago-go-client"
	operatorPb "github.com/ukama/telna/cdr/pb/gen"
	client "github.com/ukama/ukama/systems/billing/exporter/pkg/clients"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	subPb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/sirupsen/logrus"
)

const (
	simUsageRoutingKey         = "event.cloud.cdr.sim.usage"
	subscriberCreateRoutingKey = "event.cloud.registry.subscriber.create"
	handlerTimeoutFactor       = 3
)

type BillingExporterEventServer struct {
	client *client.LagoClient
	epb.UnimplementedEventNotificationServiceServer
}

func NewBillingExporterEventServer(client *client.LagoClient) *BillingExporterEventServer {
	return &BillingExporterEventServer{
		client: client,
	}
}

func (b *BillingExporterEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	case simUsageRoutingKey:
		msg, err := unmarshalOperatorSimUsage(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudCdrSimUsage(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	case subscriberCreateRoutingKey:
		msg, err := unmarshalSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudRegistrySubscriberCreate(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func unmarshalOperatorSimUsage(msg *anypb.Any) (*operatorPb.SimUsage, error) {
	p := &operatorPb.SimUsage{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func handleEventCloudCdrSimUsage(key string, msg *operatorPb.SimUsage, b *BillingExporterEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	eventInput := &lago.EventInput{
		TransactionID:      fmt.Sprintf("%s%d", msg.Id, time.Now().Unix()),
		ExternalCustomerID: msg.SubscriberId,
		// ExternalSubscriptionID: msg.SimId,
		Code:      "data_usage",
		Timestamp: time.Now().Unix(),
		Properties: map[string]string{
			"bytes_used": fmt.Sprint(msg.BytesUsed),
			"sim_id":     msg.SimId,
		},
	}

	log.Infof("Sending usage event %v to lago billing", eventInput)

	err := b.client.L.Event().Create(ctx, eventInput)
	if err != nil {
		log.Errorf("Error while sending data usage to lago instance: %v, %v, %v", err.Err, err.HTTPStatusCode, err.Msg)

		return fmt.Errorf("Failed to send data usage to lago instance: %v", err)
	}

	return nil
}

func unmarshalSubscriber(msg *anypb.Any) (*subPb.Subscriber, error) {
	p := &subPb.Subscriber{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal AddOrgRequest message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func handleEventCloudRegistrySubscriberCreate(key string, msg *subPb.Subscriber, b *BillingExporterEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, msg)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	customerInput := &lago.CustomerInput{
		ExternalID:   msg.SubscriberId,
		Name:         msg.FirstName,
		Email:        msg.Email,
		AddressLine1: msg.Address,
		Phone:        msg.PhoneNumber,
	}

	log.Infof("Sending usage event %v to billing", customerInput)

	customer, err := b.client.L.Customer().Create(ctx, customerInput)

	if err != nil {
		log.Errorf("Error while sending data usage to billing instance: %v, %v, %v", err.Err, err.HTTPStatusCode, err.Msg)

		return fmt.Errorf("Failed to send data usage to billing instance: %v", err)
	}

	log.Infof("Successfuly registered customer %v", customer)

	return nil
}
