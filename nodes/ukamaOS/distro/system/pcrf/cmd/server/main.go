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
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/cmd/version"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/client"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/controller"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/rest"
	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/providers"

	ccmd "github.com/ukama/ukama/systems/common/cmd"
)

var svcConf *pkg.Config

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()

	if svcConf.DebugMode {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Starting pcrf controller service %s", pkg.ServiceName)

	nodeClient, err := client.NewNodedClient(svcConf.HttpServices.Noded, svcConf.DebugMode)
	if err != nil {
		log.Fatalf("Failed to create node client: %v", err)
	}

	NodeId, err := nodeClient.GetNodeId()
	if err != nil {
		log.Fatalf("Failed to read node info: %v", err)
	}

	ctr, err := controller.NewController(svcConf.DB, svcConf.Bridge, svcConf.HttpServices.Policy, svcConf.SyncPeriod, NodeId, svcConf.DebugMode)
	if err != nil {
		log.Fatalf("Failed to create controller: %v", err)
	}

	go sigHandler(sigs, done, ctr)

	ac, err := providers.NewAuthClient(svcConf.Auth.AuthServerUrl, svcConf.DebugMode)
	if err != nil {
		log.Errorf("Failed to create auth client: %v", err)
	}

	metrics.StartMetricsServer(&svcConf.Metrics)

	r := rest.NewRouter(ctr, rest.NewRouterConfig(svcConf), NodeId, ac.AuthenticateUser)
	go r.Run()

	<-done
	log.Infof("Exiting service %s.", pkg.ServiceName)

}

func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
	config.LoadConfig(pkg.ServiceName, svcConf)
}

func sigHandler(sigs chan os.Signal, done chan bool, ctr *controller.Controller) {
	sig := <-sigs
	log.Infof("Starting signal handler routine for %v", sig)
	ctr.ExitController()
	done <- true
}
