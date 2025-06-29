/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package inventory

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/uuid"

	log "github.com/sirupsen/logrus"
)

const ComponentEndpoint = "/v1/components"

type ComponentInfo struct {
	Id            uuid.UUID `json:"id,omitempty"`
	Inventory     string    `json:"inventory,omitempty"`
	UserId        string    `json:"user_id,omitempty"`
	Category      string    `json:"category"`
	Type          string    `json:"type,omitempty"`
	Description   string    `json:"description,omitempty"`
	DatasheetURL  string    `json:"datasheet_url,omitempty"`
	ImagesURL     string    `json:"images_url,omitempty"`
	PartNumber    string    `json:"part_number,omitempty"`
	Manufacturer  string    `json:"manufacturer,omitempty"`
	Managed       string    `json:"managed,omitempty"`
	Warranty      uint32    `json:"warranty"`
	Specification string    `json:"specification,omitempty"`
}

type Component struct {
	ComponentInfo *ComponentInfo `json:"component"`
}

type ComponentClient interface {
	Get(Id string) (*ComponentInfo, error)
}

type componentClient struct {
	u *url.URL
	R *client.Resty
}

func NewComponentClient(h string, options ...client.Option) *componentClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse  %s url. Error: %s", h, err.Error())
	}

	return &componentClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (s *componentClient) Get(id string) (*ComponentInfo, error) {
	log.Debugf("Getting component: %v", id)

	component := Component{}

	resp, err := s.R.Get(s.u.String() + ComponentEndpoint + "/" + id)
	if err != nil {
		log.Errorf("GetComponent failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetComponent failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &component)
	if err != nil {
		log.Tracef("Failed to deserialize component info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("component info deserialization failure: %w", err)
	}

	log.Infof("Component Info: %+v", component.ComponentInfo)

	return component.ComponentInfo, nil
}
