/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/ukama"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/billing/collector/pkg/clients"
	client "github.com/ukama/ukama/systems/billing/collector/pkg/clients"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

// TODO: We need to think about retry policies for failing interaction between
// TODO: We have unmarshal methods in common/pb/gen/events for all the event messages. We should use those.
// our backend and the upstream billing service provider.

const (
	handlerTimeoutFactor      = 3
	defaultChargeModel        = "package"
	defaultCurrency           = "USD"
	postpaidBillingInterval   = "monthly"
	prepaidBillingInterval    = "yearly"
	testBillingInterval       = "weekly"
	DefaultBillableMetricCode = "data_usage"
)

type BillableMetric struct {
	Id   string
	Code string
}

type CollectorEventServer struct {
	orgName string
	orgId   string
	bMetric BillableMetric
	client  client.BillingClient
	epb.UnimplementedEventNotificationServiceServer
}

func NewCollectorEventServer(orgName, orgId string, client client.BillingClient) (*CollectorEventServer, error) {
	log.Infof("Starting billing collector for org: %s", orgName)

	bm, err := initBillingDefaults(client, DefaultBillableMetricCode, orgName, orgId)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize billable metric: %w", err)
	}

	bMetric := BillableMetric{
		Id:   bm,
		Code: DefaultBillableMetricCode,
	}

	return &CollectorEventServer{
		orgName: orgName,
		orgId:   orgId,
		client:  client,
		bMetric: bMetric,
	}, nil
}

