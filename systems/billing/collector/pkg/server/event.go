package server

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	client "github.com/ukama/ukama/systems/billing/collector/pkg/clients"
	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/ukama"
	subpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	log "github.com/sirupsen/logrus"
)

// TODO: We need to think about retry policies for failing interaction between our backend and the upstream billing service
// provider

const (
	handlerTimeoutFactor      = 3
	defaultChargeModel        = "package"
	defaultCurrency           = "USD"
	defaultBillingInterval    = "monthly"
	testBillingInterval       = "weekly"
	DefaultBillableMetricCode = "data_usage"
)

type BillableMetric struct {
	Id   string
	Code string
}

type BillingCollectorEventServer struct {
	org     string
	client  client.BillingClient
	bMetric BillableMetric
	epb.UnimplementedEventNotificationServiceServer
}

func NewBillingCollectorEventServer(org string, client client.BillingClient) *BillingCollectorEventServer {
	bm, err := initBillableMetric(client, DefaultBillableMetricCode)

	if err != nil {
		log.Fatalf("Failed to initialize billable metric: %v", err)
	}

	bMetric := BillableMetric{
		Id:   bm,
		Code: DefaultBillableMetricCode,
	}

	return &BillingCollectorEventServer{
		org:     org,
		client:  client,
		bMetric: bMetric,
	}
}

