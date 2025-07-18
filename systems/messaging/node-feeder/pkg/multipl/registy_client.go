/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

//TODO: we should use registry rest client in common instead.

package multipl

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
	nodepb "github.com/ukama/ukama/systems/registry/node/pb/gen"
)

const (
	RegistryVersion = "/v1/"
	SystemName      = "registry"
)

type RoleType int32

const (
	RoleType_OWNER RoleType = iota
	RoleType_ADMIN
	RoleType_EMPLOYEE
	RoleType_VENDOR
	RoleType_USERS
)

// Enum value maps for RoleType.
var (
	RoleType_name = map[int32]string{
		0: "OWNER",
		1: "ADMIN",
		2: "EMPLOYEE",
		3: "VENDOR",
		4: "USERS",
	}
	RoleType_value = map[string]int32{
		"OWNER":    0,
		"ADMIN":    1,
		"EMPLOYEE": 2,
		"VENDOR":   3,
		"USERS":    4,
	}
)

type RegistryProvider interface {
	GetAllNodes(org string) (*nodepb.GetNodesResponse, error)
}

type registryProvider struct {
	R       *rest.RestClient
	debug   bool
	icHost  string
	timeout time.Duration
}

type OrgMember struct {
	UserUuid string `example:"{{UserUUID}}" json:"user_uuid" validate:"required"`
	Role     string `example:"member" json:"role" validate:"required"`
}

func (r *registryProvider) GetRestyClient(org string) (*rest.RestClient, error) {
	/* Add user to member db of the org */
	url, err := ic.GetHostUrl(ic.CreateHostString(org, SystemName), r.icHost, &org, r.debug)
	if err != nil {
		log.Errorf("Failed to resolve registry address to update user as member: %v", err)
		return nil, fmt.Errorf("failed to resolve org registry address. Error: %v", err)
	}

	rc := rest.NewRestyClient(url, r.debug)

	return rc, nil
}

func NewRegistryProvider(icHost string, t int, debug bool) *registryProvider {

	r := &registryProvider{
		debug:   debug,
		icHost:  icHost,
		timeout: time.Duration(t) * time.Second,
	}

	return r
}

func (r *registryProvider) GetAllNodes(orgName string) (*nodepb.GetNodesResponse, error) {

	var err error

	nodeResp := &nodepb.GetNodesResponse{}
	/* Get Provider */
	r.R, err = r.GetRestyClient(orgName)
	if err != nil {
		return nil, err
	}

	errStatus := &rest.ErrorMessage{}

	resp, err := r.R.C.R().
		SetError(errStatus).
		Get(r.R.URL.String() + RegistryVersion + "nodes")

	if err != nil {
		log.Errorf("Failed to send api request to registry at %s . Error %s", r.R.URL.String(), err.Error())
		return nil, fmt.Errorf("api request to registry at %s failure: %v", r.R.URL.String(), err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to get nodes from registry at %s. HTTP resp code %d and Error message is %s",
			r.R.URL.String(), resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("failed to get noddes from registry at %s. Error %s", r.R.URL.String(), errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), nodeResp)
	if err != nil {
		log.Errorf("failed to decode registry response  for node list: %s", err.Error())
		return nil, fmt.Errorf("failed to decode registry response for node list")
	}

	return nodeResp, nil
}
