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
	"time"

	"github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/hub/hub/cmd/version"
	"github.com/ukama/ukama/systems/hub/hub/pkg"
	"github.com/ukama/ukama/systems/hub/hub/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()

	storage := pkg.NewMinioWrapper(&serviceConfig.Storage)
	chunker := pkg.NewChunker(&serviceConfig.Chunker, storage)

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		instanceId = uuid.NewV4().String()
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, pkg.SystemName,
		pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	r := server.NewRouter(&serviceConfig.Server, storage, chunker,
		time.Duration(serviceConfig.Storage.TimeoutSecond)*time.Second, mbClient)

	metrics.StartMetricsServer(serviceConfig.Metrics)

	go msgBusListener(mbClient)

	r.Run()
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	config.LoadConfig(pkg.ServiceName, serviceConfig)

	pkg.IsDebugMode = serviceConfig.DebugMode
}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}
	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
