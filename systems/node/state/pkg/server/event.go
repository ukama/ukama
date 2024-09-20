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
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	stm "github.com/ukama/ukama/systems/common/stateMachine"

	"github.com/ukama/ukama/systems/common/msgbus"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"

	"github.com/ukama/ukama/systems/node/state/pkg"
)
 
 type NodeStateEventServer struct {
	 s               *NodeStateServer
	 orgName         string
	 msgbus          mb.MsgBusServiceClient
	 stateRoutingKey msgbus.RoutingKeyBuilder
	 epb.UnimplementedEventNotificationServiceServer
	 stateMachine *stm.StateMachine 
 }
 

 func (s *NodeStateServer) InitializeStateMachine(configPath string) error {
	sm, err := stm.NewStateMachine(configPath)
	if err != nil {
		return err
	}
	s.stateMachine = sm
	log.Infof("Initialized state machine with config from %s", configPath)
	
	return nil
}
 func NewNodeStateEventServer(orgName string, s *NodeStateServer, msgBus mb.MsgBusServiceClient) *NodeStateEventServer {
	 return &NodeStateEventServer{
		 s:               s,
		 orgName:         orgName,
		 msgbus:          msgBus,
		 stateRoutingKey: msgbus.NewRoutingKeyBuilder().SetCloudSource().SetSystem(pkg.SystemName).SetOrgName(orgName).SetService(pkg.ServiceName),
	 }
 }
 

