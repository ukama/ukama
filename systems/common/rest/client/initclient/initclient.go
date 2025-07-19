/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package initclient

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

type InitClient struct {
	u *url.URL
	R *client.Resty
}

type SystemIPInfo struct {
	SystemId    string `json:"systemId"`
	SystemName  string `json:"systemName"`
	OrgName     string `json:"orgName"`
	Certificate string `json:"certificate"`
	Ip          string `json:"ip"`
	Port        uint   `json:"port"`
	Health      int    `json:"health"`
}

type SystemLookupReq struct {
	System string
	Org    string
}

const SYSTEM_API_VERSION = "/v1"

func NewInitClient(host string, options ...client.Option) *InitClient {
	u, err := url.Parse(host)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", host, err)
	}

	ic := &InitClient{
		u: u,
		R: client.NewResty(options...),
	}

	log.Tracef("Client created %+v for %s ", ic, ic.u.String())

	return ic
}

func GetHostUrl(host string, icHost string, org *string, debug bool) (*url.URL, error) {
	log.Infof("Getting host url from initclient for host %s and org %s", host, *org)

	// errStatus := &rest.ErrorMessage{}
	ic := NewInitClient(icHost)
	if debug {
		ic = NewInitClient(icHost, client.WithDebug())
	}

	sysIpInfo := SystemIPInfo{}
	s, err := ParseHostString(host, org)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host string: %w", err)
	}

	resp, err := ic.R.Get(ic.u.String() + SYSTEM_API_VERSION + "/orgs/" + *org + "/systems/" + s.System)
	if err != nil {
		log.Errorf("Get initclient system failure. error: %s", err.Error())

		return nil, fmt.Errorf("get initclient system failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &sysIpInfo)
	if err != nil {
		log.Tracef("Failed to deserialize system IP info. Error message is %s", err.Error())

		return nil, fmt.Errorf("system IP info deserailization failure: %w", err)
	}

	log.Infof("System IP Info: %+v", sysIpInfo)

	return CreateHTTPURL(sysIpInfo)
}

func CreateHTTPURL(s SystemIPInfo) (*url.URL, error) {
	log.Infof("Creating HTTP url for system %s and org %s", s.SystemName, s.OrgName)

	host := fmt.Sprintf("http://%s:%d", s.Ip, s.Port)
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func CreateHTTPSURL(s SystemIPInfo) (*url.URL, error) {
	log.Infof("Creating HTTPS url for system %s and org %s", s.SystemName, s.OrgName)

	host := fmt.Sprintf("https://%s:%d", s.Ip, s.Port)
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return url, nil
}

/* Host is expected to be orgname.systemname */
func ParseHostString(host string, org *string) (*SystemLookupReq, error) {
	tok := strings.Split(host, ".")
	s := &SystemLookupReq{}

	if len(tok) == 1 {
		/* If it only has system name */
		s.System = tok[0]
		if org != nil {
			s.Org = *org
		} else {
			return nil, fmt.Errorf("missing organization string for resolving host")
		}
	} else if len(tok) == 2 {
		s.System = tok[1]
		s.Org = tok[0]
	} else {
		return nil, fmt.Errorf("wrong hostname %s. Expected format orgname.systemname", host)
	}

	return s, nil
}

func CreateHostString(org string, system string) string {
	return fmt.Sprintf("%s.%s", org, system)
}
