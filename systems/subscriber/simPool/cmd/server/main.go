package main

import (
	"os"

	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/subscriber/simPool/pkg/server"

	"github.com/num30/config"
	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/systems/subscriber/simPool/pkg"

	"github.com/ukama/ukama/systems/subscriber/simPool/cmd/version"

	"github.com/ukama/ukama/systems/subscriber/simPool/pkg/db"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	generated "github.com/ukama/ukama/systems/subscriber/simPool/pb"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()
	rateDb := initDb()
	runGrpcServer(rateDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = &pkg.Config{
		DB: &uconf.Database{
			DbName: pkg.ServiceName,
		},
	}
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
	err := d.Init(&db.SimPool{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {

		srv := server.NewSimPoolServer(db.NewSimPoolRepo(gormdb))
		generated.RegisterSimPoolServiceServer(s, srv)
	})

	grpcServer.StartServer()
}
