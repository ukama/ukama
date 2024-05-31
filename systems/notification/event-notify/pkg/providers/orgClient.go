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
	orgpb "github.com/ukama/ukama/systems/nucleus/org/pb/gen"
	userspb "github.com/ukama/ukama/systems/nucleus/user/pb/gen"
)

const nucleusVersion = "/v1/"

type NucleusProvider interface {
	Whoami(uuid string) (*userspb.WhoamiResponse, error)
	GetOrg(orgName string) (*orgpb.GetByNameResponse, error)
}

type nucleusProvider struct {
	R      *rest.RestClient
	debug  bool
	icHost string
}

func (n *nucleusProvider) GetRestyClient(org string) (*rest.RestClient, error) {
	url, err := ic.GetHostUrl(ic.CreateHostString(org, "nucleus"), n.icHost, &org, n.debug)
	if err != nil {
		log.Errorf("Failed to resolve nucleus address: %v", err)
		return nil, fmt.Errorf("failed to resolve nucleus address. Error: %v", err)
	}

	rc := rest.NewRestyClient(url, n.debug)

	return rc, nil
}

func NewNucleusProvider(icHost string, debug bool) *nucleusProvider {

	r := &nucleusProvider{
		debug:  debug,
		icHost: icHost,
	}

	return r
}

func (n *nucleusProvider) Whoami(orgName string, uuid string) (*userspb.WhoamiResponse, error) {

	var err error

	/* Get Provider */
	n.R, err = n.GetRestyClient(orgName)
	if err != nil {
		return nil, err
	}

	errStatus := &rest.ErrorMessage{}

	resp, err := n.R.C.R().
		SetError(errStatus).
		Get(n.R.URL.String() + nucleusVersion + "users/whoami/" + uuid)

	if err != nil {
		log.Errorf("Failed to send api request to nucleus at %s . Error %s", n.R.URL.String(), err.Error())
		return nil, fmt.Errorf("api request to nucleus at %s failure: %v", n.R.URL.String(), err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to call whoami to nucleus at %s. HTTP resp code %d and Error message is %s", n.R.URL.String(), resp.StatusCode(), errStatus.Message)
		return nil, fmt.Errorf("failed to call whoami to nucleus at %s. Error %s", n.R.URL.String(), errStatus.Message)
	}

	whoamiResp := &userspb.WhoamiResponse{}
	err = json.Unmarshal(resp.Body(), whoamiResp)
	if err != nil {
		log.Errorf("Failed to deserialize whoami response. Error message is %s", err.Error())

		return nil, fmt.Errorf("whoami response deserialization failure: %w", err)
	}

	return whoamiResp, nil
}

func (n *nucleusProvider) GetOrg(orgName string) (*orgpb.GetByNameResponse, error) {

	var err error

	/* Get Provider */
	n.R, err = n.GetRestyClient(orgName)
	if err != nil {
		return nil, err
	}

	errStatus := &rest.ErrorMessage{}

	resp, err := n.R.C.R().
		SetError(errStatus).
		Get(n.R.URL.String() + nucleusVersion + "orgs/" + orgName)

	if err != nil {
		log.Errorf("Failed to send api request to nucleus at %s . Error %s", n.R.URL.String(), err.Error())
		return nil, fmt.Errorf("api request to nucleus at %s failure: %v", n.R.URL.String(), err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to call getOrg to nucleus at %s. HTTP resp code %d and Error message is %s", n.R.URL.String(), resp.StatusCode(), errStatus.Message)
		return nil, fmt.Errorf("failed to call getOrg to nucleus at %s. Error %s", n.R.URL.String(), errStatus.Message)
	}

	orgResponse := &orgpb.GetByNameResponse{}
	err = json.Unmarshal(resp.Body(), orgResponse)
	if err != nil {
		log.Errorf("Failed to deserialize orgByName response. Error message is %s", err.Error())

		return nil, fmt.Errorf("orgByName response deserialization failure: %w", err)
	}

	return orgResponse, nil
}
