/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ukama/ukama/nodes/apps/pcrf/cmd/version"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/client"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/controller"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/rest"
	"github.com/ukama/ukama/nodes/apps/pcrf/pkg/service"
	"github.com/ukama/ukama/systems/common/metrics"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
)

var svcConf *pkg.Config

func main() {
	var (
		err            error
		pcrfPort       int
		nodedURL       string
		initNetworkURL string
		nodeClient     interface {
			GetNodeId() (string, error)
		}
		initNetworkClient *client.InitNetworkClient
		initStatus        *client.InitNetworkStatus
		nodeId            string
		ctr               *controller.Controller
	)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	svcConf = pkg.NewConfig(pkg.ServiceName)

	if svcConf.DebugMode {
		log.SetLevel(log.DebugLevel)
	}

	pcrfPort, err = service.Port(pkg.ServiceName)
	if err != nil {
		log.Fatalf("Failed to resolve %s from /etc/services: %v",
			pkg.ServiceName, err)
	}
	svcConf.Server.Port = pcrfPort

	nodedURL, err = service.LocalURL(pkg.NodedServiceName)
	if err != nil {
		log.Fatalf("Failed to resolve %s from /etc/services: %v",
			pkg.NodedServiceName, err)
	}

	initNetworkURL, err = service.LocalURL(pkg.InitNetworkServiceName)
	if err != nil {
		log.Fatalf("Failed to resolve %s from /etc/services: %v",
			pkg.InitNetworkServiceName, err)
	}

	ukamaLocalServiceURL, err := service.LocalURL(pkg.UkamaServiceName)
	if err != nil {
		log.Fatalf("Failed to resolve %s from /etc/services: %v",
			pkg.UkamaAgentSystemName, err)
	}

	ukamaAgentURL := fmt.Sprintf("%s/%s", ukamaLocalServiceURL, pkg.UkamaAgentSystemName)

	log.Infof("Starting PCRF service %s on port %d",
		pkg.ServiceName, svcConf.Server.Port)

	initNetworkClient = client.NewInitNetworkClient(initNetworkURL)
	initStatus, err = initNetworkClient.GetStatus()
	if err != nil {
		log.Fatalf("Failed to read init-network status: %v", err)
	}

	svcConf.Bridge.Name = initStatus.Bridge.Name
	svcConf.Bridge.Ip = initStatus.Bridge.Address
	svcConf.Bridge.Management = filepath.Dir(initStatus.Bridge.ManagementSocket)

	if err = os.MkdirAll(filepath.Dir(svcConf.DB), 0755); err != nil {
		log.Fatalf("Failed to create PCRF DB directory: %v", err)
	}

	nodeClient, err = client.NewNodedClient(nodedURL, svcConf.DebugMode)
	if err != nil {
		log.Fatalf("Failed to create node client: %v", err)
	}

	nodeId, err = nodeClient.GetNodeId()
	if err != nil {
		log.Fatalf("Failed to read node info: %v", err)
	}

	ukamaAgentClient, err := client.NewRemoteControllerClient(ukamaAgentURL, svcConf.DebugMode)
	if err != nil {
		log.Fatalf("Failed to create ukama agent client: %v", err)
	}

	log.Infof("PCRF running on node %s", nodeId)

	ctr, err = controller.NewController(
		svcConf.DB,
		svcConf.Bridge,
		ukamaAgentClient,
		svcConf.SyncPeriod,
		nodeId,
		svcConf.DebugMode)
	if err != nil {
		log.Fatalf("Failed to create controller: %v", err)
	}

	go sigHandler(sigs, done, ctr)

	metrics.StartMetricsServer(&svcConf.Metrics)

	r := rest.NewRouter(ctr,
		rest.NewRouterConfig(svcConf),
		nodeId,
		rest.RuntimeStatus{
			InitNetworkURL:   initNetworkURL,
			InitNetworkReady: initStatus.Ready,
			UECidr:           initStatus.UE.Cidr,
			DBPath:           svcConf.DB,
		})
	go r.Run()

	<-done
	log.Infof("Exiting service %s.", pkg.ServiceName)
}

func sigHandler(sigs chan os.Signal, done chan bool, ctr *controller.Controller) {
	sig := <-sigs

	log.Infof("Starting signal handler routine for %v", sig)
	_ = ctr.ExitController()

	done <- true
}
