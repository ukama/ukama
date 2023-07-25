package adapters

import (
	"context"

	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/clients/providers"
)

type OperatorAgentAdaper struct {
	host    string
	isDebug bool
	client  providers.OperatorClient
}

func NewOperatorAgentAdapter(operatorAgentHost string, debug bool) (*OperatorAgentAdaper, error) {
	c, err := providers.NewOperatorClient(operatorAgentHost, debug)
	if err != nil {
		return nil, err
	}

	return &OperatorAgentAdaper{
		host:    operatorAgentHost,
		isDebug: debug,
		client:  c,
	}, nil
}

func (o *OperatorAgentAdaper) BindSim(ctx context.Context, iccid string) (any, error) {
	// think of how to use ctx with restclient
	return o.client.BindSim(iccid)
}

func (o *OperatorAgentAdaper) GetSim(ctx context.Context, iccid string) (any, error) {
	// think of how to use ctx with restclient
	return o.client.GetSimInfo(iccid)
}

func (o *OperatorAgentAdaper) ActivateSim(ctx context.Context, iccid string) error {
	// think of how to use ctx with restclient
	return o.client.ActivateSim(iccid)
}

func (o *OperatorAgentAdaper) DeactivateSim(ctx context.Context, iccid string) error {
	// think of how to use ctx with restclient
	return o.client.DeactivateSim(iccid)
}

func (o *OperatorAgentAdaper) TerminateSim(ctx context.Context, iccid string) error {
	// think of how to use ctx with restclient
	return o.client.TerminateSim(iccid)
}

func (t *OperatorAgentAdaper) Close() {
}
