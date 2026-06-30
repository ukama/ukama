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

const (
	InitApiEndpoint = "/v1/orgs"
)

var (
	HTTP_PROTOCOL = "http"
	HTTPS_PROTOCOL = "https"
	NODE_GW_INTERFACE = "node-gw"
	API_GW_INTERFACE = "api-gw"
)

type SystemIPInfo struct {
	SystemId    string `json:"systemId"`
	SystemName  string `json:"systemName"`
	OrgName     string `json:"orgName"`
	Certificate string `json:"certificate"`
	ApiGwIp     string `json:"apiGwIp"`
	ApiGwPort   uint   `json:"apiGwPort"`
	ApiGwHealth int    `json:"apiGwHealth"`
	ApiGwUrl    string `json:"apiGwUrl"`
	NodeGwIp    string `json:"nodeGwIp"`
	NodeGwPort   uint  `json:"nodeGwPort"`
	NodeGwHealth int   `json:"nodeGwHealth"`
}

type SystemLookupReq struct {
	System string
	Org    string
}

type InitClient interface {
	GetSystem(org, system string) (*SystemIPInfo, error)
	GetSystemFromHost(host string, org *string) (*SystemIPInfo, error)
}

type initClient struct {
	u *url.URL
	R *client.Resty
}

func NewInitClient(host string, options ...client.Option) *initClient {
	u, err := url.Parse(host)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", host, err)
	}

	ic := &initClient{
		u: u,
		R: client.NewResty(options...),
	}

	log.Tracef("Client created %+v for %s ", ic, ic.u.String())

	return ic
}

func (i *initClient) GetSystem(org, system string) (*SystemIPInfo, error) {
	return i.getSystem(org, system)
}

func (i *initClient) GetSystemFromHost(host string, org *string) (*SystemIPInfo, error) {
	s, err := ParseHostString(host, org)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host string: %w", err)
	}

	return i.getSystem(s.Org, s.System)
}

func (i *initClient) getSystem(org, system string) (*SystemIPInfo, error) {
	log.Debugf("Getting sysem %q from org %q", system, org)

	sysIpInfo := SystemIPInfo{}

	resp, err := i.R.Get(i.u.String() + InitApiEndpoint + "/" + org + "/systems/" + system)
	if err != nil {
		log.Errorf("Get system failure. error: %v", err)

		return nil, fmt.Errorf("get system failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &sysIpInfo)
	if err != nil {
		log.Tracef("Failed to deserialize system IP info. Error message is %v", err)

		return nil, fmt.Errorf("system IP info deserialization failure: %w", err)
	}

	log.Infof("System IP Info: %+v", sysIpInfo)

	return &sysIpInfo, nil
}

func GetHostUrl(ic InitClient, host string, org *string) (*url.URL, error) {
	log.Infof("Getting url from initclient matching host %s", host)

	sysIpInfo, err := ic.GetSystemFromHost(host, org)
	if err != nil {
		log.Errorf("Initclient GetSystem failure. error: %s", err)

		return nil, fmt.Errorf("initclient GetSystem failure: %w", err)
	}

	// Prefer the registered api-gw URL when available. A URL is typically a
	// stable hostname (DNS / service name) that the network layer re-resolves
	// to the current IP, so callers that cache it do not break when a system's
	// IP changes. Fall back to the raw ApiGwIp:ApiGwPort for backward
	// compatibility when no URL was registered.
	if strings.TrimSpace(sysIpInfo.ApiGwUrl) != "" {
		return createURLFromString(sysIpInfo.SystemName, sysIpInfo.OrgName, sysIpInfo.ApiGwUrl)
	}

	return CreateHTTPURL(*sysIpInfo, API_GW_INTERFACE)
}

func GetNodeGwHostURL(ic InitClient, host string, org *string) (*url.URL, error) {
	log.Infof("Getting url from initclient matching host %s", host)

	sysIpInfo, err := ic.GetSystemFromHost(host, org)
	if err != nil {
		log.Errorf("Initclient GetSystem failure. error: %s", err)

		return nil, fmt.Errorf("initclient GetSystem failure: %w", err)
	}

	if sysIpInfo.NodeGwIp == "" {
		return nil, fmt.Errorf("node gw ip is not found for host %s", host)
	}

	return CreateHTTPURL(*sysIpInfo, NODE_GW_INTERFACE)
}

func CreateHTTPURL(s SystemIPInfo, interfaceType string) (*url.URL, error) {
	if interfaceType == "node-gw" {
		return createURL(s.SystemName, s.OrgName, s.NodeGwIp, HTTP_PROTOCOL, s.NodeGwPort)
	}
	return createURL(s.SystemName, s.OrgName, s.ApiGwIp, HTTP_PROTOCOL, s.ApiGwPort)
}

// createURLFromString builds a URL from a pre-registered api-gw URL string.
// If the URL has no scheme, http is assumed.
func createURLFromString(name, org, rawURL string) (*url.URL, error) {
	log.Infof("Creating url from registered api-gw url %q for system %s and org %s",
		rawURL, name, org)

	rawURL = strings.TrimSpace(rawURL)
	if !strings.Contains(rawURL, "://") {
		rawURL = HTTP_PROTOCOL + "://" + rawURL
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("error while parsing api-gw url %q: %w", rawURL, err)
	}

	if u.Host == "" {
		return nil, fmt.Errorf("error while parsing api-gw url %q: missing host", rawURL)
	}

	return u, nil
}

func createURL(name, org, ip, protocol string, port uint) (*url.URL, error) {
	log.Infof("Creating %s url for system %s and org %s",
		protocol, name, org)

	//we can add more protocol validation later
	if protocol == "" {
		return nil, fmt.Errorf("error while creating url: protocol %q is not valid",
			protocol)
	}

	return url.Parse(fmt.Sprintf("%s://%s:%d", protocol, ip, port))
}

/* Host is expected to be orgname.systemname */
func ParseHostString(host string, org *string) (*SystemLookupReq, error) {
	tok := strings.Split(host, ".")
	s := &SystemLookupReq{}

	if len(tok) == 1 {
		/* If it only has system name */
		s.System = tok[0]
		if org == nil {
			return nil, fmt.Errorf("missing organization string for resolving host")
		}
		s.Org = *org
	} else if len(tok) == 2 {
		s.System = tok[1]
		s.Org = tok[0]
		if org != nil && *org != s.Org {
			return nil, fmt.Errorf("organization string for resolving host does not match: (%s != %s)",
				*org, s.Org)
		}
	} else {
		return nil, fmt.Errorf("wrong hostname %s. Expected format orgname.systemname", host)
	}

	return s, nil
}

func CreateHostString(org string, system string) string {
	return fmt.Sprintf("%s.%s", org, system)
}
