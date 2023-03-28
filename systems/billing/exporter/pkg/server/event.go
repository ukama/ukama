package server

import (
	"context"
	"fmt"
	"time"

	lago "github.com/getlago/lago-go-client"
	operatorPb "github.com/ukama/telna/cdr/pb/gen"
	client "github.com/ukama/ukama/systems/billing/exporter/pkg/clients"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	subpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	smpb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/sirupsen/logrus"
)

// TODO: We need to think about retry policies for failing interaction between our backend and the upstream billing service
// provider

const (
	handlerTimeoutFactor = 3
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
	case "event.cloud.cdr.sim.usage":
		msg, err := unmarshalOperatorSimUsage(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudCdrSimUsage(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// Create customer
	case "event.cloud.registry.subscriber.create":
		msg, err := unmarshalSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudRegistrySubscriberCreate(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// update customer
	case "event.cloud.registry.subscriber.update":
		msg, err := unmarshalSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudRegistrySubscriberUpdate(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// delete customer
	case "event.cloud.registry.subscriber.delete":
		msg, err := unmarshalSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudRegistrySubscriberDelete(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// add or update subscrition to customer
	case "event.cloud.simmanager.package.activate":
		msg, err := unmarshalSim(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleEventCloudSimManagerSetActivePackageForSim(e.RoutingKey, msg, b)
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
		log.Errorf("Failed to Unmarshal operator SimUsage message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}

func handleEventCloudCdrSimUsage(key string, simUsage *operatorPb.SimUsage, b *BillingExporterEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, simUsage)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	eventInput := &lago.EventInput{
		// To be replaced by msgClient msgId
		TransactionID: fmt.Sprintf("%s%d", simUsage.Id, time.Now().Unix()),

		ExternalCustomerID: simUsage.SubscriberId,
		// ExternalSubscriptionID: msg.SimId,
		Code:      "data_usage",
		Timestamp: time.Now().Unix(),
		Properties: map[string]string{
			"bytes_used": fmt.Sprint(simUsage.BytesUsed),
			"sim_id":     simUsage.SimId,
		},
	}

	log.Infof("Sending operator data usage event %v to billing server", eventInput)

	err := b.client.L.Event().Create(ctx, eventInput)
	if err != nil {
		log.Errorf("Error while sending operator data usage event to billing server: %v, %v, %v",
			err.Err, err.HTTPStatusCode, err.Msg)

		return fmt.Errorf("failed to send operator data usage event to billing server: %v", err)
	}

	return nil
}

func unmarshalSubscriber(msg *anypb.Any) (*subpb.Subscriber, error) {
	p := &subpb.Subscriber{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal subscriber message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func handleEventCloudRegistrySubscriberCreate(key string, subscriber *subpb.Subscriber, b *BillingExporterEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, subscriber)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	customerInput := &lago.CustomerInput{
		ExternalID:   subscriber.SubscriberId,
		Name:         subscriber.FirstName,
		Email:        subscriber.Email,
		AddressLine1: subscriber.Address,
		Phone:        subscriber.PhoneNumber,
	}

	log.Infof("Sending subscriber create event %v to billing server", customerInput)

	customer, err := b.client.L.Customer().Create(ctx, customerInput)
	if err != nil {
		log.Errorf("Error while sending subscriber creation event to billing server: %v, %v, %v",
			err.Err, err.HTTPStatusCode, err.Msg)

		return fmt.Errorf("failed to send subscriber creation event to billing server: %v", err)
	}

	log.Infof("Successfuly registered customer %v", customer)

	return nil
}

func handleEventCloudRegistrySubscriberUpdate(key string, subscriber *subpb.Subscriber, b *BillingExporterEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, subscriber)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	customerInput := &lago.CustomerInput{
		ExternalID:   subscriber.SubscriberId,
		Email:        subscriber.Email,
		AddressLine1: subscriber.Address,
		Phone:        subscriber.PhoneNumber,
	}

	log.Infof("Sending subscriber update event %v to billing", customerInput)

	customer, err := b.client.L.Customer().Update(ctx, customerInput)
	if err != nil {
		log.Errorf("Error while sending subscriber update event to billing server: %v, %v, %v",
			err.Err, err.HTTPStatusCode, err.Msg)

		return fmt.Errorf("failed to send subscriber update evetn to billing server: %v", err)
	}

	log.Infof("Successfuly updated customer %v", customer)

	return nil
}

func handleEventCloudRegistrySubscriberDelete(key string, subscriber *subpb.Subscriber, b *BillingExporterEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, subscriber)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	customer, err := b.client.L.Customer().Delete(ctx, subscriber.SubscriberId)
	if err != nil {
		log.Errorf("Error while sending subscriber deletion event to billing server: %v, %v, %v",
			err.Err, err.HTTPStatusCode, err.Msg)

		return fmt.Errorf("failed to send subscriber deletion evetn to billing server: %v", err)
	}

	log.Infof("Successfuly deleted customer %v", customer)

	return nil
}

func unmarshalSim(msg *anypb.Any) (*smpb.Sim, error) {
	p := &smpb.Sim{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal sim manager's sim message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func handleEventCloudSimManagerSetActivePackageForSim(key string, sim *smpb.Sim, b *BillingExporterEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	// Terminate the previous subscription for that sim (plan)
	_, err := b.client.L.Subscription().Terminate(ctx, sim.Id)
	if err != nil {
		log.Errorf("Error while sending subscription termination event to billing server: %v, %v, %v",
			err.Err, err.HTTPStatusCode, err.Msg)
	}

	subscriptionAt := sim.Package.StartDate.AsTime()

	// Because the Plan object does not expose an external_plan_id, we need to use
	// our backend plan_id as billing provider's plan_code
	subscriptionInput := &lago.SubscriptionInput{
		ExternalID:         sim.Id,
		ExternalCustomerID: sim.SubscriberId,
		PlanCode:           sim.Package.PlanId,
		SubscriptionAt:     &subscriptionAt,
	}

	log.Infof("Sending sim package activation event %v to billing server", subscriptionInput)

	subscription, err := b.client.L.Subscription().Create(ctx, subscriptionInput)
	if err != nil {
		log.Errorf("Error while sending subscription creation event to billing server: %v, %v, %v",
			err.Err, err.HTTPStatusCode, err.Msg)

		return fmt.Errorf("failed to send subscription creation event to billing server: %v", err)
	}

	log.Infof("Successfuly created subscription %v", subscription)

	return nil
}
