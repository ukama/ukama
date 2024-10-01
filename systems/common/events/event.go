package events

import (
	notif "github.com/ukama/ukama/systems/common/notification"
)

const (
	ScopeDefault = notif.SCOPE_ORG
	TypeDefault  = notif.TYPE_INFO
)

type EventConfig struct {
	Key         EventId
	Name        string
	Title       string
	Description string
	IsProcessRequired bool
	Scope       notif.NotificationScope
	Type        notif.NotificationType
}

type EventId int

const (
	EventInvalid EventId = iota
	EventOrgAdd
	EventUserAdd
	EventUserDeactivate
	EventUserDelete
	EventMemberCreate
	EventMemberDelete
	EventNetworkAdd
	EventNetworkDelete
	EventNodeCreate
	EventNodeUpdate
	EventNodeStateUpdate
	EventNodeDelete
	EventNodeAssign
	EventNodeRelease
	EventInviteCreate
	EventInviteDelete
	EventInviteUpdate
	EventMeshNodeOnline
	EventMeshNodeOffline
	EventSimActivate
	EventSimAllocate
	EventSimDelete
	EventSimAddPackage
	EventSimActivePackage
	EventSimRemovePackage
	EventSubscriberCreate
	EventSubscriberUpdate
	EventSubscriberDelete
	EventSimsUpload
	EventBaserateUpload
	EventPackageCreate
	EventPackageUpdate
	EventPackageDelete
	EventMarkupUpdate
	EventMarkupDelete
	EventComponentsSync
	EventAccountingSync
	EventInvoiceGenerate
	EventInvoiceDelete
	EventHealthCappStore
	EventNotificationDelete
	EventNotificationStore
)

var EventRoutingKey = [...]string{
	EventOrgAdd:             "event.cloud.local.{{ .Org}}.nucleus.org.org.add",
	EventUserAdd:            "event.cloud.local.{{ .Org}}.nucleus.user.user.add",
	EventUserDeactivate:     "event.cloud.local.{{ .Org}}.nucleus.user.user.deactivate",
	EventUserDelete:         "event.cloud.local.{{ .Org}}.nucleus.user.user.delete",
	EventMemberCreate:       "event.cloud.local.{{ .Org}}.registry.member.member.create",
	EventMemberDelete:       "event.cloud.local.{{ .Org}}.registry.member.member.delete",
	EventNetworkAdd:         "event.cloud.local.{{ .Org}}.registry.network.network.add",
	EventNetworkDelete:      "event.cloud.local.{{ .Org}}.registry.network.network.delete",
	EventNodeCreate:         "event.cloud.local.{{ .Org}}.registry.node.node.create",
	EventNodeUpdate:         "event.cloud.local.{{ .Org}}.registry.node.node.update",
	EventNodeStateUpdate:    "event.cloud.local.{{ .Org}}.registry.node.node.state.update",
	EventNodeDelete:         "event.cloud.local.{{ .Org}}.registry.node.node.delete",
	EventNodeAssign:         "event.cloud.local.{{ .Org}}.registry.node.node.assign",
	EventNodeRelease:        "event.cloud.local.{{ .Org}}.registry.node.node.release",
	EventInviteCreate:       "event.cloud.local.{{ .Org}}.registry.invitation.invite.create",
	EventInviteDelete:       "event.cloud.local.{{ .Org}}.registry.invitation.invite.delete",
	EventInviteUpdate:       "event.cloud.local.{{ .Org}}.registry.invitation.invite.update",
	EventMeshNodeOnline:     "event.cloud.local.{{ .Org}}.messaging.mesh.node.online",
	EventMeshNodeOffline:    "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline",
	EventSimActivate:        "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activate",
	EventSimAllocate:        "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate",
	EventSimDelete:          "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.delete",
	EventSimAddPackage:      "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.addpackage",
	EventSimActivePackage:   "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activepackage",
	EventSimRemovePackage:   "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.removepackage",
	EventSubscriberCreate:   "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.create",
	EventSubscriberUpdate:   "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.update",
	EventSubscriberDelete:   "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.delete",
	EventSimsUpload:         "event.cloud.local.{{ .Org}}.subscriber.simpool.sims.upload",
	EventBaserateUpload:     "event.cloud.local.{{ .Org}}.dataplan.baserate.rates.upload",
	EventPackageCreate:      "event.cloud.local.{{ .Org}}.dataplan.package.package.create",
	EventPackageUpdate:      "event.cloud.local.{{ .Org}}.dataplan.package.package.update",
	EventPackageDelete:      "event.cloud.local.{{ .Org}}.dataplan.package.package.delete",
	EventMarkupUpdate:       "event.cloud.local.{{ .Org}}.dataplan.rate.markup.update",
	EventMarkupDelete:       "event.cloud.local.{{ .Org}}.dataplan.rate.markup.delete",
	EventComponentsSync:     "event.cloud.local.{{ .Org}}.inventory.component.components.sync",
	EventAccountingSync:     "event.cloud.local.{{ .Org}}.inventory.accounting.accounting.sync",
	EventInvoiceGenerate:    "event.cloud.local.{{ .Org}}.billing.invoice.invoice.generate",
	EventInvoiceDelete:      "event.cloud.local.{{ .Org}}.billing.invoice.invoice.delete",
	EventHealthCappStore:    "event.cloud.local.{{ .Org}}.node.health.capps.store",
	EventNotificationDelete: "event.cloud.local.{{ .Org}}.notification.notify.notification.delete",
	EventNotificationStore:  "event.cloud.local.{{ .Org}}.notification.notify.notification.store",
}

