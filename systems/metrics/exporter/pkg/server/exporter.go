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

	pb "github.com/ukama/ukama/systems/metrics/exporter/pb/gen"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg"
)

type ExporterServer struct {
	pb.UnimplementedExporterServiceServer
	baseRoutingKey msgbus.RoutingKeyBuilder
	org            string
	orgName        string
	msgbus         mb.MsgBusServiceClient
}

func NewExporterServer(orgName string, org string, msgBus mb.MsgBusServiceClient) (*ExporterServer, error) {

	exp := ExporterServer{
		orgName: orgName,
		org:     org,
		msgbus:  msgBus,
	}

	if msgBus != nil {
		exp.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName)
	}

	return &exp, nil
}
