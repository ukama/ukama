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
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	stm "github.com/ukama/ukama/systems/common/stateMachine"
)
 type NodeStateEventServer struct {
	 s                      *StateServer
	 orgName                string
	 stateMachineInstance   *stm.StateMachine 
	 epb.UnimplementedEventNotificationServiceServer
 }
 
 func NewNodeStateEventServer(orgName string, s *StateServer) *NodeStateEventServer {
	sm := stm.NewStateMachine(func(event stm.Event) {
		if event.IsSubstate {
			log.Infof("Substate Transition: Event: %s, Old Substate: %s, New Substate: %s", event.Name, event.OldState, event.NewState)
		} else {
			log.Infof("Main State Transition: Event: %s, Old State: %s, New State: %s", event.Name, event.OldState, event.NewState)
		}
	})

	return &NodeStateEventServer{
		s:                      s,
		orgName:                orgName,
		stateMachineInstance:   sm,
	}
 }

//TODO: implement this method to process all node event and start a go routine to process them






