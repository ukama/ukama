package main

import (
	"os"

	"github.com/ukama/ukamaX/cloud/registry/pkg"
	"github.com/ukama/ukamaX/cloud/registry/pkg/bootstrap"

	"github.com/ukama/ukamaX/cloud/registry/cmd/version"

	generated "github.com/ukama/ukamaX/cloud/registry/pb/gen"

	"github.com/ukama/ukamaX/cloud/registry/internal/db"
	"github.com/ukama/ukamaX/cloud/registry/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
	ugrpc "github.com/ukama/ukamaX/common/grpc"
	"github.com/ukama/ukamaX/common/sql"
	"google.golang.org/grpc"
)

var svcConf *pkg.Config

const ServiceName = "registry"

func main() {
	ccmd.ProcessVersionArgument(ServiceName, os.Args, version.Version)

	initConfig()
	registryDb := initDb()
	runGrpcServer(registryDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig("", svcConf)
}
func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)
	err := d.Init(&db.Org{}, &db.Network{}, &db.Site{}, &db.Node{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	regServer := server.NewRegistryServer(db.NewOrgRepo(gormdb),
		db.NewNodeRepo(gormdb),
		bootstrap.NewBootstrapClient(svcConf.BootstrapUrl, bootstrap.NewAuthenticator(svcConf.BootstrapAuth)),
		svcConf.DeviceGatewayHost)

	grpcServer := ugrpc.NewGrpcServer(svcConf.Grpc, func(s *grpc.Server) {
		generated.RegisterRegistryServiceServer(s, regServer)
	})

	grpcServer.StartServer()
}