var EventToEventConfig = map[EventId]EventConfig{
	EventOrgAdd: {
		Key:         EventOrgAdd,
		Name:        "EventOrgAdd",
		Title:       "Organization Added",
		Description: "Organization Added",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventUserAdd: {
		Key:         EventUserAdd,
		Name:        "EventUserAdd",
		Title:       "User Added",
		Description: "User Added",
		Scope:       notif.SCOPE_USER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventUserDeactivate: {
		Key:         EventUserDeactivate,
		Name:        "EventUserDeactivate",
		Title:       "User Deactivated",
		Description: "User Deactivated",
		Scope:       notif.SCOPE_USER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventUserDelete: {
		Key:         EventUserDelete,
		Name:        "EventUserDelete",
		Title:       "User Deleted",
		Description: "User Deleted",
		Scope:       notif.SCOPE_USER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventMemberCreate: {
		Key:         EventMemberCreate,
		Name:        "EventMemberCreate",
		Title:       "Member Created",
		Description: "Member Created",
		Scope:       notif.SCOPE_USER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventMemberDelete: {
		Key:         EventMemberDelete,
		Name:        "EventMemberDelete",
		Title:       "Member Deleted",
		Description: "Member Deleted",
		Scope:       notif.SCOPE_USER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventNetworkAdd: {
		Key:         EventNetworkAdd,
		Name:        "EventNetworkAdd",
		Title:       "Network Added",
		Description: "Network Added",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventNetworkDelete: {
		Key:         EventNetworkDelete,
		Name:        "EventNetworkDelete",
		Title:       "Network Deleted",
		Description: "Network Deleted",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventNodeCreate: {
		Key:         EventNodeCreate,
		Name:        "EventNodeCreate",
		Title:       "Node Created",
		Description: "Node Created",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: true,
	},
	EventNodeUpdate: {
		Key:         EventNodeUpdate,
		Name:        "EventNodeUpdate",
		Title:       "Node Updated",
		Description: "Node Updated",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
		IsProcessRequired: true,
	},
	EventNodeStateUpdate: {
		Key:         EventNodeStateUpdate,
		Name:        "EventNodeStateUpdate",
		Title:       "Node State Updated",
		Description: "Node State Updated",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
		IsProcessRequired: true,
	},
	EventNodeDelete: {
		Key:         EventNodeDelete,
		Name:        "EventNodeDelete",
		Title:       "Node Deleted",
		Description: "Node Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: true,
	},
	EventNodeAssign: {
		Key:         EventNodeAssign,
		Name:        "EventNodeAssign",
		Title:       "Node Assigned",
		Description: "Node Assigned",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
		IsProcessRequired: true,
	},
	EventNodeRelease: {
		Key:         EventNodeRelease,
		Name:        "EventNodeRelease",
		Title:       "Node Released",
		Description: "Node Released",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
		IsProcessRequired: true,
	},
	EventInviteCreate: {
		Key:         EventInviteCreate,
		Name:        "EventInviteCreate",
		Title:       "Invite Created",
		Description: "Invite Created",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventInviteDelete: {
		Key:         EventInviteDelete,
		Name:        "EventInviteDelete",
		Title:       "Invite Deleted",
		Description: "Invite Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventInviteUpdate: {
		Key:         EventInviteUpdate,
		Name:        "EventInviteUpdate",
		Title:       "Invite Updated",
		Description: "Invite Updated",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventMeshNodeOnline: {
		Key:         EventMeshNodeOnline,
		Name:        "EventMeshNodeOnline",
		Title:       "Mesh Node Online",
		Description: "Mesh Node Online",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventMeshNodeOffline: {
		Key:         EventMeshNodeOffline,
		Name:        "EventMeshNodeOffline",
		Title:       "Mesh Node Offline",
		Description: "Mesh Node Offline",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventSimActivate: {
		Key:         EventSimActivate,
		Name:        "EventSimActivate",
		Title:       "Sim Activated",
		Description: "Sim Activated",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventSimAllocate: {
		Key:         EventSimAllocate,
		Name:        "EventSimAllocate",
		Title:       "Sim Allocated",
		Description: "Sim Allocated",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventSimDelete: {
		Key:         EventSimDelete,
		Name:        "EventSimDelete",
		Title:       "Sim Deleted",
		Description: "Sim Deleted",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventSimAddPackage: {
		Key:         EventSimAddPackage,
		Name:        "EventSimAddPackage",
		Title:       "Sim Package Added",
		Description: "Sim Package Added",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventSimActivePackage: {
		Key:         EventSimActivePackage,
		Name:        "EventSimActivePackage",
		Title:       "Sim Active Package",
		Description: "Sim Active Package",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventSimRemovePackage: {
		Key:         EventSimRemovePackage,
		Name:        "EventSimRemovePackage",
		Title:       "Sim Package Removed",
		Description: "Sim Package Removed",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventSubscriberCreate: {
		Key:         EventSubscriberCreate,
		Name:        "EventSubscriberCreate",
		Title:       "Subscriber Created",
		Description: "Subscriber Created",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventSubscriberUpdate: {
		Key:         EventSubscriberUpdate,
		Name:        "EventSubscriberUpdate",
		Title:       "Subscriber Updated",
		Description: "Subscriber Updated",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventSubscriberDelete: {
		Key:         EventSubscriberDelete,
		Name:        "EventSubscriberDelete",
		Title:       "Subscriber Deleted",
		Description: "Subscriber Deleted",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventSimsUpload: {
		Key:         EventSimsUpload,
		Name:        "EventSimsUpload",
		Title:       "Sim Uploaded",
		Description: "Sim Uploaded",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventBaserateUpload: {
		Key:         EventBaserateUpload,
		Name:        "EventBaserateUpload",
		Title:       "Baserate uploaded",
		Description: "Baserate uploaded",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventPackageCreate: {
		Key:         EventPackageCreate,
		Name:        "EventPackageCreate",
		Title:       "Package Created",
		Description: "Package Created",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventPackageUpdate: {
		Key:         EventPackageUpdate,
		Name:        "EventPackageUpdate",
		Title:       "Package Updated",
		Description: "Package Updated",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventPackageDelete: {
		Key:         EventPackageDelete,
		Name:        "EventPackageDelete",
		Title:       "Package Deleted",
		Description: "Package Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventMarkupUpdate: {
		Key:         EventMarkupUpdate,
		Name:        "EventMarkupUpdate",
		Title:       "Markup Updated",
		Description: "Markup Updated",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventMarkupDelete: {
		Key:         EventMarkupDelete,
		Name:        "EventMarkupDelete",
		Title:       "Markup Deleted",
		Description: "Markup Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventComponentsSync: {
		Key:         EventComponentsSync,
		Name:        "EventComponentsSync",
		Title:       "Components Sync",
		Description: "Components Sync",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventAccountingSync: {
		Key:         EventAccountingSync,
		Name:        "EventAccountingSync",
		Title:       "Accounting Sync",
		Description: "Accounting Sync",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventInvoiceGenerate: {
		Key:         EventInvoiceGenerate,
		Name:        "EventInvoiceGenerate",
		Title:       "Invoice Generated",
		Description: "Invoice Generated",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventInvoiceDelete: {
		Key:         EventInvoiceDelete,
		Name:        "EventInvoiceDelete",
		Title:       "Invoice Deleted",
		Description: "Invoice Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventHealthCappStore: {
		Key:         EventHealthCappStore,
		Name:        "EventHealthCappStore",
		Title:       "Health CAPP Store",
		Description: "Health CAPP Store",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventNotificationDelete: {
		Key:         EventNotificationDelete,
		Name:        "EventNotificationDelete",
		Title:       "Notification Deleted",
		Description: "Notification Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
	EventNotificationStore: {
		Key:         EventNotificationStore,
		Name:        "EventNotificationStore",
		Title:       "Notification Stored",
		Description: "Notification Stored",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
		IsProcessRequired: false,
	},
}
