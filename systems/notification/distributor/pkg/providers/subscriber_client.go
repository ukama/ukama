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

	log "github.com/sirupsen/logrus"
	ic "github.com/ukama/ukama/systems/common/initclient"
	"github.com/ukama/ukama/systems/common/rest"
	spb "github.com/ukama/ukama/systems/subscriber/registry/pb/gen"
)

const subscriberVersion = "/v1/"

type SubscriberProvider interface {
	GetSubscriber(orgName string, subscriberId string) (*spb.GetSubscriberResponse, error)
}

type subscriberProvider struct {
	R      *rest.RestClient
	debug  bool
	icHost string
}

func (r *subscriberProvider) GetRestyClient(org string) (*rest.RestClient, error) {
	/* Add user to member db of the org */
	url, err := ic.GetHostUrl(ic.CreateHostString(org, "subscriber"), r.icHost, &org, r.debug)
	if err != nil {
		log.Errorf("Failed to resolve subscriber address to getSubscriber by subId: %v", err)
		return nil, fmt.Errorf("failed to resolve subscriber address. Error: %v", err)
	}

	rc := rest.NewRestyClient(url, r.debug)

	return rc, nil
}

func NewSubscriberProvider(icHost string, debug bool) *registryProvider {

	r := &registryProvider{
		debug:  debug,
		icHost: icHost,
	}

	return r
}

func (r *registryProvider) GetSubscriber(orgName string, subscriberId string) (*spb.GetSubscriberResponse, error) {

	var err error

	/* Get Provider */
	r.R, err = r.GetRestyClient(orgName)
	if err != nil {
		return nil, err
	}

	errStatus := &rest.ErrorMessage{}

	resp, err := r.R.C.R().
		SetError(errStatus).
		Get(r.R.URL.String() + subscriberVersion + "subscriber/" + subscriberId)

	if err != nil {
		log.Errorf("Failed to send api request to subscriber at %s . Error %s", r.R.URL.String(), err.Error())
		return nil, fmt.Errorf("api request to subscriber at %s failure: %v", r.R.URL.String(), err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to get subscriber to subscriber at %s. HTTP resp code %d and Error message is %s", r.R.URL.String(), resp.StatusCode(), errStatus.Message)
		return nil, fmt.Errorf("failed to get subscriber to subscriber at %s. Error %s", r.R.URL.String(), errStatus.Message)
	}

	subResp := &spb.GetSubscriberResponse{}
	err = json.Unmarshal(resp.Body(), subResp)
	if err != nil {
		log.Errorf("Failed to deserialize subscriber response. Error message is %s", err.Error())

		return nil, fmt.Errorf("subscriber response deserialization failure: %w", err)
	}

	return subResp, nil
}
