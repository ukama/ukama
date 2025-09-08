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
	"github.com/ukama/ukama/systems/api/api-gateway/pkg/rest"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/rest/client/auth"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
	cclient "github.com/ukama/ukama/systems/common/rest/client"
	cdplan "github.com/ukama/ukama/systems/common/rest/client/dataplan"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	csub "github.com/ukama/ukama/systems/common/rest/client/subscriber"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	networkClient := client.NewNetworkClientSet(creg.NewNetworkClient(svcConf.HttpServices.RegistryHost))
	packageClient := client.NewPackageClientSet(cdplan.NewPackageClient(svcConf.HttpServices.DataPlanHost))
	simClient := client.NewSimClientSet(csub.NewSimClient(svcConf.HttpServices.SubscriberHost),
		csub.NewSubscriberClient(svcConf.HttpServices.SubscriberHost))
	nodeClient := client.NewNodeClientSet(creg.NewNodeClient(svcConf.HttpServices.RegistryHost))

	router := rest.NewRouter(networkClient, packageClient, simClient, nodeClient, rest.NewRouterConfig(svcConf),
		auth.NewAuthClient(svcConf.Auth.AuthServerUrl, cclient.WithDebug(svcConf.DebugMode)).AuthenticateUser)
	router.Run()
}

func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
	config.LoadConfig(pkg.ServiceName, svcConf)
}
