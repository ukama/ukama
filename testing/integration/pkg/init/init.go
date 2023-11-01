/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package init

import (
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	api "github.com/ukama/ukama/systems/init/api-gateway/pkg/rest"
	lpb "github.com/ukama/ukama/systems/init/lookup/pb/gen"
	"github.com/ukama/ukama/testing/integration/pkg/utils"
	"k8s.io/apimachinery/pkg/util/json"
)

type InitClient struct {
	u *url.URL
	r utils.Resty
}

func NewInitClient(h string) *InitClient {
	u, _ := url.Parse(h)
	return &InitClient{
		u: u,
		r: *utils.NewResty(),
	}

}
func (s *InitClient) InitAddOrg(req api.AddOrgRequest) (*lpb.AddOrgResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &lpb.AddOrgResponse{}

	resp, err := s.r.Put(s.u.String()+"/v1/orgs/"+req.OrgName, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *InitClient) InitGetOrg(req api.GetOrgRequest) (*lpb.GetOrgResponse, error) {

	rsp := &lpb.GetOrgResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/orgs/" + req.OrgName)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *InitClient) InitAddSystem(req api.AddSystemRequest) (*lpb.AddSystemResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &lpb.AddSystemResponse{}

	resp, err := s.r.Put(s.u.String()+"/v1/orgs/"+req.OrgName+"/systems/"+req.SysName, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *InitClient) InitGetSystem(req api.GetSystemRequest) (*lpb.GetSystemResponse, error) {

	rsp := &lpb.GetSystemResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/orgs/" + req.OrgName + "/systems/" + req.SysName)

	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *InitClient) InitAddNode(req api.AddNodeRequest) (*lpb.AddNodeResponse, error) {

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error. error: %s", err.Error())
	}
	rsp := &lpb.AddNodeResponse{}

	resp, err := s.r.Put(s.u.String()+"/v1/orgs/"+req.OrgName+"/nodes/"+req.NodeId, b)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}

func (s *InitClient) InitGetNode(req api.GetNodeRequest) (*lpb.GetNodeResponse, error) {

	rsp := &lpb.GetNodeResponse{}

	resp, err := s.r.Get(s.u.String() + "/v1/orgs/" + req.OrgName + "/nodes/" + req.NodeId)
	if err != nil {
		log.Errorf("Failed to send api request. error %s", err.Error())
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), rsp)
	if err != nil {
		return nil, fmt.Errorf("response unmarshal error. error: %s", err.Error())
	}

	return rsp, nil
}
