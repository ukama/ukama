package main

import (
	"os"

	bootstrap "github.com/ukama/ukama/services/bootstrap/client"
	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/common/msgbus"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	"github.com/ukama/ukama/systems/registry/org/pkg/server"

	"github.com/ukama/ukama/systems/registry/org/pkg"
	"gopkg.in/yaml.v2"

	"github.com/ukama/ukama/systems/registry/org/cmd/version"

	"github.com/ukama/ukama/systems/registry/org/pkg/db"

	"github.com/num30/config"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	uconf "github.com/ukama/ukama/systems/common/config"

	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")
	initConfig()
	orgDb := initDb()
	metrics.StartMetricsServer(svcConf.Metrics)

	runGrpcServer(orgDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	svcConf = &pkg.Config{
		DB: &uconf.Database{DbName: pkg.ServiceName},
	}
	// We change only DB name. Rest of the fields is set by default.
	svcConf.DB.DbName = pkg.ServiceName

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
	err := d.Init(&db.Org{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {

	bootstrapCl := bootstrap.NewBootstrapClient(svcConf.BootstrapUrl, bootstrap.NewAuthenticator(*svcConf.BootstrapAuth))
	if svcConf.Debug.DisableBootstrap {
		bootstrapCl = bootstrap.DummyBootstrapClient{}
	}

	pub, err := msgbus.NewQPub(svcConf.Queue.Uri, pkg.ServiceName, pkg.InstanceId)
	if err != nil {
		log.Fatalf("Failed to create publisher. Error: %v", err)
	}

	regServer := server.NewOrgServer(db.NewOrgRepo(gormdb),
		bootstrapCl,
		svcConf.DeviceGatewayHost,
		pub)

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		pb.RegisterOrgServiceServer(s, regServer)
	})

	grpcServer.StartServer()
}
