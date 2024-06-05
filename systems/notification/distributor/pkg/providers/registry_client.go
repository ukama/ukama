/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package providers

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	ic "github.com/ukama/ukama/systems/common/initclient"
	"github.com/ukama/ukama/systems/common/rest"
	mpb "github.com/ukama/ukama/systems/registry/member/pb/gen"
	netpb "github.com/ukama/ukama/systems/registry/network/pb/gen"
	"google.golang.org/protobuf/encoding/protojson"
)

const registryVersion = "/v1/"

type RegistryProvider interface {
	GetMember(orgName string, uuid string) (*mpb.MemberResponse, error)
	GetNetwork(orgName string, netID string) (*netpb.GetResponse, error)
}

type registryProvider struct {
	R      *rest.RestClient
	debug  bool
	icHost string
}

func (r *registryProvider) GetRestyClient(org string) (*rest.RestClient, error) {
	/* Add user to member db of the org */
	url, err := ic.GetHostUrl(ic.CreateHostString(org, "registry"), r.icHost, &org, r.debug)
	if err != nil {
		log.Errorf("Failed to resolve registry address to getMember by userId: %v", err)
		return nil, fmt.Errorf("failed to resolve registry address. Error: %v", err)
	}

	rc := rest.NewRestyClient(url, r.debug)

	return rc, nil
}

func NewRegistryProvider(icHost string, debug bool) *registryProvider {

	r := &registryProvider{
		debug:  debug,
		icHost: icHost,
	}

	return r
}

func (r *registryProvider) GetMember(orgName string, userId string) (*mpb.MemberResponse, error) {

	var err error

	/* Get Provider */
	r.R, err = r.GetRestyClient(orgName)
	if err != nil {
		return nil, err
	}

	errStatus := &rest.ErrorMessage{}

	resp, err := r.R.C.R().
		SetError(errStatus).
		Get(r.R.URL.String() + registryVersion + "members/" + userId)

	if err != nil {
		log.Errorf("Failed to send api request to registry at %s . Error %s", r.R.URL.String(), err.Error())
		return nil, fmt.Errorf("api request to registry at %s failure: %v", r.R.URL.String(), err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to get member to registry at %s. HTTP resp code %d and Error message is %s", r.R.URL.String(), resp.StatusCode(), errStatus.Message)
		return nil, fmt.Errorf("failed to get memeber to registry at %s. Error %s", r.R.URL.String(), errStatus.Message)
	}

	memResp := &mpb.MemberResponse{}
	err = protojson.Unmarshal(resp.Body(), memResp)
	if err != nil {
		log.Errorf("Failed to deserialize member response. Error message is %s", err.Error())

		return nil, fmt.Errorf("member response deserialization failure: %w", err)
	}

	return memResp, nil
}

func (r *registryProvider) GetNetwork(orgName string, networkId string) (*netpb.GetResponse, error) {

	var err error

	/* Get Provider */
	r.R, err = r.GetRestyClient(orgName)
	if err != nil {
		return nil, err
	}

	errStatus := &rest.ErrorMessage{}

	resp, err := r.R.C.R().
		SetError(errStatus).
		Get(r.R.URL.String() + registryVersion + "networks/" + networkId)

	if err != nil {
		log.Errorf("Failed to send api request to registry at %s . Error %s", r.R.URL.String(), err.Error())
		return nil, fmt.Errorf("api request to registry at %s failure: %v", r.R.URL.String(), err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to get network from registry at %s. HTTP resp code %d and Error message is %s", r.R.URL.String(), resp.StatusCode(), errStatus.Message)
		return nil, fmt.Errorf("failed to get network from registry at %s. Error %s", r.R.URL.String(), errStatus.Message)
	}

	netResp := &netpb.GetResponse{}
	err = protojson.Unmarshal(resp.Body(), netResp)
	if err != nil {
		log.Errorf("Failed to deserialize network response. Error message is %s", err.Error())

		return nil, fmt.Errorf("network response deserialization failure: %w", err)
	}

	return netResp, nil
}
