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

const SimEndpoint = "/v1/sims"

type SimInfo struct {
	Id            string    `json:"id,omitempty"`
	NetworkId     string    `json:"network_id,omitempty"`
	SubscriberId  string    `json:"subscriber_id,omitempty"`
	Iccid         string    `json:"iccid,omitempty"`
	Msisdn        string    `json:"msisdn,omitempty"`
	Imsi          string    `json:"imsi,omitempty"`
	Status        string    `json:"status,omitempty"`
	SimType       string    `json:"sim_type,omitempty"`
	IsPhysical    bool      `json:"is_physical,omitempty"`
	TrafficPolicy uint32    `json:"traffic_policy"`
	SyncStatus    string    `json:"sync_status,omitempty"`
	AllocatedAt   time.Time `json:"allocated_at,omitempty"`
}

type Sim struct {
	SimInfo *SimInfo `json:"sim"`
}

type AddSimRequest struct {
	SubscriberId  string `json:"subscriber_id" validate:"required"`
	NetworkId     string `json:"network_id" validate:"required"`
	PackageId     string `json:"package_id" validate:"required"`
	SimType       string `json:"sim_type"`
	SimToken      string `json:"sim_token"`
	TrafficPolicy uint32 `json:"traffic_policy"`
}

type SimClient interface {
	Get(Id string) (*SimInfo, error)
	Add(req AddSimRequest) (*SimInfo, error)
}

type simClient struct {
	u *url.URL
	R *Resty
}

func NewSimClient(h string) *simClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &simClient{
		u: u,
		R: NewResty(),
	}
}

// TODO check upstream returns payload
func (s *simClient) Add(req AddSimRequest) (*SimInfo, error) {
	log.Debugf("Adding sim: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}

	sim := Sim{}

	resp, err := s.R.Post(s.u.String()+SimEndpoint, b)
	if err != nil {
		log.Errorf("AddSim failure. error: %s", err.Error())

		return nil, fmt.Errorf("AddSim failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &sim)
	if err != nil {
		log.Tracef("Failed to deserialize sim info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("sim info deserailization failure: %w", err)
	}

	log.Infof("Sim Info: %+v", sim.SimInfo)

	return sim.SimInfo, nil
}

func (s *simClient) Get(id string) (*SimInfo, error) {
	log.Debugf("Getting sim: %v", id)

	sim := Sim{}

	resp, err := s.R.Get(s.u.String() + SimEndpoint + "/" + id)
	if err != nil {
		log.Errorf("GetSim failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetSim failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &sim)
	if err != nil {
		log.Tracef("Failed to deserialize sim info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("sim info deserailization failure: %w", err)
	}

	log.Infof("Sim Info: %+v", sim.SimInfo)

	return sim.SimInfo, nil
}
