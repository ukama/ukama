package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	pb "github.com/ukama/ukama/systems/registry/org/pb/gen"
	provider "github.com/ukama/ukama/systems/registry/org/pkg/providers"
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

func initConfig() {
	svcConf = &pkg.Config{
		DB: &uconf.Database{
			DbName: pkg.ServiceName,
		},
		Grpc: &uconf.Grpc{
			Port: 9091,
		},
		Metrics: &uconf.Metrics{
			Port: 10251,
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

	err := d.Init(&db.Org{}, &db.User{}, &db.OrgUser{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}

	return d
}

func runGrpcServer(gormdb sql.Db) {
	regServer := server.NewOrgServer(db.NewOrgRepo(gormdb),
		db.NewUserRepo(gormdb),
		provider.NewUserClientProvider(svcConf.UsersHost),
	)

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		pb.RegisterOrgServiceServer(s, regServer)
	})

	grpcServer.StartServer()
}
