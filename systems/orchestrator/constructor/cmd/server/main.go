package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/ukama/ukama/systems/orchestrator/constructor/cmd/version"
	"github.com/ukama/ukama/systems/orchestrator/constructor/pkg"
	"github.com/ukama/ukama/systems/orchestrator/constructor/pkg/db"
	"github.com/ukama/ukama/systems/orchestrator/constructor/pkg/server"

	"github.com/num30/config"
	"gopkg.in/yaml.v3"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egenerated "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/common/sql"
	generated "github.com/ukama/ukama/systems/orchestrator/constructor/pb/gen"
	"google.golang.org/grpc"
)

var serviceConfig = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	log.Infof("Starting the base-rate service")

	initConfig()

	db := initDb()

	runGrpcServer(db)
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Org{}, &db.Deployment{}, &db.Config{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func initConfig() {
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

func runGrpcServer(d sql.Db) {
	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	mbClient := mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, pkg.SystemName,
		pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		srv := server.NewConstructorServer(db.NewOrgRepo(d), db.NewDeploymentRepo(d), db.NewConfigRepo(d), mbClient, serviceConfig.DebugMode)
		nSrv := server.NewConstructorEventServer(srv)
		generated.RegisterConstructorServiceServer(s, srv)
		egenerated.RegisterEventNotificationServiceServer(s, nSrv)
	})

	go msgBusListener(mbClient)

	grpcServer.StartServer()

}

func msgBusListener(m mb.MsgBusServiceClient) {
	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}
	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
