package main

import (
	"github.com/ukama/ukama/services/cloud/node/pkg/server"
	"github.com/ukama/ukama/services/common/msgbus"
	"os"

	"github.com/num30/config"
	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/services/cloud/node/pkg"

	"github.com/ukama/ukama/services/cloud/node/cmd/version"

	"github.com/ukama/ukama/services/cloud/node/pkg/db"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	generated "github.com/ukama/ukama/services/cloud/node/pb/gen"
	ccmd "github.com/ukama/ukama/services/common/cmd"
	ugrpc "github.com/ukama/ukama/services/common/grpc"
	"github.com/ukama/ukama/services/common/sql"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()
	nodeDb := initDb()
	runGrpcServer(nodeDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = &pkg.Config{}
	serviceConfig.DB.DbName = pkg.ServiceName
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			logrus.Infof("Config:\n%s", string(b))
		}
	}

	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Node{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	instanceId := os.Getenv("POD_NAME")
	grpcServer := ugrpc.NewGrpcServer(serviceConfig.Grpc, func(s *grpc.Server) {
		pub, err := msgbus.NewQPub(serviceConfig.Queue.Uri, pkg.ServiceName, instanceId)
		if err != nil {
			log.Fatalf("Failed to create publisher. Error: %v", err)
		}

		srv := server.NewNodeServer(db.NewNodeRepo(gormdb), pub)
		generated.RegisterNodeServiceServer(s, srv)
	})

	grpcServer.StartServer()
}
