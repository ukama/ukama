/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/test-sys/node-gateway/cmd/version"
	"github.com/ukama/ukama/systems/test-sys/node-gateway/pkg"
	"github.com/ukama/ukama/systems/test-sys/node-gateway/pkg/rest"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
)
 
 var svcConf = pkg.NewConfig()
 
 func main() {
	 ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	 initializeNotificationConfig()
 
	 router := rest.NewRouter(rest.NewRouterConfig(svcConf))
	 router.Run()
 }
 
 func initializeNotificationConfig() {
	 svcConf = pkg.NewConfig()
	 config.LoadConfig(pkg.ServiceName, svcConf)
 }
 