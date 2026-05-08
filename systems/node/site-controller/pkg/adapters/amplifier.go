package adapters

import (
	"context"
	"encoding/json"
)

type AmplifierAdapter struct{ cmd NodeCommandAdapter }

func NewAmplifierAdapter(cmd NodeCommandAdapter) *AmplifierAdapter {
	return &AmplifierAdapter{cmd: cmd}
}
func (a *AmplifierAdapter) SetRadio(ctx context.Context, nodeID, state string) error {
	b, _ := json.Marshal(map[string]string{"state": state})
	return a.cmd.Send(ctx, nodeID, "POST", "/device/v1/radio", b)
}
