/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package operatoragent

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/validation"

	log "github.com/sirupsen/logrus"
)

const (
	OperatorSimsEndpoint   = "/v1/sims"
	OperatorUsagesEndpoint = "/v1/usages"
)

type OperatorAgentClient interface {
	BindSim(iccid string) (*OperatorSimInfo, error)
	GetSimInfo(iccid string) (*OperatorSimInfo, error)
	GetUsages(iccid, cdrType, from, to, region string) (map[string]any, map[string]any, error)
	ActivateSim(iccid string) error
	DeactivateSim(iccid string) error
	TerminateSim(iccid string) error
}

type operatorAgentClient struct {
	u *url.URL
	R *client.Resty
}

func NewOperatorAgentClient(h string, options ...client.Option) *operatorAgentClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &operatorAgentClient{
		u: u,
		R: client.NewResty(options...),
	}
}

// Bind sim is a no-op for operator sims for now
func (o *operatorAgentClient) BindSim(iccid string) (*OperatorSimInfo, error) {
	return &OperatorSimInfo{}, nil
}

func (o *operatorAgentClient) GetSimInfo(iccid string) (*OperatorSimInfo, error) {
	log.Debugf("Getting operator sim info: %v", iccid)

	sim := OperatorSim{}

	resp, err := o.R.Get(o.u.String() + OperatorSimsEndpoint + "/" + iccid)
	if err != nil {
		log.Errorf("GetSimInfo failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetSimInfo failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &sim)
	if err != nil {
		log.Tracef("Failed to deserialize operator sim info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("operator sim info deserialization failure: %w", err)
	}

	log.Infof("Operator Sim Info: %+v", sim.SimInfo)

	return sim.SimInfo, nil
}

func (o *operatorAgentClient) GetUsages(iccid, cdrType, from, to, region string) (map[string]any, map[string]any, error) {
	log.Debugf("Getting operator sim usages: %v", iccid)

	usage := OperatorSimUsage{}

	_, err := validation.FromString(from)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid format for from: %s. Error: %s", from, err)
	}

	_, err = validation.FromString(to)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid format for to: %s. Error: %s", to, err)
	}

	resp, err := o.R.Get(o.u.String() + OperatorUsagesEndpoint +
		fmt.Sprintf("?iccid=%s&cdr_type=%s&from=%s&to=%s&region=%s", iccid, cdrType, from, to, region))
	if err != nil {
		log.Errorf("GetSim usages failure. error: %s", err.Error())

		return nil, nil, fmt.Errorf("GetSim usages failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &usage)
	if err != nil {
		log.Tracef("Failed to deserialize operator sim info. Error message is: %s", err.Error())

		return nil, nil, fmt.Errorf("operator sim info deserialization failure: %w", err)
	}

	log.Infof("Operator data usage (of type %T): %+v", usage.Usage, usage.Usage)
	log.Infof("Operator data cost (of type %T): %+v", usage.Cost, usage.Cost)

	return usage.Usage, usage.Cost, nil
}

func (o *operatorAgentClient) ActivateSim(iccid string) error {
	log.Debugf("Activationg operator sim: %v", iccid)

	_, err := o.R.Put(o.u.String()+OperatorSimsEndpoint+"/"+iccid, nil)
	if err != nil {
		log.Errorf("ActivateSim failure. error: %s", err.Error())

		return fmt.Errorf("ActivateSim failure: %w", err)
	}

	return nil
}

func (o *operatorAgentClient) DeactivateSim(iccid string) error {
	log.Debugf("Deactivating operator sim: %v", iccid)

	_, err := o.R.Patch(o.u.String()+OperatorSimsEndpoint+"/"+iccid, nil)
	if err != nil {
		log.Errorf("DeactivateSim failure. error: %s", err.Error())

		return fmt.Errorf("DeactivateSim failure: %w", err)
	}

	return nil
}

func (o *operatorAgentClient) TerminateSim(iccid string) error {
	log.Debugf("Terminating operator sim: %v", iccid)

	_, err := o.R.Delete(o.u.String() + OperatorSimsEndpoint + "/" + iccid)
	if err != nil {
		log.Errorf("TerminateSim failure. error: %s", err.Error())

		return fmt.Errorf("TerminateSim failure: %w", err)
	}

	return nil
}

type OperatorSimInfo struct {
	Iccid string `json:"iccid"`
	Imsi  string `json:"imsi"`
}

type OperatorSim struct {
	SimInfo *OperatorSimInfo `json:"sim"`
}

type OperatorSimUsage struct {
	Usage map[string]any `json:"usage"`
	Cost  map[string]any `json:"cost"`
}
