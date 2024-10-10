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
	"net/http"
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
	c *HttpClient
}

func NewPaymentsClient(h string) *paymentsClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	headers := map[string]string{"Content-Type": "application/json"}

	return &paymentsClient{
		u: u,
		R: client.NewResty(client.WithError(&Err{}),
			client.WithDebug(), client.WithContentTypeJSON()),
		// client.WithContentTypeJSON()),

		c: NewHttpClient(WithHeaders(headers)),
	}
}

// func (p *paymentsClient) ListPayments2(queryString string) ([]*PaymentInfo, error) {
// log.Infof("Listing payments matching: %v", queryString)

// payments := &Payment{}

// resp, err := p.R.Get(p.u.String() + PaymentEndpoint + queryString)
// if err != nil {
// log.Errorf("ListPayments failure. error: %s", err.Error())

// return nil, fmt.Errorf("ListPayments failure: %w", err)
// }

// log.Infof("Header: %v", resp.Header())
// log.Infof("IsSuccess: %v", resp.IsSuccess())
// log.Infof("IsError: %v", resp.IsError())
// log.Infof("Size: %v", resp.Size())
// log.Infof("Status: %v", resp.Status())
// log.Infof("String: %v", resp.String())

// b := resp.RawBody()
// defer b.Close()

// data, err := io.ReadAll(b)
// if err != nil {
// return nil, fmt.Errorf("read response body failure: %w", err)
// }

// if data == nil {
// return nil, fmt.Errorf("data is nil after reading response body: %v", data)
// }

// // log.Infof("Listing response body: %v", resp.Body())

// err = json.Unmarshal(data, &payments)
// if err != nil {
// log.Tracef(deserializeLogMsg, "payments", err.Error())

// return nil, fmt.Errorf(deserializeErrorMsg, "payments", err)
// }

// log.Infof(resourceLogMsg, "Payments", payments.Payments)

// return payments.Payments, nil
// }

func (p *paymentsClient) ListPayments(queryString string) ([]*PaymentInfo, error) {
	log.Infof("Listing payments matching: %v", queryString)

	resp, err := p.c.Get(p.u.String() + PaymentEndpoint + queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}

	if !((resp.StatusCode >= http.StatusOK) && resp.StatusCode < http.StatusBadRequest) {
		errResp := &ErrorResponse{}

		err = DecodeJSONResponse(resp, errResp)
		if err != nil {
			log.Errorf("fail to unmarshal error response: %v", err)

			return nil, fmt.Errorf("fail to unmarshal error response: %w", err)
		}

		log.Errorf("rest api GET failure with error: %v", errResp)

		return nil, fmt.Errorf("rest api GET failure with error: %w", errResp)
	}

	payments := &Payment{}
	err = DecodeJSONResponse(resp, payments)
	if err != nil {
		log.Errorf("fail to unmarshal payment response: %v", err)

		return nil, fmt.Errorf("fail to unmarshal payment response: %w", err)
	}

	log.Infof(resourceLogMsg, "Payments", payments.Payments)

	return payments.Payments, nil
}

type Payment struct {
	Payments []*PaymentInfo `json:"payments"`
}

type PaymentInfo struct {
	Id              string  `json:"id,omitempty"`
	ItemId          string  `json:"item_id,omitempty"`
	ItemType        string  `json:"item_type,omitempty"`
	Amount          float64 `json:"amount,string,omitempty"`
	DepositedAmount float64 `json:"deposited_amount,string,omitempty"`
	Currency        string  `json:"currency,omitempty"`
	PaymentMethod   string  `json:"payment_method,omitempty"`
	PayerName       string  `json:"payer_name,omitempty"`
	PayerEmail      string  `json:"payer_email,omitempty"`
	PayerPhone      string  `json:"payer_phone,omitempty"`
	Correspondent   string  `json:"correspondent,omitempty"`
	Country         string  `json:"country,omitempty"`
	Description     string  `json:"description,omitempty"`
	Status          string  `json:"status,omitempty"`
	FailureReason   string  `json:"failure_reason,omitempty"`
	// CreatedAt     time.Time `json:"created_at,omitempty"`
	// PaidAt               *time.Time `json:"paid_at,omitempty"`
}

type ErrorResponse struct {
	Err    string `json:"error,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("{error: %s, reason: %s}", e.Err, e.Reason)
}
