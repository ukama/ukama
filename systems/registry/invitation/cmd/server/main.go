package main

import (
	"os"

	generated "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/systems/common/msgBusServiceClient"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/uuid"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/registry/invitation/cmd/version"
	"github.com/ukama/ukama/systems/registry/invitation/pkg"

	"github.com/ukama/ukama/systems/registry/invitation/pkg/db"
	"github.com/ukama/ukama/systems/registry/invitation/pkg/providers"
	"github.com/ukama/ukama/systems/registry/invitation/pkg/server"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	invitationDb := initDb()
	runGrpcServer(invitationDb)
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

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Invitation{})
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

	notificationClient, err := providers.NewNotificationClient(serviceConfig.NotificationHost, pkg.IsDebugMode)
	if err != nil {
		logrus.Fatalf("Notification Client initilization failed. Error: %v", err.Error())
	}
	nucleusP := providers.NewNucleusClientProvider(serviceConfig.OrgRegistryHost, serviceConfig.DebugMode)

	mbClient := msgBusServiceClient.NewMsgBusClient(serviceConfig.MsgClient.Timeout, serviceConfig.OrgName, pkg.SystemName, pkg.ServiceName, instanceId, serviceConfig.Queue.Uri, serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange, serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue, serviceConfig.MsgClient.RetryCount, serviceConfig.MsgClient.ListenerRoutes)
	invitationServer := server.NewInvitationServer(db.NewInvitationRepo(gormdb),
		serviceConfig.InvitationExpiryTime, serviceConfig.AuthLoginbaseURL,
		notificationClient, nucleusP, mbClient, serviceConfig.OrgName, serviceConfig.TemplateName)
	log.Debugf("MessageBus Client is %+v", mbClient)
	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterInvitationServiceServer(s, invitationServer)
	})

	go grpcServer.StartServer()

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
