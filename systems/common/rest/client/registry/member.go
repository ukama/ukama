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

type MemberInfoResponse struct {
	Member MemberInfo `json:"member"`
}

type MemberClient interface {
	GetByUserId(Id string) (*MemberInfoResponse, error)
}

type memberClient struct {
	u *url.URL
	R *client.Resty
}

func NewMemberClient(h string, options ...client.Option) *memberClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
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

		return nil, fmt.Errorf("member info deserailization failure: %w", err)
	}

	log.Infof("Member Info: %+v", mem)

	return &mem, nil
}
