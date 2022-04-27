package main

import (
	"os"

	"github.com/ukama/ukama/services/cloud/hss/pb/gen"
	"github.com/ukama/ukama/services/cloud/hss/pkg"
	"github.com/ukama/ukama/services/cloud/hss/pkg/server"

	"github.com/ukama/ukama/services/cloud/hss/cmd/version"

	"github.com/ukama/ukama/services/cloud/hss/pkg/db"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
	ugrpc "github.com/ukama/ukamaX/common/grpc"

	"github.com/ukama/ukamaX/common/sql"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()
	hssDb := initDb()
	runGrpcServer(hssDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Org{}, &db.Imsi{}, &db.Guti{}, &db.Tai{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	reqGenerator, err := pkg.NewDeviceFeederReqGenerator(serviceConfig.Queue.Uri)
	if err != nil {
		log.Fatalf("Failed to create device feeder request generator. Error: %v", err)
	}

	// hss service
	subs := server.NewHssEventsSubscribers(pkg.NewHssNotifications(serviceConfig.Queue), reqGenerator)
	imsiService := server.NewImsiService(db.NewImsiRepo(gormdb), db.NewGutiRepo(gormdb), subs)

	// users service

	rpcServer := ugrpc.NewGrpcServer(serviceConfig.Grpc, func(s *grpc.Server) {
		gen.RegisterImsiServiceServer(s, imsiService)
	})
	rpcServer.StartServer()

}
