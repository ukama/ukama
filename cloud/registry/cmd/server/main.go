package main

import (
	db2 "github.com/ukama/ukamaX/cloud/registry/pkg/db"
	"os"

	"github.com/ukama/ukamaX/cloud/registry/pkg/bootstrap"

	"github.com/ukama/ukamaX/cloud/registry/cmd/version"
	"github.com/ukama/ukamaX/cloud/registry/pkg"

	generated "github.com/ukama/ukamaX/cloud/registry/pb/gen"

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
	config.LoadConfig(ServiceName, svcConf)
}
func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)
	err := d.Init(&db2.Org{}, &db2.Network{}, &db2.Site{}, &db2.Node{}, &db2.NodeIp{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	bootstrapCl := bootstrap.NewBootstrapClient(svcConf.BootstrapUrl, bootstrap.NewAuthenticator(svcConf.BootstrapAuth))
	if svcConf.Debug.DisableBootstrap {
		bootstrapCl = bootstrap.DummyBootstrapClient{}
	}

	regServer := server.NewRegistryServer(db2.NewOrgRepo(gormdb),
		db2.NewNodeRepo(gormdb),
		db2.NewNetRepo(gormdb),
		bootstrapCl,
		svcConf.DeviceGatewayHost)

	grpcServer := ugrpc.NewGrpcServer(svcConf.Grpc, func(s *grpc.Server) {
		generated.RegisterRegistryServiceServer(s, regServer)
	})

	grpcServer.StartServer()
}
