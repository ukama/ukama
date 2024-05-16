/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"log"
	"time"

	"github.com/ukama/ukama/systems/common/config"
	uconf "github.com/ukama/ukama/systems/common/config"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/notification/event-notify/pkg/db"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"3s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	Service          *uconf.Service
	OrgName          string
	OrgId            string
}

const (
	EventOrgAdd             = "event.cloud.local.{{ .Org}}.nucleus.org.org.add"
	EventUserAdd            = "event.cloud.local.{{ .Org}}.nucleus.user.user.add"
	EventUserDeactivate     = "event.cloud.local.{{ .Org}}.nucleus.user.user.deactivate"
	EventUserDelete         = "event.cloud.local.{{ .Org}}.nucleus.user.user.delete"
	EventMemberCreate       = "event.cloud.local.{{ .Org}}.registry.member.member.create"
	EventMemberDelete       = "event.cloud.local.{{ .Org}}.registry.member.member.delete"
	EventNetworkAdd         = "event.cloud.local.{{ .Org}}.registry.network.network.add"
	EventNetworkDelete      = "event.cloud.local.{{ .Org}}.registry.network.network.delete"
	EventNodeCreate         = "event.cloud.local.{{ .Org}}.registry.node.node.create"
	EventNodeUpdate         = "event.cloud.local.{{ .Org}}.registry.node.node.update"
	EventNodeStateUpdate    = "event.cloud.local.{{ .Org}}.registry.node.node.state.update"
	EventNodeDelete         = "event.cloud.local.{{ .Org}}.registry.node.node.delete"
	EventNodeAssign         = "event.cloud.local.{{ .Org}}.registry.node.node.assign"
	EventNodeRelease        = "event.cloud.local.{{ .Org}}.registry.node.node.release"
	EventInviteCreate       = "event.cloud.local.{{ .Org}}.registry.invitation.invite.create"
	EventInviteDelete       = "event.cloud.local.{{ .Org}}.registry.invitation.invite.delete"
	EventInviteUpdate       = "event.cloud.local.{{ .Org}}.registry.invitation.invite.update"
	EventMeshNodeOnline     = "event.cloud.local.{{ .Org}}.messaging.mesh.node.online"
	EventMeshNodeOffline    = "event.cloud.local.{{ .Org}}.messaging.mesh.node.offline"
	EventSimActivate        = "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activate"
	EventSimAllocate        = "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.allocate"
	EventSimDelete          = "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.delete"
	EventSimAddPackage      = "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.addpackage"
	EventSimActivePackage   = "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.activepackage"
	EventSimRemovePackage   = "event.cloud.local.{{ .Org}}.subscriber.simmanager.sim.removepackage"
	EventSubscriberCreate   = "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.create"
	EventSubscriberUpdate   = "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.update"
	EventSubscriberDelete   = "event.cloud.local.{{ .Org}}.subscriber.registry.subscriber.delete"
	EventSimsUpload         = "event.cloud.local.{{ .Org}}.subscriber.simpool.sims.upload"
	EventBaserateUpload     = "event.cloud.local.{{ .Org}}.dataplan.baserate.rates.upload"
	EventPackageCreate      = "event.cloud.local.{{ .Org}}.dataplan.package.package.create"
	EventPackageUpdate      = "event.cloud.local.{{ .Org}}.dataplan.package.package.update"
	EventPackageDelete      = "event.cloud.local.{{ .Org}}.dataplan.package.package.delete"
	EventMarkupUpdate       = "event.cloud.local.{{ .Org}}.dataplan.rate.markup.update"
	EventMarkupDelete       = "event.cloud.local.{{ .Org}}.dataplan.rate.markup.delete"
	EventComponentsSync     = "event.cloud.local.{{ .Org}}.inventory.component.components.sync"
	EventAccountingSync     = "event.cloud.local.{{ .Org}}.inventory.accounting.accounting.sync"
	EventInvoiceGenerate    = "event.cloud.local.{{ .Org}}.billing.invoice.invoice.generate"
	EventInvoiceDelete      = "event.cloud.local.{{ .Org}}.billing.invoice.invoice.delete"
	EventHealthCAPPStore    = "event.cloud.local.{{ .Org}}.node.health.capps.store"
	EventNotificationDelete = "event.cloud.local.{{ .Org}}.notification.notify.notification.delete"
	EventNotificationStore  = "event.cloud.local.{{ .Org}}.notification.notify.notification.store"
)

