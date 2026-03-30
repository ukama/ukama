/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package server

import (
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/messaging/broadcaster/pkg"
)

type BroadcasterServer struct {
	broadcasterRoutingKey msgbus.RoutingKeyBuilder
	msgbus               mb.MsgBusServiceClient
	debug                bool
	orgName              string
}

func NewBroadcasterServer(orgName string, msgBus mb.MsgBusServiceClient, debug bool) *BroadcasterServer {
	return &BroadcasterServer{
		debug:                debug,
		msgbus:               msgBus,
		orgName:              orgName,
		broadcasterRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}
