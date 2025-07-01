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

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const SiteEndpoint = "/v1/sites"

type SiteInfo struct {
	Id          string  `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	NetworkId   string  `json:"network_id,omitempty"`
	BackhaulId  string  `json:"backhaul_id,omitempty"`
	PowerId     string  `json:"power_id,omitempty"`
	AccessId    string  `json:"access_id,omitempty"`
	SwitchId    string  `json:"switch_id,omitempty"`
	SpectrumId  string  `json:"spectrum_id,omitempty"`
	Deactivated bool    `json:"deactivated,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	InstallDate string  `json:"install_date,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	Location    string  `json:"location,omitempty"`
}

type Site struct {
	SiteInfo *SiteInfo `json:"site"`
}

type SiteClient interface {
	Get(id string) (*SiteInfo, error)
}

type siteClient struct {
	u *url.URL
	R *client.Resty
}

func NewSiteClient(h string, options ...client.Option) *siteClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &siteClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (s *siteClient) Get(id string) (*SiteInfo, error) {
	log.Debugf("Getting site: %v", id)

	site := Site{}

	resp, err := s.R.Get(s.u.String() + SiteEndpoint + "/" + id)
	if err != nil {
		log.Errorf("GetSite failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetSite failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &site)
	if err != nil {
		log.Tracef("Failed to deserialize site info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("site info deserialization failure: %w", err)
	}

	log.Infof("Site Info: %+v", site.SiteInfo)

	return site.SiteInfo, nil
}
