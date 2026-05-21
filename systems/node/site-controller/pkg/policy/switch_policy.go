/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package policy

import (
	"encoding/json"
	"time"

	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
)

const SourceSiteController = "site-controller"

type SwitchPolicy struct {
	SiteID    string       `json:"site_id"`
	Source    string       `json:"source"`
	UpdatedAt string       `json:"updated_at"`
	Ports     []PolicyPort `json:"ports"`
}
type PolicyPort struct {
	Port   int    `json:"port"`
	Role   string `json:"role"`
	NodeID string `json:"node_id,omitempty"`
	Class  string `json:"class"`
	Policy string `json:"policy"`
}

func BuildSwitchPolicy(siteID string, ports []db.SitePortMap) (*SwitchPolicy, error) {
	if err := ValidatePortMap(ports); err != nil {
		return nil, err
	}
	p := &SwitchPolicy{SiteID: siteID, Source: SourceSiteController, UpdatedAt: time.Now().UTC().Format(time.RFC3339), Ports: make([]PolicyPort, 0, len(ports))}
	for _, port := range ports {
		p.Ports = append(p.Ports, PolicyPort{Port: port.Port, Role: port.Role, NodeID: port.NodeID, Class: port.Class, Policy: port.Policy})
	}
	return p, nil
}
func Marshal(policy *SwitchPolicy) ([]byte, error) { return json.Marshal(policy) }
