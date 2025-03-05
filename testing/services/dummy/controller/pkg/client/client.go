package client

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	cenums "github.com/ukama/ukama/testing/common/enums"
)

 type DNodeClient struct {
	 baseURL    string
	 httpClient *http.Client
 }
 
 func NewDNodeClient(baseURL string, timeout time.Duration) *DNodeClient {
	 return &DNodeClient{
		 baseURL: baseURL,
		 httpClient: &http.Client{
			 Timeout: timeout,
		 },
	 }
 }
 
 func (c *DNodeClient) UpdateNodeScenario(nodeID string, scenario cenums.SCENARIOS, profile cenums.Profile) error {
	scenarioStr := string(scenario)
	var profileStr string
	switch profile {
	case cenums.PROFILE_NORMAL:
		profileStr = "normal"
	case cenums.PROFILE_MIN:
		profileStr = "min"
	case cenums.PROFILE_MAX:
		profileStr = "max"
	default:
		profileStr = "normal"
	}
	
	data := url.Values{}
	data.Add("nodeid", nodeID)
	data.Add("profile", profileStr)
	data.Add("scenario", scenarioStr)
	
	fullURL := fmt.Sprintf("%s/update", c.baseURL)
	
	log.Infof("Sending update request to dnode for node %s: scenario=%s, profile=%s", nodeID, scenarioStr, profileStr)
	
	req, err := http.NewRequest(http.MethodPost, fullURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to dnode: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("dnode server returned non-OK status: %s", resp.Status)
	}
	
	log.Infof("Successfully updated node %s with scenario %s and profile %s", nodeID, scenarioStr, profileStr)
	return nil
}
 
 func (c *DNodeClient) NotifyNodePowerDown(nodeID string) error {
	 return c.UpdateNodeScenario(nodeID, cenums.SCENARIO_POWER_DOWN, cenums.PROFILE_NORMAL)
 }
 
 func (c *DNodeClient) NotifyNodeBackhaulDown(nodeID string) error {
	 return c.UpdateNodeScenario(nodeID, cenums.SCENARIO_BACKHAUL_DOWN, cenums.PROFILE_NORMAL)
 }
 
