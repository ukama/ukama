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

const SubscribersEndpoint = "/v1/subscribers"

/* TODO-verify: response shape against subscriber api-gateway. */

type SubscriberRecord struct {
	SubscriberId string `json:"subscriber_id"`
	NetworkId    string `json:"network_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	CreatedAt    string `json:"created_at"`
}

type subscribersResponse struct {
	Subscribers []SubscriberRecord `json:"subscribers"`
}

type SubscriberClient interface {
	GetSubscribers() ([]SubscriberRecord, error)
}

type subscriberClient struct {
	u *url.URL
	R *client.Resty
}

func NewSubscriberClient(h string, options ...client.Option) SubscriberClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &subscriberClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (c *subscriberClient) GetSubscribers() ([]SubscriberRecord, error) {
	resp, err := c.R.Get(c.u.String() + SubscribersEndpoint)
	if err != nil {
		return nil, fmt.Errorf("GetSubscribers failure: %w", err)
	}

	out := subscribersResponse{}

	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		return nil, fmt.Errorf("subscribers deserialization failure: %w", err)
	}

	return out.Subscribers, nil
}
