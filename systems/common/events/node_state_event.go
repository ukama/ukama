package events

type NodeEventId int

const (
	NodeEventInvalid NodeEventId = iota
	NodeEventCreate
	NodeEventAssign
	NodeEventRelease
	NodeEventOnline
	NodeEventOffline
)

var NodeEventRoutingKey = map[NodeEventId]string{
	NodeEventCreate:  "event.cloud.local.{{ .Org}}.registry.node.node.create",
	NodeEventAssign:  "event.cloud.local.{{ .Org}}.registry.node.node.assign",
	NodeEventRelease: "event.cloud.local.{{ .Org}}.registry.node.node.release",
	NodeEventOnline:  "event.cloud.local.{{ .Org}}.messaging.mesh.node.online",
	NodeEventOffline: "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline",
}

var NodeEventToEventConfig = map[NodeEventId]NodeEventConfig{
	NodeEventCreate: {
		Key:        NodeEventCreate,
		Name:       "online",
		RoutingKey: NodeEventRoutingKey[NodeEventCreate],
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

type NodeEventConfig struct {
	Key        NodeEventId
	Name       string
	RoutingKey string
}
