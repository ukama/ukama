/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/testing/services/dummy/dsubscriber/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy/dsubscriber/pkg"
)

type DsubscriberServer struct {
	pb.UnimplementedDsubscriberServiceServer
	orgName        string
	msgbus         mb.MsgBusServiceClient
	baseRoutingKey msgbus.RoutingKeyBuilder
}

func NewDsubscriberServer(orgName string, msgBus mb.MsgBusServiceClient) *DsubscriberServer {
	return &DsubscriberServer{
		orgName:        orgName,
		msgbus:         msgBus,
		baseRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	}
}
