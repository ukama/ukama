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
	Scope       notif.NotificationScope
	Type        notif.NotificationType
}

type EventId int 

const (
	EventInvalid EventId = iota
	EventOrgAdd
	EventSiteCreate
	EventSiteUpdate
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
	EventNodeOnline
	EventNodeOffline
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
	EventPaymentSuccess
	EventPaymentFailed
	EventNodeStateTransition
)

var EventRoutingKey = [...]string{
	EventOrgAdd:             "event.cloud.local.{{ .Org}}.nucleus.org.org.add",
	EventSiteCreate:         "event.cloud.local.{{ .Org}}.registry.site.site.create",
	EventSiteUpdate:         "event.cloud.local.{{ .Org}}.registry.site.site.update",
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
	EventNodeOnline:         "event.cloud.local.{{ .Org}}.messaging.mesh.node.online",
	EventNodeOffline:        "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline",
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
	EventInvoiceGenerate:    "event.cloud.local.{{ .Org}}.billing.report.invoice.generate",
	EventInvoiceDelete:      "event.cloud.local.{{ .Org}}.billing.invoice.invoice.delete",
	EventHealthCappStore:    "event.cloud.local.{{ .Org}}.node.health.capps.store",
	EventNotificationDelete: "event.cloud.local.{{ .Org}}.notification.notify.notification.delete",
	EventNotificationStore:  "event.cloud.local.{{ .Org}}.notification.notify.notification.store",
	EventPaymentSuccess:     "event.cloud.local.{{ .Org}}.payments.processor.payment.success",
	EventPaymentFailed:      "event.cloud.local.{{ .Org}}.payments.processor.payment.failed",
	EventNodeStateTransition: "event.cloud.local.{{ .Org}}.node.state.node.transition",
}

