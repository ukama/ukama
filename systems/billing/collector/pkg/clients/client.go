/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"context"
	"time"
)

type BillingClient interface {
	GetBillableMetricId(context.Context, string) (string, error)
	CreateBillableMetric(context.Context, BillableMetric) (string, error)
	CreatePlan(context.Context, Plan) (string, error)
	CreateCustomer(context.Context, Customer) (string, error)
	UpdateCustomer(context.Context, Customer) (string, error)
	DeleteCustomer(context.Context, string) (string, error)
	CreateSubscription(context.Context, Subscription) (string, error)
	TerminateSubscription(context.Context, string) (string, error)
	AddUsageEvent(context.Context, Event) error
}

type BillableMetric struct {
	Name        string
	Code        string
	Description string
	FieldName   string
}

type Event struct {
	TransactionId        string
	CustomerId           string
	SubscriptionId       string
	Code                 string
	SentAt               time.Time
	AdditionalProperties map[string]string
}

type Plan struct {
	Name              string
	Code              string
	Interval          string
	PayInAdvance      bool
	AmountCents       int
	AmountCurrency    string
	BillChargeMonthly bool
	TrialPeriod       float32

	BillableMetricID     string
	ChargeModel          string
	ChargeAmountCents    string
	ChargeAmountCurrency string
	FreeUnits            int
	PackageSize          int
}

type Customer struct {
	Id      string
	Name    string
	Email   string
	Address string
	Phone   string
}

type Subscription struct {
	Id             string
	CustomerId     string
	PlanCode       string
	SubscriptionAt *time.Time
}
