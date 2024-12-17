/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package clients

import (
	"fmt"

	"github.com/stripe/stripe-go/v78/client"
	"github.com/ukama/ukama/systems/common/util/payments"

	log "github.com/sirupsen/logrus"
	stripelib "github.com/stripe/stripe-go/v78"
)

type StripeClient interface {
	GetPaymentIntent(id string) (*payments.Intent, error)
}

type StripeClientWrapper struct {
	k string
	c *client.API
}

func NewStripeClient(key string, options ...Option) *StripeClientWrapper {
	strpBackends := &stripelib.Backends{
		API:     stripelib.GetBackend(stripelib.APIBackend),
		Connect: stripelib.GetBackend(stripelib.ConnectBackend),
		Uploads: stripelib.GetBackend(stripelib.UploadsBackend),
	}

	c := &StripeClientWrapper{
		k: key,
		c: client.New(key, strpBackends),
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (s *StripeClientWrapper) GetPaymentIntent(id string) (*payments.Intent, error) {
	log.Infof("Getting payment internt %v", id)

	intentParams := &stripelib.PaymentIntentParams{}

	paymentIntent, err := s.c.PaymentIntents.Get(id, intentParams)
	if err != nil {
		log.Errorf("Failed to get a stripe payment intent: %v", err)

		return nil, fmt.Errorf("failed to get a stripe payment intent: %w", err)
	}

	log.Infof("Successfuly got a stripe payment intent request %s from payment provider",
		paymentIntent.ID)

	return &payments.Intent{paymentIntent}, nil
}

type Option func(*StripeClientWrapper)

func WithCustomBackends(backends *stripelib.Backends) Option {
	return func(s *StripeClientWrapper) {
		s.c = client.New(s.k, backends)
	}
}
