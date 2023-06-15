package main

import (
	"os"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/num30/config"

	"github.com/ukama/ukama/systems/notification/mailer/pkg/server"
	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/systems/notification/mailer/pkg"

	"github.com/ukama/ukama/systems/notification/mailer/cmd/version"

	pb "github.com/ukama/ukama/systems/notification/mailer/pb/gen"
	"github.com/ukama/ukama/systems/notification/mailer/pkg/db"

	"github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	initConfig()
	mailerDb := initDb()
	runGrpcServer(mailerDb)
}

func initConfig() {

	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		logrus.Fatal("Error reading config ", err)
	} else if serviceConfig.DebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			logrus.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	logrus.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Mailing{})

	if err != nil {
		logrus.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {

	instanceId := os.Getenv("POD_NAME")
	if instanceId == "" {
		inst := uuid.NewV4()
		instanceId = inst.String()
	}

	
	srv := server.NewMaillingServer(db.NewMaillingRepo(gormdb))

	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		pb.RegisterMaillingServiceServer(s, srv)
	})

	grpcServer.StartServer()
}


