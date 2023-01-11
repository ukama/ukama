package main

import (
	"context"
	"os"

	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/server"

	"github.com/num30/config"
	"gopkg.in/yaml.v3"

	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg"

	"github.com/ukama/ukama/systems/data-plan/base-rate/cmd/version"

	"github.com/ukama/ukama/systems/data-plan/base-rate/pkg/db"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	generated "github.com/ukama/ukama/systems/data-plan/base-rate/pb/gen"
	msgClient "github.com/ukama/ukama/systems/init/msgClient/pb/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serviceConfig *pkg.Config
var host = "localhost"
var port = 50051

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
	err := d.Init(&db.Rate{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	grpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {

		srv := server.NewBaseRateServer(db.NewBaseRateRepo(gormdb))
		generated.RegisterBaseRatesServiceServer(s, srv)
	})

	grpcServer.StartServer()
	conn, err := grpc.Dial("localhost:7070", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial msg server: %v", err)
	}
	defer conn.Close()
	client := msgClient.NewMsgClientServiceClient(conn)

	res, err := client.RegisterService(context.Background(), &msgClient.RegisterServiceReq{
		SystemName:  pkg.SystemName,
		ServiceName: pkg.ServiceName,
		InstanceId:  pkg.InstanceId,
	})
	if err != nil {
		log.Fatalf("Error while Registering service: %v", err)
	}
	log.Printf("res msg data: %v", res)

}
