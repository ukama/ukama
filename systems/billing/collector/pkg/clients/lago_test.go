/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package clients_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"

	"github.com/ukama/ukama/systems/billing/collector/mocks"
	"github.com/ukama/ukama/systems/billing/collector/pkg/clients"

	lago "github.com/getlago/lago-go-client"
	guuid "github.com/google/uuid"
)

func TestLagoClient_AddUsaeEvent(t *testing.T) {
	t.Run("UsageEventNotSent", func(t *testing.T) {
		e := &mocks.LagoEvent{}
		l := clients.NewLagoClientFromClients(nil, nil, e, nil, nil, nil)

		e.On("Create", mock.Anything, mock.Anything).Return(&lago.Event{}, &lago.Error{})

		err := l.AddUsageEvent(context.TODO(), clients.Event{})

		assert.Error(t, err)
	})

	t.Run("UsageEventSent", func(t *testing.T) {
		e := &mocks.LagoEvent{}
		l := clients.NewLagoClientFromClients(nil, nil, e, nil, nil, nil)

		e.On("Create", mock.Anything, mock.Anything).Return(&lago.Event{}, nil)

		err := l.AddUsageEvent(context.TODO(), clients.Event{})

		assert.NoError(t, err)
	})
}

func TestLagoClient_GetBillableMetric(t *testing.T) {
	t.Run("BillableMetricNotFound", func(t *testing.T) {
		b := &mocks.LagoBillableMetric{}
		l := clients.NewLagoClientFromClients(b, nil, nil, nil, nil, nil)

		b.On("Get", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		bm, err := l.GetBillableMetricId(context.TODO(), "SomeBillableMectricCode")

		assert.Error(t, err)
		assert.Empty(t, bm)
	})

	t.Run("BillableMetricFound", func(t *testing.T) {
		b := &mocks.LagoBillableMetric{}
		l := clients.NewLagoClientFromClients(b, nil, nil, nil, nil, nil)

		b.On("Get", mock.Anything, mock.Anything).Return(&lago.BillableMetric{}, nil)

		bm, err := l.GetBillableMetricId(context.TODO(), "SomeBillableMectricCode")

		assert.NoError(t, err)
		assert.NotEmpty(t, bm)
	})
}

func TestLagoClient_CreateBillableMetric(t *testing.T) {
	t.Run("BillableMetricNotCreated", func(t *testing.T) {
		b := &mocks.LagoBillableMetric{}
		l := clients.NewLagoClientFromClients(b, nil, nil, nil, nil, nil)

		b.On("Create", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		bm, err := l.CreateBillableMetric(context.TODO(), clients.BillableMetric{})

		assert.Error(t, err)
		assert.Empty(t, bm)
	})

	t.Run("BillableMetricCreated", func(t *testing.T) {
		b := &mocks.LagoBillableMetric{}
		l := clients.NewLagoClientFromClients(b, nil, nil, nil, nil, nil)

		b.On("Create", mock.Anything, mock.Anything).Return(&lago.BillableMetric{}, nil)

		bm, err := l.CreateBillableMetric(context.TODO(), clients.BillableMetric{})

		assert.NoError(t, err)
		assert.NotEmpty(t, bm)
	})
}

func TestLagoClient_GetPlan(t *testing.T) {
	t.Run("PlanNotFound", func(t *testing.T) {
		p := &mocks.LagoPlan{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, p, nil, nil)

		p.On("Get", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		plan, err := l.GetPlan(context.TODO(), "SomePlanCode")

		assert.Error(t, err)
		assert.Empty(t, plan)
	})

	t.Run("PlanFound", func(t *testing.T) {
		p := &mocks.LagoPlan{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, p, nil, nil)

		p.On("Get", mock.Anything, mock.Anything).Return(&lago.Plan{}, nil)

		plan, err := l.GetPlan(context.TODO(), "SomePlanCode")

		assert.NoError(t, err)
		assert.NotEmpty(t, plan)
	})
}

func TestLagoClient_CreatePlan(t *testing.T) {
	t.Run("InvalidPlanChargeBillableMetricId", func(t *testing.T) {
		p := &mocks.LagoPlan{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, p, nil, nil)

		charges := []clients.PlanCharge{
			clients.PlanCharge{
				BillableMetricID: "lol",
			},
		}

		plan, err := l.CreatePlan(context.TODO(), clients.Plan{}, charges...)

		assert.Error(t, err)
		assert.Empty(t, plan)
	})
	t.Run("PlanNotCreated", func(t *testing.T) {
		p := &mocks.LagoPlan{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, p, nil, nil)

		bmId := guuid.New().String()

		charges := []clients.PlanCharge{
			clients.PlanCharge{
				BillableMetricID: bmId,
			},
		}

		p.On("Create", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		plan, err := l.CreatePlan(context.TODO(), clients.Plan{}, charges...)

		assert.Error(t, err)
		assert.Empty(t, plan)
	})

	t.Run("PlanCreated", func(t *testing.T) {
		p := &mocks.LagoPlan{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, p, nil, nil)

		p.On("Create", mock.Anything, mock.Anything).Return(&lago.Plan{}, nil)

		plan, err := l.CreatePlan(context.TODO(), clients.Plan{})

		assert.NoError(t, err)
		assert.NotEmpty(t, plan)
	})
}

func TestLagoClient_TerminatePlan(t *testing.T) {
	t.Run("PlanNotTerminated", func(t *testing.T) {
		p := &mocks.LagoPlan{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, p, nil, nil)

		p.On("Delete", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		plan, err := l.TerminatePlan(context.TODO(), "SomePlanCode")

		assert.Error(t, err)
		assert.Empty(t, plan)
	})

	t.Run("PlanTerminated", func(t *testing.T) {
		p := &mocks.LagoPlan{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, p, nil, nil)

		p.On("Delete", mock.Anything, mock.Anything).Return(&lago.Plan{}, nil)

		plan, err := l.TerminatePlan(context.TODO(), "SomePlanCode")

		assert.NoError(t, err)
		assert.NotEmpty(t, plan)
	})
}

func TestLagoClient_GetCustomer(t *testing.T) {
	t.Run("CustomerNotFound", func(t *testing.T) {
		c := &mocks.LagoCustomer{}
		l := clients.NewLagoClientFromClients(nil, c, nil, nil, nil, nil)

		c.On("Get", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		customer, err := l.GetCustomer(context.TODO(), "SomeCustomerId")

		assert.Error(t, err)
		assert.Empty(t, customer)
	})

	t.Run("CustomerFound", func(t *testing.T) {
		c := &mocks.LagoCustomer{}
		l := clients.NewLagoClientFromClients(nil, c, nil, nil, nil, nil)

		c.On("Get", mock.Anything, mock.Anything).Return(&lago.Customer{}, nil)

		customer, err := l.GetCustomer(context.TODO(), "SomeCustomerId")

		assert.NoError(t, err)
		assert.NotEmpty(t, customer)
	})
}

func TestLagoClient_CreateCustomer(t *testing.T) {
	t.Run("CustomerNotCreated", func(t *testing.T) {
		c := &mocks.LagoCustomer{}
		l := clients.NewLagoClientFromClients(nil, c, nil, nil, nil, nil)

		c.On("Create", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		customer, err := l.CreateCustomer(context.TODO(), clients.Customer{
			Type: "company",
		})

		assert.Error(t, err)
		assert.Empty(t, customer)
	})

	t.Run("CustomerCreated", func(t *testing.T) {
		c := &mocks.LagoCustomer{}
		l := clients.NewLagoClientFromClients(nil, c, nil, nil, nil, nil)

		c.On("Create", mock.Anything, mock.Anything).Return(&lago.Customer{}, nil)

		customer, err := l.CreateCustomer(context.TODO(), clients.Customer{
			Type: "individual",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, customer)
	})
}

func TestLagoClient_UpdateCustomer(t *testing.T) {
	t.Run("CustomerNotUpdated", func(t *testing.T) {
		c := &mocks.LagoCustomer{}
		l := clients.NewLagoClientFromClients(nil, c, nil, nil, nil, nil)

		c.On("Update", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		customer, err := l.UpdateCustomer(context.TODO(), clients.Customer{})

		assert.Error(t, err)
		assert.Empty(t, customer)
	})

	t.Run("CustomerUpdated", func(t *testing.T) {
		c := &mocks.LagoCustomer{}
		l := clients.NewLagoClientFromClients(nil, c, nil, nil, nil, nil)

		c.On("Update", mock.Anything, mock.Anything).Return(&lago.Customer{}, nil)

		customer, err := l.UpdateCustomer(context.TODO(), clients.Customer{})

		assert.NoError(t, err)
		assert.NotEmpty(t, customer)
	})
}

func TestLagoClient_DeleteCustomer(t *testing.T) {
	t.Run("CustomerNotDeleted", func(t *testing.T) {
		c := &mocks.LagoCustomer{}
		l := clients.NewLagoClientFromClients(nil, c, nil, nil, nil, nil)

		c.On("Delete", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		customer, err := l.DeleteCustomer(context.TODO(), "SomeCustomerId")

		assert.Error(t, err)
		assert.Empty(t, customer)
	})

	t.Run("CustomerDeleted", func(t *testing.T) {
		c := &mocks.LagoCustomer{}
		l := clients.NewLagoClientFromClients(nil, c, nil, nil, nil, nil)

		c.On("Delete", mock.Anything, mock.Anything).Return(&lago.Customer{}, nil)

		customer, err := l.DeleteCustomer(context.TODO(), "SomeCustomerId")

		assert.NoError(t, err)
		assert.NotEmpty(t, customer)
	})
}

func TestLagoClient_CreateSubscription(t *testing.T) {
	t.Run("SubscriptionNotCreated", func(t *testing.T) {
		s := &mocks.LagoSubscription{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, nil, s, nil)

		s.On("Create", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		subscription, err := l.CreateSubscription(context.TODO(), clients.Subscription{})

		assert.Error(t, err)
		assert.Empty(t, subscription)
	})

	t.Run("SubscriptionCreated", func(t *testing.T) {
		s := &mocks.LagoSubscription{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, nil, s, nil)

		s.On("Create", mock.Anything, mock.Anything).Return(&lago.Subscription{}, nil)

		subscription, err := l.CreateSubscription(context.TODO(), clients.Subscription{})

		assert.NoError(t, err)
		assert.NotEmpty(t, subscription)
	})
}

func TestLagoClient_TerminateSubscription(t *testing.T) {
	t.Run("SubscriptionNotTerminated", func(t *testing.T) {
		s := &mocks.LagoSubscription{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, nil, s, nil)

		s.On("Terminate", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		subscription, err := l.TerminateSubscription(context.TODO(), "SomeSubscriptionId")

		assert.Error(t, err)
		assert.Empty(t, subscription)
	})

	t.Run("SubscriptionTerminated", func(t *testing.T) {
		s := &mocks.LagoSubscription{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, nil, s, nil)

		s.On("Terminate", mock.Anything, mock.Anything).Return(&lago.Subscription{}, nil)

		subscription, err := l.TerminateSubscription(context.TODO(), "SomeSubscriptionId")

		assert.NoError(t, err)
		assert.NotEmpty(t, subscription)
	})
}

func TestLagoClient_CreateWebhook(t *testing.T) {
	t.Run("WebhookNotCreated", func(t *testing.T) {
		w := &mocks.LagoWebhookEndpoint{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, nil, nil, w)

		w.On("Create", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		webhhok, err := l.CreateWebhook(context.TODO(), clients.WebhookEndpoint{
			SignatureAlgo: "hmac",
		})

		assert.Error(t, err)
		assert.Empty(t, webhhok)
	})

	t.Run("WebhookCreated", func(t *testing.T) {
		w := &mocks.LagoWebhookEndpoint{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, nil, nil, w)

		w.On("Create", mock.Anything, mock.Anything).Return(&lago.WebhookEndpoint{}, nil)

		webhook, err := l.CreateWebhook(context.TODO(), clients.WebhookEndpoint{
			SignatureAlgo: "hmac",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, webhook)
	})
}

func TestLagoClient_ListWebhook(t *testing.T) {
	t.Run("WebhookNotFound", func(t *testing.T) {
		w := &mocks.LagoWebhookEndpoint{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, nil, nil, w)

		w.On("GetList", mock.Anything, mock.Anything).Return(nil, &lago.Error{})

		webhooks, err := l.ListWebhooks(context.TODO())

		assert.Error(t, err)
		assert.Empty(t, webhooks)
	})

	t.Run("WebhookFound", func(t *testing.T) {
		w := &mocks.LagoWebhookEndpoint{}
		l := clients.NewLagoClientFromClients(nil, nil, nil, nil, nil, w)

		w.On("GetList", mock.Anything, mock.Anything).Return(&lago.WebhookEndpointResult{
			WebhookEndpoints: []lago.WebhookEndpoint{
				lago.WebhookEndpoint{},
			},
		}, nil)

		webhook, err := l.ListWebhooks(context.TODO())

		assert.NoError(t, err)
		assert.NotEmpty(t, webhook)
	})
}
