package main

import (
	"os"

	"github.com/cloudflare/cfssl/log"
	"github.com/num30/config"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	"github.com/ukama/ukama/systems/node/controller/pkg/db"
	"github.com/ukama/ukama/systems/node/controller/pkg/providers"
	"github.com/ukama/ukama/systems/node/controller/pkg/server"
	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/systems/node/controller/pkg"

	"github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/uuid"
	"github.com/ukama/ukama/systems/node/controller/cmd/version"
	pb "github.com/ukama/ukama/systems/node/controller/pb/gen"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	nodeLogDb := initDb()
	runGrpcServer(nodeLogDb)
	logrus.Infof("Starting %s", pkg.ServiceName)
	runGrpcServer(nodeLogDb)
}
func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.NodeLog{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}


func initConfig() {
	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		logrus.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			logrus.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = serviceConfig.DebugMode
}

func runGrpcServer(gormdb sql.Db) {

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		/* used on local machines */
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	reg := providers.NewRegistryProvider(serviceConfig.RegistryHost, serviceConfig.DebugMode)
	mbClient := msgBusServiceClient.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName, pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri, serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange, serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue, serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)
	controllerServer := server.NewControllerServer(mbClient, reg, pkg.IsDebugMode, serviceConfig.OrgName,db.NewNodeLogRepo(gormdb) )
	controllerEventServer := server.NewControllerEventServer(serviceConfig.OrgName, controllerServer)

	logrus.Debugf("MessageBus Client is %+v", mbClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		pb.RegisterControllerServiceServer(s, controllerServer)
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
