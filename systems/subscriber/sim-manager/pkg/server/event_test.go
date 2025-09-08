/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server_test

import (
	"testing"
)

const (
	OrgName = "testOrg"
	simId   = "592f7a8e-f318-4d3a-aab8-8d4187cde7f9"
	// webhookurl = "http://webhooks:8080/reports"

	// sessionId = 22
	// bmId      = "e044081b-fbbe-45e9-8f78-0f9c0f112977"
	// custId    = "e231a7cd-03f6-470a-9e8c-e02f54f9b415"
	// planId    = "0f8be763-8bd6-406d-9d82-158d7f1a2140"
)

func TesSimManagerEventServer_HandleSimManagerSimAllocateEvent(t *testing.T) {
	// s := server.NewSimManagerEventServer(OrgName, &mocks.SimManagerServiceServer{})

	// t.Run("NewSimAllocated", func(t *testing.T) {
	// routingKey := msgbus.PrepareRoute(OrgName,
	// "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate")

	// allocatedSim := epb.EventSimAllocation{
	// Id: simId,
	// }

	// anyE, err := anypb.New(&allocatedSim)
	// assert.NoError(t, err)

	// msg := &epb.Event{
	// RoutingKey: routingKey,
	// Msg:        anyE,
	// }

	// _, err = s.EventNotification(context.TODO(), msg)

	// assert.Error(t, err)
	// })
}
