/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package server

import (
	log "github.com/sirupsen/logrus"
	stm "github.com/ukama/ukama/systems/common/stateMachine"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/node/state/pb/gen"
	"github.com/ukama/ukama/systems/node/state/pkg/db"
)
 
 type NodeStateServer struct {
	 pb.UnimplementedNodeStateServiceServer
	 sRepo                db.NodeStateRepo
	 nodeStateRoutingKey  msgbus.RoutingKeyBuilder
	 msgbus               mb.MsgBusServiceClient
	 stateMachine *stm.StateMachine 
	 debug                bool
	 orgName              string
 }
 

 func NewNodeStateServer(orgName string, sRepo db.NodeStateRepo, msgBus mb.MsgBusServiceClient, debug bool, configPath string) *NodeStateServer {
	ns := &NodeStateServer{
		sRepo:            sRepo,
		orgName:          orgName,
		msgbus:           msgBus,
		debug:            debug,
	}

	if err := ns.InitializeStateMachine(configPath); err != nil {
		log.Fatalf("Failed to initialize state machine: %v", err)
	}

	return ns
}
 

 
