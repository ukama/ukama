package adapters

import (
	"context"
	"encoding/json"
	"fmt"
)

type CNodeAdapter struct{ cmd NodeCommandAdapter }

func NewCNodeAdapter(cmd NodeCommandAdapter) *CNodeAdapter { return &CNodeAdapter{cmd: cmd} }

func (a *CNodeAdapter) RequestSwitchPolicy(ctx context.Context, nodeID string) error {
	return a.cmd.Send(ctx, nodeID, "GET", "/switch/v1/ports/policy", nil)
}

func (a *CNodeAdapter) SetPortPoe(ctx context.Context, nodeID string, port int, on bool, reason string) error {
	b, _ := json.Marshal(map[string]interface{}{
		"on":     on,
		"source": "site-controller",
		"reason": reason,
	})
	return a.cmd.Send(ctx, nodeID, "POST", fmt.Sprintf("/switch/v1/ports/%d/poe", port), b)
}

func (a *CNodeAdapter) PowerCyclePort(ctx context.Context, nodeID string, port int, reason string) error {
	b, _ := json.Marshal(map[string]string{
		"source": "site-controller",
		"reason": reason,
	})
	return a.cmd.Send(ctx, nodeID, "POST", fmt.Sprintf("/switch/v1/ports/%d/poe/cycle", port), b)
}
