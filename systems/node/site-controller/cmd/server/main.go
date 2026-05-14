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
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	epb "github.com/ukama/ukama/systems/common/pb/gen/events"

	"github.com/ukama/ukama/systems/common/rest/client"
	ic "github.com/ukama/ukama/systems/common/rest/client/initclient"
	creg "github.com/ukama/ukama/systems/common/rest/client/registry"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/site-controller/cmd/version"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"github.com/ukama/ukama/systems/node/site-controller/pkg"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/adapters"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/reconciler"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/server"
	"github.com/ukama/ukama/systems/node/site-controller/providers"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

var svcConf *pkg.Config
const registrySystemName = "registry"

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")
	initConfig()
	gormdb := initDb()
	runGrpcServer(gormdb)
	waitForExit()
}
func initConfig() {
	svcConf = pkg.NewConfig(pkg.ServiceName)
	if err := config.NewConfReader(pkg.ServiceName).Read(svcConf); err != nil {
		log.Fatal("Error reading config ", err)
	} else if svcConf.DebugMode {
		if b, err := yaml.Marshal(svcConf); err == nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = svcConf.DebugMode
}

func initDb() sql.Db {
	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)
	if err := d.Init(&db.Site{}, &db.SiteIntent{}, &db.SiteState{}, &db.SiteComponent{}, &db.SitePortMap{}); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		instanceId = uuid.NewV4().String()
	}
	_ = instanceId

	mbClient := mb.NewMsgBusClient(svcConf.MsgClient.Timeout, svcConf.OrgName, pkg.SystemName, pkg.ServiceName, instanceId, svcConf.Queue.Uri, svcConf.Service.Uri, svcConf.MsgClient.Host, svcConf.MsgClient.Exchange, svcConf.MsgClient.ListenQueue, svcConf.MsgClient.PublishQueue, svcConf.MsgClient.RetryCount, svcConf.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	cmdAdapter, err := adapters.NewControllerAdapter(svcConf.Service.Host, svcConf.Timeout)
	if err != nil {
		log.Fatalf("failed to connect controller: %v", err)
	}
	r := reconciler.New(db.NewIntentRepo(gormdb), db.NewStateRepo(gormdb), db.NewPortMapRepo(gormdb), db.NewComponentRepo(gormdb), adapters.NewTowerAdapter(cmdAdapter), adapters.NewAmplifierAdapter(cmdAdapter), adapters.NewCNodeAdapter(cmdAdapter))

	regUrl, err := ic.GetHostUrl(ic.NewInitClient(svcConf.Http.InitClient, client.WithDebug(svcConf.DebugMode)),
		ic.CreateHostString(svcConf.OrgName, registrySystemName), &svcConf.OrgName)
	if err != nil {
		log.Errorf("Failed to resolve registry address: %v", err)
	}

	nodeClient := creg.NewNodeClient(regUrl.String())

	srv := server.NewSiteControllerServer(svcConf.OrgName, r, mbClient, nodeClient, providers.NewHealthClientProvider(svcConf.HealthHost))
	eventServer := server.NewSiteControllerEventServer(svcConf.OrgName, srv)
	
	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		pb.RegisterSiteControllerServiceServer(s, srv)
		epb.RegisterEventNotificationServiceServer(s, eventServer)

	})
	go grpcServer.StartServer()
	go msgBusListener(mbClient)
	log.Infof("Starting %s", pkg.ServiceName)
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
	go func() { sig := <-sigs; log.Info(sig); done <- true }()
	<-done
}
