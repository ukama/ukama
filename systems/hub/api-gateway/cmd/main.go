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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/metrics"

	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/hub/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/hub/api-gateway/pkg"
	server "github.com/ukama/ukama/systems/hub/api-gateway/pkg/rest"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
)

var svcConf = pkg.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	clientSet := server.NewClientsSet(&svcConf.Services)
	ac, err := providers.NewAuthClient(svcConf.Auth.AuthServerUrl, svcConf.BaseConfig.DebugMode)
	if err != nil {
		log.Errorf("Failed to create auth client: %v", err)
	}

	metrics.StartMetricsServer(&svcConf.Metrics)
	r := server.NewRouter(clientSet, server.NewRouterConfig(svcConf), ac.AuthenticateUser)
	r.Run()

}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}
