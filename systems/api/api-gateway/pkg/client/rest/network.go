/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package rest

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

const NetworkEndpoint = "/v1/networks"

type NetworkInfo struct {
	Id               string    `json:"id,omitempty"`
	Name             string    `json:"name,omitempty"`
	OrgId            string    `json:"org_id,omitempty"`
	IsDeactivated    bool      `json:"is_deactivated,omitempty"`
	AllowedCountries []string  `json:"allowed_countries"`
	AllowedNetworks  []string  `json:"allowed_networks"`
	Budget           float64   `json:"budget"`
	Overdraft        float64   `json:"overdraft"`
	TrafficPolicy    uint32    `json:"traffic_policy"`
	PaymentLinks     bool      `json:"payment_links"`
	SyncStatus       string    `json:"sync_status,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
}

type Network struct {
	NetworkInfo *NetworkInfo `json:"network"`
}

type AddNetworkRequest struct {
	OrgName          string   `json:"org" validate:"required"`
	NetName          string   `json:"network_name" validate:"required"`
	AllowedCountries []string `json:"allowed_countries"`
	AllowedNetworks  []string `json:"allowed_networks"`
	Budget           float64  `json:"budget"`
	Overdraft        float64  `json:"overdraft"`
	TrafficPolicy    uint32   `json:"traffic_policy"`
	PaymentLinks     bool     `json:"payment_links"`
}

type NetworkClient interface {
	Get(Id string) (*NetworkInfo, error)
	Add(req AddNetworkRequest) (*NetworkInfo, error)
}

type networkClient struct {
	u *url.URL
	R *Resty
}

func NewNetworkClient(h string) *networkClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &networkClient{
		u: u,
		R: NewResty(),
	}
}

// TODO check upstream returns payload
func (n *networkClient) Add(req AddNetworkRequest) (*NetworkInfo, error) {
	log.Debugf("Adding network: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	ntwk := Network{}

	resp, err := n.R.Post(n.u.String()+NetworkEndpoint, b)
	if err != nil {
		log.Errorf("AddNetwork failure. error: %s", err.Error())

		return nil, fmt.Errorf("AddNetwork failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &ntwk)
	if err != nil {
		log.Tracef("Failed to deserialize network info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("network info deserailization failure: %w", err)
	}

	log.Infof("Network Info: %+v", ntwk.NetworkInfo)

	return ntwk.NetworkInfo, nil
}

func (n *networkClient) Get(id string) (*NetworkInfo, error) {
	log.Debugf("Getting network: %v", id)

	ntwk := Network{}

	resp, err := n.R.Get(n.u.String() + NetworkEndpoint + "/" + id)
	if err != nil {
		log.Errorf("GetNetwork failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetNetwork failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &ntwk)
	if err != nil {
		log.Tracef("Failed to deserialize network info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("network info deserailization failure: %w", err)
	}

	log.Infof("Network Info: %+v", ntwk.NetworkInfo)

	return ntwk.NetworkInfo, nil
}
