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
	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/metrics/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg/rest"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg/sanitizer"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ic "github.com/ukama/ukama/systems/common/initclient"
)

const (
	registrySystem = "registry"
)

var svcConf = pkg.NewConfig()

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	log.Infof("Config %+v", svcConf.MetricsConfig)
	clientSet := rest.NewClientsSet(&svcConf.Services, svcConf.MetricsStore, svcConf.DebugMode)

	metrics.StartMetricsServer(&svcConf.MetricsServer)

	// Looking up registry system's host from initClient
	registrySystemUrl, err := ic.GetHostUrl(ic.CreateHostString(svcConf.OrgName, registrySystem),
		svcConf.Http.InitClient, &svcConf.OrgName, svcConf.DebugMode)
	if err != nil {
		log.Fatalf("Failed to resolve %s system address from initClient: %v", registrySystem, err)
	}

	m, err := pkg.NewMetrics(svcConf.MetricsConfig)
	if err != nil {
		panic("Error creating NodeMetrics. Error: " + err.Error())
	}

	sanitizer := sanitizer.NewMetricSanitizer(m, registrySystemUrl.String(), svcConf.Period)
	defer sanitizer.Stop()

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
