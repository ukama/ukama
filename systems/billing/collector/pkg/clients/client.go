package client

import (
	"context"
	"time"
)

type BillingClient interface {
	CreateCustomer(context.Context, Customer) (string, error)
	UpdateCustomer(context.Context, Customer) (string, error)
	DeleteCustomer(context.Context, string) (string, error)
	CreateSubscription(context.Context, Subscription) (string, error)
	TerminateSubscription(context.Context, string) (string, error)
	AddUsageEvent(context.Context, Event) error
}

type Event struct {
	TransactionId        string
	CustomerId           string
	SubscriptionId       string
	Code                 string
	SentAt               time.Time
	AdditionalProperties map[string]string
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
