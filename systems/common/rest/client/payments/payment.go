/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package payments

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const PaymentEndpoint = "/v1/payments"

type PaymentInfo struct {
	Id              string `json:"id"`
	ItemId          string `json:"item_id"`
	ItemType        string `json:"item_type"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	PaymentMethod   string `json:"payment_method"`
	DepositedAmount string `json:"deposited_amount"`
	PaidAt          string `json:"paid_at"`
	PayerName       string `json:"payer_name"`
	PayerEmail      string `json:"payer_email"`
	PayerPhone      string `json:"payer_phone"`
	Country         string `json:"country"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	ExternalId      string `json:"external_id"`
	Metadata        string `json:"metadata"`
	CreatedAt       string `json:"created_at"`
}

type Payment struct {
	PaymentInfo *PaymentInfo `json:"payment"`
}

// AddPaymentRequest mirrors the payments api-gateway POST /v1/payments body.
// Metadata is a plain map (e.g. {"sim": "<id>"}); the api-gateway marshals it to
// the payment's bytes metadata.
type AddPaymentRequest struct {
	ItemId        string            `json:"item_id"`
	ItemType      string            `json:"item_type"`
	Amount        string            `json:"amount"`
	Currency      string            `json:"currency"`
	PaymentMethod string            `json:"payment_method"`
	PayerName     string            `json:"payer_name"`
	PayerPhone    string            `json:"payer_phone,omitempty"`
	PayerEmail    string            `json:"payer_email,omitempty"`
	Country       string            `json:"country,omitempty"`
	Description   string            `json:"description,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

type PaymentClient interface {
	Add(req AddPaymentRequest) (*PaymentInfo, error)
}

type paymentClient struct {
	u *url.URL
	R *client.Resty
}

func NewPaymentClient(h string, options ...client.Option) *paymentClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &paymentClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (p *paymentClient) Add(req AddPaymentRequest) (*PaymentInfo, error) {
	log.Debugf("Adding payment: %+v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %w", err)
	}

	payment := Payment{}

	resp, err := p.R.Post(p.u.String()+PaymentEndpoint, b)
	if err != nil {
		log.Errorf("AddPayment failure. error: %s", err.Error())

		return nil, fmt.Errorf("AddPayment failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &payment)
	if err != nil {
		log.Tracef("Failed to deserialize payment info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("payment info deserialization failure: %w", err)
	}

	log.Infof("Payment Info: %+v", payment.PaymentInfo)

	return payment.PaymentInfo, nil
}
