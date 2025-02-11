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
        log.Fatalf("Error reading config: %v", err)
    }
    
    if serviceConfig.DebugMode {
        b, err := yaml.Marshal(serviceConfig)
        if err != nil {
            log.Warnf("Error marshaling config: %v", err)
        } else {
            log.Infof("Config:\n%s", string(b))
        }
    }
    pkg.IsDebugMode = serviceConfig.DebugMode
}

func runGrpcServer() {
  
    controllerServer := server.NewControllerServer(
        serviceConfig.OrgName,
    )

    grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
        generated.RegisterMetricsControllerServer(s, controllerServer)
    })

	go grpcServer.StartServer()


	waitForExit()
}
func waitForExit() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		log.Info(sig)
		done <- true
	}()

	log.Debug("awaiting terminate/interrrupt signal")
	<-done
	log.Infof("exiting service %s", pkg.ServiceName)
}