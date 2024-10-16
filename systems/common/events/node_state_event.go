/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package events

import (
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	pb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
)

type NodeStateEventId int

type NodeStateEventsConfig struct {
	Key        NodeStateEventId
	Name       string
	RoutingKey string
	body  interface{} 
}

const (
	NodeStateEventInvalid NodeStateEventId = iota
	NodeStateEventCreate
	NodeStateEventAssign
	NodeStateEventRelease
	NodeStateEventOnline
	NodeStateEventOffline
	NodeStateEventUpdate
	NodeStateEventConfig
	NodeStateEventReady
)

var NodeStateEventRoutingKey = map[NodeStateEventId]string{
	NodeStateEventCreate:  "event.cloud.local.{{ .Org}}.registry.node.node.create",
	NodeStateEventAssign:  "event.cloud.local.{{ .Org}}.registry.node.node.assign",
	NodeStateEventRelease: "event.cloud.local.{{ .Org}}.registry.node.node.release",
	NodeStateEventOnline:  "event.cloud.local.{{ .Org}}.messaging.mesh.node.online",
	NodeStateEventOffline: "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline",
	NodeStateEventUpdate:  "event.node.local.{{ .Org}}.messaging.mesh.config.create",
	NodeStateEventConfig:  " event.cloud.local.{{ .Org}}.node.configurator.node.publish",
	NodeStateEventReady:   "event.cloud.local.{{ .Org}}.messaging.mesh.node.ready",
}

var NodeEventToEventConfig = map[NodeStateEventId]NodeStateEventsConfig{
	NodeStateEventCreate: {
		Key:        NodeStateEventCreate,
		Name:       "online",
		RoutingKey: NodeStateEventRoutingKey[NodeStateEventCreate],
		body:epb.NodeCreatedEvent{},
	},
	NodeStateEventUpdate: {
		Key:        NodeStateEventUpdate,
		Name:       "update",
		RoutingKey: NodeStateEventRoutingKey[NodeStateEventUpdate],
		 body:epb.NodeUpdatedEvent{},
	},
	NodeStateEventAssign: {
		Key:        NodeStateEventAssign,
		Name:       "onboarding",
		RoutingKey: NodeStateEventRoutingKey[NodeStateEventAssign],
		body:epb.NodeAssignedEvent{},
	},
	NodeStateEventRelease: {
		Key:        NodeStateEventRelease,
		Name:       "offboarding",
		RoutingKey: NodeStateEventRoutingKey[NodeStateEventRelease],
		body:epb.NodeReleasedEvent{},
	},
	NodeStateEventOffline: {
		Key:        NodeStateEventOffline,
		Name:       "offline",
		RoutingKey: NodeStateEventRoutingKey[NodeStateEventOffline],
		body:epb.NodeOfflineEvent{},
	},
	NodeStateEventOnline: {
		Key:        NodeStateEventOnline,
		Name:       "online",
		RoutingKey: NodeStateEventRoutingKey[NodeStateEventOnline],
		body:epb.NodeOnlineEvent{},
	},
	NodeStateEventConfig: {
		Key:        NodeStateEventConfig,
		Name:       "config",
		RoutingKey: NodeStateEventRoutingKey[NodeStateEventConfig],
		 body:pb.NodeFeederMessage{},
	},
	NodeStateEventReady: {
		Key:        NodeStateEventReady,
		Name:       "ready",
		RoutingKey: NodeStateEventRoutingKey[NodeStateEventReady],
		// body:epb.NodeReadyEvent{},
	},
}
