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
	"time"

	log "github.com/sirupsen/logrus"
)

const OrgEndpoint = "/v1/orgs"

type OrgInfo struct {
	Id            string    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Owner         string    `json:"owner,omitempty"`
	Certificate   string    `json:"certificate,omitempty"`
	IsDeactivated bool      `json:"isDeactivated,omitempty"`
	CreatedAt     time.Time `json:"created_AT,omitempty"`
}

type Org struct {
	OrgInfo *OrgInfo `json:"org"`
}

type OrgClient interface {
	Get(name string) (*OrgInfo, error)
}

type orgClient struct {
	u *url.URL
	R *Resty
}

func NewOrgClient(h string) *orgClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &orgClient{
		u: u,
		R: NewResty(),
	}
}

func (n *orgClient) Get(name string) (*OrgInfo, error) {
	log.Debugf("Getting org: %v", name)

	org := Org{}

	resp, err := n.R.Get(n.u.String() + OrgEndpoint + "/" + name)
	if err != nil {
		log.Errorf("GetOrg failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetOrg failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &org)
	if err != nil {
		log.Tracef("Failed to deserialize org info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("org info deserailization failure: %w", err)
	}

	log.Infof("Org Info: %+v", org.OrgInfo)

	return org.OrgInfo, nil
}