const (
	ScopeDefault = 0
	TypeDefault  = 0
)

func newPbEvent(event proto.Message) *anypb.Any {
	pb, err := anypb.New(event)
	if err != nil {
		log.Fatalf("Failed to create protobuf event: %v", err)
	}
	return pb
}

type Event struct {
	Key         string
	Title       string
	Description string
	Scope       uint8
	Type        uint8
	PB          proto.Message
}

var RoleToScopeMap = map[db.RoleType][]db.NotificationScope{
	db.OWNER:  {db.ORG, db.NETWORK, db.SITE, db.SUBSCRIBER, db.USER, db.NODE},
	db.ADMIN:  {db.ORG, db.NETWORK, db.SITE, db.SUBSCRIBER, db.USER, db.NODE},
	db.VENDOR: {db.NETWORK},
	db.USERS:  {db.USER},
}

var EventsSTMapping = map[string]Event{
	"EventOrgAdd": {
		Key:         EventOrgAdd,
		Title:       "Organization Added",
		Description: "Organization Added",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},

	"EventUserAdd": {
		Key:         EventUserAdd,
		Title:       "User Added",
		Description: "User Added",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventUserDeactivate": {
		Key:         EventUserDeactivate,
		Title:       "User Deactivated",
		Description: "User Deactivated",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventUserDelete": {
		Key:         EventUserDelete,
		Title:       "User Deleted",
		Description: "User Deleted",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventMemberCreate": {
		Key:         EventMemberCreate,
		Title:       "Member Created",
		Description: "Member Created",
		Scope:       uint8(db.ORG),
		Type:        uint8(db.INFO),
		PB:          &epb.AddMemberEventRequest{},
	},
	"EventMemberDelete": {
		Key:         EventMemberDelete,
		Title:       "Member Deleted",
		Description: "Member Deleted",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventNetworkAdd": {
		Key:         EventNetworkAdd,
		Title:       "Network Added",
		Description: "Network Added",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventNetworkDelete": {
		Key:         EventNetworkDelete,
		Title:       "Network Deleted",
		Description: "Network Deleted",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventNodeCreate": {
		Key:         EventNodeCreate,
		Title:       "Node Created",
		Description: "Node Created",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventNodeUpdate": {
		Key:         EventNodeUpdate,
		Title:       "Node Updated",
		Description: "Node Updated",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventNodeDelete": {
		Key:         EventNodeDelete,
		Title:       "Node Deleted",
		Description: "Node Deleted",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventNodeAssign": {
		Key:         EventNodeAssign,
		Title:       "Node Assigned",
		Description: "Node Assigned",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventNodeRelease": {
		Key:         EventNodeRelease,
		Title:       "Node Released",
		Description: "Node Released",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventInviteCreate": {
		Key:         EventInviteCreate,
		Title:       "Invite Created",
		Description: "Invite Created",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventInviteDelete": {
		Key:         EventInviteDelete,
		Title:       "Invite Deleted",
		Description: "Invite Deleted",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventInviteUpdate": {
		Key:         EventInviteUpdate,
		Title:       "Invite Updated",
		Description: "Invite Updated",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventMeshNodeOnline": {
		Key:         EventMeshNodeOnline,
		Title:       "Mesh Node Online",
		Description: "Mesh Node Online",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventMeshNodeOffline": {
		Key:         EventMeshNodeOffline,
		Title:       "Mesh Node Offline",
		Description: "Mesh Node Offline",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventSimActivate": {
		Key:         EventSimActivate,
		Title:       "Sim Activated",
		Description: "Sim Activated",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventSimAllocate": {
		Key:         EventSimAllocate,
		Title:       "Sim Allocated",
		Description: "Sim Allocated",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventSimDelete": {
		Key:         EventSimDelete,
		Title:       "Sim Deleted",
		Description: "Sim Deleted",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventSimAddPackage": {
		Key:         EventSimAddPackage,
		Title:       "Sim Package Added",
		Description: "Sim Package Added",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventSimActivePackage": {
		Key:         EventSimActivePackage,
		Title:       "Sim Active Package",
		Description: "Sim Active Package",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventSimRemovePackage": {
		Key:         EventSimRemovePackage,
		Title:       "Sim Package Removed",
		Description: "Sim Package Removed",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventSubscriberCreate": {
		Key:         EventSubscriberCreate,
		Title:       "Subscriber Created",
		Description: "Subscriber Created",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventSubscriberUpdate": {
		Key:         EventSubscriberUpdate,
		Title:       "Subscriber Updated",
		Description: "Subscriber Updated",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventSubscriberDelete": {
		Key:         EventSubscriberDelete,
		Title:       "Subscriber Deleted",
		Description: "Subscriber Deleted",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventSimUpload": {
		Key:         EventSimsUpload,
		Title:       "Sim Uploaded",
		Description: "Sim Uploaded",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventBaserateUpload": {
		Key:         EventBaserateUpload,
		Title:       "Baserate uploaded",
		Description: "Baserate uploaded",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventPackageCreate": {
		Key:         EventPackageCreate,
		Title:       "Package Created",
		Description: "Package Created",
		Scope:       uint8(db.ORG),
		Type:        uint8(db.INFO),
	},
	"EventPackageUpdate": {
		Key:         EventPackageUpdate,
		Title:       "Package Updated",
		Description: "Package Updated",
		Scope:       uint8(db.ORG),
		Type:        uint8(db.INFO),
	},
	"EventPackageDelete": {
		Key:         EventPackageDelete,
		Title:       "Package Deleted",
		Description: "Package Deleted",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventMarkupUpdate": {
		Key:         EventMarkupUpdate,
		Title:       "Markup Updated",
		Description: "Markup Updated",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventMarkupDelete": {
		Key:         EventMarkupDelete,
		Title:       "Markup Deleted",
		Description: "Markup Deleted",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventComponentsSync": {
		Key:         EventComponentsSync,
		Title:       "Components Sync",
		Description: "Components Sync",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventAccountingSync": {
		Key:         EventAccountingSync,
		Title:       "Accounting Sync",
		Description: "Accounting Sync",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventInvoiceGenerate": {
		Key:         EventInvoiceGenerate,
		Title:       "Invoice Generated",
		Description: "Invoice Generated",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventInvoiceDelete": {
		Key:         EventInvoiceDelete,
		Title:       "Invoice Deleted",
		Description: "Invoice Deleted",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventHealthCAPPStore": {
		Key:         EventHealthCAPPStore,
		Title:       "Health CAPP Store",
		Description: "Health CAPP Store",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventNotificationDelete": {
		Key:         EventNotificationDelete,
		Title:       "Notification Deleted",
		Description: "Notification Deleted",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
	"EventNotificationStore": {
		Key:         EventNotificationStore,
		Title:       "Notification Stored",
		Description: "Notification Stored",
		Scope:       ScopeDefault,
		Type:        TypeDefault,
	},
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &config.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 7 * time.Second,
			ListenerRoutes: []string{
				EventOrgAdd,
				EventUserAdd,
				EventUserDeactivate,
				EventUserDelete,
				EventMemberCreate,
				EventMemberDelete,
				EventNetworkAdd,
				EventNetworkDelete,
				EventNodeCreate,
				EventNodeUpdate,
				EventNodeDelete,
				EventNodeAssign,
				EventNodeRelease,
				EventInviteCreate,
				EventInviteDelete,
				EventInviteUpdate,
				EventMeshNodeOnline,
				EventMeshNodeOffline,
				EventSimActivate,
				EventSimAllocate,
				EventSimDelete,
				EventSimAddPackage,
				EventSimActivePackage,
				EventSimRemovePackage,
				EventSubscriberCreate,
				EventSubscriberUpdate,
				EventSubscriberDelete,
				EventSimsUpload,
				EventBaserateUpload,
				EventPackageCreate,
				EventPackageUpdate,
				EventPackageDelete,
				EventMarkupUpdate,
				EventMarkupDelete,
				EventComponentsSync,
				EventAccountingSync,
				EventInvoiceGenerate,
				EventInvoiceDelete,
				EventHealthCAPPStore,
				EventNotificationDelete,
				EventNotificationStore,
			},
		},
	}
}
