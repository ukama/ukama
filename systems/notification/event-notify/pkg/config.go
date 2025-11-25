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

	uconf "github.com/ukama/ukama/systems/common/config"
	evt "github.com/ukama/ukama/systems/common/events"
)

type Config struct {
	uconf.BaseConfig `mapstructure:",squash"`
	DB               *uconf.Database  `default:"{}"`
	Grpc             *uconf.Grpc      `default:"{}"`
	Queue            *uconf.Queue     `default:"{}"`
	Timeout          time.Duration    `default:"3s"`
	MsgClient        *uconf.MsgClient `default:"{}"`
	Service          *uconf.Service
	Http             HttpServices
	OrgName          string
	OrgId            string
	OwnerId          string
}

type HttpServices struct {
	NucleusClient string `default:"api-gateway-nucleus:8080"`
	InitClient    string `default:"api-gateway-init:8080"`
}

func NewConfig(name string) *Config {
	return &Config{
		DB: &uconf.Database{
			DbName: name,
		},
		Service: uconf.LoadServiceHostConfig(name),
		MsgClient: &uconf.MsgClient{
			Timeout: 7 * time.Second,
			ListenerRoutes: []string{
				evt.EventRoutingKey[evt.EventOrgAdd],
				evt.EventRoutingKey[evt.EventUserAdd],
				evt.EventRoutingKey[evt.EventUserDeactivate],
				evt.EventRoutingKey[evt.EventUserDelete],
				evt.EventRoutingKey[evt.EventMemberCreate],
				evt.EventRoutingKey[evt.EventMemberDelete],
				evt.EventRoutingKey[evt.EventNetworkAdd],
				evt.EventRoutingKey[evt.EventNetworkDelete],
				evt.EventRoutingKey[evt.EventNodeCreate],
				evt.EventRoutingKey[evt.EventNodeUpdate],
				evt.EventRoutingKey[evt.EventNodeDelete],
				evt.EventRoutingKey[evt.EventNodeAssign],
				evt.EventRoutingKey[evt.EventNodeRelease],
				evt.EventRoutingKey[evt.EventInviteCreate],
				evt.EventRoutingKey[evt.EventInviteDelete],
				evt.EventRoutingKey[evt.EventInviteUpdate],
				evt.EventRoutingKey[evt.EventNodeOnline],
				evt.EventRoutingKey[evt.EventNodeOffline],
				evt.EventRoutingKey[evt.EventSimActivate],
				evt.EventRoutingKey[evt.EventSimAllocate],
				evt.EventRoutingKey[evt.EventSimDelete],
				evt.EventRoutingKey[evt.EventSimAddPackage],
				evt.EventRoutingKey[evt.EventSimActivePackage],
				evt.EventRoutingKey[evt.EventSimRemovePackage],
				evt.EventRoutingKey[evt.EventSubscriberCreate],
				evt.EventRoutingKey[evt.EventSubscriberUpdate],
				evt.EventRoutingKey[evt.EventSubscriberDelete],
				evt.EventRoutingKey[evt.EventSimsUpload],
				evt.EventRoutingKey[evt.EventBaserateUpload],
				evt.EventRoutingKey[evt.EventPackageCreate],
				evt.EventRoutingKey[evt.EventPackageUpdate],
				evt.EventRoutingKey[evt.EventPackageDelete],
				evt.EventRoutingKey[evt.EventMarkupUpdate],
				evt.EventRoutingKey[evt.EventMarkupDelete],
				evt.EventRoutingKey[evt.EventComponentsSync],
				evt.EventRoutingKey[evt.EventAccountingSync],
				evt.EventRoutingKey[evt.EventInvoiceGenerate],
				evt.EventRoutingKey[evt.EventInvoiceDelete],
				evt.EventRoutingKey[evt.EventHealthCappStore],
				evt.EventRoutingKey[evt.EventNotificationDelete],
				evt.EventRoutingKey[evt.EventNotificationStore],
				evt.EventRoutingKey[evt.EventPaymentSuccess],
				evt.EventRoutingKey[evt.EventPaymentFailed],
				evt.EventRoutingKey[evt.EventNodeStateTransition],
			}},
	}
}
