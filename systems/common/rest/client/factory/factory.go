/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

package factory

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const FactoryEndpoint = "/v1/nodefactory"

type NodeFactoryInfo struct {
	Id            string    `json:"id,omitempty"`
	Type          string    `json:"type,omitempty"`
	OrgName       string    `json:"orgName,omitempty"`
	IsProvisioned bool      `json:"isProvisioned,omitempty"`
	ProvisionedAt time.Time `json:"provisionedAt,omitempty"`
}

type Nodes struct {
	Nodes []*NodeFactoryInfo `json:"nodes"`
}

type NodeFactoryClient interface {
	Get(Id string) (*NodeFactoryInfo, error)
	List(nodeType string, orgName string, isProvisioned bool) (*Nodes, error)
}

type nodeFactoryClient struct {
	u *url.URL
	R *client.Resty
}

func NewNodeFactoryClient(h string, options ...client.Option) *nodeFactoryClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &nodeFactoryClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (s *nodeFactoryClient) Get(id string) (*NodeFactoryInfo, error) {
	log.Debugf("Getting node by id from factory: %v", id)

	nodeFactory := NodeFactoryInfo{}

	resp, err := s.R.Get(s.u.String() + FactoryEndpoint + "/node/" + id)
	if err != nil {
		log.Errorf("Get node failure. error: %s", err.Error())

		return nil, fmt.Errorf("getNodeFactory failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &nodeFactory)
	if err != nil {
		log.Tracef("Failed to deserialize node info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("node info deserialization failure: %w", err)
	}

	log.Infof("Node Factory Info: %+v", nodeFactory)

	return &nodeFactory, nil
}

func (s *nodeFactoryClient) List(nodeType string, orgName string, isProvisioned bool) (*Nodes, error) {
	log.Debugf("Listing nodes from factory. nodeType: %v, orgName: %v, isProvisioned: %v", nodeType, orgName, isProvisioned)

	nodes := Nodes{}

	resp, err := s.R.Get(s.u.String() + FactoryEndpoint + "/nodes?type=" + nodeType + "&orgName=" + orgName + "&isProvisioned=" + strconv.FormatBool(isProvisioned))
	if err != nil {
		log.Errorf("ListNodes failure. error: %s", err.Error())

		return nil, fmt.Errorf("ListNodes failure: %w", err)
	}

	err = json.Unmarshal(resp.Body(), &nodes)
	if err != nil {
		log.Tracef("Failed to deserialize node info. Error message is: %s", err.Error())

		return nil, fmt.Errorf("node info deserialization failure: %w", err)
	}

	return &nodes, nil
}
