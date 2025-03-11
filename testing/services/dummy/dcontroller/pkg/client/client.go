/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	cenums "github.com/ukama/ukama/testing/common/enums"
)

 type DNodeClient struct {
	 baseURL    string
	 httpClient *http.Client
 }
 
 func NewDNodeClient(baseURL string, timeout time.Duration) *DNodeClient {
	 if timeout == 0 {
		 timeout = 10 * time.Second
	 }
 
	 return &DNodeClient{
		 baseURL: baseURL,
		 httpClient: &http.Client{
			 Timeout: timeout,
		 },
	 }
 }
 
 func (c *DNodeClient) UpdateNodeScenario(nodeId string, scenario cenums.SCENARIOS, profile cenums.Profile) error {
	 updateURL := fmt.Sprintf("%s/update", c.baseURL)
 
	 params := url.Values{}
	 params.Add("nodeid", nodeId)
	 params.Add("scenario", string(scenario))
	 params.Add("profile", string(profile))
 
	 requestURL := fmt.Sprintf("%s?%s", updateURL, params.Encode())
 
	 req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	 if err != nil {
		 return fmt.Errorf("failed to create request: %w", err)
	 }
 
	 log.Infof("Sending request to dnode: %s", requestURL)
 
	 resp, err := c.httpClient.Do(req)
	 if err != nil {
		 return fmt.Errorf("failed to send request to dnode: %w", err)
	 }
	 defer resp.Body.Close()
 
	 if resp.StatusCode != http.StatusOK {
		 return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	 }
 
	 log.Infof("Successfully updated node %s scenario to %s with profile %s", nodeId, scenario, string(profile))
	 return nil
 }
 
 func (c *DNodeClient) SetBackhaulDown(nodeId string, profile cenums.Profile) error {
	 return c.UpdateNodeScenario(nodeId, cenums.SCENARIO_BACKHAUL_DOWN, profile)
 }
 
 // SetNodeOff notifies the DNode that the node is off
 func (c *DNodeClient) SetNodeOff(nodeId string, profile cenums.Profile) error {
	 return c.UpdateNodeScenario(nodeId, cenums.SCENARIO_NODE_OFF, profile)
 }
 
 // SetDefault sets the node back to the default scenario
 func (c *DNodeClient) SetDefault(nodeId string, profile cenums.Profile) error {
	 return c.UpdateNodeScenario(nodeId, cenums.SCENARIO_DEFAULT, profile)
 }