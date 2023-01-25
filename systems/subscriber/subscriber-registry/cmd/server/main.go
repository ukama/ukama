package main

import (
	"os"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/subscriber/subscriber-registry/pkg/server"

	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/systems/subscriber/subscriber-registry/pkg"

	"github.com/ukama/ukama/systems/subscriber/subscriber-registry/cmd/version"

	"github.com/ukama/ukama/systems/subscriber/subscriber-registry/pkg/db"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mbc "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/sql"
	generated "github.com/ukama/ukama/systems/subscriber/subscriber-registry/pb/gen"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()
	packageDb := initDb()
	runGrpcServer(packageDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			logrus.Infof("Config:\n%s", string(b))
		}
	}

	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Subscriber{})

	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	mbClient := mbc.NewMsgBusClient(serviceConfig.MsgClient.Timeout,
		pkg.SystemName,
		pkg.ServiceName,
		"susbcriber-registry",
		serviceConfig.Queue.Uri,
		"localhost:9090",
		serviceConfig.MsgClient.Host,
		serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue,
		serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)
	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {

		srv := server.NewSubscriberServer(db.NewSubscriberRepo(gormdb))
		generated.RegisterSubscriberRegistryServiceServer(s, srv)
	})
	go msgBusListener(mbClient)

	grpcServer.StartServer()
}

func msgBusListener(m *mbc.MsgBusClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
