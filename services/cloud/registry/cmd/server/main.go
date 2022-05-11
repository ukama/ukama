package main

import (
	"github.com/ukama/ukama/services/common/msgbus"
	"os"

	db2 "github.com/ukama/ukama/services/cloud/registry/pkg/db"

	"github.com/ukama/ukama/services/cloud/registry/pkg/bootstrap"

	"github.com/ukama/ukama/services/cloud/registry/cmd/version"
	"github.com/ukama/ukama/services/cloud/registry/pkg"

	generated "github.com/ukama/ukama/services/cloud/registry/pb/gen"

	"github.com/ukama/ukama/services/cloud/registry/pkg/server"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/services/common/config"
	ugrpc "github.com/ukama/ukama/services/common/grpc"
	"github.com/ukama/ukama/services/common/sql"
	"google.golang.org/grpc"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()
	registryDb := initDb()
	runGrpcServer(registryDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	svcConf = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, svcConf)
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)
	err := d.Init(&db2.Org{}, &db2.Network{}, &db2.Site{}, &db2.Node{})
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

	pub, err := msgbus.NewQPub(svcConf.Queue.Uri)
	if err != nil {
		log.Fatalf("Failed to create publisher. Error: %v", err)
	}

	regServer := server.NewRegistryServer(db2.NewOrgRepo(gormdb),
		db2.NewNodeRepo(gormdb),
		db2.NewNetRepo(gormdb),
		bootstrapCl,
		svcConf.DeviceGatewayHost,
		pub)

	grpcServer := ugrpc.NewGrpcServer(svcConf.Grpc, func(s *grpc.Server) {
		generated.RegisterRegistryServiceServer(s, regServer)
	})

	grpcServer.StartServer()
}
