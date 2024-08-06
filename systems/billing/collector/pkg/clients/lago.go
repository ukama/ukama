/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package clients

import (
	"context"
	"fmt"

	lago "github.com/getlago/lago-go-client"
	guuid "github.com/google/uuid"
)

type lagoClient struct {
	b LagoBillableMetric
	c LagoCustomer
	e LagoEvent
	p LagoPlan
	s LagoSubscription
}

func NewLagoClient(APIKey, Host string, Port uint) BillingClient {
	lagoBaseURL := fmt.Sprintf("http://%s:%d", Host, Port)
	c := lago.New().SetBaseURL(lagoBaseURL).SetApiKey(APIKey).SetDebug(true)

	return &lagoClient{
		b: c.BillableMetric(), c: c.Customer(),
		e: c.Event(), p: c.Plan(), s: c.Subscription(),
	}
}

func NewLagoClientFromClients(b LagoBillableMetric, c LagoCustomer,
	e LagoEvent, p LagoPlan, s LagoSubscription) BillingClient {
	return &lagoClient{
		b: b,
		c: c,
		e: e,
		p: p,
		s: s,
	}

}

func (l *lagoClient) AddUsageEvent(ctx context.Context, ev Event) error {
	eventInput := &lago.EventInput{
		TransactionID:          ev.TransactionId,
		ExternalCustomerID:     ev.CustomerId,
		ExternalSubscriptionID: ev.SubscriptionId,
		Code:                   ev.Code,
		Timestamp:              ev.SentAt.Unix(),
		Properties:             ev.AdditionalProperties,
	}

	err := l.e.Create(ctx, eventInput)
	if err != nil {
		return fmt.Errorf("error while sending sim usage event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return nil
}

func (l *lagoClient) GetBillableMetricId(ctx context.Context, code string) (string, error) {
	bm, pErr := l.b.Get(ctx, code)
	if pErr != nil {
		return "", fmt.Errorf("error while getting billable metrict Id: %s. code: %d. %w",
			pErr.Msg, pErr.HTTPStatusCode, pErr.Err)
	}

	return bm.LagoID.String(), nil
}

func (l *lagoClient) CreateBillableMetric(ctx context.Context, bMetric BillableMetric) (string, error) {
	bmInput := &lago.BillableMetricInput{
		Name:            bMetric.Name,
		Code:            bMetric.Code,
		Description:     bMetric.Description,
		AggregationType: lago.SumAggregation,
		FieldName:       bMetric.FieldName,
	}

	bm, pErr := l.b.Create(ctx, bmInput)
	if pErr != nil {
		return "", fmt.Errorf("error while creating billable metrict: %s. code: %d. %w",
			pErr.Msg, pErr.HTTPStatusCode, pErr.Err)
	}

	return bm.LagoID.String(), nil
}

func (l *lagoClient) GetPlan(ctx context.Context, planCode string) (string, error) {
	plan, pErr := l.p.Get(ctx, planCode)
	if pErr != nil {
		return "", fmt.Errorf("error while getting plan with Id %s: %s. code: %d. %w",
			planCode, pErr.Msg, pErr.HTTPStatusCode, pErr.Err)
	}

	return plan.LagoID.String(), nil
}

func (l *lagoClient) CreatePlan(ctx context.Context, pl Plan, charges ...PlanCharge) (string, error) {
	newPlan := &lago.PlanInput{
		Name:           pl.Name,
		Code:           pl.Code,
		Interval:       lago.PlanInterval(pl.Interval),
		PayInAdvance:   pl.PayInAdvance,
		AmountCents:    pl.AmountCents,
		AmountCurrency: lago.Currency(pl.AmountCurrency),
	}

	// Processing each charge, if any
	for _, charge := range charges {
		bMetricId, err := guuid.Parse(charge.BillableMetricID)
		if err != nil {
			return "", fmt.Errorf("fail to parse billable metric Id: %w", err)
		}

		props := make(map[string]interface{})

		props["amount"] = charge.ChargeAmountCents
		props["free_units"] = charge.FreeUnits
		props["package_size"] = charge.PackageSize

		newCharge := lago.PlanChargeInput{
			BillableMetricID: bMetricId,
			ChargeModel:      lago.ChargeModel(charge.ChargeModel),
			AmountCurrency:   lago.Currency(pl.AmountCurrency),
			// PayInAdvance:     true,
			Properties: props,
		}

		// Appending charge to plan
		newPlan.Charges = append(newPlan.Charges, newCharge)
	}

	plan, pErr := l.p.Create(ctx, newPlan)
	if pErr != nil {
		return "", fmt.Errorf("error while sending plan creation event: %s. code: %d. %w",
			pErr.Msg, pErr.HTTPStatusCode, pErr.Err)
	}

	return plan.LagoID.String(), nil
}

func (l *lagoClient) TerminatePlan(ctx context.Context, planCode string) (string, error) {
	plan, err := l.p.Delete(ctx, planCode)

	if err != nil {
		return "", fmt.Errorf("error while sending plan delete event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return plan.LagoID.String(), nil
}

func (l *lagoClient) GetCustomer(ctx context.Context, custId string) (string, error) {
	customer, pErr := l.c.Get(ctx, custId)
	if pErr != nil {
		return "", fmt.Errorf("error while getting customer with Id %s: %s. code: %d. %w",
			custId, pErr.Msg, pErr.HTTPStatusCode, pErr.Err)
	}

	return customer.LagoID.String(), nil
}

func (l *lagoClient) CreateCustomer(ctx context.Context, cust Customer) (string, error) {
	newCust := &lago.CustomerInput{
		ExternalID:   cust.Id,
		Name:         cust.Name,
		Email:        cust.Email,
		AddressLine1: cust.Address,
		Phone:        cust.Phone,
	}

	customer, err := l.c.Create(ctx, newCust)
	if err != nil {
		return "", fmt.Errorf("error while sending customer creation event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return customer.LagoID.String(), nil
}

func (l *lagoClient) UpdateCustomer(ctx context.Context, cust Customer) (string, error) {
	newCust := &lago.CustomerInput{
		Name:         cust.Name,
		Email:        cust.Email,
		AddressLine1: cust.Address,
		Phone:        cust.Phone,
	}

	customer, err := l.c.Update(ctx, newCust)
	if err != nil {
		return "", fmt.Errorf("error while sending customer update event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return customer.LagoID.String(), nil
}

func (l *lagoClient) DeleteCustomer(ctx context.Context, custId string) (string, error) {
	customer, err := l.c.Delete(ctx, custId)
	if err != nil {
		return "", fmt.Errorf("error while sending customer delete event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return customer.LagoID.String(), nil
}

func (l *lagoClient) CreateSubscription(ctx context.Context, sub Subscription) (string, error) {
	newSub := &lago.SubscriptionInput{
		ExternalID:         sub.Id,
		ExternalCustomerID: sub.CustomerId,
		PlanCode:           sub.PlanCode,
		SubscriptionAt:     sub.SubscriptionAt,
	}

	subscription, err := l.s.Create(ctx, newSub)
	if err != nil {
		return "", fmt.Errorf("error while sending subscription create event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return subscription.LagoID.String(), nil
}

func (l *lagoClient) TerminateSubscription(ctx context.Context, subscritionId string) (string, error) {
	subscription, err := l.s.Terminate(ctx, subscritionId)
	if err != nil {
		return "", fmt.Errorf("error while sending subscription termination event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return subscription.LagoID.String(), nil
}

type LagoBillableMetric interface {
	Get(context.Context, string) (*lago.BillableMetric, *lago.Error)
	Create(context.Context, *lago.BillableMetricInput) (*lago.BillableMetric, *lago.Error)
}

type LagoCustomer interface {
	Get(context.Context, string) (*lago.Customer, *lago.Error)
	Create(context.Context, *lago.CustomerInput) (*lago.Customer, *lago.Error)
	Update(context.Context, *lago.CustomerInput) (*lago.Customer, *lago.Error)
	Delete(context.Context, string) (*lago.Customer, *lago.Error)
}

type LagoEvent interface {
	Create(context.Context, *lago.EventInput) *lago.Error
}

type LagoPlan interface {
	Get(context.Context, string) (*lago.Plan, *lago.Error)
	Create(context.Context, *lago.PlanInput) (*lago.Plan, *lago.Error)
	Delete(context.Context, string) (*lago.Plan, *lago.Error)
}

type LagoSubscription interface {
	Create(context.Context, *lago.SubscriptionInput) (*lago.Subscription, *lago.Error)
	Terminate(context.Context, string) (*lago.Subscription, *lago.Error)
}