func (c *CollectorEventServer) EventNotification(ctx context.Context, e *epb.Event) (*epb.EventResponse, error) {
	log.Infof("Received a message with Routing key %s and Message %+v", e.RoutingKey, e.Msg)

	switch e.RoutingKey {
	// Update org subscription
	case msgbus.PrepareRoute(c.orgName, "event.cloud.global.{{ .Org}}.inventory.accounting.accounting.sync"): // or from orchestrator spin
		msg, err := unmarshalOrgSubscription(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleOrgSubscriptionEvent(e.RoutingKey, msg, c)
		if err != nil {
			return nil, err
		}

	// Create plan
	case msgbus.PrepareRoute(c.orgName, "event.cloud.local.{{ .Org}}.dataplan.package.package.create"):
		msg, err := unmarshalPackage(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleDataPlanPackageCreateEvent(e.RoutingKey, msg, c)
		if err != nil {
			return nil, err
		}

	// Create customer
	case msgbus.PrepareRoute(c.orgName, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.create"):
		msg, err := unmarshalAddSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleRegistrySubscriberCreateEvent(e.RoutingKey, msg, c)
		if err != nil {
			return nil, err
		}

	// Update customer
	case msgbus.PrepareRoute(c.orgName, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.update"):
		msg, err := unmarshalUpdateSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleRegistrySubscriberUpdateEvent(e.RoutingKey, msg, c)
		if err != nil {
			return nil, err
		}

	// Delete customer
	case msgbus.PrepareRoute(c.orgName, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.delete"):
		msg, err := unmarshalRemoveSubscriber(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleRegistrySubscriberDeleteEvent(e.RoutingKey, msg, c)
		if err != nil {
			return nil, err
		}

	// add subscrition to customer
	case msgbus.PrepareRoute(c.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate"):
		msg, err := unmarshalSimAllocation(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleSimManagerAllocateSimEvent(e.RoutingKey, msg, c)
		if err != nil {
			return nil, err
		}

	// update subscrition for customer
	case msgbus.PrepareRoute(c.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activepackage"):
		msg, err := unmarshalSimAcivePackage(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleSimManagerSetActivePackageForSimEvent(e.RoutingKey, msg, c)
		if err != nil {
			return nil, err
		}

		// Send usage event
	case msgbus.PrepareRoute(c.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.usage"):
		msg, err := unmarshalSimUsage(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleSimUsageEvent(e.RoutingKey, msg, c)
		if err != nil {
			return nil, err
		}

	// Terminate subscription
	case msgbus.PrepareRoute(c.orgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.expirepackage"):
		msg, err := unmarshalSimPackageExpire(e.Msg)
		if err != nil {
			return nil, err
		}

		err = handleSimManagerSimPackageExpireEvent(e.RoutingKey, msg, c)
		if err != nil {
			return nil, err
		}

	default:
		log.Errorf("No handler routing key %s", e.RoutingKey)
	}

	return &epb.EventResponse{}, nil
}

func handleOrgSubscriptionEvent(key string, usrAccountItems *epb.UserAccountingEvent,
	b *CollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, usrAccountItems)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	for _, accountItem := range usrAccountItems.Accounting {
		// Do we already have a plan with the same code for that org
		_, err := b.client.GetPlan(ctx, accountItem.Id)
		if err == nil {
			// The plan, therefore the subscription exists. Remove those
			log.Warnf("Plan with similar code %s was found", accountItem.Id)
			log.Infof("Removing plan and subscription associated with code: %s", accountItem.Id)

			_, err := b.client.TerminateSubscription(ctx, accountItem.Id)
			if err != nil {
				return fmt.Errorf("fail to terminate org subscription from item %s: %w",
					accountItem.Id, err)
			}

			_, err = b.client.TerminatePlan(ctx, accountItem.Id)
			if err != nil {
				return fmt.Errorf("fail to terminate org plan from item %s: %w",
					accountItem.Id, err)
			}
		}

		amount, err := strconv.ParseFloat(strings.TrimSpace(accountItem.OpexFee), 64)
		if err != nil {
			return fmt.Errorf("fail to parse opex fees %s for org subscription line with account item %s: %w",
				accountItem.OpexFee, accountItem.Id, err)
		}

		// Then we recreate the plan and the subscription
		newPlan := client.Plan{
			Name:        accountItem.Item + ": " + accountItem.Id,
			Code:        accountItem.Id,
			Interval:    postpaidBillingInterval,
			AmountCents: int(amount * 100),

			//TODO: update currency to pkg.Currency when the discussiion about currency is definetly settled.
			AmountCurrency: defaultCurrency,
			PayInAdvance:   false,
		}

		log.Infof("Sending plan create event %v with no charge to billing", newPlan)

		plan, err := b.client.CreatePlan(ctx, newPlan)
		if err != nil {
			return fmt.Errorf("fail to create org plan: %w", err)
		}

		log.Infof("New billing plan: %q", plan)
		log.Infof("Successfuly created org item plan from account item  %q", accountItem.Id)

		subscriptionInput := client.Subscription{
			Id:         accountItem.Id,
			CustomerId: b.orgId,
			PlanCode:   accountItem.Id,
		}

		log.Infof("Sending org subscription event %v to billing server", subscriptionInput)

		subscriptionId, err := b.client.CreateSubscription(ctx, subscriptionInput)
		if err != nil {
			return fmt.Errorf("fail to create org subscripton: %w", err)
		}

		log.Infof("New subscription created on billing server:  %q", subscriptionId)
		log.Infof("Successfuly created new subscription org item from account item: %q", accountItem.Id)
	}

	return nil
}

func handleSimUsageEvent(key string, simUsage *epb.EventSimUsage, b *CollectorEventServer) error {
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

		AdditionalProperties: map[string]any{
			"bytes_used": fmt.Sprint(simUsage.BytesUsed),
			"sim_id":     simUsage.SimId,
		},
	}

	log.Infof("Sending data usage event %v to billing server", event)

	return b.client.AddUsageEvent(ctx, event)
}

func handleDataPlanPackageCreateEvent(key string, pkg *epb.CreatePackageEvent, b *CollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, pkg)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	// Get the cost of the package per byte
	dataUnit := ukama.ParseDataUnitType(pkg.DataUnit)
	if dataUnit == ukama.DataUnitTypeUnknown {
		return fmt.Errorf("invalid data unit type: %s", pkg.DataUnit)
	}

	// Get the type of the package
	pkgType := ukama.ParsePackageType(pkg.Type)
	if pkgType == ukama.PackageTypeUnknown {
		return fmt.Errorf("invalid package type: %s", pkg.DataUnit)
	}

	var amount string

	pkgIntervall := prepaidBillingInterval

	switch pkgType {
	case ukama.PackageTypePostpaid:
		pkgIntervall = postpaidBillingInterval
		amount = strconv.FormatFloat(pkg.DataUnitCost, 'f', 2, 64)

	case ukama.PackageTypePrepaid:
		dataUnitCost := pkg.Amount / float64(pkg.DataVolume)
		amount = strconv.FormatFloat(dataUnitCost, 'f', 2, 64)
	}

	billableDataSize := math.Pow(1024, float64(dataUnit-1))

	charge := client.PlanCharge{
		BillableMetricID: b.bMetric.Id,
		ChargeModel:      defaultChargeModel,
		ChargeAmount:     amount,

		//TODO: update currency to pkg.Currency when the discussiion about currency is definetly settled.
		ChargeAmountCurrency: defaultCurrency,
		PackageSize:          int(billableDataSize),
	}

	newPlan := client.Plan{
		Name:        "Plan: " + pkg.Uuid,
		Code:        pkg.Uuid,
		Interval:    pkgIntervall,
		AmountCents: 0,

		//TODO: update currency to pkg.Currency when the discussiion about currency is definetly settled.
		AmountCurrency: defaultCurrency,
		PayInAdvance:   false,
	}

	log.Infof("Sending plan create event %v with charges %v to billing", newPlan, charge)

	plan, err := b.client.CreatePlan(ctx, newPlan, charge)
	if err != nil {
		return fmt.Errorf("fail to create subscriber plan: %w", err)
	}

	log.Infof("New billing plan: %q", plan)
	log.Infof("Successfuly created plan from package  %q", pkg.Uuid)

	return nil
}

func handleRegistrySubscriberCreateEvent(key string, subscriber *epb.AddSubscriber,
	b *CollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, subscriber)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	customer := client.Customer{
		Id:      subscriber.Subscriber.SubscriberId,
		Name:    subscriber.Subscriber.FirstName,
		Email:   subscriber.Subscriber.Email,
		Address: subscriber.Subscriber.Address,
		Phone:   subscriber.Subscriber.PhoneNumber,
		Type:    client.IndividualCustomerType,
	}

	log.Infof("Sending subscriber create event %v to billing server", customer)

	customerBillingId, err := b.client.CreateCustomer(ctx, customer)
	if err != nil {
		return fmt.Errorf("fail to create subscriber: %w", err)
	}

	log.Infof("New billing customer: %q", customerBillingId)
	log.Infof("Successfuly registered subscriber %q as billing customer", subscriber.Subscriber.SubscriberId)

	return nil
}

func handleRegistrySubscriberUpdateEvent(key string, subscriber *epb.UpdateSubscriber,
	b *CollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, subscriber)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	customer := client.Customer{
		Id:      subscriber.Subscriber.SubscriberId,
		Name:    subscriber.Subscriber.FirstName,
		Email:   subscriber.Subscriber.Email,
		Address: subscriber.Subscriber.Address,
		Phone:   subscriber.Subscriber.PhoneNumber,
	}

	log.Infof("Sending subscriber update event %v to billing", customer)

	customerBillingId, err := b.client.UpdateCustomer(ctx, customer)
	if err != nil {
		return fmt.Errorf("fail to update subscriber: %w", err)
	}

	log.Infof("Updated billing customer: %q",
		customerBillingId)

	log.Infof("Successfuly updated subscriber %q",
		subscriber.Subscriber.SubscriberId)

	return nil
}

func handleRegistrySubscriberDeleteEvent(key string, subscriber *epb.RemoveSubscriber,
	b *CollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, subscriber)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	customerBillingId, err := b.client.DeleteCustomer(ctx, subscriber.Subscriber.SubscriberId)
	if err != nil {
		return fmt.Errorf("fail to delete subscriber: %w", err)
	}

	log.Infof("Successfuly deleted customer %v", customerBillingId)

	return nil
}

func handleSimManagerAllocateSimEvent(key string, sim *epb.EventSimAllocation,
	b *CollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	// Because the Plan object does not expose an external_plan_id, we need to use
	// our backend plan_id as billing provider's plan_code
	subscriptionInput := client.Subscription{
		Id:         sim.Id,
		CustomerId: sim.SubscriberId,
		PlanCode:   sim.DataPlanId,
	}

	log.Infof("Sending subscripton creation event %v to billing server",
		subscriptionInput)

	subscriptionId, err := b.client.CreateSubscription(ctx, subscriptionInput)
	if err != nil {
		return fmt.Errorf("fail to create subscriber subscripton: %w", err)
	}

	log.Infof("New subscription created on billing server:  %q",
		subscriptionId)

	log.Infof("Successfuly created new subscription from sim: %q",
		sim.Id)

	return nil
}

func handleSimManagerSetActivePackageForSimEvent(key string, sim *epb.EventSimActivePackage,
	b *CollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	subscriptionId, err := b.client.TerminateSubscription(ctx, sim.Id)
	if err != nil {
		var e *clients.Error
		if !errors.As(err, &e) || (errors.As(err, &e) && e.Code != http.StatusNotFound) {
			return fmt.Errorf("fail to terminate subscriber's current subscripton: %w",
				err)
		}
	} else {
		log.Infof("Successfuly terminated previous subscription: %q",
			subscriptionId)
	}

	subscriptionInput := client.Subscription{
		Id:         sim.Id,
		CustomerId: sim.SubscriberId,
		PlanCode:   sim.PlanId,
	}

	log.Infof("Sending sim package activation event %v to billing server",
		subscriptionInput)

	subscriptionId, err = b.client.CreateSubscription(ctx, subscriptionInput)
	if err != nil {
		return fmt.Errorf("fail to create subscriber subscripton: %w", err)
	}

	log.Infof("New subscription created on billing server:  %q",
		subscriptionId)

	log.Infof("Successfuly created new subscription from sim: %q",
		sim.Id)

	return nil
}

func handleSimManagerSimPackageExpireEvent(key string, sim *epb.EventSimPackageExpire,
	b *CollectorEventServer) error {
	log.Infof("Keys %s and Proto is: %+v", key, sim)

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	subscriptionId, err := b.client.TerminateSubscription(ctx, sim.Id)
	if err != nil {
		return fmt.Errorf("fail to terminate subscriber subscripton: %w", err)
	}

	log.Infof("Successfuly terminated current subscription: %q",
		subscriptionId)

	return nil
}

func unmarshalOrgSubscription(msg *anypb.Any) (*epb.UserAccountingEvent, error) {
	p := &epb.UserAccountingEvent{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal SimUsage message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalSimAcivePackage(msg *anypb.Any) (*epb.EventSimActivePackage, error) {
	p := &epb.EventSimActivePackage{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal sim active package message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalPackage(msg *anypb.Any) (*epb.CreatePackageEvent, error) {
	p := &epb.CreatePackageEvent{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal package  message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalAddSubscriber(msg *anypb.Any) (*epb.AddSubscriber, error) {
	p := &epb.AddSubscriber{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal subscriber message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalUpdateSubscriber(msg *anypb.Any) (*epb.UpdateSubscriber, error) {
	p := &epb.UpdateSubscriber{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal subscriber message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalRemoveSubscriber(msg *anypb.Any) (*epb.RemoveSubscriber, error) {
	p := &epb.RemoveSubscriber{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("failed to Unmarshal subscriber message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalSimUsage(msg *anypb.Any) (*epb.EventSimUsage, error) {
	p := &epb.EventSimUsage{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal SimUsage message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalSimAllocation(msg *anypb.Any) (*epb.EventSimAllocation, error) {
	p := &epb.EventSimAllocation{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal SimAllocation message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}

func unmarshalSimPackageExpire(msg *anypb.Any) (*epb.EventSimPackageExpire, error) {
	p := &epb.EventSimPackageExpire{}

	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("Failed to Unmarshal SimAllocation message with : %+v. Error %s.",
			msg, err.Error())

		return nil, err
	}

	return p, nil
}

func initBillingDefaults(clt client.BillingClient, bmCode, orgName, orgId string) (string, error) {
	log.Infof("Initializing billing defaults")

	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	_, err := clt.GetCustomer(ctx, orgId)
	if err != nil {
		log.Warnf("Error while getting org billable account: %v", err)
		log.Infof("Creating org billable account: %s", orgId)

		_, err = createOrgCustomer(clt, orgId, orgName)
		if err != nil {
			return "", err
		}
	}

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

func createOrgCustomer(clt client.BillingClient, orgId, OrgName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), handlerTimeoutFactor*time.Second)
	defer cancel()

	customer := client.Customer{
		Id:   orgId,
		Name: OrgName,
		Type: client.CompanyCustomerType,

		// TODO: we might need additional fields such as Email, Address, Phone.
	}

	log.Infof("Sending org customer create event %v to billing server",
		customer)

	customerBillingId, err := clt.CreateCustomer(ctx, customer)
	if err != nil {
		return "", err
	}

	log.Infof("New org customer: %q",
		customerBillingId)

	log.Infof("Successfuly registered org %q as billing customer",
		orgId)

	return customerBillingId, nil
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

	log.Infof("Sending create request for billable metric %q to billing server",
		bMetric)

	bm, err := clt.CreateBillableMetric(ctx, bMetric)
	if err != nil {
		return "", err
	}

	log.Infof("Successfuly created billable metric. Id: %s", bm)

	return bm, nil
}
