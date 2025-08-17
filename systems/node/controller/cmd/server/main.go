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
	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/systems/common/rest/client"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/controller/cmd/version"
	"github.com/ukama/ukama/systems/node/controller/pkg"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"
	"github.com/ukama/ukama/systems/node/controller/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
)

const registrySystemName = "registry"

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")
	initConfig()
	cDb := initDb()
	runGrpcServer(cDb)
	log.Infof("Starting %s", pkg.ServiceName)
}

func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(svcConf)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if svcConf.DebugMode {
		b, err := yaml.Marshal(svcConf)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = svcConf.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)
	err := d.Init(&db.NodeLog{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}
func runGrpcServer(gormdb sql.Db) {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	//TODO: We should do initclient resolution on demand, in order to avoid systems url changes side effects
	regUrl, err := ic.GetHostUrl(ic.NewInitClient(svcConf.Http.InitClient, client.WithDebug(svcConf.DebugMode)),
		ic.CreateHostString(svcConf.OrgName, registrySystemName), &svcConf.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve registry address: %v", err)
	}

	cnet := creg.NewNetworkClient(regUrl.String())
	csite := creg.NewSiteClient(regUrl.String())
	cnode := creg.NewNodeClient(regUrl.String())

	mbClient := mb.NewMsgBusClient(svcConf.MsgClient.Timeout, svcConf.OrgName, pkg.SystemName, pkg.ServiceName, instanceId, svcConf.Queue.Uri, svcConf.Service.Uri, svcConf.MsgClient.Host, svcConf.MsgClient.Exchange, svcConf.MsgClient.ListenQueue, svcConf.MsgClient.PublishQueue, svcConf.MsgClient.RetryCount, svcConf.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)
	contServer := server.NewControllerServer(svcConf.OrgName, db.NewNodeLogRepo(gormdb),
		mbClient, cnet, csite, cnode, svcConf.DebugMode)
	controllerEventServer := server.NewControllerEventServer(svcConf.OrgName, contServer)

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		pb.RegisterControllerServiceServer(s, contServer)
		epb.RegisterEventNotificationServiceServer(s, controllerEventServer)

	})

	go grpcServer.StartServer()

	go msgBusListener(mbClient)

	waitForExit()
}

func msgBusListener(m mb.MsgBusServiceClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
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
