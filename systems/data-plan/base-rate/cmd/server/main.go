package main

import (
	"os"

	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/server"

	"github.com/num30/config"
	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg"

	"github.com/ukama/ukama/systems/data-plan/base-rate/cmd/version"

	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mbc "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	"github.com/ukama/ukama/systems/common/sql"
	generated "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"

	"google.golang.org/grpc"
)

var svcConf = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()
	rateDb := initDb()
	runGrpcServer(rateDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	svcConf = &pkg.Config{
		DB: &uconf.Database{
			DbName: pkg.ServiceName,
		},
	}
	err := config.NewConfReader(pkg.ServiceName).Read(svcConf)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if svcConf.DebugMode {
		b, err := yaml.Marshal(svcConf)
		if err != nil {
			logrus.Infof("Config:\n%s", string(b))
		}
	}

	pkg.IsDebugMode = svcConf.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)
	err := d.Init(&db.Rate{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	// instanceId := os.Getenv("POD_NAME")

	// mbClient := mbc.NewMsgBusClient(serviceConfig.MsgClient.Timeout, pkg.SystemName,
	// 	pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
	// 	serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
	// 	serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
	// 	serviceConfig.MsgClient.RetryCount,
	// 	serviceConfig.MsgClient.ListenerRoutes)

	// log.Debugf("MessageBus Client is %+v", mbClient)

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {

		srv := server.NewBaseRateServer(db.NewBaseRateRepo(gormdb))
		generated.RegisterBaseRatesServiceServer(s, srv)
	})

	// go msgBusListener(mbClient)

	grpcServer.StartServer()

}

func msgBusListener(m mbc.MsgBusServiceClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
