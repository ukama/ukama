/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package main

import (
	"os"

	"github.com/num30/config"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/systems/init/reflector/cmd/version"
	"github.com/ukama/ukama/systems/init/reflector/pkg"
	"github.com/ukama/ukama/systems/init/reflector/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	pb "github.com/ukama/ukama/systems/init/reflector/pb/gen"
)

var FactorySystem = "factory"
var MessagingSystem = "messaging"

var svcConf = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")
	initConfig()
	runGrpcServer()
	log.Infof("Starting %s", pkg.ServiceName)
}

func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(svcConf)
	if err != nil {
		log.Fatal("Error reading config ", err)
	}

	if svcConf.DebugMode {
		b, err := yaml.Marshal(svcConf)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = svcConf.DebugMode
}

func runGrpcServer() {
	reflectorServer := server.NewReflectorServer(svcConf)

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		pb.RegisterReflectorServiceServer(s, reflectorServer)
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

	log.Debug("awaiting terminate/interrupt signal")
	<-done
	log.Infof("exiting service %s", pkg.ServiceName)
}