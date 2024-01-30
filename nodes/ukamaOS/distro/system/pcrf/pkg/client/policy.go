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

	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/rest"

	log "github.com/sirupsen/logrus"
)

const PolicyEndpoint = "/v1/policy/imsi/"
const CDREndpoint = "/v1/cdr/imsi/"

type PolicyController interface {
	GetPolicy(imsi string) (*rest.Policy, error)
	PushCdr(cdr rest.CDR) error
}

type policyControllerClient struct {
	u *url.URL
	R *Resty
}

func NewPolicyControllerClient(h string) *policyControllerClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &policyControllerClient{
		u: u,
		R: NewResty(),
	}
}

func (p *policyControllerClient) PushCdr(req rest.CDR) error {
	log.Debugf("Posting  CDR: %v", req)

	url := p.u.String() + "/" + CDREndpoint + req.Imsi

	b, err := json.Marshal(req)
	if err != nil {
		log.Errorf("Error marshalling CDR. error: %s", err.Error())
		return fmt.Errorf("Marshal CDR request failure for imsi : %d. Error %s", req.Imsi, err.Error())
	}

	resp, err := p.R.Post(url, b)
	if err != nil {
		log.Errorf("Post CDR failure. error: %s", err.Error())
		return fmt.Errorf("Post CDR failure: %w", err)
	}

	return nil
}

func (p *policyControllerClient) GetPolicy(imsi string) (*rest.Policy, error) {
	log.Debugf("Getting policy for ismi: %s", imsi)

	policy := &rest.Policy{}
	resp, err := p.R.GetPolicy(p.u.String() + PolicyEndpoint + "/" + imsi)
	if err != nil {
		log.Errorf("GetPolicy failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetPolicy failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &policy)
	if err != nil {
		log.Tracef("Failed to deserialize policy info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("policy info deserailization failure: %w", err)
	}

	log.Infof("Policy Info: %+v", policy)

	return policy, nil
}
