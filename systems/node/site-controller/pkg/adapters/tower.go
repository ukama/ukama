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
)

type TowerAdapter struct{ cmd NodeCommandAdapter }

func NewTowerAdapter(cmd NodeCommandAdapter) *TowerAdapter { return &TowerAdapter{cmd: cmd} }
func (a *TowerAdapter) SetService(ctx context.Context, nodeID, state string) error {
	b, _ := json.Marshal(map[string]string{"state": state})
	return a.cmd.Send(ctx, nodeID, "POST", "/device/v1/service", b)
}
