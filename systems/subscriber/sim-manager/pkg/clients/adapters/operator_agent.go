/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

//TODO: add unit tests.

package adapters

import (
	"context"

	"github.com/ukama/ukama/systems/common/rest/client"

	cop "github.com/ukama/ukama/systems/common/rest/client/operatoragent"
)

type OperatorAgentAdaper struct {
	host    string
	isDebug bool
	client  cop.OperatorAgentClient
}

func NewOperatorAgentAdapter(operatorAgentHost string, debug bool) (*OperatorAgentAdaper, error) {
	c := cop.NewOperatorAgentClient(operatorAgentHost)

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

func (o *OperatorAgentAdaper) GetUsages(ctx context.Context, iccid, cdrType, from, to, region string) (any, any, error) {
	// think of how to use ctx with restclient
	return o.client.GetUsages(iccid, cdrType, from, to, region)
}

func (o *OperatorAgentAdaper) ActivateSim(ctx context.Context, req client.AgentRequestData) error {
	// think of how to use ctx with restclient
	return o.client.ActivateSim(req.Iccid)
}

func (o *OperatorAgentAdaper) DeactivateSim(ctx context.Context, req client.AgentRequestData) error {
	// think of how to use ctx with restclient
	return o.client.DeactivateSim(req.Iccid)
}

func (o *OperatorAgentAdaper) UpdatePackage(ctx context.Context, req client.AgentRequestData) error {
	// think of how to use ctx with restclient
	return nil
}

func (o *OperatorAgentAdaper) TerminateSim(ctx context.Context, iccid string) error {
	// think of how to use ctx with restclient
	return o.client.TerminateSim(iccid)
}

func (t *OperatorAgentAdaper) Close() {
}
