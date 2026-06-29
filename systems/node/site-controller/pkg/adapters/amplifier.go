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
	"github.com/ukama/ukama/systems/node/site-controller/providers"
)

type AmplifierAdapter struct {
	cmd providers.ControllerClientProvider
}

func NewAmplifierAdapter(cmd providers.ControllerClientProvider) *AmplifierAdapter {
	return &AmplifierAdapter{cmd: cmd}
}
func (a *AmplifierAdapter) SetRadio(ctx context.Context, nodeID, state string) error {
	b, _ := json.Marshal(map[string]string{"state": state})
	client, err := a.cmd.GetClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}
	_, err = client.SendNodeCommand(ctx, &crpc.SendNodeCommandRequest{NodeId: nodeID, Method: "POST", Path: "/device/v1/radio", Body: b})
	if err != nil {
		return fmt.Errorf("failed to send node command: %w", err)
	}
	return nil
}
