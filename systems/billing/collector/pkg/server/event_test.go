/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ukama/ukama/systems/billing/collector/mocks"
	"github.com/ukama/ukama/systems/billing/collector/pkg/server"
	"github.com/ukama/ukama/systems/common/msgbus"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

const (
	OrgName = "testOrg"
	OrgId   = "592f7a8e-f318-4d3a-aab8-8d4187cde7f9"

	sessionId = 22
	bmId      = "e044081b-fbbe-45e9-8f78-0f9c0f112977"
	custId    = "e231a7cd-03f6-470a-9e8c-e02f54f9b415"
	planId    = "0f8be763-8bd6-406d-9d82-158d7f1a2140"
)

func TestCollectorEventServer_NewCollectorEventServer(t *testing.T) {
	t.Run("CreateCustomerFaillure", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}

		billingClient.On("GetCustomer", mock.Anything,
			OrgId).Return("", errors.New("fail to get customer")).Once()

		billingClient.On("CreateCustomer", mock.Anything,
			mock.Anything).Return("", errors.New("fail to create customer")).Once()

		_, err := server.NewCollectorEventServer(OrgName, OrgId, billingClient)
		assert.Error(t, err)
	})

	t.Run("CreateBillableMetricFaillure", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}

		billingClient.On("GetCustomer", mock.Anything,
			OrgId).Return("", errors.New("faillure")).Once()

		billingClient.On("CreateCustomer", mock.Anything,
			mock.Anything).Return(custId, nil).Once()

		billingClient.On("GetBillableMetricId", mock.Anything,
			server.DefaultBillableMetricCode).Return("", errors.New("faillure")).Once()

		billingClient.On("CreateBillableMetric", mock.Anything,
			mock.Anything).Return("", errors.New("faillure")).Once()

		_, err := server.NewCollectorEventServer(OrgName, OrgId, billingClient)
		assert.Error(t, err)
	})

	t.Run("BillableMetricCreated", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}
		routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.foo.bar.baz.buz")

		billingClient.On("GetCustomer", mock.Anything,
			OrgId).Return("", errors.New("faillure")).Once()

		billingClient.On("CreateCustomer", mock.Anything,
			mock.Anything).Return(custId, nil).Once()

		billingClient.On("GetBillableMetricId", mock.Anything,
			server.DefaultBillableMetricCode).Return("", errors.New("faillure")).Once()

		billingClient.On("CreateBillableMetric", mock.Anything,
			mock.Anything).Return(bmId, nil).Once()

		s, err := server.NewCollectorEventServer(OrgName, OrgId, billingClient)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        nil,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})
}

