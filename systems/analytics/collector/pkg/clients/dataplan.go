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

const PackagesEndpoint = "/v1/packages"

/* TODO-verify: response shape against dataplan api-gateway. */

type DataplanPackage struct {
	Uuid       string  `json:"uuid"`
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	Duration   uint64  `json:"duration"`
	DataVolume int64   `json:"data_volume"`
	DataUnit   string  `json:"data_unit"`
	IsActive   bool    `json:"active"`
}

type packagesResponse struct {
	Packages []DataplanPackage `json:"packages"`
}

type DataplanClient interface {
	GetPackages() ([]DataplanPackage, error)
}

type dataplanClient struct {
	u *url.URL
	R *client.Resty
}

func NewDataplanClient(h string, options ...client.Option) DataplanClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &dataplanClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (c *dataplanClient) GetPackages() ([]DataplanPackage, error) {
	resp, err := c.R.Get(c.u.String() + PackagesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("GetPackages failure: %w", err)
	}

	out := packagesResponse{}

	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		return nil, fmt.Errorf("packages deserialization failure: %w", err)
	}

	return out.Packages, nil
}
