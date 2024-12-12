/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package clients

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const (
	PaymentEndpoint = "/v1/payments"
)

type PaymentsClient interface {
	ListPayments(string) ([]*PaymentInfo, error)
}

type paymentsClient struct {
	u *url.URL
	R *client.Resty
}

func NewPaymentsClient(h string) *paymentsClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &paymentsClient{
		u: u,
		R: client.NewResty(client.WithError(&Err{}),
			client.WithDebug(), client.WithContentTypeJSON()),
	}
}

func (p *paymentsClient) ListPayments(queryString string) ([]*PaymentInfo, error) {
	log.Infof("Listing payments matching: %v", queryString)
	payments := &Payment{}

	resp, err := p.R.Get(p.u.String() + PaymentEndpoint + queryString)
	if err != nil {
		log.Errorf("ListPayments failure. error: %s", err.Error())

		return nil, fmt.Errorf("ListPayments failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &payments)
	log.Infof("Error %v", err)
	if err != nil {
		log.Tracef(deserializeLogMsg, "payments", err.Error())

		return nil, fmt.Errorf(deserializeErrorMsg, "payments", err)
	}

	log.Infof(resourceLogMsg, "Payments", payments.Payments)

	return payments.Payments, nil
}

type Payment struct {
	Payments []*PaymentInfo `json:"payments"`
}

type PaymentInfo struct {
	Id                   string `json:"id,omitempty"`
	ItemId               string `json:"item_id,omitempty"`
	ItemType             string `json:"item_type,omitempty"`
	AmountCents          int64  `json:"amount_cents,omitempty"`
	DepositedAmountCents int64  `json:"deposited_amount_cents,omitempty"`
	Currency             string `json:"currency,omitempty"`
	PaymentMethod        string `json:"payment_method,omitempty"`
	PaidAt               string `json:"paid_at,omitempty"`
	PayerName            string `json:"payer_name,omitempty"`
	PayerEmail           string `json:"payer_email,omitempty"`
	PayerPhone           string `json:"payer_phone,omitempty"`
	Correspondent        string `json:"correspondent,omitempty"`
	Country              string `json:"country,omitempty"`
	Description          string `json:"description,omitempty"`
	Status               string `json:"status,omitempty"`
	FailureReason        string `json:"faillure_reason,omitempty"`
	ExternalId           string `json:"externa_id,omitempty"`
	CreatedAt            string `json:"created_at,omitempty"`
	Metadata             []byte `json:"metadata,omitempty"`
}
