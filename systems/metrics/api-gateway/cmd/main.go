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

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/metrics"

	"github.com/ukama/ukama/systems/metrics/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg/rest"
	"github.com/ukama/ukama/systems/common/providers"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
)

var svcConf = pkg.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	log.Infof("Config %+v", svcConf.MetricsConfig)
	clientSet := rest.NewClientsSet(&svcConf.Services, svcConf.MetricsStore, svcConf.DebugMode)

	metrics.StartMetricsServer(&svcConf.MetricsServer)

	m, err := pkg.NewMetrics(svcConf.MetricsConfig)
	if err != nil {
		panic("Error creating NodeMetrics. Error: " + err.Error())
	}

	ac, err := providers.NewAuthClient(svcConf.Auth.AuthServerUrl, svcConf.DebugMode)
	if err != nil {
		log.Errorf("Failed to create auth client: %v", err)
	}
	r := rest.NewRouter(clientSet, rest.NewRouterConfig(svcConf), m, ac.AuthenticateUser)
	r.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}
