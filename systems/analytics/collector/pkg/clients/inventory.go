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

const ComponentsEndpoint = "/v1/components"

/* TODO-verify: response shape against inventory api-gateway. */

type InventoryComponent struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	PartNumber string `json:"part_number"`
	Inventory  string `json:"inventory"`
}

type componentsResponse struct {
	Components []InventoryComponent `json:"components"`
}

type InventoryClient interface {
	GetComponents() ([]InventoryComponent, error)
}

type inventoryClient struct {
	u *url.URL
	R *client.Resty
}

func NewInventoryClient(h string, options ...client.Option) InventoryClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &inventoryClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (c *inventoryClient) GetComponents() ([]InventoryComponent, error) {
	resp, err := c.R.Get(c.u.String() + ComponentsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("GetComponents failure: %w", err)
	}

	out := componentsResponse{}

	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		return nil, fmt.Errorf("components deserialization failure: %w", err)
	}

	return out.Components, nil
}
