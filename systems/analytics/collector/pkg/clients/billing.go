/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package clients

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const BillingAccountEndpoint = "/v1/account"

/* TODO-verify: response shape against billing api-gateway. */

type BillingAccount struct {
	Balance             float64 `json:"balance"`
	Currency            string  `json:"currency"`
	PaymentMethodStatus string  `json:"payment_method_status"`
	LastInvoiceAt       string  `json:"last_invoice_at"`
}

type billingAccountResponse struct {
	Account BillingAccount `json:"account"`
}

type BillingClient interface {
	GetAccount() (*BillingAccount, error)
}

type billingClient struct {
	u *url.URL
	R *client.Resty
}

func NewBillingClient(h string, options ...client.Option) BillingClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &billingClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (c *billingClient) GetAccount() (*BillingAccount, error) {
	resp, err := c.R.Get(c.u.String() + BillingAccountEndpoint)
	if err != nil {
		return nil, fmt.Errorf("GetAccount failure: %w", err)
	}

	out := billingAccountResponse{}

	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		return nil, fmt.Errorf("billing account deserialization failure: %w", err)
	}

	return &out.Account, nil
}
