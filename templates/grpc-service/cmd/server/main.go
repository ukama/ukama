package main

import (
	"os"
	"ukamaX/service-boilerplate/pkg"

	"github.com/ukama/ukamaX/cloud/foo/pkg"

	"github.com/ukama/ukamaX/cloud/foo/cmd/version"

	"github.com/ukama/ukamaX/cloud/foo/pkg/db"

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
	fooDb := initDb()
	runGrpcServer(fooDb)
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
	err := d.Init(&db.Foo{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {

	grpcServer := ugrpc.NewGrpcServer(serviceConfig.Grpc, func(s *grpc.Server) {
		log.Fatal("Register your service here")
	})

	grpcServer.StartServer()
}
