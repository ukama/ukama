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

const MetricsEndpoint = "/v1/metrics"

/* TODO-verify: response shape against metrics api-gateway. */

type MetricValue struct {
	Metric       string  `json:"metric"`
	ResourceType string  `json:"resource_type"`
	ResourceId   string  `json:"resource_id"`
	Value        float64 `json:"value"`
	Unit         string  `json:"unit"`
	Timestamp    int64   `json:"timestamp"`
}

type metricsResponse struct {
	Metrics []MetricValue `json:"metrics"`
}

type MetricsClient interface {
	GetLatestMetrics() ([]MetricValue, error)
}

type metricsClient struct {
	u *url.URL
	R *client.Resty
}

func NewMetricsClient(h string, options ...client.Option) MetricsClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &metricsClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (c *metricsClient) GetLatestMetrics() ([]MetricValue, error) {
	resp, err := c.R.Get(c.u.String() + MetricsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("GetLatestMetrics failure: %w", err)
	}

	out := metricsResponse{}

	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		return nil, fmt.Errorf("metrics deserialization failure: %w", err)
	}

	return out.Metrics, nil
}
