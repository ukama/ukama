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

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type ErrorMessage struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type InitClient struct {
	C   *resty.Client
	URL *url.URL
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

func NewInitClient(host string, debug bool) (*InitClient, error) {
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	c := resty.New()
	c.SetDebug(debug)
	rc := &InitClient{
		C:   c,
		URL: url,
	}
	log.Tracef("Client created %+v for %s ", rc, rc.URL.String())
	return rc, nil
}

func GetHostUrl(host string, icHost string, org *string, debug bool) (*url.URL, error) {
	i, err := NewInitClient(icHost, debug)
	if err != nil {
		log.Errorf("Failed to get rest client.Error %+v", err)
		return nil, err
	}

	errStatus := &ErrorMessage{}

	pkg := SystemIPInfo{}
	s, err := ParseHostString(host, org)
	if err != nil {
		return nil, err
	}

	resp, err := i.C.R().
		SetError(errStatus).
		Get(i.URL.String() + SYSTEM_API_VERSION + "/orgs/" + *org + "/systems/" + s.System)

	if err != nil {
		log.Errorf("Failed to send api request to InitClient. Error %s", err.Error())

		return nil, fmt.Errorf("api request to initclient system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Tracef("Failed to fetch %s host info. HTTP resp code %d and Error message is %s", icHost, resp.StatusCode(), errStatus.Message)

		return nil, fmt.Errorf("host Info failure %s", errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), &pkg)
	if err != nil {
		log.Tracef("Failed to deserialize system IP info. Error message is %s", err.Error())

		return nil, fmt.Errorf("system IP info deserailization failure: %w", err)
	}

	log.Infof("System IP Info: %+v", pkg)

	return CreateHTTPURL(pkg)

}

func CreateHTTPURL(s SystemIPInfo) (*url.URL, error) {
	host := fmt.Sprintf("http://%s:%d", s.Ip, s.Port)
	url, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func CreateHTTPSURL(s SystemIPInfo) (*url.URL, error) {
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
