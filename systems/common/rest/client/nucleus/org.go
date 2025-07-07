/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package nucleus

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ukama/ukama/systems/common/rest/client"

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
	AddUser(orgId string, userId string) error
	RemoveUser(orgId string, userId string) error
}

type orgClient struct {
	u *url.URL
	R *client.Resty
}

func NewOrgClient(h string, options ...client.Option) *orgClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &orgClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (o *orgClient) Get(name string) (*OrgInfo, error) {
	log.Debugf("Getting org: %v", name)

	org := Org{}

	resp, err := o.R.Get(o.u.String() + OrgEndpoint + "/" + name)
	if err != nil {
		log.Errorf("GetOrg failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetOrg failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &org)
	if err != nil {
		log.Tracef("Failed to deserialize org info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("org info deserialization failure: %w", err)
	}

	log.Infof("Org Info: %+v", org.OrgInfo)

	return org.OrgInfo, nil
}

func (o *orgClient) AddUser(orgId string, userId string) error {
	log.Debugf("Adding user %q to org %q", userId, orgId)

	_, err := o.R.Put(o.u.String()+OrgEndpoint+"/"+orgId+"/users/"+userId, nil)
	if err != nil {
		log.Errorf("AddUser failure. error: %s", err.Error())

		return fmt.Errorf("AddUser failure: %w", err)
	}

	return nil
}

func (o *orgClient) RemoveUser(orgId string, userId string) error {
	log.Debugf("Removing user %q from org %q", userId, orgId)

	_, err := o.R.Delete(o.u.String() + OrgEndpoint + "/" + orgId + "/users/" + userId)
	if err != nil {
		log.Errorf("RemoveUser failure. error: %s", err.Error())

		return fmt.Errorf("RemoveUser failure: %w", err)
	}

	return nil
}
