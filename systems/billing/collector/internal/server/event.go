package server

import (
	"context"
	"fmt"
	"time"

	client "github.com/ukama/ukama/systems/billing/collector/internal/clients"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	subpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	simpb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/sirupsen/logrus"
)

// TODO: We need to think about retry policies for failing interaction between our backend and the upstream billing service
// provider

const (
	handlerTimeoutFactor = 3
)

type BillingCollectorEventServer struct {
	client client.BillingClient
	epb.UnimplementedEventNotificationServiceServer
}

func NewBillingCollectorEventServer(client client.BillingClient) *BillingCollectorEventServer {
	return &BillingCollectorEventServer{
		client: client,
	}
}

func (b *BillingCollectorEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {

	// Send usage event
	case "event.cloud.cdr.sim.usage":
		msg, err := unmarshalSimUsage(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleSimUsageEvent(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// Create customer
	case "event.cloud.registry.subscriber.create":
		msg, err := unmarshalSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleRegistrySubscriberCreateEvent(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// Update customer
	case "event.cloud.registry.subscriber.update":
		msg, err := unmarshalSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleRegistrySubscriberUpdateEvent(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// Delete customer
	case "event.cloud.registry.subscriber.delete":
		msg, err := unmarshalSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleRegistrySubscriberDeleteEvent(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// add or update subscrition to customer
	case "event.cloud.simmanager.package.activate":
		msg, err := unmarshalSim(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleSimManagerSetActivePackageForSimEvent(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func handleSimUsageEvent(key string, simUsage *epb.SimUsage, b *BillingCollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, simUsage)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	event := client.Event{
		//TODO: To be replaced by msgClient msgId
		TransactionId: fmt.Sprintf("%s%d", simUsage.Id, time.Now().Unix()),

		CustomerId:     simUsage.SubscriberId,
		SubscriptionId: simUsage.SimId,
		Code:           "data_usage",
		SentAt:         time.Now(),

		AdditionalProperties: map[string]string{
			"bytes_used": fmt.Sprint(simUsage.BytesUsed),
			"sim_id":     simUsage.SimId,
		},
	}

	log.Infof("Sending data usage event %v to billing server", event)

	return b.client.AddUsageEvent(ctx, event)
}

func handleRegistrySubscriberCreateEvent(key string, subscriber *subpb.Subscriber,
	b *BillingCollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, subscriber)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	customer := client.Customer{
		Id:      subscriber.SubscriberId,
		Name:    subscriber.FirstName,
		Email:   subscriber.Email,
		Address: subscriber.Address,
		Phone:   subscriber.PhoneNumber,
	}

	log.Infof("Sending subscriber create event %v to billing server", customer)

	customerBillingId, err := b.client.CreateCustomer(ctx, customer)
	if err != nil {
		return err
	}

	log.Infof("Successfuly registered customer. Id: %s", customerBillingId)

	return nil
}

func handleRegistrySubscriberUpdateEvent(key string, subscriber *subpb.Subscriber,
	b *BillingCollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, subscriber)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	customer := client.Customer{
		Name:    subscriber.FirstName,
		Email:   subscriber.Email,
		Address: subscriber.Address,
		Phone:   subscriber.PhoneNumber,
	}

	log.Infof("Sending subscriber update event %v to billing", customer)

	customerBillingId, err := b.client.UpdateCustomer(ctx, customer)
	if err != nil {
		return err
	}

	log.Infof("Successfuly updated customer %v", customerBillingId)

	return nil
}

func handleRegistrySubscriberDeleteEvent(key string, subscriber *subpb.Subscriber,
	b *BillingCollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, subscriber)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	customerBillingId, err := b.client.DeleteCustomer(ctx, subscriber.SubscriberId)
	if err != nil {
		return err
	}

	log.Infof("Successfuly deleted customer %v", customerBillingId)

	return nil
}

func handleSimManagerSetActivePackageForSimEvent(key string, sim *simpb.Sim,
	b *BillingCollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	subscriptionId, err := b.client.TerminateSubscription(ctx, sim.Id)
	if err != nil {
		return err
	}

	log.Infof("Successfuly terminated previous subscription %v", subscriptionId)

	subscriptionAt := sim.Package.StartDate.AsTime()

	// Because the Plan object does not expose an external_plan_id, we need to use
	// our backend plan_id as billing provider's plan_code
	subscriptionInput := client.Subscription{
		Id:         sim.Id,
		CustomerId: sim.SubscriberId,
		// PlanCode:       sim.Package.PlanId,
		SubscriptionAt: &subscriptionAt,
	}

	log.Infof("Sending sim package activation event %v to billing server", subscriptionInput)

	subscriptionId, err = b.client.CreateSubscription(ctx, subscriptionInput)
	if err != nil {
		return err
	}

	log.Infof("Successfuly created new subscription %v", subscriptionId)

	return nil
}

func unmarshalSim(msg *anypb.Any) (*simpb.Sim, error) {
	p := &simpb.Sim{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal sim manager's sim message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
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

func unmarshalSimUsage(msg *anypb.Any) (*epb.SimUsage, error) {
	p := &epb.SimUsage{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal SimUsage message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}
