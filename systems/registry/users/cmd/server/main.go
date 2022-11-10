package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/registry/users/pkg/server"

	"github.com/ukama/ukama/systems/registry/users/pkg"

	"github.com/ukama/ukama/systems/registry/users/cmd/version"

	"github.com/ukama/ukama/systems/registry/users/pkg/db"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/registry/users/pb/gen"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()

	usersDb := initDb()

	metrics.StartMetricsServer(serviceConfig.Metrics)

	runGrpcServer(usersDb)
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

	err := d.Init(&db.User{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	return d
}

func runGrpcServer(gormdb sql.Db) {
	userService := server.NewUserService(db.NewUserRepo(gormdb))
	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		gen.RegisterUserServiceServer(s, userService)
	})

	grpcServer.StartServer()
}
