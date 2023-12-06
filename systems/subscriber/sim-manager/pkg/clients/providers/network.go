/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package providers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
)

const netEndpoint = "/v1/networks"

type NetworkClientProvider interface {
	GetNetwork(Id string) (*Network, error)
}

type nucleusInfoClient struct {
	R *rest.RestClient
}

type Network struct {
	Id            string    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	OrgId         string    `json:"org_id,omitempty"`
	Overdraft     float64   `json:"overdraft"`
	TrafficPolicy uint32    `json:"traffic_policy"`
	IsDeactivated bool      `json:"is_deactivated,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

type NetworkInfos struct {
	Network *Network `json:"network"`
}

func NewNetworkClientProvider(url string, debug bool) NetworkClientProvider {
	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		log.Fatalf("Can't connect to %s url. Error %s", url, err.Error())
	}

	n := &nucleusInfoClient{
		R: f,
	}

	return n
}

func (p *nucleusInfoClient) GetNetwork(id string) (*Network, error) {
	errStatus := &rest.ErrorMessage{}

	pkg := Network{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + netEndpoint + "/" + id)

	if err != nil {
		log.Errorf("Failed to send api request to nucleus/network. Error %s", err.Error())

		return nil, fmt.Errorf("api request to nucleus system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch network info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("Network Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &pkg)
	if err != nil {
		log.Tracef("Failed to deserialize org info. Error message is %s", err.Error())

		return nil, fmt.Errorf("network info deserailization failure: %w", err)
	}

	log.Infof("Network Info: %+v", pkg)

	return &pkg, nil
}
