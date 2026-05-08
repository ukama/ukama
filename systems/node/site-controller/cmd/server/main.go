package main

import (
	"os"

	"github.com/num30/config"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/site-controller/cmd/version"
	pb "github.com/ukama/ukama/systems/node/site-controller/pb/gen"
	"github.com/ukama/ukama/systems/node/site-controller/pkg"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/adapters"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/db"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/reconciler"
	"github.com/ukama/ukama/systems/node/site-controller/pkg/server"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")
	initConfig()
	siteDb := initDb()
	runGrpcServer(siteDb)
	log.Infof("Starting %s", pkg.ServiceName)
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
	log.Infof("Initializing Database")
	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)
	if err := d.Init(&db.SiteIntent{}, &db.SiteState{}, &db.SiteComponent{}, &db.SiteSwitchPolicy{}); err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		instanceId = uuid.NewV4().String()
	}

	mbClient := mb.NewMsgBusClient(
		svcConf.MsgClient.Timeout,
		svcConf.OrgName,
		pkg.SystemName,
		pkg.ServiceName,
		instanceId,
		svcConf.Queue.Uri,
		svcConf.Service.Uri,
		svcConf.MsgClient.Host,
		svcConf.MsgClient.Exchange,
		svcConf.MsgClient.ListenQueue,
		svcConf.MsgClient.PublishQueue,
		svcConf.MsgClient.RetryCount,
		svcConf.MsgClient.ListenerRoutes,
	)

	cmdAdapter, err := adapters.NewControllerAdapter(svcConf.Services.Controller, svcConf.Services.Timeout)
	if err != nil {
		log.Fatalf("failed to connect controller: %v", err)
	}

	r := reconciler.New(
		db.NewIntentRepo(gormdb),
		db.NewStateRepo(gormdb),
		db.NewSwitchPolicyRepo(gormdb),
		adapters.NewTowerAdapter(cmdAdapter),
		adapters.NewAmplifierAdapter(cmdAdapter),
		adapters.NewCNodeAdapter(cmdAdapter),
	)

	srv := server.NewSiteControllerServer(r)
	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		pb.RegisterSiteControllerServiceServer(s, srv)
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
		log.Fatalf("Failed to start Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
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
	log.Debug("awaiting terminate/interrupt signal")
	<-done
	log.Infof("exiting service %s", pkg.ServiceName)
}
