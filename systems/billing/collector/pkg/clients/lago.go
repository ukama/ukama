package client

import (
	"context"
	"fmt"

	lago "github.com/getlago/lago-go-client"
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

func (l *lagoClient) AddUsageEvent(ctx context.Context, ev Event) error {
	eventInput := &lago.EventInput{
		TransactionID:      ev.TransactionId,
		ExternalCustomerID: ev.CustomerId,
		// ExternalSubscriptionID: ev.SubscriptionId,
		Code:       ev.Code,
		Timestamp:  ev.SentAt.Unix(),
		Properties: ev.AdditionalProperties,
	}

	err := l.c.Event().Create(ctx, eventInput)

	if err != nil {
		return fmt.Errorf("error while sending sim usage event: %s. code: %d. %w",
			err.Msg, err.HTTPStatusCode, err.Err)
	}

	return nil
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
