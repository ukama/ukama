package main

import (
	"os"

	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/data_plan/base_rate/pkg/server"

	"github.com/num30/config"
	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/systems/data_plan/base_rate/pkg"

	"github.com/ukama/ukama/systems/data_plan/base_rate/cmd/version"

	"github.com/ukama/ukama/systems/data_plan/base_rate/pkg/db"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mbc "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/sql"
	generated "github.com/ukama/ukama/systems/data_plan/base_rate/pb/gen"

	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config
var host = "localhost"
var port = 50051

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()
	rateDb := initDb()
	runGrpcServer(rateDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = &pkg.Config{
		DB: &uconf.Database{
			DbName: pkg.ServiceName,
		},
	}
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
	err := d.Init(&db.Rate{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	instanceId := os.Getenv("POD_NAME")

	mbClient := mbc.NewMsgBusClient(serviceConfig.MsgClient.Timeout, pkg.SystemName,
		pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
		serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
		serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
		serviceConfig.MsgClient.RetryCount,
		serviceConfig.MsgClient.ListenerRoutes)

	log.Debugf("MessageBus Client is %+v", mbClient)

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {

		srv := server.NewBaseRateServer(db.NewBaseRateRepo(gormdb))
		generated.RegisterBaseRatesServiceServer(s, srv)
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
