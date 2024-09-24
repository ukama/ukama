package events

type NodeEventId int

type NodeEventConfig struct {
	Key        NodeEventId
	Name       string
	RoutingKey string
}

const (
	NodeEventInvalid NodeEventId = iota
	NodeEventCreate
	NodeEventAssign
	NodeEventRelease
	NodeEventOnline
	NodeEventOffline
	NodeEventConfigUpdate
)

var NodeEventRoutingKey = map[NodeEventId]string{
	NodeEventCreate:  "event.cloud.local.{{ .Org}}.registry.node.node.create",
	NodeEventAssign:  "event.cloud.local.{{ .Org}}.registry.node.node.assign",
	NodeEventRelease: "event.cloud.local.{{ .Org}}.registry.node.node.release",
	NodeEventOnline:  "event.cloud.local.{{ .Org}}.messaging.mesh.node.online",
	NodeEventOffline: "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline",
	NodeEventConfigUpdate: "event.node.local.{{ .Org}}.messaging.mesh.config.create",
}

var NodeEventToEventConfig = map[NodeEventId]NodeEventConfig{
	NodeEventCreate: {
		Key:        NodeEventCreate,
		Name:       "online",
		RoutingKey: NodeEventRoutingKey[NodeEventCreate],
	},
	NodeEventConfigUpdate: {
		Key:        NodeEventConfigUpdate,
		Name:       "config",
		RoutingKey: NodeEventRoutingKey[NodeEventConfigUpdate],
	},
	NodeEventAssign: {
		Key:        NodeEventAssign,
		Name:       "onboarding",
		RoutingKey: NodeEventRoutingKey[NodeEventAssign],
	},
	NodeEventRelease: {
		Key:        NodeEventRelease,
		Name:       "offboarding",
		RoutingKey: NodeEventRoutingKey[NodeEventRelease],
	},
	NodeEventOffline: {
		Key:        NodeEventOffline,
		Name:       "offline",
		RoutingKey: NodeEventRoutingKey[NodeEventOffline],
	},
	NodeEventOnline: {
		Key:        NodeEventOnline,
		Name:       "online",
		RoutingKey: NodeEventRoutingKey[NodeEventOnline],
	},
}

