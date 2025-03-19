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

	"github.com/ukama/ukama/systems/common/providers"
	"github.com/ukama/ukama/systems/metrics/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg"
	"github.com/ukama/ukama/systems/metrics/api-gateway/pkg/rest"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
	ic "github.com/ukama/ukama/systems/common/initclient"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
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

	regUrl, err := ic.GetHostUrl(ic.CreateHostString(svcConf.OrgName, "registry"), svcConf.Http.InitClient, &svcConf.OrgName, svcConf.DebugMode)
	if err != nil {
		log.Errorf("Failed to resolve registry address: %v", err)
	}

	subUrl, err := ic.GetHostUrl(ic.CreateHostString(svcConf.OrgName, "subscriber"), svcConf.Http.InitClient, &svcConf.OrgName, svcConf.DebugMode)
	if err != nil {
		log.Errorf("Failed to resolve subscriber address: %v", err)
	}

	nodeClient := creg.NewNodeClient(regUrl.String())
	networkClient := creg.NewNetworkClient(regUrl.String())
	siteClient := creg.NewSiteClient(regUrl.String())
	subClient := csub.NewSubscriberClient(subUrl.String())
	simClient := csub.NewSimClient(subUrl.String())

	r := rest.NewRouter(clientSet, rest.NewRouterConfig(svcConf), m, ac.AuthenticateUser, networkClient, siteClient, nodeClient, subClient, simClient)
	r.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}
