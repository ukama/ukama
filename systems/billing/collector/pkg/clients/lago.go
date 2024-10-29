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

	return &lagoClient{b: b, c: c, e: e, p: p, s: s}
}

func (l *lagoClient) AddUsageEvent(ctx context.Context, ev Event) error {
	eventInput := &lago.EventInput{
		TransactionID:          ev.TransactionId,
		ExternalSubscriptionID: ev.SubscriptionId,
		Code:                   ev.Code,
		Timestamp:              ev.SentAt.String(),
		Properties:             ev.AdditionalProperties,
	}

	_, err := l.e.Create(ctx, eventInput)
	if err != nil {
		msg := "error while sending sim usage event"

		return unpackLagoError(msg, err)
	}

	return nil
}

func (l *lagoClient) GetBillableMetricId(ctx context.Context, code string) (string, error) {
	bm, err := l.b.Get(ctx, code)
	if err != nil {
		msg := fmt.Sprintf("error while getting billable metrict Id (%s)",
			code)

		return "", unpackLagoError(msg, err)
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

	bm, err := l.b.Create(ctx, bmInput)
	if err != nil {
		msg := "error while creating billable metric"

		return "", unpackLagoError(msg, err)
	}

	return bm.LagoID.String(), nil
}

func (l *lagoClient) GetPlan(ctx context.Context, planCode string) (string, error) {
	plan, err := l.p.Get(ctx, planCode)
	if err != nil {
		msg := fmt.Sprintf("error while getting plan with Id (%s)",
			planCode)

		return "", unpackLagoError(msg, err)
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
			return "", fmt.Errorf("fail to parse billable metric Id: %w",
				err)
		}

		props := make(map[string]interface{})

		props["amount"] = charge.ChargeAmount
		props["free_units"] = charge.FreeUnits
		props["package_size"] = charge.PackageSize

		newCharge := lago.PlanChargeInput{
			BillableMetricID: bMetricId,
			ChargeModel:      lago.ChargeModel(charge.ChargeModel),
			AmountCurrency:   lago.Currency(pl.AmountCurrency),
			Properties:       props,
		}

		// Appending charge to plan
		newPlan.Charges = append(newPlan.Charges, newCharge)
	}

	plan, err := l.p.Create(ctx, newPlan)
	if err != nil {
		msg := "error while sending plan creation event"

		return "", unpackLagoError(msg, err)
	}

	return plan.LagoID.String(), nil
}

func (l *lagoClient) TerminatePlan(ctx context.Context, planCode string) (string, error) {
	plan, err := l.p.Delete(ctx, planCode)
	if err != nil {
		msg := fmt.Sprintf("error while sending plan delete event with code (%s)",
			planCode)

		return "", unpackLagoError(msg, err)
	}

	return plan.LagoID.String(), nil
}

func (l *lagoClient) GetCustomer(ctx context.Context, custId string) (string, error) {
	customer, err := l.c.Get(ctx, custId)
	if err != nil {
		msg := fmt.Sprintf("error while getting customer with Id (%s)",
			custId)

		return "", unpackLagoError(msg, err)
	}

	return customer.LagoID.String(), nil
}

func (l *lagoClient) CreateCustomer(ctx context.Context, cust Customer) (string, error) {
	var customerType lago.CustomerType = IndividualCustomerType

	if cust.Type == CompanyCustomerType {
		customerType = lago.CompanyCustomerType
	}

	newCust := &lago.CustomerInput{
		ExternalID:   cust.Id,
		Name:         cust.Name,
		Email:        cust.Email,
		AddressLine1: cust.Address,
		Phone:        cust.Phone,
		CustomerType: customerType,
	}

	customer, err := l.c.Create(ctx, newCust)
	if err != nil {
		msg := "error while sending customer creation event"

		return "", unpackLagoError(msg, err)
	}

	return customer.LagoID.String(), nil
}

func (l *lagoClient) UpdateCustomer(ctx context.Context, cust Customer) (string, error) {
	newCust := &lago.CustomerInput{
		ExternalID:   cust.Id,
		Name:         cust.Name,
		Email:        cust.Email,
		AddressLine1: cust.Address,
		Phone:        cust.Phone,
	}

	customer, err := l.c.Update(ctx, newCust)
	if err != nil {
		msg := fmt.Sprintf("error while sending customer update event with Id (%s)",
			cust.Id)

		return "", unpackLagoError(msg, err)
	}

	return customer.LagoID.String(), nil
}

func (l *lagoClient) DeleteCustomer(ctx context.Context, custId string) (string, error) {
	customer, err := l.c.Delete(ctx, custId)
	if err != nil {
		msg := fmt.Sprintf("error while sending customer delete event witch Id (%s)",
			custId)

		return "", unpackLagoError(msg, err)
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
		msg := "error while sending subscription create event"

		return "", unpackLagoError(msg, err)
	}

	return subscription.LagoID.String(), nil
}

func (l *lagoClient) TerminateSubscription(ctx context.Context, subscriptionId string) (string, error) {
	subscriptionTerminateInput := lago.SubscriptionTerminateInput{
		ExternalID: subscriptionId,
	}

	subscription, err := l.s.Terminate(ctx, subscriptionTerminateInput)
	if err != nil {
		msg := fmt.Sprintf("error while sending subscription termination event with Id (%s)",
			subscriptionTerminateInput.ExternalID)

		return "", unpackLagoError(msg, err)
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
	Create(context.Context, *lago.EventInput) (*lago.Event, *lago.Error)
}

type LagoPlan interface {
	Get(context.Context, string) (*lago.Plan, *lago.Error)
	Create(context.Context, *lago.PlanInput) (*lago.Plan, *lago.Error)
	Delete(context.Context, string) (*lago.Plan, *lago.Error)
}

type LagoSubscription interface {
	Create(context.Context, *lago.SubscriptionInput) (*lago.Subscription, *lago.Error)
	Terminate(context.Context, lago.SubscriptionTerminateInput) (*lago.Subscription, *lago.Error)
}

func unpackLagoError(msg string, err *lago.Error) error {
	cltError := &Error{
		Code: err.HTTPStatusCode,
		Msg:  err.Message,
		Err:  err.Err,
	}

	return fmt.Errorf(msg+": %s. code: %d. %w",
		err.Message, err.HTTPStatusCode, cltError)
}
