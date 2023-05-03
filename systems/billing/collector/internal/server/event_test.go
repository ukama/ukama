package server_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/billing/collector/internal/server"
	"github.com/ukama/ukama/systems/billing/collector/mocks"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	operatorpb "github.com/ukama/telna/cdr/pb/gen"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	subpb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
	simpb "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"
)

func TestBillingCollectorEventServer_HandleCdrSimUsageEvent(t *testing.T) {
	t.Run("SimUsageEventSent", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}
		routingKey := "event.cloud.cdr.sim.usage"

		billingClient.On("AddUsageEvent", mock.Anything, mock.Anything).Return(nil).Once()

		simUsage := operatorpb.SimUsage{
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

		s := server.NewBillingCollectorEventServer(billingClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("SimUsageEventNotSent", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}
		routingKey := "event.cloud.cdr.sim.usage"

		billingClient.On("AddUsageEvent", mock.Anything, mock.Anything).
			Return(errors.New("failed to send sim usage")).Once()

		simUsage := operatorpb.SimUsage{}

		anyE, err := anypb.New(&simUsage)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewBillingCollectorEventServer(billingClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestBillingCollectorEventServer_HandleRegistrySubscriberCreateEvent(t *testing.T) {
	t.Run("CreateCustomerEventSent", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}
		routingKey := "event.cloud.registry.subscriber.create"

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

		s := server.NewBillingCollectorEventServer(billingClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("CreateCustomerEventNotSent", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}
		routingKey := "event.cloud.registry.subscriber.create"

		billingClient.On("CreateCustomer", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send create customer event")).Once()

		subscriber := subpb.Subscriber{}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewBillingCollectorEventServer(billingClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestBillingCollectorEventServer_HandleRegistrySubscriberUpdateEvent(t *testing.T) {
	t.Run("UpdateCustomerEventSent", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}
		routingKey := "event.cloud.registry.subscriber.update"

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

		s := server.NewBillingCollectorEventServer(billingClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("UpdateCustomerEventNotSent", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}
		routingKey := "event.cloud.registry.subscriber.update"

		billingClient.On("UpdateCustomer", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send update customer event")).Once()

		subscriber := subpb.Subscriber{}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewBillingCollectorEventServer(billingClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestBillingCollectorEventServer_HandleRegistrySubscriberDeleteEvent(t *testing.T) {
	t.Run("DeleteCustomerEventSent", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}
		routingKey := "event.cloud.registry.subscriber.delete"

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

		s := server.NewBillingCollectorEventServer(billingClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("DeleteCustomerEventNotSent", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}
		routingKey := "event.cloud.registry.subscriber.delete"

		billingClient.On("DeleteCustomer", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send delete customer event")).Once()

		subscriber := subpb.Subscriber{}

		anyE, err := anypb.New(&subscriber)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewBillingCollectorEventServer(billingClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}

func TestBillingCollectorEventServer_HandleSimManagerSetActivePackageForSimEvent(t *testing.T) {
	t.Run("SetActivePackageEventSent", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}
		routingKey := "event.cloud.simmanager.package.activate"

		billingClient.On("TerminateSubscription", mock.Anything, mock.Anything).Return("9fd07299-2826-4f8b-aea9-69da56440bec", nil).Once()
		billingClient.On("CreateSubscription", mock.Anything, mock.Anything).Return("75ec112a-8745-49f9-ab64-1a37edade794", nil).Once()

		pkg := &simpb.Package{
			// PlanId:    "9fd07299-2826-4f8b-aea9-69da56440bec",
			StartDate: timestamppb.New(time.Now()),
		}

		sim := simpb.Sim{
			Id:           "b20c61f1-1c5a-4559-bfff-cd00f746697d",
			SubscriberId: "c214f255-0ed6-4aa1-93e7-e333658c7318",
			Package:      pkg,
		}

		anyE, err := anypb.New(&sim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewBillingCollectorEventServer(billingClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.NoError(t, err)
	})

	t.Run("SetActivePackageEventNotSent", func(t *testing.T) {
		billingClient := &mocks.BillingClient{}
		routingKey := "event.cloud.simmanager.package.activate"

		billingClient.On("TerminateSubscription", mock.Anything, mock.Anything).
			Return("", errors.New("failed to send terminate subscription event")).Once()

		sim := simpb.Sim{}

		anyE, err := anypb.New(&sim)
		assert.NoError(t, err)

		msg := &epb.Event{
			RoutingKey: routingKey,
			Msg:        anyE,
		}

		s := server.NewBillingCollectorEventServer(billingClient)

		_, err = s.EventNotification(context.TODO(), msg)

		assert.Error(t, err)
	})
}
