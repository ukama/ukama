/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package node

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ukama/ukama/systems/common/rest/client"

	log "github.com/sirupsen/logrus"
)

const StateEndpoint = "/v1/state"

type StateInfo struct {
	Id              string    `json:"id,omitempty"`
	NodeId          string    `json:"nodeId,omitempty"`
	PreviousStateId string    `json:"previousStateId,omitempty"`
	CurrentState    string    `json:"currentState,omitempty"`
	SubState        []string  `json:"subState,omitempty"`
	Events          []string  `json:"events,omitempty"`
	NodeType        string    `json:"nodeType,omitempty"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
	UpdatedAt       time.Time `json:"updatedAt,omitempty"`
}

type GetLatestStateResponse struct {
	State *StateInfo `json:"State"`
}

type NodeStateClient interface {
	GetLatestState(nodeId string) (*StateInfo, error)
}

type nodeStateClient struct {
	u *url.URL
	R *client.Resty
}

func NewNodeStateClient(h string, options ...client.Option) *nodeStateClient {
	u, err := url.Parse(h)

	if err != nil {
		log.Fatalf("Can't parse %s url. Error: %v", h, err)
	}

	return &nodeStateClient{
		u: u,
		R: client.NewResty(options...),
	}
}

func (s *nodeStateClient) GetLatestState(nodeId string) (*StateInfo, error) {
	log.Debugf("Getting latest state for node: %v", nodeId)

	resp, err := s.R.Get(s.u.String() + StateEndpoint + "/" + nodeId + "/latest")
	if err != nil {
		log.Errorf("GetLatestState failure. error: %s", err.Error())

		return nil, fmt.Errorf("GetLatestState failure: %w", err)
	}

	out := &GetLatestStateResponse{}

	err = json.Unmarshal(resp.Body(), out)
	if err != nil {
		log.Tracef("Failed to deserialize latest state response. Error message is: %s", err.Error())

		return nil, fmt.Errorf("latest state response deserialization failure: %w", err)
	}

	log.Infof("Latest state: %+v", out.State)

	return out.State, nil
}
