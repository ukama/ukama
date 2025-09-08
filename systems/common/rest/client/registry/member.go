/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package registry

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ukama/ukama/systems/common/pb/gen/ukama"
	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const MemberEndpoint = "/v1/members"

type MemberInfo struct {
	MemberId      string    `json:"member_id,omitempty"`
	UserId        string    `json:"user_id,omitempty"`
	Role          string    `json:"role"`
	IsDeactivated bool      `json:"is_deactivated"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

type OrgMember struct {
	UserUuid string `example:"{{UserUUID}}" json:"user_uuid" validate:"required"`
	Role     string `example:"member" json:"role" validate:"required"`
}

type MemberInfoResponse struct {
	Member MemberInfo `json:"member"`
}

type MemberClient interface {
	GetByUserId(Id string) (*MemberInfoResponse, error)
	AddMember(uuid string) (*MemberInfoResponse, error)
}

type memberClient struct {
	u *url.URL
	R *client.Resty
}

func NewMemberClient(h string, options ...client.Option) *memberClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &memberClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (m *memberClient) GetByUserId(id string) (*MemberInfoResponse, error) {
	log.Debugf("Getting member: %v", id)

	mem := MemberInfoResponse{}

	resp, err := m.R.Get(m.u.String() + MemberEndpoint + "/user/" + id)
	if err != nil {
		log.Errorf("GetMember failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetMember failure: %w", err)
	}
	err = json.Unmarshal(resp.Body(), &mem)
	if err != nil {
		log.Tracef("Failed to deserialize member info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("member info deserialization failure: %w", err)
	}

	log.Infof("Member Info: %+v", mem)

	return &mem, nil
}

func (m *memberClient) AddMember(uuid string) (*MemberInfoResponse, error) {

	log.Debugf("Adding member: %v", uuid)

	memberRes := MemberInfoResponse{}
	req := OrgMember{
		UserUuid: uuid,
		Role:     ukama.RoleType_ROLE_USER.String(),
	}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("request marshal error: %w", err)
	}

	resp, err := m.R.Post(m.u.String()+MemberEndpoint, b)

	if err != nil {
		log.Errorf("Failed to send api request to registry at %s . Error %s", m.u.String(), err.Error())
		return nil, fmt.Errorf("api request to registry at %s failure: %v", m.u.String(), err)
	}

	err = json.Unmarshal(resp.Body(), &memberRes)
	if err != nil {
		log.Errorf("Failed to deserialize member info. Error message is: %s", err.Error())
		return nil, fmt.Errorf("member info deserialization failure: %w", err)
	}

	return &memberRes, nil
}