var EventToEventConfig = map[EventId]EventConfig{
	EventOrgAdd: {
		Key:         EventOrgAdd,
		Name:        "EventOrgAdd",
		Title:       "Organization Added",
		Description: "Organization Added",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventNodeStateTransition: {
		Key:         EventNodeStateTransition,
		Name:        "EventNodeStateTransition",
		Title:       "Node State Transition",
		Description: "Node State Transition",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
	},
	EventSiteCreate: {
		Key:         EventSiteCreate,
		Name:        "EventSiteCreate",
		Title:       "Site Added",
		Description: "Site Added",
		Scope:       notif.SCOPE_SITE,
		Type:        TypeDefault,
	},
	EventSiteUpdate: {
		Key:         EventSiteUpdate,
		Name:        "EventSiteUpdate",
		Title:       "Site Update",
		Description: "Site Update",
		Scope:       notif.SCOPE_SITE,
		Type:        TypeDefault,
	},
	EventUserAdd: {
		Key:         EventUserAdd,
		Name:        "EventUserAdd",
		Title:       "User Added",
		Description: "User Added",
		Scope:       notif.SCOPE_USER,
		Type:        TypeDefault,
	},
	EventUserDeactivate: {
		Key:         EventUserDeactivate,
		Name:        "EventUserDeactivate",
		Title:       "User Deactivated",
		Description: "User Deactivated",
		Scope:       notif.SCOPE_USER,
		Type:        TypeDefault,
	},
	EventUserDelete: {
		Key:         EventUserDelete,
		Name:        "EventUserDelete",
		Title:       "User Deleted",
		Description: "User Deleted",
		Scope:       notif.SCOPE_USER,
		Type:        TypeDefault,
	},
	EventMemberCreate: {
		Key:         EventMemberCreate,
		Name:        "EventMemberCreate",
		Title:       "Member Created",
		Description: "Member Created",
		Scope:       notif.SCOPE_USER,
		Type:        TypeDefault,
	},
	EventMemberDelete: {
		Key:         EventMemberDelete,
		Name:        "EventMemberDelete",
		Title:       "Member Deleted",
		Description: "Member Deleted",
		Scope:       notif.SCOPE_USER,
		Type:        TypeDefault,
	},
	EventNetworkAdd: {
		Key:         EventNetworkAdd,
		Name:        "EventNetworkAdd",
		Title:       "Network Added",
		Description: "Network Added",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
	},
	EventNetworkDelete: {
		Key:         EventNetworkDelete,
		Name:        "EventNetworkDelete",
		Title:       "Network Deleted",
		Description: "Network Deleted",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
	},
	EventNodeCreate: {
		Key:         EventNodeCreate,
		Name:        "EventNodeCreate",
		Title:       "Node added",
		Description: "Node added",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventNodeUpdate: {
		Key:         EventNodeUpdate,
		Name:        "EventNodeUpdate",
		Title:       "Node Updated",
		Description: "Node Updated",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
	},
	EventNodeStateUpdate: {
		Key:         EventNodeStateUpdate,
		Name:        "EventNodeStateUpdate",
		Title:       "Node State Updated",
		Description: "Node State Updated",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
	},
	EventNodeDelete: {
		Key:         EventNodeDelete,
		Name:        "EventNodeDelete",
		Title:       "Node Deleted",
		Description: "Node Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventNodeAssign: {
		Key:         EventNodeAssign,
		Name:        "EventNodeAssign",
		Title:       "Node Assigned",
		Description: "Node Assigned",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
	},
	EventNodeRelease: {
		Key:         EventNodeRelease,
		Name:        "EventNodeRelease",
		Title:       "Node Released",
		Description: "Node Released",
		Scope:       notif.SCOPE_NETWORK,
		Type:        TypeDefault,
	},
	EventInviteCreate: {
		Key:         EventInviteCreate,
		Name:        "EventInviteCreate",
		Title:       "Invite Created",
		Description: "Invite Created",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventInviteDelete: {
		Key:         EventInviteDelete,
		Name:        "EventInviteDelete",
		Title:       "Invite Deleted",
		Description: "Invite Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventInviteUpdate: {
		Key:         EventInviteUpdate,
		Name:        "EventInviteUpdate",
		Title:       "Invite Updated",
		Description: "Invite Updated",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventNodeOnline: {
		Key:         EventNodeOnline,
		Name:        "EventNodeOnline",
		Title:       "Node Online",
		Description: "Node Online",
		Scope:       notif.SCOPE_ORG,
		Type:        notif.TYPE_INFO,
	},
	EventNodeOffline: {
		Key:         EventNodeOffline,
		Name:        "EventNodeOffline",
		Title:       "Node Offline",
		Description: "Node Offline",
		Scope:       notif.SCOPE_ORG,
		Type:        notif.TYPE_INFO,
	},
	EventSimActivate: {
		Key:         EventSimActivate,
		Name:        "EventSimActivate",
		Title:       "Sim Activated",
		Description: "Sim Activated",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
	},
	EventSimAllocate: {
		Key:         EventSimAllocate,
		Name:        "EventSimAllocate",
		Title:       "Sim Allocated",
		Description: "Sim Allocated",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
	},
	EventSimDelete: {
		Key:         EventSimDelete,
		Name:        "EventSimDelete",
		Title:       "Sim Deleted",
		Description: "Sim Deleted",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
	},
	EventSimAddPackage: {
		Key:         EventSimAddPackage,
		Name:        "EventSimAddPackage",
		Title:       "Sim Package Added",
		Description: "Sim Package Added",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
	},
	EventSimActivePackage: {
		Key:         EventSimActivePackage,
		Name:        "EventSimActivePackage",
		Title:       "Sim Active Package",
		Description: "Sim Active Package",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
	},
	EventSimRemovePackage: {
		Key:         EventSimRemovePackage,
		Name:        "EventSimRemovePackage",
		Title:       "Sim Package Removed",
		Description: "Sim Package Removed",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
	},
	EventSubscriberCreate: {
		Key:         EventSubscriberCreate,
		Name:        "EventSubscriberCreate",
		Title:       "Subscriber Created",
		Description: "Subscriber Created",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
	},
	EventSubscriberUpdate: {
		Key:         EventSubscriberUpdate,
		Name:        "EventSubscriberUpdate",
		Title:       "Subscriber Updated",
		Description: "Subscriber Updated",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
	},
	EventSubscriberDelete: {
		Key:         EventSubscriberDelete,
		Name:        "EventSubscriberDelete",
		Title:       "Subscriber Deleted",
		Description: "Subscriber Deleted",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
	},
	EventSimsUpload: {
		Key:         EventSimsUpload,
		Name:        "EventSimsUpload",
		Title:       "Sim Uploaded",
		Description: "Sim Uploaded",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventBaserateUpload: {
		Key:         EventBaserateUpload,
		Name:        "EventBaserateUpload",
		Title:       "Baserate uploaded",
		Description: "Baserate uploaded",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventPackageCreate: {
		Key:         EventPackageCreate,
		Name:        "EventPackageCreate",
		Title:       "Package Created",
		Description: "Package Created",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventPackageUpdate: {
		Key:         EventPackageUpdate,
		Name:        "EventPackageUpdate",
		Title:       "Package Updated",
		Description: "Package Updated",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventPackageDelete: {
		Key:         EventPackageDelete,
		Name:        "EventPackageDelete",
		Title:       "Package Deleted",
		Description: "Package Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventMarkupUpdate: {
		Key:         EventMarkupUpdate,
		Name:        "EventMarkupUpdate",
		Title:       "Markup Updated",
		Description: "Markup Updated",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventMarkupDelete: {
		Key:         EventMarkupDelete,
		Name:        "EventMarkupDelete",
		Title:       "Markup Deleted",
		Description: "Markup Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventComponentsSync: {
		Key:         EventComponentsSync,
		Name:        "EventComponentsSync",
		Title:       "Components Sync",
		Description: "Components Sync",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventAccountingSync: {
		Key:         EventAccountingSync,
		Name:        "EventAccountingSync",
		Title:       "Accounting Sync",
		Description: "Accounting Sync",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventInvoiceGenerate: {
		Key:         EventInvoiceGenerate,
		Name:        "EventInvoiceGenerate",
		Title:       "Invoice Generated",
		Description: "Invoice Generated",
		Scope:       notif.SCOPE_ORG,
		Type:        notif.TYPE_ACTIONABLE_INFO,
	},
	EventInvoiceDelete: {
		Key:         EventInvoiceDelete,
		Name:        "EventInvoiceDelete",
		Title:       "Invoice Deleted",
		Description: "Invoice Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventHealthCappStore: {
		Key:         EventHealthCappStore,
		Name:        "EventHealthCappStore",
		Title:       "Health CAPP Store",
		Description: "Health CAPP Store",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventNotificationDelete: {
		Key:         EventNotificationDelete,
		Name:        "EventNotificationDelete",
		Title:       "Notification Deleted",
		Description: "Notification Deleted",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventNotificationStore: {
		Key:         EventNotificationStore,
		Name:        "EventNotificationStore",
		Title:       "Notification Stored",
		Description: "Notification Stored",
		Scope:       notif.SCOPE_ORG,
		Type:        TypeDefault,
	},
	EventPaymentSuccess: {
		Key:         EventPaymentSuccess,
		Name:        "EventPaymentSuccess",
		Title:       "Payment Success",
		Description: "Payment Success",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        TypeDefault,
	},
	EventPaymentFailed: {
		Key:         EventPaymentFailed,
		Name:        "EventPaymentFailed",
		Title:       "Payment Failed",
		Description: "Payment Failed",
		Scope:       notif.SCOPE_SUBSCRIBER,
		Type:        notif.TYPE_ERROR,
	},
}
