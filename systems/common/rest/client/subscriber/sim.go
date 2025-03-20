/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package subscriber

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ukama/ukama/systems/common/rest/client"

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
	Type          string    `json:"type,omitempty"`
	IsPhysical    bool      `json:"is_physical,omitempty"`
	TrafficPolicy uint32    `json:"traffic_policy"`
	SyncStatus    string    `json:"sync_status,omitempty"`
	AllocatedAt   time.Time `json:"allocated_at,omitempty"`
	Package       *Pacakge  `json:"package,omitempty"`
}

type Pacakge struct {
	Id        string `json:"id,omitempty"`
	PackageId string `json:"package_id,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	IsActive  bool   `json:"is_active,omitempty"`
	AsExpired bool   `json:"as_expired,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type SimList struct {
	Sims []*SimInfo `json:"sims"`
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
type ListSimRequest struct {
	ICCID         string `json:"iccid,omitempty"`
	Imsi          string `json:"imsi,omitempty"`
	SubscriberId  string `json:"subscriber_id,omitempty"`
	NetworkId     string `json:"network_id,omitempty"`
	SimType       string `json:"sim_type,omitempty"`
	SimStatus     string `json:"sim_status,omitempty"`
	TrafficPolicy uint32 `json:"traffic_policy,omitempty"`
	IsPhysical    bool   `json:"is_physical,omitempty"`
	Count         uint32 `json:"count,omitempty"`
	Sort          bool   `json:"sort,omitempty"`
}

type SimClient interface {
	List(req ListSimRequest) (SimList, error)
	Get(Id string) (*SimInfo, error)
	Add(req AddSimRequest) (*SimInfo, error)
}

type simClient struct {
	u *url.URL
	R *client.Resty
}

func NewSimClient(h string, options ...client.Option) *simClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &simClient{
		u: u,
		R: client.NewResty(options...),
	}
}

// TODO check upstream returns payload
func (s *simClient) Add(req AddSimRequest) (*SimInfo, error) {
	log.Debugf("Adding sim: %v", req)

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %w", err)
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

func (s *simClient) List(req ListSimRequest) (SimList, error) {
	log.Debugf("Listing sims: %v", req)

	simList := SimList{}

	q := s.u.Query()
	if req.ICCID != "" {
		q.Set("iccid", req.ICCID)
	}
	if req.Imsi != "" {
		q.Set("imsi", req.Imsi)
	}
	if req.SubscriberId != "" {
		q.Set("subscriber_id", req.SubscriberId)
	}
	if req.NetworkId != "" {
		q.Set("network_id", req.NetworkId)
	}
	if req.SimType != "" {
		q.Set("sim_type", req.SimType)
	}
	if req.SimStatus != "" {
		q.Set("sim_status", req.SimStatus)
	}
	if req.TrafficPolicy != 0 {
		q.Set("traffic_policy", fmt.Sprintf("%d", req.TrafficPolicy))
	}
	if req.IsPhysical {
		q.Set("is_physical", fmt.Sprintf("%t", req.IsPhysical))
	}
	if req.Count != 0 {
		q.Set("count", fmt.Sprintf("%d", req.Count))
	}
	if req.Sort {
		q.Set("sort", fmt.Sprintf("%t", req.Sort))
	}

	fullURL := s.u.String() + "/v1/sim" + "?" + q.Encode()

	resp, err := s.R.Get(fullURL)
	if err != nil {
		log.Errorf("ListSim failure. error: %s", err.Error())

		return simList, fmt.Errorf("ListSim failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &simList)
	if err != nil {
		log.Tracef("Failed to deserialize sim list. Error message is: %s", err.Error())

		return simList, fmt.Errorf("sim list deserailization failure: %w", err)
	}

	log.Infof("Sim List: %+v", simList)

	return simList, nil
}
