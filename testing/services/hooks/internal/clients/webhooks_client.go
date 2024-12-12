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
	"github.com/ukama/ukama/testing/services/hooks/util"

	log "github.com/sirupsen/logrus"
)

const (
	MopayHooksEndpoint  = "/v1/pawapay"
	StripeHooksEndpoint = "/v1/stripe"
)

type WebhooksClient interface {
	PostDepositHook(*util.Deposit) (*WebhookInfo, error)
	PostPaymentIntentHook(*util.PaymentIntent) (*WebhookInfo, error)
}

type webhooksClient struct {
	u *url.URL
	R *client.Resty
}

func NewWebhooksClient(h string) *webhooksClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &webhooksClient{
		u: u,
		R: client.NewResty(client.WithError(&Err{}),
			client.WithDebug(), client.WithContentTypeJSON()),
	}
}

func (p *webhooksClient) PostDepositHook(req *util.Deposit) (*WebhookInfo, error) {
	log.Debugf("Posting deposit webhook response: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf(requestMarshalErrorMsg, err)
	}

	webhook := Webhook{}

	resp, err := p.R.Post(p.u.String()+MopayHooksEndpoint, b)
	if err != nil {
		log.Errorf("PostDepositHook failure. error: %s", err.Error())

		return nil, fmt.Errorf("PostDepositHook failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &webhook)
	if err != nil {
		log.Tracef(deserializeLogMsg, "webhook", err.Error())

		return nil, fmt.Errorf(deserializeErrorMsg, "webhook", err)
	}

	log.Infof(resourceLogMsg, "Webhook", webhook.WebhookInfo)

	return webhook.WebhookInfo, nil
}

func (p *webhooksClient) PostPaymentIntentHook(req *util.PaymentIntent) (*WebhookInfo, error) {
	log.Debugf("Posting payment intent webhook response: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf(requestMarshalErrorMsg, err)
	}

	webhook := Webhook{}

	resp, err := p.R.Post(p.u.String()+StripeHooksEndpoint, b)
	if err != nil {
		log.Errorf("PostPaymentIntentHook failure. error: %s", err.Error())

		return nil, fmt.Errorf("PostPaymentIntentHook failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &webhook)
	if err != nil {
		log.Tracef(deserializeLogMsg, "webhook", err.Error())

		return nil, fmt.Errorf(deserializeErrorMsg, "webhook", err)
	}

	log.Infof(resourceLogMsg, "Webhook", webhook.WebhookInfo)

	return webhook.WebhookInfo, nil
}

type Webhook struct {
	WebhookInfo *WebhookInfo `json:"webhook"`
}

type WebhookInfo struct {
	Id      string `json:"id,omitempty"`
	OrgName string `json:"org_name,omitempty"`
	Payload string `json:"payload"`
}
