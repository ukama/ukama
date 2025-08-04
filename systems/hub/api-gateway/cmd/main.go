/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/auth"
	"github.com/ukama/ukama/systems/hub/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/hub/api-gateway/pkg"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	server "github.com/ukama/ukama/systems/hub/api-gateway/pkg/rest"
)

var svcConf = pkg.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	clientSet := server.NewClientsSet(&svcConf.Services)

	metrics.StartMetricsServer(&svcConf.Metrics)
	r := server.NewRouter(clientSet, server.NewRouterConfig(svcConf),
		auth.NewAuthClient(svcConf.Auth.AuthServerUrl, client.WithDebug()).AuthenticateUser)
	r.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}
