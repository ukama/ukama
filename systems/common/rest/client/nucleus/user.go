/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.user/MPL/2.0/.
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

const UserEndpoint = "/v1/users"

type UserInfo struct {
	Id              string    `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Email           string    `json:"email,omitempty"`
	Phone           string    `json:"phone,omitempty"`
	IsDeactivated   bool      `json:"is_deactivated,omitempty"`
	AuthId          string    `json:"auth_id,omitempty"`
	RegisteredSince time.Time `json:"registered_since,omitempty"`
}

type User struct {
	UserInfo *UserInfo `json:"user"`
}

type UserClient interface {
	GetById(userId string) (*UserInfo, error)
	GetByEmail(email string) (*UserInfo, error)
}

type userClient struct {
	u *url.URL
	R *client.Resty
}

func NewUserClient(h string, options ...client.Option) *userClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &userClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (u *userClient) GetById(id string) (*UserInfo, error) {
	log.Debugf("Getting user: %v", id)

	user := User{}

	resp, err := u.R.Get(u.u.String() + UserEndpoint + "/" + id)
	if err != nil {
		log.Errorf("GetUser failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetUser failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &user)
	if err != nil {
		log.Tracef("Failed to deserialize user info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("user info deserialization failure: %w", err)
	}

	log.Infof("User Info: %+v", user.UserInfo)

	return user.UserInfo, nil
}

func (u *userClient) GetByEmail(email string) (*UserInfo, error) {
	log.Debugf("Getting user: %v", email)

	user := User{}

	resp, err := u.R.Get(u.u.String() + UserEndpoint + "/email/" + email)
	if err != nil {
		log.Errorf("GetUser failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetUser failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &user)
	if err != nil {
		log.Tracef("Failed to deserialize user info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("user info deserialization failure: %w", err)
	}

	log.Infof("User Info: %+v", user.UserInfo)

	return user.UserInfo, nil
}
