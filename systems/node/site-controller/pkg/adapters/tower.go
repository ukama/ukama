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