func (b *BillingCollectorEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {

	// Send usage event
	case msgbus.PrepareRoute(b.org, "event.cloud.local.{{ .Org}}.operator.cdr.sim.usage"):
		msg, err := unmarshalSimUsage(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleSimUsageEvent(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// Create plan
	case msgbus.PrepareRoute(b.org, "event.cloud.local.{{ .Org}}.dataplan.package.create"):
		msg, err := unmarshalPackage(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleDataPlanPackageCreateEvent(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// Create customer
	case msgbus.PrepareRoute(b.org, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.create"):
		msg, err := unmarshalSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleRegistrySubscriberCreateEvent(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// Update customer
	case msgbus.PrepareRoute(b.org, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.update"):
		msg, err := unmarshalSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleRegistrySubscriberUpdateEvent(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// Delete customer
	case msgbus.PrepareRoute(b.org, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.delete"):
		msg, err := unmarshalSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleRegistrySubscriberDeleteEvent(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// add subscrition to customer
	case msgbus.PrepareRoute(b.org, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate"):
		msg, err := unmarshalSimAllocation(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleSimManagerAllocateSimEvent(e.RoutingKey, msg, b)
		if err != nil {
			return nil, err
		}

	// update subscrition to customer
	case msgbus.PrepareRoute(b.org, "event.cloud.local.{{ .Org}}.subscriber.simmanager.package.activate"):
		msg, err := unmarshalSimAcivePackage(e.Msg)
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
		Code:           b.bMetric.Code,
		SentAt:         time.Now(),

		AdditionalProperties: map[string]string{
			"bytes_used": fmt.Sprint(simUsage.BytesUsed),
			"sim_id":     simUsage.SimId,
		},
	}

	log.Infof("Sending data usage event %v to billing server", event)

	return b.client.AddUsageEvent(ctx, event)
}

func handleDataPlanPackageCreateEvent(key string, pkg *epb.CreatePackageEvent, b *BillingCollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, pkg)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	// TODO: upstream billing provider fails on a DB constraint when pay in advance
	// is set to false (postpaid). Somwhow, false bool value from go is sent as null
	// to upstream DB. Need to investigate this between upstream go client and DB.
	// TODO updates: It seems like 0, false values are not sent by go client.

	// payAdvance := false
	// if ukama.ParsePackageType(pkg.Type) == ukama.PackageTypePrepaid {
	// payAdvance = true
	// }

	// Get the cost of the package per bytke
	dataUnit := ukama.ParseDataUnitType(pkg.DataUnit)
	if dataUnit == ukama.DataUnitTypeUnknown {
		return fmt.Errorf("invalid data unit type")
	}

	billableDataSize := math.Pow(1024, float64(dataUnit-1))
	amountCents := strconv.Itoa(int(pkg.DataUnitCost * 100))

	newPlan := client.Plan{
		Name:     "Plan " + pkg.Uuid,
		Code:     pkg.Uuid,
		Interval: testBillingInterval,

		// 0 values are not sent by the upstream billing provider client. see above Todos
		AmountCents: 1,

		AmountCurrency: defaultCurrency,

		// fails on false (postpaid). See abouve Todos
		PayInAdvance: true,

		BillableMetricID:     b.bMetric.Id,
		ChargeModel:          defaultChargeModel,
		ChargeAmountCents:    amountCents,
		ChargeAmountCurrency: defaultCurrency,
		PackageSize:          int(billableDataSize),
	}

	log.Infof("Sending plan create event %v to billing", newPlan)

	plan, err := b.client.CreatePlan(ctx, newPlan)
	if err != nil {
		return err
	}

	log.Infof("New billing plan: %q", plan)
	log.Infof("Successfuly created plan from package  %q", pkg.Uuid)

	return nil
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

	log.Infof("New billing customer: %q", customerBillingId)
	log.Infof("Successfuly registered subscriber %q as billing customer", subscriber.SubscriberId)

	return nil
}

func handleRegistrySubscriberUpdateEvent(key string, subscriber *subpb.Subscriber,
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

	log.Infof("Sending subscriber update event %v to billing", customer)

	customerBillingId, err := b.client.UpdateCustomer(ctx, customer)
	if err != nil {
		return err
	}

	log.Infof("Updated billing customer: %q", customerBillingId)
	log.Infof("Successfuly updated subscriber %q", subscriber.SubscriberId)

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

func handleSimManagerAllocateSimEvent(key string, sim *epb.SimAllocation,
	b *BillingCollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	// subscriptionAt := time.Now()

	// Because the Plan object does not expose an external_plan_id, we need to use
	// our backend plan_id as billing provider's plan_code
	subscriptionInput := client.Subscription{
		Id:         sim.Id,
		CustomerId: sim.SubscriberId,
		PlanCode:   sim.DataPlanId,
		// SubscriptionAt: &subscriptionAt,
	}

	log.Infof("Sending subscripton creation event %v to billing server", subscriptionInput)

	subscriptionId, err := b.client.CreateSubscription(ctx, subscriptionInput)
	if err != nil {
		return err
	}

	log.Infof("New subscription created on billing server:  %q", subscriptionId)
	log.Infof("Successfuly created new subscription from sim: %q", sim.Id)

	return nil
}

func handleSimManagerSetActivePackageForSimEvent(key string, sim *epb.SimActivePackage,
	b *BillingCollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	subscriptionId, err := b.client.TerminateSubscription(ctx, sim.Id)
	if err != nil {
		return err
	}

	log.Infof("Successfuly terminated previous subscription: %q", subscriptionId)

	// subscriptionAt := sim.PackageStartDate.AsTime()

	subscriptionInput := client.Subscription{
		Id:         sim.Id,
		CustomerId: sim.SubscriberId,
		PlanCode:   sim.PlanId,
		// SubscriptionAt: &subscriptionAt,
	}

	log.Infof("Sending sim package activation event %v to billing server", subscriptionInput)

	subscriptionId, err = b.client.CreateSubscription(ctx, subscriptionInput)
	if err != nil {
		return err
	}

	log.Infof("New subscription created on billing server:  %q", subscriptionId)
	log.Infof("Successfuly created new subscription from sim: %q", sim.Id)

	return nil
}

func unmarshalSimAcivePackage(msg *anypb.Any) (*epb.SimActivePackage, error) {
	p := &epb.SimActivePackage{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal sim active package message with : %+v. Error %s.", msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalPackage(msg *anypb.Any) (*epb.CreatePackageEvent, error) {
	p := &epb.CreatePackageEvent{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal package  message with : %+v. Error %s.", msg, err.Error())

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

func unmarshalSimAllocation(msg *anypb.Any) (*epb.SimAllocation, error) {
	p := &epb.SimAllocation{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal SimAllocation message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}

func initBillableMetric(clt client.BillingClient, bmCode string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	log.Infof("Sending get request for billable metric %q to billing server", bmCode)

	bmId, err := clt.GetBillableMetricId(ctx, bmCode)
	if err != nil {
		log.Warnf("Error while getting default billable metric: %v", err)
		log.Infof("Creating default billable metric: %s", bmCode)

		bmId, err = createBillableMetric(clt)
		if err != nil {
			return "", err
		}
	}

	log.Infof("Successfuly returning billable metric. Id: %s", bmId)
	return bmId, nil
}

func createBillableMetric(clt client.BillingClient) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	bMetric := client.BillableMetric{
		Name:        "Data Usage",
		Code:        DefaultBillableMetricCode,
		Description: "Data Usage Billable Metric",
		FieldName:   "bytes_used",
	}

	log.Infof("Sending create request for billable metric %q to billing server", bMetric)

	bm, err := clt.CreateBillableMetric(ctx, bMetric)
	if err != nil {
		return "", err
	}

	log.Infof("Successfuly created billable metric. Id: %s", bm)

	return bm, nil
}
