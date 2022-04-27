package main

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/ukama/ukama/services/cloud/users/pkg/server"
	"github.com/ukama/ukama/services/cloud/users/pkg/sims"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ukama/ukama/services/cloud/users/pkg"

	"github.com/ukama/ukama/services/cloud/users/cmd/version"

	"github.com/ukama/ukama/services/cloud/users/pkg/db"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/services/cloud/users/pb/gen"
	pbclient "github.com/ukama/ukama/services/cloud/users/pb/gen/simmgr"
	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/services/common/config"
	ugrpc "github.com/ukama/ukama/services/common/grpc"
	"github.com/ukama/ukama/services/common/sql"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()
	usersDb := initDb()
	runGrpcServer(usersDb)
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
	err := d.Init(&db.Org{}, &db.Simcard{}, &db.SimPool{}, &db.User{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {

	client, conn := newSimManagerClient()
	defer conn.Close()

	simPoolRepo := db.NewIccidpoolRepo(gormdb)
	userService := server.NewUserService(db.NewUserRepo(gormdb),
		pkg.NewImsiClientProvider(serviceConfig.HssHost),
		db.NewSimcardRepo(gormdb),
		sims.NewSimProvider(serviceConfig.SimTokenKey, simPoolRepo),
		client,
		serviceConfig.SimManager.Name+":"+serviceConfig.SimManager.Host)

	grpcServer := ugrpc.NewGrpcServer(serviceConfig.Grpc, func(s *grpc.Server) {

		gen.RegisterUserServiceServer(s, userService)
	})

	grpcServer.StartServer()
}

func newSimManagerClient() (client pbclient.SimManagerServiceClient, connection io.Closer) {
	var conn *grpc.ClientConn

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Infoln("Connecting to Sim Manager service ", serviceConfig.SimManager.Host)

	conn, err := grpc.DialContext(ctx, serviceConfig.SimManager.Host, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(
			grpc.ConnectParams{
				MinConnectTimeout: time.Second * 5,
			}))
	if err != nil {
		log.Fatalf("Failed to connect to service %s. Error: %v", serviceConfig.SimManager.Host, err)
	}

	return pbclient.NewSimManagerServiceClient(conn), conn
}
