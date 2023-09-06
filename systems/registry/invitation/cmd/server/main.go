package main

import (
	"os"

	generated "github.com/ukama/ukama/systems/registry/invitation/pb/gen"
	"gopkg.in/yaml.v2"

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
	networkDb := initDb()
	runGrpcServer(networkDb)
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

	notificationClient, err := providers.NewNotificationClient(serviceConfig.NotificationHost, pkg.IsDebugMode)
	if err != nil {
		logrus.Fatalf("Notification Client initilization failed. Error: %v", err.Error())
	}
	nucleusP := providers.NewNucleusClientProvider(serviceConfig.OrgRegistryHost, serviceConfig.DebugMode)

	invitationServer := server.NewInvitationServer(db.NewInvitationRepo(gormdb), serviceConfig.InvitationExpiryTime, serviceConfig.AuthLoginbaseURL, notificationClient, nucleusP)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		generated.RegisterInvitationServiceServer(s, invitationServer)
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
