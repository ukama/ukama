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
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/rest/client/auth"
	"github.com/ukama/ukama/systems/report/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/report/api-gateway/internal/rest"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	internal "github.com/ukama/ukama/systems/report/api-gateway/internal"
)

var svcConf = internal.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(internal.ServiceName, os.Args, version.Version)
	initConfig()

	clientSet := rest.NewClientsSet(&svcConf.Services, &svcConf.HttpServices, svcConf.DebugMode)

	metrics.StartMetricsServer(&svcConf.Metrics)

	r := rest.NewRouter(clientSet, rest.NewRouterConfig(svcConf),
		auth.NewAuthClient(svcConf.Auth.AuthServerUrl, client.WithDebug(svcConf.DebugMode)).AuthenticateUser)
	r.Run()
}

func initConfig() {
	svcConf = internal.NewConfig()
	config.LoadConfig(internal.ServiceName, svcConf)
}
