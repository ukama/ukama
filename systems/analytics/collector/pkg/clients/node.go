/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package clients

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const NodesEndpoint = "/v1/nodes"

/* TODO-verify: response shape against registry/node api-gateway. */

type NodeRecord struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	SiteId       string `json:"site_id"`
	NetworkId    string `json:"network_id"`
	State        string `json:"state"`
	Connectivity string `json:"connectivity"`
}

type nodesResponse struct {
	Nodes []NodeRecord `json:"nodes"`
}

type NodeClient interface {
	GetNodes() ([]NodeRecord, error)
}

type nodeClient struct {
	u *url.URL
	R *client.Resty
}

func NewNodeClient(h string, options ...client.Option) NodeClient {
	u, err := url.Parse(h)
	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &nodeClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (c *nodeClient) GetNodes() ([]NodeRecord, error) {
	resp, err := c.R.Get(c.u.String() + NodesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("GetNodes failure: %w", err)
	}

	out := nodesResponse{}

	err = json.Unmarshal(resp.Body(), &out)
	if err != nil {
		return nil, fmt.Errorf("nodes deserialization failure: %w", err)
	}

	return out.Nodes, nil
}
