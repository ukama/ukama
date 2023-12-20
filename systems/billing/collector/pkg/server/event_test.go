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

	"github.com/ukama/ukama/systems/billing/collector/mocks"
	"github.com/ukama/ukama/systems/billing/collector/pkg/server"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	subpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
)

const OrgName = "testOrg"
const bmId = "e044081b-fbbe-45e9-8f78-0f9c0f112977"

func TestBillingCollectorEventServer_HandleCdrSimUsageEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.operator.cdr.sim.usage")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	s := server.NewBillingCollectorEventServer(OrgName, billingClient)

	t.Run("SimUsageEventSent", func(t *testing.T) {
		billingClient.On("AddUsageEvent", mock.Anything, mock.Anything).Return(nil).Once()

		simUsage := epb.SimUsage{
			Id:           "b20c61f1-1c5a-4559-bfff-cd00f746697d",
			SubscriberId: "c214f255-0ed6-4aa1-93e7-e333658c7318",
			NetworkId:    "9fd07299-2826-4f8b-aea9-69da56440bec",
			OrgId:        "75ec112a-8745-49f9-ab64-1a37edade794",
			Type:         "test_simple",
			BytesUsed:    uint64(rand.Int63n(4096000)),
			SessionId:    uuid.NewV4().String(),
			StartTime:    time.Now().Unix() - int64(rand.Intn(30000)),
			EndTime:      time.Now().Unix(),
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

	t.Run("SimUsageEventNotSent", func(t *testing.T) {
		billingClient.On("AddUsageEvent", mock.Anything, mock.Anything).
			Return(errors.New("failed to send sim usage")).Once()

		simUsage := epb.SimUsage{}

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

func TestBillingCollectorEventServer_HandleDataPlanPackageCreateEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.dataplan.package.package.create")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	s := server.NewBillingCollectorEventServer(OrgName, billingClient)

	t.Run("PackageCreateEventSent", func(t *testing.T) {
		billingClient.On("CreatePlan", mock.Anything, mock.Anything).
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

		assert.NoError(t, err)
	})

	t.Run("CreatePackageEventNotSent", func(t *testing.T) {
		billingClient.On("CreatePlan", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send create package event")).Once()

		pkg := epb.CreatePackageEvent{}

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

func TestBillingCollectorEventServer_HandleRegistrySubscriberCreateEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.create")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	s := server.NewBillingCollectorEventServer(OrgName, billingClient)

	t.Run("CreateCustomerEventSent", func(t *testing.T) {

		billingClient.On("CreateCustomer", mock.Anything, mock.Anything).
			Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

		subscriber := subpb.Subscriber{
			SubscriberId: "c214f255-0ed6-4aa1-93e7-e333658c7318",
			FirstName:    "John Doe",
			Email:        "john.doe@example.com",
			Address:      "This is my address",
			PhoneNumber:  "000111222",
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

	t.Run("CreateCustomerEventNotSent", func(t *testing.T) {
		billingClient.On("CreateCustomer", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send create customer event")).Once()

		subscriber := subpb.Subscriber{}

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

func TestBillingCollectorEventServer_HandleRegistrySubscriberUpdateEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.update")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	s := server.NewBillingCollectorEventServer(OrgName, billingClient)

	t.Run("UpdateCustomerEventSent", func(t *testing.T) {
		billingClient.On("UpdateCustomer", mock.Anything, mock.Anything).
			Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

		subscriber := subpb.Subscriber{
			FirstName:   "Fox Doe",
			Email:       "Fox.doe@example.com",
			Address:     "This is my address",
			PhoneNumber: "000111222",
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

	t.Run("UpdateCustomerEventNotSent", func(t *testing.T) {
		billingClient.On("UpdateCustomer", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send update customer event")).Once()

		subscriber := subpb.Subscriber{}

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

func TestBillingCollectorEventServer_HandleRegistrySubscriberDeleteEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.delete")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	s := server.NewBillingCollectorEventServer(OrgName, billingClient)

	t.Run("DeleteCustomerEventSent", func(t *testing.T) {
		billingClient.On("DeleteCustomer", mock.Anything, mock.Anything).
			Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

		subscriber := subpb.Subscriber{
			SubscriberId: "c214f255-0ed6-4aa1-93e7-e333658c7318",
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

	t.Run("DeleteCustomerEventNotSent", func(t *testing.T) {

		billingClient.On("DeleteCustomer", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send delete customer event")).Once()

		subscriber := subpb.Subscriber{}

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

func TestBillingCollectorEventServer_HandleSimManagerSimAllocationEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	s := server.NewBillingCollectorEventServer(OrgName, billingClient)

	t.Run("AllocateSimEventSent", func(t *testing.T) {
		billingClient.On("CreateSubscription", mock.Anything, mock.Anything).
			Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

		planId := "f1ad4204-ab9e-4574-b6bb-bffcc104f8f9"

		sim := epb.SimAllocation{
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

	t.Run("AllocateSimEventNotSent", func(t *testing.T) {
		billingClient.On("CreateSubscription", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send create subscription event")).Once()

		sim := epb.SimAllocation{}

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

func TestBillingCollectorEventServer_HandleSimManagerSetActivePackageForSimEvent(t *testing.T) {
	billingClient := &mocks.BillingClient{}
	routingKey := msgbus.PrepareRoute(OrgName, "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activepackage")

	billingClient.On("GetBillableMetricId", mock.Anything,
		server.DefaultBillableMetricCode).Return(bmId, nil).Once()

	s := server.NewBillingCollectorEventServer(OrgName, billingClient)

	t.Run("SetActivePackageEventSent", func(t *testing.T) {
		billingClient.On("TerminateSubscription", mock.Anything, "b20c61f1-1c5a-4559-bfff-cd00f746697d").
			Return("9fd07299-2826-4f8b-aea9-69da56440bec", nil).Once()

		billingClient.On("CreateSubscription", mock.Anything, mock.Anything).
			Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

		sim := epb.SimActivePackage{
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

	t.Run("SetActivePackageEventNotSent", func(t *testing.T) {
		billingClient.On("TerminateSubscription", mock.Anything, "16befbda-250f-4a68-9cb4-31cccce3005e").
			Return("", errors.New("failed to send terminate subscription event")).Once()

		sim := epb.SimActivePackage{
			Id: "16befbda-250f-4a68-9cb4-31cccce3005e",
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
}
