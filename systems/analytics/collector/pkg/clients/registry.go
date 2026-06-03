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

const (
	NetworksEndpoint = "/v1/networks"
	SitesEndpoint    = "/v1/sites"
)

/* TODO-verify: response shapes against registry api-gateway. */

type RegistryNetwork struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	IsDeactivated bool   `json:"is_deactivated"`
}

type RegistrySite struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	NetworkId     string `json:"network_id"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	IsDeactivated bool   `json:"is_deactivated"`
}

type networksResponse struct {
	Networks []RegistryNetwork `json:"networks"`
}

type sitesResponse struct {
	Sites []RegistrySite `json:"sites"`
}

type RegistryClient interface {
	GetNetworks() ([]RegistryNetwork, error)
	GetSites() ([]RegistrySite, error)
}

type registryClient struct {
	u *url.URL
	R *client.Resty
}

func NewRegistryClient(h string, options ...client.Option) RegistryClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &registryClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (c *registryClient) GetNetworks() ([]RegistryNetwork, error) {
	resp, err := c.R.Get(c.u.String() + NetworksEndpoint)
	if err != nil {
		return nil, fmt.Errorf("GetNetworks failure: %w", err)
	}

	out := networksResponse{}

	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		return nil, fmt.Errorf("networks deserialization failure: %w", err)
	}

	return out.Networks, nil
}

func (c *registryClient) GetSites() ([]RegistrySite, error) {
	resp, err := c.R.Get(c.u.String() + SitesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("GetSites failure: %w", err)
	}

	out := sitesResponse{}

	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		return nil, fmt.Errorf("sites deserialization failure: %w", err)
	}

	return out.Sites, nil
}
