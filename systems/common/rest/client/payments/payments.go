/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package payments

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const PaymentEndpoint = "/v1/payments"

type PaymentListReq struct {
	PaymentMethod string `json:"payment_method,omitempty"`
	ItemType      string `json:"item_type,omitempty"`
	Status        string `json:"status,omitempty"`
	ItemId        string `json:"item_id,omitempty"`
	Count         uint32 `json:"count,omitempty"`
	Sort          bool   `json:"sort,omitempty"`
}

type PaymentsRes struct {
	Payments []*Payment `json:"payments"`
}

type Payment struct {
	Id              string `json:"id,omitempty"`
	ItemId          string `json:"item_id,omitempty"`
	ItemType        string `json:"item_type,omitempty"`
	Amount          string `json:"amount,omitempty"`
	Currency        string `json:"currency,omitempty"`
	PaymentMethod   string `json:"payment_method,omitempty"`
	DepositedAmount string `json:"deposited_amount,omitempty"`
	PaidAt          string `json:"paid_at,omitempty"`
	TransactionId   string `json:"transaction_id,omitempty"`
	PayerName       string `json:"payer_name,omitempty"`
	PayerEmail      string `json:"payer_email,omitempty"`
	PayerPhone      string `json:"payer_phone,omitempty"`
	Correspondent   string `json:"correspondent,omitempty"`
	Country         string `json:"country,omitempty"`
	Description     string `json:"description,omitempty"`
	Status          string `json:"status,omitempty"`
	FailureReason   string `json:"failure_reason,omitempty"`
	ExternalId      string `json:"external_id,omitempty"`
	Extra           string `json:"extra,omitempty"`
	Metadata        string `json:"metadata,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
}

type PaymentClient interface {
	List(req PaymentListReq) (PaymentsRes, error)
}

type paymentClient struct {
	u *url.URL
	R *client.Resty
}

func NewPaymentClient(h string, options ...client.Option) *paymentClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &paymentClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (p *paymentClient) List(req PaymentListReq) (PaymentsRes, error) {
	log.Debugf("Listing payments: %v", req)

	paymentRes := PaymentsRes{}

	q := p.u.Query()
	if req.ItemId != "" {
		q.Set("item_id", req.ItemId)
	}
	if req.ItemType != "" {
		q.Set("item_type", req.ItemType)
	}
	if req.PaymentMethod != "" {
		q.Set("payment_method", req.PaymentMethod)
	}
	if req.Status != "" {
		q.Set("status", req.Status)
	}
	if req.Count != 0 {
		q.Set("count", strconv.Itoa(int(req.Count)))
	}
	if req.Sort {
		q.Set("sort", "true")
	}

	fullURL := p.u.String() + PaymentEndpoint + "?" + q.Encode()

	resp, err := p.R.Get(fullURL)
	if err != nil {
		log.Errorf("ListPayments failure. error: %s", err.Error())

		return paymentRes, fmt.Errorf("ListPayments failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &paymentRes)
	if err != nil {
		log.Tracef("Failed to deserialize payments list. Error message is: %s", err.Error())

		return paymentRes, fmt.Errorf("payments list deserailization failure: %w", err)
	}

	log.Infof("Payments List: %+v", paymentRes)

	return paymentRes, nil
}
