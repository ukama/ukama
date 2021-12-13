package main

import (
	"github.com/ukama/ukamaX/cloud/hss/pb/gen"
	"github.com/ukama/ukamaX/cloud/hss/pkg/server"
	"os"

	"github.com/ukama/ukamaX/cloud/hss/pkg"

	"github.com/ukama/ukamaX/cloud/hss/cmd/version"

	"github.com/ukama/ukamaX/cloud/hss/pkg/db"

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
	config.LoadConfig("", serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Org{}, &db.Imsi{}, &db.User{}, &db.Guti{}, &db.Tai{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	imsiService := server.NewImsiService(db.NewImsiRepo(gormdb), db.NewGutiRepo(gormdb), pkg.NewHssQueue(serviceConfig.Queue))
	userService := server.NewUserService(db.NewUserRepo(gormdb), db.NewImsiRepo(gormdb))

	grpcServer := ugrpc.NewGrpcServer(serviceConfig.Grpc, func(s *grpc.Server) {
		gen.RegisterImsiServiceServer(s, imsiService)
		gen.RegisterUserServiceServer(s, userService)
	})

	grpcServer.StartServer()
}
