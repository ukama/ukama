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
	 Grpc             *uconf.Grpc      `default:"{}"` 
	 Queue            *uconf.Queue     `default:"{}"` 
	 Timeout          time.Duration    `default:"3s"` 
	 Service          *uconf.Service
	 OrgName          string
	 MsgClient        *uconf.MsgClient `default:"{}"` 
	 Port             string           `default:"2112"`
	 DNodeURL         string           `default:"http://dnode:8085"` 
	 RegistryHost string `default:"http://api-gateway-registry:8080"`
 }


 
 func NewConfig(name string) *Config {
	 return &Config{
		 Service: uconf.LoadServiceHostConfig(name),
		 MsgClient: &uconf.MsgClient{
			 Timeout: 5 * time.Second,
			 ListenerRoutes: []string{
				 evt.EventRoutingKey[evt.EventNodeAssign],
				 "request.cloud.local.{{ .Org}}.node.controller.nodefeeder.publish",
			 },
		 },
	 }
 }