/*
* This Source Code Form is subject to the terms of the Mozilla Public
* License, v. 2.0. If a copy of the MPL was not distributed with this
* file, You can obtain one at https://mozilla.org/MPL/2.0/.
*
* Copyright (c) 2026-present, Ukama Inc.
 */

package adapters

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ukama/ukama/systems/node/site-controller/pkg/policy"
)

type CNodeAdapter struct{ cmd NodeCommandAdapter }

func NewCNodeAdapter(cmd NodeCommandAdapter) *CNodeAdapter { return &CNodeAdapter{cmd: cmd} }
func (a *CNodeAdapter) ApplySwitchPolicy(ctx context.Context, nodeID string, p *policy.SwitchPolicy) error {
	b, err := policy.Marshal(p)
	if err != nil {
		return err
	}
	return a.cmd.Send(ctx, nodeID, "PUT", "/switch/v1/ports/policy", b)
}
func (a *CNodeAdapter) SetPortPoe(ctx context.Context, nodeID string, port int, on bool, reason string) error {
	b, _ := json.Marshal(map[string]interface{}{"on": on, "source": "site-controller", "reason": reason})
	return a.cmd.Send(ctx, nodeID, "POST", fmt.Sprintf("/switch/v1/ports/%d/poe", port), b)
}
func (a *CNodeAdapter) PowerCyclePort(ctx context.Context, nodeID string, port int, reason string) error {
	b, _ := json.Marshal(map[string]string{"source": "site-controller", "reason": reason})
	return a.cmd.Send(ctx, nodeID, "POST", fmt.Sprintf("/switch/v1/ports/%d/poe/cycle", port), b)
}