func TestCollectorEventServer_HandleOrgSubscriptionEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.global.{{ .Org}}.inventory.accounting.accounting.sync")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	billingClient.On("GetCustomer", mock.Anything,
		OrgId).Return(custId, nil).Once()

	s, err := server.NewCollectorEventServer(OrgName, OrgId, billingClient)
	assert.NoError(t, err)

	t.Run("PlanFoundButSubscriptionNotTerminated", func(t *testing.T) {
		billingClient.On("GetPlan", mock.Anything, mock.Anything).Return("", nil).Once()
		billingClient.On("TerminateSubscription", mock.Anything, mock.Anything).
			Return("", errors.New("some unexpected error")).Once()

		userAccountItems := epb.UserAccountingEvent{
			UserId: OrgId,
			Accounting: []*epb.UserAccounting{
				&epb.UserAccounting{
					Id: planId,
				},
			},
		}

		anyE, err := anypb.New(&userAccountItems)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("PlanFoundAndNotTerminated", func(t *testing.T) {
		billingClient.On("GetPlan", mock.Anything, mock.Anything).Return("", nil).Once()
		billingClient.On("TerminateSubscription", mock.Anything, mock.Anything).
			Return("", nil).Once()
		billingClient.On("TerminatePlan", mock.Anything, mock.Anything).
			Return("", errors.New("some unexpected error")).Once()

		userAccountItems := epb.UserAccountingEvent{
			UserId: OrgId,
			Accounting: []*epb.UserAccounting{
				&epb.UserAccounting{
					Id: planId,
				},
			},
		}

		anyE, err := anypb.New(&userAccountItems)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("OpexFeeNotValid", func(t *testing.T) {
		billingClient.On("GetPlan", mock.Anything, mock.Anything).Return("", errors.New("Not found")).Once()

		userAccountItems := epb.UserAccountingEvent{
			UserId: OrgId,
			Accounting: []*epb.UserAccounting{
				&epb.UserAccounting{
					Id:      planId,
					OpexFee: "lol",
				},
			},
		}

		anyE, err := anypb.New(&userAccountItems)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("NewPlanNotCreated", func(t *testing.T) {
		billingClient.On("GetPlan", mock.Anything, mock.Anything).Return("", errors.New("Not found")).Once()
		billingClient.On("CreatePlan", mock.Anything, mock.Anything).
			Return("", errors.New("some unexpected error occured")).Once()

		userAccountItems := epb.UserAccountingEvent{
			UserId: OrgId,
			Accounting: []*epb.UserAccounting{
				&epb.UserAccounting{
					Id:      planId,
					Item:    "SAS subscription",
					OpexFee: "199",
				},
			},
		}

		anyE, err := anypb.New(&userAccountItems)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("NewSubscriptionNotCreated", func(t *testing.T) {
		billingClient.On("GetPlan", mock.Anything, mock.Anything).Return("", errors.New("Not found")).Once()
		billingClient.On("CreatePlan", mock.Anything, mock.Anything).
			Return("", nil).Once()

		billingClient.On("CreateSubscription", mock.Anything, mock.Anything).
			Return("", errors.New("some unexpected error occured")).Once()

		userAccountItems := epb.UserAccountingEvent{
			UserId: OrgId,
			Accounting: []*epb.UserAccounting{
				&epb.UserAccounting{
					Id:      planId,
					Item:    "SAS subscription",
					OpexFee: "199",
				},
			},
		}

		anyE, err := anypb.New(&userAccountItems)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("OrgSubscriptionCreated", func(t *testing.T) {
		billingClient.On("GetPlan", mock.Anything, mock.Anything).Return("", errors.New("Not found")).Once()
		billingClient.On("CreatePlan", mock.Anything, mock.Anything).
			Return("", nil).Once()
		billingClient.On("CreateSubscription", mock.Anything, mock.Anything).
			Return("", nil).Once()

		userAccountItems := epb.UserAccountingEvent{
			UserId: OrgId,
			Accounting: []*epb.UserAccounting{
				&epb.UserAccounting{
					Id:      planId,
					Item:    "SAS subscription",
					OpexFee: "199",
				},
			},
		}

		anyE, err := anypb.New(&userAccountItems)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("OrgSubscriptionEventNotSent", func(t *testing.T) {
		userAccountItems := epb.Notification{}

		anyE, err := anypb.New(&userAccountItems)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestCollectorEventServer_HandleCdrSimUsageEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.usage")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	billingClient.On("GetCustomer", mock.Anything,
		OrgId).Return(custId, nil).Once()

	s, err := server.NewCollectorEventServer(OrgName, OrgId, billingClient)

	assert.NoError(t, err)

	t.Run("SimUsageEventSent", func(t *testing.T) {
		billingClient.On("AddUsageEvent", mock.Anything, mock.Anything).Return(nil).Once()

		simUsage := epb.EventSimUsage{
			Id:           "b20c61f1-1c5a-4559-bfff-cd00f746697d",
			SubscriberId: "c214f255-0ed6-4aa1-93e7-e333658c7318",
			NetworkId:    "9fd07299-2826-4f8b-aea9-69da56440bec",
			OrgId:        "75ec112a-8745-49f9-ab64-1a37edade794",
			Type:         "test_simple",
			BytesUsed:    uint64(rand.Int63n(4096000)),
			SessionId:    sessionId,
			StartTime:    uint64(time.Now().Unix() - int64(rand.Intn(30000))),
			EndTime:      uint64(time.Now().Unix()),
		}

		anyE, err := anypb.New(&simUsage)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("SimUsageEventNotAdded", func(t *testing.T) {
		billingClient.On("AddUsageEvent", mock.Anything, mock.Anything).
			Return(errors.New("failed to send sim usage")).Once()

		simUsage := epb.EventSimUsage{}

		anyE, err := anypb.New(&simUsage)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("SimUsageEventNotSent", func(t *testing.T) {
		simUsage := epb.Notification{}

		anyE, err := anypb.New(&simUsage)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestCollectorEventServer_HandleDataPlanPackageCreateEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.dataplan.package.package.create")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	billingClient.On("GetCustomer", mock.Anything,
		OrgId).Return(custId, nil).Once()

	s, err := server.NewCollectorEventServer(OrgName, OrgId, billingClient)

	assert.NoError(t, err)

	t.Run("PackageCreateEventSent", func(t *testing.T) {
		billingClient.On("CreatePlan", mock.Anything, mock.Anything, mock.Anything).
			Return("da337d0e-5678-446f-95c3-e94ac27a93b3", nil).Once()

		pkg := epb.CreatePackageEvent{
			Uuid:        "b20c61f1-1c5a-4559-bfff-cd00f746697d",
			SimType:     "operator_data",
			OrgId:       "75ec112a-8745-49f9-ab64-1a37edade794",
			OwnerId:     "c214f255-0ed6-4aa1-93e7-e333658c7318",
			SmsVolume:   1000,
			DataVolume:  5000000,
			VoiceVolume: 500,
			Type:        "prepaid",
			DataUnit:    "MegaBytes",
			Flatrate:    true,
			Country:     "USA",
			Provider:    "ukama",
		}

		anyE, err := anypb.New(&pkg)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("DataunintNotValid", func(t *testing.T) {
		billingClient.On("CreatePlan", mock.Anything, mock.Anything, mock.Anything).
			Return("", errors.New("failed to send create package event")).Once()

		pkg := epb.CreatePackageEvent{
			Uuid:        "b20c61f1-1c5a-4559-bfff-cd00f746697d",
			SimType:     "operator_data",
			OrgId:       "75ec112a-8745-49f9-ab64-1a37edade794",
			OwnerId:     "c214f255-0ed6-4aa1-93e7-e333658c7318",
			SmsVolume:   1000,
			DataVolume:  5000000,
			VoiceVolume: 500,
			Type:        "prepaid",
			DataUnit:    "lol",
			Flatrate:    false,
			Country:     "USA",
			Provider:    "ukama",
		}

		anyE, err := anypb.New(&pkg)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("CreatePackageFaillure", func(t *testing.T) {
		billingClient.On("CreatePlan", mock.Anything, mock.Anything, mock.Anything).
			Return("", errors.New("failed to send create package event")).Once()

		pkg := epb.CreatePackageEvent{
			Uuid:        "b20c61f1-1c5a-4559-bfff-cd00f746697d",
			SimType:     "operator_data",
			OrgId:       "75ec112a-8745-49f9-ab64-1a37edade794",
			OwnerId:     "c214f255-0ed6-4aa1-93e7-e333658c7318",
			SmsVolume:   1000,
			DataVolume:  5000000,
			VoiceVolume: 500,
			Type:        "prepaid",
			DataUnit:    "MegaBytes",
			Flatrate:    false,
			Country:     "USA",
			Provider:    "ukama",
		}

		anyE, err := anypb.New(&pkg)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("CreatePackageEventNotSent", func(t *testing.T) {
		pkg := epb.Notification{}

		anyE, err := anypb.New(&pkg)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestCollectorEventServer_HandleRegistrySubscriberCreateEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.create")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	billingClient.On("GetCustomer", mock.Anything,
		OrgId).Return(custId, nil).Once()

	s, err := server.NewCollectorEventServer(OrgName, OrgId, billingClient)

	assert.NoError(t, err)

	t.Run("CreateCustomerEventSent", func(t *testing.T) {

		billingClient.On("CreateCustomer", mock.Anything, mock.Anything).
			Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

		subs := &upb.Subscriber{
			SubscriberId: "c214f255-0ed6-4aa1-93e7-e333658c7318",
			FirstName:    "John Doe",
			Email:        "john.doe@example.com",
			Address:      "This is my address",
			PhoneNumber:  "000111222",
		}

		subscriber := epb.AddSubscriber{
			Subscriber: subs,
		}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("CreateCustomerFaillure", func(t *testing.T) {
		billingClient.On("CreateCustomer", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send create customer event")).Once()

		subs := &upb.Subscriber{}

		subscriber := epb.AddSubscriber{
			Subscriber: subs,
		}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("CreateCustomerEventNotSent", func(t *testing.T) {
		subscriber := epb.Notification{}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestCollectorEventServer_HandleRegistrySubscriberUpdateEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.update")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	billingClient.On("GetCustomer", mock.Anything,
		OrgId).Return(custId, nil).Once()

	s, err := server.NewCollectorEventServer(OrgName, OrgId, billingClient)

	assert.NoError(t, err)

	t.Run("UpdateCustomerEventSent", func(t *testing.T) {
		billingClient.On("UpdateCustomer", mock.Anything, mock.Anything).
			Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

		subs := &upb.Subscriber{
			FirstName:   "Fox Doe",
			Email:       "Fox.doe@example.com",
			Address:     "This is my address",
			PhoneNumber: "000111222",
		}

		subscriber := epb.UpdateSubscriber{
			Subscriber: subs,
		}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("UpdateCustomerFaillure", func(t *testing.T) {
		billingClient.On("UpdateCustomer", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send update customer event")).Once()

		subs := &upb.Subscriber{}

		subscriber := epb.UpdateSubscriber{
			Subscriber: subs,
		}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("UpdateCustomerEventNotSent", func(t *testing.T) {
		subscriber := epb.Notification{}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestCollectorEventServer_HandleRegistrySubscriberDeleteEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.delete")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	billingClient.On("GetCustomer", mock.Anything,
		OrgId).Return(custId, nil).Once()

	s, err := server.NewCollectorEventServer(OrgName, OrgId, billingClient)

	assert.NoError(t, err)

	t.Run("DeleteCustomerEventSent", func(t *testing.T) {
		billingClient.On("DeleteCustomer", mock.Anything, mock.Anything).
			Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

		subs := &upb.Subscriber{
			SubscriberId: "c214f255-0ed6-4aa1-93e7-e333658c7318",
		}

		subscriber := epb.RemoveSubscriber{
			Subscriber: subs,
		}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("DeleteCustomerFaillure", func(t *testing.T) {
		billingClient.On("DeleteCustomer", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send delete customer event")).Once()

		subs := &upb.Subscriber{}

		subscriber := epb.RemoveSubscriber{
			Subscriber: subs,
		}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("DeleteCustomerEventNotSent", func(t *testing.T) {
		subscriber := epb.Notification{}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestCollectorEventServer_HandleSimManagerSimAllocationEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	billingClient.On("GetCustomer", mock.Anything,
		OrgId).Return(custId, nil).Once()

	s, err := server.NewCollectorEventServer(OrgName, OrgId, billingClient)

	assert.NoError(t, err)

	t.Run("AllocateSimEventSent", func(t *testing.T) {
		billingClient.On("CreateSubscription", mock.Anything, mock.Anything).
			Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

		planId := "f1ad4204-ab9e-4574-b6bb-bffcc104f8f9"

		sim := epb.EventSimAllocation{
			Id:           "b20c61f1-1c5a-4559-bfff-cd00f746697d",
			SubscriberId: "c214f255-0ed6-4aa1-93e7-e333658c7318",
			DataPlanId:   planId,
		}

		anyE, err := anypb.New(&sim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("AllocateSimInvalidEventSent", func(t *testing.T) {
		billingClient.On("CreateSubscription", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send create subscription event")).Once()

		sim := epb.EventSimAllocation{}

		anyE, err := anypb.New(&sim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("AllocateSimEventNotSent", func(t *testing.T) {
		billingClient.On("CreateSubscription", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send create subscription event")).Once()

		sim := epb.Notification{}

		anyE, err := anypb.New(&sim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestCollectorEventServer_HandleSimManagerSetActivePackageForSimEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activepackage")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	billingClient.On("GetCustomer", mock.Anything,
		OrgId).Return(custId, nil).Once()

	s, err := server.NewCollectorEventServer(OrgName, OrgId, billingClient)

	assert.NoError(t, err)

	t.Run("SetActivePackageEventSent", func(t *testing.T) {
		billingClient.On("TerminateSubscription", mock.Anything, "b20c61f1-1c5a-4559-bfff-cd00f746697d").
			Return("9fd07299-2826-4f8b-aea9-69da56440bec", nil).Once()

		billingClient.On("CreateSubscription", mock.Anything, mock.Anything).
			Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

		sim := epb.EventSimActivePackage{
			Id:               "b20c61f1-1c5a-4559-bfff-cd00f746697d",
			SubscriberId:     "c214f255-0ed6-4aa1-93e7-e333658c7318",
			PackageId:        "3c353228-34ce-42ac-8ce4-0d4abb90bd8e",
			PackageStartDate: timestamppb.New(time.Now()),
		}

		anyE, err := anypb.New(&sim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("TerminateSubscriptionFaillure", func(t *testing.T) {
		billingClient.On("TerminateSubscription", mock.Anything, "16befbda-250f-4a68-9cb4-31cccce3005e").
			Return("", errors.New("failed to send terminate subscription event")).Once()

		sim := epb.EventSimActivePackage{
			Id:               "16befbda-250f-4a68-9cb4-31cccce3005e",
			SubscriberId:     "c214f255-0ed6-4aa1-93e7-e333658c7318",
			PackageId:        "3c353228-34ce-42ac-8ce4-0d4abb90bd8e",
			PackageStartDate: timestamppb.New(time.Now()),
		}

		anyE, err := anypb.New(&sim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("CreateSubscriptionFaillure", func(t *testing.T) {
		billingClient.On("TerminateSubscription", mock.Anything, "b20c61f1-1c5a-4559-bfff-cd00f746697d").
			Return("9fd07299-2826-4f8b-aea9-69da56440bec", nil).Once()

		billingClient.On("CreateSubscription", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send create subscription event")).Once()

		sim := epb.EventSimActivePackage{
			Id:               "b20c61f1-1c5a-4559-bfff-cd00f746697d",
			SubscriberId:     "c214f255-0ed6-4aa1-93e7-e333658c7318",
			PackageId:        "3c353228-34ce-42ac-8ce4-0d4abb90bd8e",
			PackageStartDate: timestamppb.New(time.Now()),
		}

		anyE, err := anypb.New(&sim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})

	t.Run("WrongEventSent", func(t *testing.T) {
		sim := epb.Notification{}

		anyE, err := anypb.New(&sim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}
