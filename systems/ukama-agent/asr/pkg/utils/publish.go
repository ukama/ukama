/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package utils

import (
	"github.com/ukama/ukama/systems/common/msgbus"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg"
	"google.golang.org/protobuf/reflect/protoreflect"

	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
)

func PublishEvent(orgname, action, object string, msg protoreflect.ProtoMessage,
	msgBus mb.MsgBusServiceClient) error {
	var err error
	if msgBus != nil {
		baseRoutingKey := msgbus.NewRoutingKeyBuilder().SetEventType().SetCloudSource().
			SetSystem(pkg.SystemName).SetOrgName("ukama").SetService(pkg.ServiceName)

		route := baseRoutingKey.SetObject(object).SetAction(action).MustBuild()
		err = msgBus.PublishRequest(route, msg)
	}
	return err
}
