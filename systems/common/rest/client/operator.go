/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
)

const OperatorEndpoint = "/v1/sims"

type OperatorSimInfo struct {
	Iccid string `json:"iccid"`
	Imsi  string `json:"imsi"`
}

type OperatorSim struct {
	SimInfo *OperatorSimInfo `json:"Sim"`
}

type OperatorClient interface {
	BindSim(iccid string) (*OperatorSimInfo, error)
	GetSimInfo(iccid string) (*OperatorSimInfo, error)
	ActivateSim(iccid string) error
	DeactivateSim(iccid string) error
	TerminateSim(iccid string) error
}

type operatorClient struct {
	u *url.URL
	R *Resty
}

func NewOperatorClient(h string) *operatorClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &operatorClient{
		u: u,
		R: NewResty(),
	}
}

// Bind sim is a no-op for operator sims for now
func (o *operatorClient) BindSim(iccid string) (*OperatorSimInfo, error) {
	return &OperatorSimInfo{}, nil
}

func (o *operatorClient) GetSimInfo(iccid string) (*OperatorSimInfo, error) {
	log.Debugf("Getting operator sim info: %v", iccid)

	sim := OperatorSim{}

	resp, err := o.R.Get(o.u.String() + OperatorEndpoint + "/" + iccid)
	if err != nil {
		log.Errorf("GetSimInfo failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetSimInfo failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &sim)
	if err != nil {
		log.Tracef("Failed to deserialize operator sim info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("operator sim info deserailization failure: %w", err)
	}

	log.Infof("Operator Sim Info: %+v", sim.SimInfo)

	return sim.SimInfo, nil
}

func (o *operatorClient) ActivateSim(iccid string) error {
	log.Debugf("Activationg operator sim: %v", iccid)

	_, err := o.R.Put(o.u.String()+OperatorEndpoint+"/"+iccid, nil)
	if err != nil {
		log.Errorf("ActivateSim failure. error: %s", err.Error())

		return fmt.Errorf("ActivateSim failure: %w", err)
	}

	return nil
}

func (o *operatorClient) DeactivateSim(iccid string) error {
	log.Debugf("Deactivating operator sim: %v", iccid)

	_, err := o.R.Patch(o.u.String()+OperatorEndpoint+"/"+iccid, nil)
	if err != nil {
		log.Errorf("DeactivateSim failure. error: %s", err.Error())

		return fmt.Errorf("DeactivateSim failure: %w", err)
	}

	return nil
}

func (o *operatorClient) TerminateSim(iccid string) error {
	log.Debugf("Terminating operator sim: %v", iccid)

	_, err := o.R.Delete(o.u.String() + OperatorEndpoint + "/" + iccid)
	if err != nil {
		log.Errorf("TerminateSim failure. error: %s", err.Error())

		return fmt.Errorf("TerminateSim failure: %w", err)
	}

	return nil
}
