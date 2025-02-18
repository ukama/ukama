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

	cuk "github.com/ukama/ukama/systems/common/rest/client/ukamaagent"
)

type UkamaAgentAdaper struct {
	host    string
	isDebug bool
	client  cuk.UkamaAgentClient
}

func NewUkamaAgentAdapter(ukamaAgentHost string, debug bool) (*UkamaAgentAdaper, error) {
	c := cuk.NewUkamaAgentClient(ukamaAgentHost)

	return &UkamaAgentAdaper{
		host:    ukamaAgentHost,
		isDebug: debug,
		client:  c,
	}, nil
}

func (u *UkamaAgentAdaper) BindSim(ctx context.Context, iccid string) (any, error) {
	// think of how to use ctx with restclient
	return u.client.BindSim(iccid)
}

func (u *UkamaAgentAdaper) GetSim(ctx context.Context, iccid string) (any, error) {
	// think of how to use ctx with restclient
	return u.client.GetSimInfo(iccid)
}

func (u *UkamaAgentAdaper) GetUsages(ctx context.Context, iccid, cdrType, from, to, region string) (any, any, error) {
	// think of how to use ctx with restclient
	return u.client.GetUsages(iccid, cdrType, from, to, region)
}

func (u *UkamaAgentAdaper) ActivateSim(ctx context.Context, req client.AgentRequestData) error {
	// think of how to use ctx with restclient
	return u.client.ActivateSim(req)
}

func (u *UkamaAgentAdaper) DeactivateSim(ctx context.Context, req client.AgentRequestData) error {
	// think of how to use ctx with restclient
	return u.client.DeactivateSim(req)
}

func (u *UkamaAgentAdaper) UpdatePackage(ctx context.Context, req client.AgentRequestData) error {
	// think of how to use ctx with restclient
	return u.client.UpdatePackage(req)
}

func (u *UkamaAgentAdaper) TerminateSim(ctx context.Context, iccid string) error {
	// think of how to use ctx with restclient
	return u.client.TerminateSim(iccid)
}

func (t *UkamaAgentAdaper) Close() {
}
