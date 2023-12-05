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

	"github.com/ukama/ukama/systems/api/api-gateway/cmd/version"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client"
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/client/rest"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/providers"

	log "github.com/sirupsen/logrus"
	prest "github.com/ukama/ukama/systems/api/api-gateway/pkg/rest"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	networkClient := client.NewNetworkClientSet(rest.NewNetworkClient(svcConf.HttpServices.RegistryHost))
	packageClient := client.NewPackageClientSet(rest.NewPackageClient(svcConf.HttpServices.DataPlanHost))
	simClient := client.NewSimClientSet(rest.NewSimClient(svcConf.HttpServices.SubscriberHost),
		rest.NewSubscriberClient(svcConf.HttpServices.SubscriberHost))
	nodeClient := client.NewNodeClientSet(rest.NewNodeClient(svcConf.HttpServices.RegistryHost))

	ac, err := providers.NewAuthClient(svcConf.Auth.AuthServerUrl, svcConf.DebugMode)
	if err != nil {
		log.Errorf("Failed to create auth client: %v", err)
	}

	router := prest.NewRouter(networkClient, packageClient, simClient, nodeClient,
		prest.NewRouterConfig(svcConf), ac.AuthenticateUser)
	router.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
	config.LoadConfig(pkg.ServiceName, svcConf)
}
