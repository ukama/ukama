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
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/api"
	"github.com/ukama/ukama/systems/common/rest"
)

const ProfileEndpoint = "/v1/policy/imsi"
const CDREndpoint = "/v1/cdr/imsi"

type RemoteController interface {
	GetSubscriberProfile(imsi string) (*api.Spr, error)
	PushCdr(cdr *api.CDR) error
}

type remoteControllerClient struct {
	u     *url.URL
	R     *rest.RestClient
	debug bool
}

func NewRemoteControllerClient(h string, debug bool) (*remoteControllerClient, error) {
	u, err := url.Parse(h)

	if err != nil {
		log.Errorf("Can't parse  %s url. Error: %s", h, err.Error())
		return nil, err
	}

	return &remoteControllerClient{
		u:     u,
		R:     rest.NewRestyClient(u, debug),
		debug: debug,
	}, nil
}

func (r *remoteControllerClient) PushCdr(req *api.CDR) error {
	log.Debugf("Posting  CDR: %v", req)

	url := r.u.String() + "/" + CDREndpoint + req.Imsi

	b, err := json.Marshal(req)
	if err != nil {
		log.Errorf("Error marshalling CDR. error: %s", err.Error())
		return fmt.Errorf("marshal CDR request failure for imsi : %s. Error %s", req.Imsi, err.Error())
	}

	_, err = r.R.C.R().
		SetBody(b).
		Post(url)
	if err != nil {
		log.Errorf("Post CDR failure. error: %s", err.Error())
		return fmt.Errorf("post CDR failure: %s", err.Error())
	}

	return nil
}

func (r *remoteControllerClient) GetSubscriberProfile(imsi string) (*api.Spr, error) {
	log.Debugf("Getting policy for ismi: %s", imsi)

	spr := &api.Spr{}
	resp, err := r.R.C.R().Get(r.u.String() + ProfileEndpoint + "/" + imsi)
	if err != nil {
		log.Errorf("GetPolicy failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetPolicy failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &spr)
	if err != nil {
		log.Tracef("Failed to deserialize policy info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("policy info deserailization failure: %w", err)
	}

	log.Infof("SPR Info: %+v", spr)

	return spr, nil
}
