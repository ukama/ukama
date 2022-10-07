package main

import (
	"os"

	// bootstrap "github.com/ukama/ukama/services/bootstrap/client"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/msgbus"

	db2 "github.com/ukama/ukama/systems/registry/network/pkg/db"

	"github.com/ukama/ukama/systems/registry/network/cmd/version"
	"github.com/ukama/ukama/systems/registry/network/pkg"

	generated "github.com/ukama/ukama/systems/registry/network/pb/gen"

	"github.com/ukama/ukama/systems/registry/network/pkg/server"

	confr "github.com/num30/config"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")
	initConfig()
	metrics.StartMetricsServer(&svcConf.Metrics)
	networkDb := initDb()
	runGrpcServer(networkDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	svcConf = &pkg.Config{
		DB: config.Database{
			DbName: pkg.ServiceName,
		},
	}
	err := confr.NewConfReader(pkg.ServiceName).Read(svcConf)
	if err != nil {
		log.Fatalf("Failed to read config. Error: %v", err)
	}
	pkg.IsDebugMode = svcConf.DebugMode
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
	// bootstrapCl := bootstrap.NewBootstrapClient(svcConf.BootstrapUrl, bootstrap.NewAuthenticator(svcConf.BootstrapAuth))
	if svcConf.Debug.DisableBootstrap {
		// bootstrapCl = bootstrap.DummyBootstrapClient{}
	}

	pub, err := msgbus.NewQPub(svcConf.Queue.Uri, pkg.ServiceName, pkg.InstanceId)
	if err != nil {
		log.Fatalf("Failed to create publisher. Error: %v", err)
	}

	regServer := server.NewNetworkServer(db2.NewOrgRepo(gormdb),
		db2.NewNodeRepo(gormdb),
		db2.NewNetRepo(gormdb),
		// bootstrapCl,
		pub)

	grpcServer := ugrpc.NewGrpcServer(svcConf.Grpc, func(s *grpc.Server) {
		generated.RegisterNetworkServiceServer(s, regServer)
	})

	grpcServer.StartServer()
}
