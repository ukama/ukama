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
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller"
)

const PolicyEndpoint = "/v1/policy/imsi/"
const CDREndpoint = "/v1/cdr/imsi/"

type RemoteController interface {
	GetPolicy(imsi string) (*rest.Policy, error)
	PushCdr(cdr rest.CDR) error
}

type remoteControllerClient struct {
	u *url.URL
	R *Resty
}

func NewRemoteControllerClient(h string) (*remoteControllerClient, error) {
	u, err := url.Parse(h)

	if err != nil {
		log.Errorf("Can't parse  %s url. Error: %s", h, err.Error())
		return nil, err
	}

	return &remoteControllerClient{
		u: u,
		R: NewResty(),
	}, nil
}

func (r *remoteControllerClient) PushCdr(req controller.CDR) error {
	log.Debugf("Posting  CDR: %v", req)

	url := r.u.String() + "/" + CDREndpoint + req.Imsi

	b, err := json.Marshal(req)
	if err != nil {
		log.Errorf("Error marshalling CDR. error: %s", err.Error())
		return fmt.Errorf("Marshal CDR request failure for imsi : %d. Error %s", req.Imsi, err.Error())
	}

	resp, err := r.R.Post(url, b)
	if err != nil {
		log.Errorf("Post CDR failure. error: %s", err.Error())
		return fmt.Errorf("Post CDR failure: %w", err)
	}

	return nil
}

func (r *remoteControllerClient) GetPolicy(imsi string) (*rest.Policy, error) {
	log.Debugf("Getting policy for ismi: %s", imsi)

	policy := &rest.Policy{}
	resp, err := r.R.GetPolicy(r.u.String() + PolicyEndpoint + "/" + imsi)
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
