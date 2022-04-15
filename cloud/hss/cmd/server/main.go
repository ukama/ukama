package main

import (
	"github.com/ukama/ukamaX/cloud/hss/pb/gen"
	pbclient "github.com/ukama/ukamaX/cloud/hss/pb/gen/simmgr"
	"github.com/ukama/ukamaX/cloud/hss/pkg/server"
	"github.com/ukama/ukamaX/cloud/hss/pkg/sims"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"time"

	"github.com/ukama/ukamaX/cloud/hss/pkg"

	"github.com/ukama/ukamaX/cloud/hss/cmd/version"

	"github.com/ukama/ukamaX/cloud/hss/pkg/db"

	"context"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
	ugrpc "github.com/ukama/ukamaX/common/grpc"

	"github.com/ukama/ukamaX/common/sql"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()
	hssDb := initDb()
	runGrpcServer(hssDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}

func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(serviceConfig.DB, serviceConfig.DebugMode)
	err := d.Init(&db.Org{}, &db.Imsi{}, &db.User{}, &db.Guti{}, &db.Tai{}, &db.Simcard{}, &db.SimPool{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	reqGenerator, err := pkg.NewDeviceFeederReqGenerator(serviceConfig.Queue.Uri)
	if err != nil {
		log.Fatalf("Failed to create device feeder request generator. Error: %v", err)
	}

	subs := server.NewHssEventsSubscribers(pkg.NewHssNotifications(serviceConfig.Queue), reqGenerator)
	client, conn := newSimManagerClient()
	defer conn.Close()

	simPoolRepo := db.NewIccidpoolRepo(gormdb)

	imsiService := server.NewImsiService(db.NewImsiRepo(gormdb), db.NewGutiRepo(gormdb), subs)
	userService := server.NewUserService(db.NewUserRepo(gormdb),
		db.NewImsiRepo(gormdb),
		db.NewSimcardRepo(gormdb),
		sims.NewSimProvider(serviceConfig.SimTokenKey, simPoolRepo),
		client,
		serviceConfig.SimManager.Name+":"+serviceConfig.SimManager.Host)

	grpcServer := ugrpc.NewGrpcServer(serviceConfig.Grpc, func(s *grpc.Server) {
		gen.RegisterImsiServiceServer(s, imsiService)
		gen.RegisterUserServiceServer(s, userService)
	})

	grpcServer.StartServer()
}

func newSimManagerClient() (client pbclient.SimManagerServiceClient, connection io.Closer) {
	var conn *grpc.ClientConn
	if serviceConfig.SimManager.Disabled {
		return &pkg.SimManagerStub{}, &pkg.CloserStub{}
	}

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Infoln("Connecting to service ", serviceConfig.SimManager.Host)

	conn, err := grpc.DialContext(ctx, serviceConfig.SimManager.Host, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to service %s. Error: %v", serviceConfig.SimManager.Host, err)
	}

	return pbclient.NewSimManagerServiceClient(conn), conn
}
