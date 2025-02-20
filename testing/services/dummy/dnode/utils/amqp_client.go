/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ukama/ukama/testing/services/dummy/dnode/config"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	publishApiEndpoint = "/api/exchanges/%s/%s/publish"
)

type AmqpClient interface {
	PublishMessage(vhost, exchange, route string, payload *anypb.Any) (any, error)
}

type amqpClient struct {
	u *url.URL
	c *HttpClient
}

func NewAmqpClient(amqpConf config.AmqpConfig, timeout time.Duration) AmqpClient {
	u, err := url.ParseRequestURI(amqpConf.Uri)
	if err != nil {
		logrus.Errorf("Can't parse  %s url. Error: %s",
			amqpConf.Uri, err.Error())
	}

	headers := map[string]string{"Content-Type": "application/json"}
	return &amqpClient{
		u: u,
		c: NewHttpClient(WithBasicAuth(amqpConf.Username, amqpConf.Password),
			WithHeaders(headers), WithTimeout(timeout)),
	}
}

func (a *amqpClient) PublishMessage(vhost, exchange, route string, payload *anypb.Any) (any, error) {
	fullURL := a.u.JoinPath(fmt.Sprintf(publishApiEndpoint, vhost, exchange)).String()

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload message: %v. Error: %w",
			payload, err)
	}

	encodedPayload := base64.StdEncoding.EncodeToString(payloadBytes)

	m := &Message{
		RoutingKey:      route,
		Payload:         encodedPayload,
		PayloadEncoding: "base64",
	}

	p, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal publish message: %v. Error: %w", m, err)
	}

	resp, err := a.c.Post(fullURL, bytes.NewBuffer(p))
	if err != nil {
		return nil, fmt.Errorf("failed to post message to amqp server: %w", err)
	}

	if !((resp.StatusCode >= http.StatusOK) && resp.StatusCode < http.StatusBadRequest) {
		errResp := &ErrorResponse{}

		err = DecodeJSONResponse(resp, errResp)
		if err != nil {
			return nil, fmt.Errorf("fail to unmarshal error response: %w", err)

		}

		return nil, fmt.Errorf("rest api POST failure with error: %w", errResp)
	}

	succResp := &SuccessResponse{}
	err = DecodeJSONResponse(resp, succResp)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal success response: %w", err)

	}

	return succResp, nil
}

type Message struct {
	Properties      struct{} `json:"properties"`
	RoutingKey      string   `json:"routing_key"`
	Payload         string   `json:"payload"`
	PayloadEncoding string   `json:"payload_encoding"`
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
