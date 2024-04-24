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
}

const (
	EventOrgAdd             = "nucleus.org.org.add"
	EventUserAdd            = "nucleus.user.user.add"
	EventUserDeactivate     = "nucleus.user.user.deactivate"
	EventUserDelete         = "nucleus.user.user.delete"
	EventMemberCreate       = "registry.member.member.create"
	EventMemberDelete       = "registry.member.member.delete"
	EventNetworkAdd         = "registry.network.network.add"
	EventNetworkDelete      = "registry.network.network.delete"
	EventNodeCreate         = "registry.node.node.create"
	EventNodeUpdate         = "registry.node.node.update"
	EventNodeDelete         = "registry.node.node.delete"
	EventNodeAssign         = "registry.node.node.assign"
	EventNodeRelease        = "registry.node.node.release"
	EventInviteCreate       = "registry.invitation.invite.create"
	EventInviteDelete       = "registry.invitation.invite.delete"
	EventInviteUpdate       = "registry.invitation.invite.update"
	EventMeshNodeOnline     = "messaging.mesh.node.online"
	EventMeshNodeOffline    = "messaging.mesh.node.offline"
	EventSimActivate        = "subscriber.simmanager.sim.activate"
	EventSimAllocate        = "subscriber.simmanager.sim.allocate"
	EventSimDelete          = "subscriber.simmanager.sim.delete"
	EventSimAddPackage      = "subscriber.simmanager.sim.addpackage"
	EventSimActivePackage   = "subscriber.simmanager.sim.activepackage"
	EventSimRemovePackage   = "subscriber.simmanager.sim.removepackage"
	EventSubscriberCreate   = "subscriber.registry.subscriber.create"
	EventSubscriberUpdate   = "subscriber.registry.subscriber.update"
	EventSubscriberDelete   = "subscriber.registry.subscriber.delete"
	EventSimUpload          = "subscriber.simpool.sim.upload"
	EventRateUpdate         = "dataplan.baserate.rate.update"
	EventPackageCreate      = "dataplan.package.package.create"
	EventPackageUpdate      = "dataplan.package.package.update"
	EventPackageDelete      = "dataplan.package.package.delete"
	EventMarkupUpdate       = "dataplan.rate.markup.update"
	EventMarkupDelete       = "dataplan.rate.markup.delete"
	EventComponentsSync     = "inventory.component.components.sync"
	EventAccountingSync     = "inventory.accounting.accounting.sync"
	EventInvoiceGenerate    = "billing.invoice.invoice.generate"
	EventInvoiceDelete      = "billing.invoice.invoice.delete"
	EventHealthCAPPStore    = "node.health.capps.store"
	EventNotificationDelete = "notification.notify.notification.delete"
	EventNotificationStore  = "notification.notify.notification.store"
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
