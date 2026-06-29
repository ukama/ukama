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

	crpc "github.com/ukama/ukama/systems/node/controller/pb/gen"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/policy"
	"github.com/ukama/ukama/systems/node/site-controller/providers"
)

type CNodeAdapter struct {
	cmd providers.ControllerClientProvider
}

func NewCNodeAdapter(cmd providers.ControllerClientProvider) *CNodeAdapter {
	return &CNodeAdapter{cmd: cmd}
}
func (a *CNodeAdapter) ApplySwitchPolicy(ctx context.Context, nodeID string, p *policy.SwitchPolicy) error {
	b, err := policy.Marshal(p)
	if err != nil {
		return err
	}

	client, err := a.cmd.GetClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}
	_, err = client.SendNodeCommand(ctx, &crpc.SendNodeCommandRequest{NodeId: nodeID, Method: "PUT", Path: "/v1/ports/policy", Body: b})
	if err != nil {
		return fmt.Errorf("failed to send node command: %w", err)
	}
	return nil
}

func (a *CNodeAdapter) SetPortPoe(ctx context.Context, nodeID string, port int, on bool, reason string) error {
	b, _ := json.Marshal(map[string]interface{}{"on": on, "source": "site-controller", "reason": reason})
	client, err := a.cmd.GetClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}
	_, err = client.SendNodeCommand(ctx, &crpc.SendNodeCommandRequest{NodeId: nodeID, Method: "POST", Path: fmt.Sprintf("/v1/ports/%d/poe", port), Body: b})
	if err != nil {
		return fmt.Errorf("failed to send node command: %w", err)
	}
	return nil
}

func (a *CNodeAdapter) PowerCyclePort(ctx context.Context, nodeID string, port int, reason string) error {
	b, _ := json.Marshal(map[string]string{"source": "site-controller", "reason": reason})
	client, err := a.cmd.GetClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}
	_, err = client.SendNodeCommand(ctx, &crpc.SendNodeCommandRequest{NodeId: nodeID, Method: "POST", Path: fmt.Sprintf("/v1/ports/%d/poe/cycle", port), Body: b})
	if err != nil {
		return fmt.Errorf("failed to send node command: %w", err)
	}
	return nil
}
