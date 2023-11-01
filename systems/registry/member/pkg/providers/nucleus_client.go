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
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const orgEndpoint = "/v1/orgs"
const userEndpoint = "/v1/users"

type NucleusClientProvider interface {
	GetOrgByName(name string) (*OrgInfo, error)
	GetUserById(userId string) (*UserInfo, error)
	UpdateOrgToUser(orgId string, userId string) error
	RemoveOrgFromUser(orgId string, userId string) error
}

type nucleusInfoClient struct {
	R *rest.RestClient
}

type OrgInfo struct {
	Org Org `json:"org,omitempty"`
}

type Org struct {
	Id            string    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Owner         string    `json:"owner,omitempty"`
	Certificate   string    `json:"certificate,omitempty"`
	IsDeactivated bool      `json:"is_deactivated,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

type UserInfo struct {
	User User `json:"user,omitempty"`
}

type User struct {
	Id              string    `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Email           string    `json:"email,omitempty"`
	Phone           string    `json:"phone,omitempty"`
	IsDeactivated   bool      `json:"is_deactivated,omitempty"`
	AuthId          string    `json:"auth_id,omitempty"`
	RegisteredSince time.Time `json:"registered_since,omitempty"`
}

type UserOrgRequest struct {
	UserId string
	OrgId  string
}

func NewNucleusClientProvider(url string, debug bool) NucleusClientProvider {
	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		log.Fatalf("Can't connect to %s url. Error %s", url, err.Error())
	}

	n := &nucleusInfoClient{
		R: f,
	}

	return n
}

func (p *nucleusInfoClient) GetOrgByName(name string) (*OrgInfo, error) {
	errStatus := &rest.ErrorMessage{}

	pkg := OrgInfo{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + orgEndpoint + "/" + name)

	if err != nil {
		log.Errorf("Failed to send api request to nucleus/org. Error %s", err.Error())

		return nil, fmt.Errorf("api request to org system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch org info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("org info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &pkg)
	if err != nil {
		log.Tracef("Failed to deserialize org info. Error message is %s", err.Error())

		return nil, fmt.Errorf("org info deserailization failure: %w", err)
	}

	log.Infof("Org Info: %+v", pkg)

	return &pkg, nil
}

func (p *nucleusInfoClient) GetUserById(userId string) (*UserInfo, error) {
	errStatus := &rest.ErrorMessage{}

	pkg := UserInfo{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Get(p.R.URL.String() + userEndpoint + "/" + userId)

	if err != nil {
		log.Errorf("Failed to send api request to nucleus/user. Error %s", err.Error())

		return nil, fmt.Errorf("api request to user system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch org info. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("User Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &pkg)
	if err != nil {
		log.Tracef("Failed to deserialize user info. Error message is %s", err.Error())

		return nil, fmt.Errorf("user info deserailization failure: %w", err)
	}

	log.Infof("User Info: %+v", pkg)

	return &pkg, nil
}

func (p *nucleusInfoClient) UpdateOrgToUser(orgId string, userId string) error {
	errStatus := &rest.ErrorMessage{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Put(p.R.URL.String() + orgEndpoint + "/" + orgId + "/users/" + userId)

	if err != nil {
		log.Errorf("Failed to send api request to nucleus/org. Error %s", err.Error())

		return fmt.Errorf("api request to nucleus/org system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to updated org to user. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return fmt.Errorf("failed to update org to user. error: %s", errStatus.Message)
	}

	return nil
}

func (p *nucleusInfoClient) RemoveOrgFromUser(orgId string, userId string) error {
	errStatus := &rest.ErrorMessage{}

	resp, err := p.R.C.R().
		SetError(errStatus).
		Delete(p.R.URL.String() + orgEndpoint + "/" + orgId + "/users" + userId)

	if err != nil {
		log.Errorf("Failed to send api request to nucleus/org. Error %s", err.Error())

		return fmt.Errorf("api request to nucleus/org system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to delete org from user. HTTP resp code %d and Error message is %s", resp.StatusCode(), errStatus.Message)

		return fmt.Errorf("failed to delete org from user. error: %s", errStatus.Message)
	}

	return nil
}
