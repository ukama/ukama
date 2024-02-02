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

	cop "github.com/ukama/ukama/systems/common/rest/client/operator"
)

type OperatorAgentAdaper struct {
	host    string
	isDebug bool
	client  cop.OperatorClient
}

func NewOperatorAgentAdapter(operatorAgentHost string, debug bool) (*OperatorAgentAdaper, error) {
	c := cop.NewOperatorClient(operatorAgentHost)

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

func (o *OperatorAgentAdaper) GetUsages(ctx context.Context, iccid, cdrType, from, to string) (any, error) {
	// think of how to use ctx with restclient
	return o.client.GetUsages(iccid, cdrType, from, to)
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
