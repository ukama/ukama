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
	"time"

	log "github.com/sirupsen/logrus"
	cenums "github.com/ukama/ukama/testing/common/enums"
)

// DNodeClient represents a client for interacting with the dnode service
 type DNodeClient struct {
	 baseURL    string
	 httpClient *http.Client
 }
 
 // NewDNodeClient creates a new dnode client
 func NewDNodeClient(baseURL string, timeout time.Duration) *DNodeClient {
	 return &DNodeClient{
		 baseURL: baseURL,
		 httpClient: &http.Client{
			 Timeout: timeout,
		 },
	 }
 }
 
 // UpdateNodeScenario sends a scenario update request to the dnode service
 func (c *DNodeClient) UpdateNodeScenario(nodeID string, scenario cenums.SCENARIOS, profile cenums.Profile) error {
	 // Convert scenario and profile to string
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
	 
	 // Construct the query parameters
	 params := url.Values{}
	 params.Add("nodeid", nodeID)
	 params.Add("profile", profileStr)
	 params.Add("scenario", scenarioStr)
 
	 // Construct the full URL with query parameters
	 fullURL := fmt.Sprintf("%s/update?%s", c.baseURL, params.Encode())
	 
	 log.Infof("Sending update request to dnode for node %s: scenario=%s, profile=%s", nodeID, scenarioStr, profileStr)
	 
	 // Create a new request
	 req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	 if err != nil {
		 return fmt.Errorf("failed to create request: %w", err)
	 }
 
	 // Send the request
	 resp, err := c.httpClient.Do(req)
	 if err != nil {
		 return fmt.Errorf("failed to send request to dnode: %w", err)
	 }
	 defer resp.Body.Close()
 
	 // Check the response status
	 if resp.StatusCode != http.StatusOK {
		 return fmt.Errorf("dnode server returned non-OK status: %s", resp.Status)
	 }
 
	 log.Infof("Successfully updated node %s with scenario %s and profile %s", nodeID, scenarioStr, profileStr)
	 return nil
 }
 
 // NotifyNodePowerDown notifies the dnode that a node should be in power down scenario
 func (c *DNodeClient) NotifyNodePowerDown(nodeID string) error {
	 return c.UpdateNodeScenario(nodeID, cenums.SCENARIO_POWER_DOWN, cenums.PROFILE_NORMAL)
 }
 
 // NotifyNodeBackhaulDown notifies the dnode that a node should be in backhaul down scenario
 func (c *DNodeClient) NotifyNodeBackhaulDown(nodeID string) error {
	 return c.UpdateNodeScenario(nodeID, cenums.SCENARIO_BACKHAUL_DOWN, cenums.PROFILE_NORMAL)
 }
 
 // ResetNodeScenario resets a node to the default scenario
 func (c *DNodeClient) ResetNodeScenario(nodeID string) error {
	 return c.UpdateNodeScenario(nodeID, cenums.SCENARIO_DEFAULT, cenums.PROFILE_NORMAL)
 }