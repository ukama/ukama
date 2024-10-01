/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
)

const (
	publishApiEndpoint = "/api/exchanges/%s/%s/publish"
)

type AmqpClient interface {
	PublishMessage(vhost, exchange, route string, payload *epb.Event) (any, error)
}

type amqpClient struct {
	u *url.URL
	c *httpClient
}

func NewAmqpClient(h, usr, pwd string, timeout time.Duration) AmqpClient {
	u, err := url.ParseRequestURI(h)
	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s",
			h, err.Error())
	}

	headers := map[string]string{"Content-Type": "application/json"}
	return &amqpClient{
		u: u,
		c: NewHttpClient(WithBasicAuth(usr, pwd),
			WithHeaders(headers), WithTimeout(timeout)),
	}
}

func (a *amqpClient) PublishMessage(vhost, exchange, route string, payload *epb.Event) (any, error) {
	fullUrl := a.u.JoinPath(fmt.Sprintf(publishApiEndpoint, vhost, exchange)).String()

	// fmt.Printf("publishing to path: %q\n", fullUrl)
	// fmt.Printf("with routing key : %q\n", route)

	m := &Message{
		RoutingKey:      route,
		Payload:         payload,
		PayloadEncoding: "string",
	}

	p, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal publish message payload: %v. Error: %w", m, err)
	}

	resp, err := a.c.Post(fullUrl, bytes.NewBuffer(p))
	if err != nil {
		return nil, fmt.Errorf("failed to post message to amqp server: %w", err)
	}

	if !((resp.StatusCode >= http.StatusOK) && resp.StatusCode < http.StatusBadRequest) {
		errResp := &ErrorResponse{}

		err = ResponseToJson(resp, errResp)
		if err != nil {
			return nil, fmt.Errorf("fail to unmarshal error response: %w", err)

		}

		return nil, fmt.Errorf("rest api POST failure with error: %w", errResp)
	}

	succResp := &SuccessResponse{}
	err = ResponseToJson(resp, succResp)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal success response: %w", err)

	}

	return succResp, nil
}

type Message struct {
	RoutingKey      string     `json:"routing_key,omitempty"`
	Payload         *epb.Event `json:"payload,omitempty"`
	PayloadEncoding string     `json:"payload_encoding,omitempty"`
}

type SuccessResponse struct {
	Routed bool `json:"routed,omitempty"`
}

type ErrorResponse struct {
	Err    string `json:"error,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("{error: %s, reason: %s}", e.Err, e.Reason)
}

// curl -X POST -u guest:guest 'http://localhost:15672/api/exchanges/%2F/amq.topic/publish' -d '{"properties":{},"routing_key": "event.cloud.local.ukamatestorg.subscriber.registry.subscriber.delete" ,"payload":"{\"SubscriberId\": \"0000\"}","payload_encoding":"string"}' | jq
