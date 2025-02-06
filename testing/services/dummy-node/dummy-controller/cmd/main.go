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

	"github.com/num30/config"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"

	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/cmd/version"
	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/internal/backhaul"
	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/internal/battery"
	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/internal/solar"

	generated "github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/pb/gen"
	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/pkg"
	"github.com/ukama/ukama/testing/services/dummy-node/dummy-controller/pkg/server"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	runGrpcServer()
}
func initConfig() {
	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = serviceConfig.DebugMode
}

func runGrpcServer() {
    // Initialize providers
    solarProvider := solar.NewSolarProvider()
    backhaulProvider := backhaul.NewBackhaulProvider()
    batteryProvider := battery.NewMockBatteryProvider()  // Add this if not already defined

    controllerServer := server.NewControllerServer(
        serviceConfig.OrgName,
        solarProvider,
        backhaulProvider,
        batteryProvider,
    )

    grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
        generated.RegisterMetricsControllerServer(s, controllerServer)
    })

    go grpcServer.StartServer()
}