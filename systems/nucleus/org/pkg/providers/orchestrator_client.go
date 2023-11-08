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
	"github.com/ukama/ukama/systems/common/rest"
)

// OrchestratorClientProvider creates a local client to interact with
// a remote instance of  Users service.
type OrchestratorProvider interface {
	DeployOrg(req DeployOrgRequest) (*DeployOrgResponse, error)
	DestroyOrg(req DestroyOrgRequest) (*DestroyOrgResponse, error)
}

type orchestratorProvider struct {
	R *rest.RestClient
}

type System struct {
	Name     string   `json:"system" validate:"required"`
	KeyValue []string `json:"keyvalue" validate:"required"`
}

type DeployOrgRequest struct {
	OrgId   string   `path:"org_id" validate:"required"`
	OrgName string   `json:"org_name" validate:"required"`
	OwnerId string   `json:"owner_id" validate:"required"`
	Systems []System `json:"systems"`
}

type DeployOrgResponse struct {
}

type DestroyOrgRequest struct {
	OrgId   string `path:"org_id" validate:"required"`
	OwnerId string `json:"owner_id" validate:"required"`
}

type DestroyOrgResponse struct {
}

const ORCH_PATH = "/v1/orchestrator"

func NewOrchestratorProvider(orchestratorHost string, debug bool) OrchestratorProvider {

	f, err := rest.NewRestClient(orchestratorHost, debug)
	if err != nil {
		log.Fatalf("Can't connect to %s url. Error %s", orchestratorHost, err.Error())
	}

	n := &orchestratorProvider{
		R: f,
	}

	return n

}

func (p *orchestratorProvider) DeployOrg(req DeployOrgRequest) (*DeployOrgResponse, error) {
	errStatus := &rest.ErrorMessage{}

	dResp := &DeployOrgResponse{}

	url := fmt.Sprintf("%s%s%s%s", p.R.URL.String(), ORCH_PATH, "/deploy/orgs/", req.OrgId)
	resp, err := p.R.C.R().
		SetError(errStatus).
		SetBody(req).
		Put(url)

	if err != nil {
		log.Errorf("Failed to send api request to orchestrator. Error %s", err.Error())

		return nil, fmt.Errorf("api request to orchestrator system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to deploy org. URL %s HTTP resp code %d and Error message is %s", url, resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("orchestrator deploy org request failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), dResp)
	if err != nil {
		log.Errorf("Failed to deserialize orchestartor response. Error message is %s", err.Error())

		return nil, fmt.Errorf("orchestartor response deserialization failure: %w", err)
	}

	log.Infof("Deploy Org Response: %+v", dResp)

	return dResp, nil

}

func (p *orchestratorProvider) DestroyOrg(req DestroyOrgRequest) (*DestroyOrgResponse, error) {

	errStatus := &rest.ErrorMessage{}

	dResp := &DestroyOrgResponse{}
	resp, err := p.R.C.R().
		SetError(errStatus).
		SetBody(req).
		Delete(p.R.URL.String() + ORCH_PATH + "/orgs/" + req.OrgId)

	if err != nil {
		log.Errorf("Failed to send api request to orchestrator. Error %s", err.Error())

		return nil, fmt.Errorf("api request to orchestrator system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to destroy org. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("orchestrator destroy org request failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), dResp)
	if err != nil {
		log.Errorf("Failed to deserialize orchestartor response. Error message is %s", err.Error())

		return nil, fmt.Errorf("orchestartor response deserialization failure: %w", err)
	}

	log.Infof("Destroy Org Response: %+v", dResp)

	return dResp, nil

}
