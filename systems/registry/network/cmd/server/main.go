package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	"gopkg.in/yaml.v2"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/registry/network/cmd/version"
	"github.com/ukama/ukama/systems/registry/network/pkg"

	generated "github.com/ukama/ukama/systems/registry/network/pb/gen"

	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/registry/network/pkg/db"
	"github.com/ukama/ukama/systems/registry/network/pkg/providers"
	"github.com/ukama/ukama/systems/registry/network/pkg/server"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()

	metrics.StartMetricsServer(svcConf.Metrics)

	networkDb := initDb()

	runGrpcServer(networkDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	svcConf = &pkg.Config{
		DB: &uconf.Database{
			DbName: pkg.ServiceName,
		},
		Grpc: &uconf.Grpc{
			Port: 9093,
		},
		Metrics: &uconf.Metrics{
			Port: 10253,
		},
	}

	err := config.NewConfReader(pkg.ServiceName).Read(svcConf)
	if err != nil {
		log.Fatalf("Error reading config file. Error: %v", err)
	} else if svcConf.DebugMode {
		// output config in debug mode
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

	err := d.Init(&db.Org{}, &db.Network{}, &db.Site{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	return d
}

func runGrpcServer(gormdb sql.Db) {
	networkServer := server.NewNetworkServer(db.NewNetRepo(gormdb),
		db.NewOrgRepo(gormdb), db.NewSiteRepo(gormdb),
		providers.NewOrgClientProvider(svcConf.OrgHost))

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		generated.RegisterNetworkServiceServer(s, networkServer)
	})

	grpcServer.StartServer()
}
