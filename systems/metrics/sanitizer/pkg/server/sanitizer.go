/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	"context"

	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/metrics/sanitizer/pkg"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	pb "github.com/ukama/ukama/systems/metrics/sanitizer/pb/gen"
)

type SanitizerServer struct {
	pb.UnimplementedSanitizerServiceServer
	baseRoutingKey msgbus.RoutingKeyBuilder
	org            string
	orgName        string
	msgbus         mb.MsgBusServiceClient
}

func NewSanitizerServer(orgName string, org string, msgBus mb.MsgBusServiceClient) (*SanitizerServer, error) {
	exp := SanitizerServer{
		orgName: orgName,
		org:     org,
		msgbus:  msgBus,
	}

	if msgBus != nil {
		exp.baseRoutingKey = msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName)
	}

	return &exp, nil
}

func (s *SanitizerServer) Sanitize(ctx context.Context, req *pb.SanitizeRequest) (*pb.SanitizeResponse, error) {
	return nil, nil
}
