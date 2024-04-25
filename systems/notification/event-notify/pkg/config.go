/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package pkg

import (
	"time"

	"github.com/ukama/ukama/systems/common/config"
	uconf "github.com/ukama/ukama/systems/common/config"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"3s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	OrgHost          string           `default:"org:9090"`
	UserHost         string           `default:"user:9090"`
	Service          *uconf.Service
	OrgName          string
	OrgId            string
}

const (
	EventOrgAdd             = "event.cloud.local.{{ .Org}}.event.cloud.local.{{ .Org}}.nucleus.org.org.add"
	EventUserAdd            = "event.cloud.local.{{ .Org}}.nucleus.user.user.add"
	EventUserDeactivate     = "event.cloud.local.{{ .Org}}.nucleus.user.user.deactivate"
	EventUserDelete         = "event.cloud.local.{{ .Org}}.nucleus.user.user.delete"
	EventMemberCreate       = "event.cloud.local.{{ .Org}}.registry.member.member.create"
	EventMemberDelete       = "event.cloud.local.{{ .Org}}.registry.member.member.delete"
	EventNetworkAdd         = "event.cloud.local.{{ .Org}}.registry.network.network.add"
	EventNetworkDelete      = "event.cloud.local.{{ .Org}}.registry.network.network.delete"
	EventNodeCreate         = "event.cloud.local.{{ .Org}}.registry.node.node.create"
	EventNodeUpdate         = "event.cloud.local.{{ .Org}}.registry.node.node.update"
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
	EventSimUpload          = "event.cloud.local.{{ .Org}}.subscriber.simpool.sim.upload"
	EventRateUpdate         = "event.cloud.local.{{ .Org}}.dataplan.baserate.rate.update"
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
				EventSimUpload,
				EventRateUpdate,
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
