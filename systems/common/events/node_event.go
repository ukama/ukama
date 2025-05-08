/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package events
 
 type NodeStateEventId int
 
 type NodeStateEventsConfig struct {
	 Key        NodeStateEventId
	 Name       string
	 RoutingKey string
 }
 
 const (
	 NodeStateEventInvalid NodeStateEventId = iota
	 NodeStateEventAssign
	 NodeStateEventRelease
	 NodeStateEventOnline
	 NodeStateEventOffline
	 NodeStateTransition
 )
 
 var NodeStateEventRoutingKey = map[NodeStateEventId]string{
	 NodeStateEventAssign:  "event.cloud.local.{{ .Org}}.registry.node.node.assign",
	 NodeStateEventRelease: "event.cloud.local.{{ .Org}}.registry.node.node.release",
	 NodeStateEventOnline:  "event.cloud.local.{{ .Org}}.messaging.mesh.node.online",
	 NodeStateEventOffline: "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline",
	 NodeStateTransition: "event.cloud.local.{{ .Org}}.node.state.node.transition",
 }
 
 var NodeEventToEventConfig = map[NodeStateEventId]NodeStateEventsConfig{
	
	 NodeStateEventAssign: {
		 Key:        NodeStateEventAssign,
		 Name:       "onboarding",
		 RoutingKey: NodeStateEventRoutingKey[NodeStateEventAssign],
	 },
	 NodeStateEventRelease: {
		 Key:        NodeStateEventRelease,
		 Name:       "offboarding",
		 RoutingKey: NodeStateEventRoutingKey[NodeStateEventRelease],
	 },
	 NodeStateEventOffline: {
		 Key:        NodeStateEventOffline,
		 Name:       "offline",
		 RoutingKey: NodeStateEventRoutingKey[NodeStateEventOffline],
	 },
	 NodeStateEventOnline: {
		 Key:        NodeStateEventOnline,
		 Name:       "online",
		 RoutingKey: NodeStateEventRoutingKey[NodeStateEventOnline],
	 },
	 NodeStateTransition: {
		 Key:        NodeStateTransition,
		 Name:       "transition",
		 RoutingKey: NodeStateEventRoutingKey[NodeStateTransition],
	 },
 }