/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const InitNetworkStatusEndpoint = "/v1/status"

type InitNetworkStatus struct {
	Ready  bool   `json:"ready"`
	State  string `json:"state"`
	Reason string `json:"reason"`

	Bridge struct {
		Name             string `json:"name"`
		Address          string `json:"address"`
		Cidr             string `json:"cidr"`
		ManagementSocket string `json:"managementSocket"`
		Openflow         string `json:"openflow"`
	} `json:"bridge"`

	UE struct {
		Cidr        string `json:"cidr"`
		DefaultDrop bool   `json:"defaultDrop"`
	} `json:"ue"`
}

type InitNetworkClient struct {
	baseURL string
	client  *http.Client
}

func NewInitNetworkClient(baseURL string) *InitNetworkClient {
	return &InitNetworkClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (c *InitNetworkClient) GetStatus() (*InitNetworkStatus, error) {
	var status InitNetworkStatus

	if c == nil || c.baseURL == "" {
		return nil, fmt.Errorf("init-network client not configured")
	}

	resp, err := c.client.Get(c.baseURL + InitNetworkStatusEndpoint)
	if err != nil {
		return nil, fmt.Errorf("get init-network status: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Warnf("Fail to gracefuly close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("init-network status returned http %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("decode init-network status: %w", err)
	}

	if !status.Ready {
		return nil, fmt.Errorf("init-network not ready state=%s reason=%s",
			status.State, status.Reason)
	}

	if status.Bridge.Name == "" {
		return nil, fmt.Errorf("init-network status missing bridge.name")
	}

	if status.Bridge.ManagementSocket == "" {
		return nil, fmt.Errorf("init-network status missing bridge.managementSocket")
	}

	if status.UE.Cidr == "" {
		return nil, fmt.Errorf("init-network status missing ue.cidr")
	}

	return &status, nil
}
