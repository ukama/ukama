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
	"fmt"

	lago "github.com/getlago/lago-go-client"
	guuid "github.com/google/uuid"
)

type lagoClient struct {
	c *lago.Client
}

func NewLagoClient(APIKey, Host string, Port uint) BillingClient {
	lagoBaseURL := fmt.Sprintf("http://%s:%d", Host, Port)

	return &lagoClient{
		c: lago.New().SetBaseURL(lagoBaseURL).SetApiKey(APIKey).SetDebug(true),
	}
}

func (l *lagoClient) GetBillableMetricId(ctx context.Context, code string) (string, error) {
	bm, pErr := l.c.BillableMetric().Get(ctx, code)
	if pErr != nil {
		return "", fmt.Errorf("error while getting billable metrict ID: %s. code: %d. %w",
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

	bm, pErr := l.c.BillableMetric().Create(ctx, bmInput)
	if pErr != nil {
		return "", fmt.Errorf("error while creating billable metrict: %s. code: %d. %w",
			pErr.Msg, pErr.HTTPStatusCode, pErr.Err)
	}

	return bm.LagoID.String(), nil
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

	err := l.c.Event().Create(ctx, eventInput)

	if err != nil {
		return fmt.Errorf("error while sending sim usage event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return nil
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
			return "", fmt.Errorf("fail to parse billable metric ID: %w", err)
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

	plan, pErr := l.c.Plan().Create(ctx, newPlan)
	if pErr != nil {
		return "", fmt.Errorf("error while sending plan creation event: %s. code: %d. %w",
			pErr.Msg, pErr.HTTPStatusCode, pErr.Err)
	}

	return plan.LagoID.String(), nil
}

func (l *lagoClient) GetCustomer(ctx context.Context, id string) (string, error) {
	customer, pErr := l.c.Customer().Get(ctx, id)
	if pErr != nil {
		return "", fmt.Errorf("error while getting customer with ID %s: %s. code: %d. %w",
			id, pErr.Msg, pErr.HTTPStatusCode, pErr.Err)
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

	customer, err := l.c.Customer().Create(ctx, newCust)

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

	customer, err := l.c.Customer().Update(ctx, newCust)

	if err != nil {
		return "", fmt.Errorf("error while sending customer update event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return customer.LagoID.String(), nil
}

func (l *lagoClient) DeleteCustomer(ctx context.Context, custId string) (string, error) {
	customer, err := l.c.Customer().Delete(ctx, custId)

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

	subscription, err := l.c.Subscription().Create(ctx, newSub)

	if err != nil {
		return "", fmt.Errorf("error while sending subscription create event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return subscription.LagoID.String(), nil
}

func (l *lagoClient) TerminateSubscription(ctx context.Context, subscritionId string) (string, error) {
	subscription, err := l.c.Subscription().Terminate(ctx, subscritionId)

	if err != nil {
		return "", fmt.Errorf("error while sending subscription termination event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return subscription.LagoID.String(), nil
}
