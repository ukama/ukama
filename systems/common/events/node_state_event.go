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
	NodeStateEventCreate
	NodeStateEventAssign
	NodeStateEventRelease
	NodeStateEventOnline
	NodeStateEventOffline
	NodeStateEventUpdate
	NodeStateEventConfig
	NodeStateEventDowngrade
	NodeStateEventUpgrade
)

var NodeEventRoutingKey = map[NodeStateEventId]string{
	NodeStateEventCreate:  "event.cloud.local.{{ .Org}}.registry.node.node.create",
	NodeStateEventAssign:  "event.cloud.local.{{ .Org}}.registry.node.node.assign",
	NodeStateEventRelease: "event.cloud.local.{{ .Org}}.registry.node.node.release",
	NodeStateEventOnline:  "event.cloud.local.{{ .Org}}.messaging.mesh.node.online",
	NodeStateEventOffline: "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline",
	NodeStateEventUpdate: "event.node.local.{{ .Org}}.messaging.mesh.config.create",
	NodeStateEventConfig: "event.cloud.local.{{ .Org}}.messaging.mesh.config.create",
	NodeStateEventDowngrade: "event.cloud.local.{{ .Org}}.messaging.mesh.config.downgrade",
	NodeStateEventUpgrade: "event.cloud.local.{{ .Org}}.messaging.mesh.config.upgrade",
}

var NodeEventToEventConfig = map[NodeStateEventId]NodeStateEventsConfig{
	NodeStateEventCreate: {
		Key:        NodeStateEventCreate,
		Name:       "online",
		RoutingKey: NodeEventRoutingKey[NodeStateEventCreate],
	},
	NodeStateEventUpdate: {
		Key:        NodeStateEventUpdate,
		Name:       "update",
		RoutingKey: NodeEventRoutingKey[NodeStateEventUpdate],
	},
	NodeStateEventAssign: {
		Key:        NodeStateEventAssign,
		Name:       "onboarding",
		RoutingKey: NodeEventRoutingKey[NodeStateEventAssign],
	},
	NodeStateEventRelease: {
		Key:        NodeStateEventRelease,
		Name:       "offboarding",
		RoutingKey: NodeEventRoutingKey[NodeStateEventRelease],
	},
	NodeStateEventOffline: {
		Key:        NodeStateEventOffline,
		Name:       "offline",
		RoutingKey: NodeEventRoutingKey[NodeStateEventOffline],
	},
	NodeStateEventOnline: {
		Key:        NodeStateEventOnline,
		Name:       "online",
		RoutingKey: NodeEventRoutingKey[NodeStateEventOnline],
	},
	NodeStateEventConfig: {
		Key:        NodeStateEventConfig,
		Name:       "config",
		RoutingKey: NodeEventRoutingKey[NodeStateEventConfig],
	},
	NodeStateEventDowngrade: {
		Key:        NodeStateEventDowngrade,
		Name:       "downgrade",
		RoutingKey: NodeEventRoutingKey[NodeStateEventDowngrade],
	},
	NodeStateEventUpgrade: {
		Key:        NodeStateEventUpgrade,
		Name:       "upgrade",
		RoutingKey: NodeEventRoutingKey[NodeStateEventUpgrade],
	},
}

