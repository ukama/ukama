package main

import (
	"context"
	"github.com/ukama/ukama/services/cloud/users/pkg/server"
	"github.com/ukama/ukama/services/cloud/users/pkg/sims"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"

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
	err := d.Init(&db.Org{}, &db.Simcard{}, &db.User{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {

	simMgr, conn := newSimManagerClient(serviceConfig.SimManager)
	defer conn.Close()

	simPool, pcon := NewIccidPool(serviceConfig.SimManager)
	defer pcon.Close()

	userService := server.NewUserService(db.NewUserRepo(gormdb),
		pkg.NewImsiClientProvider(serviceConfig.HssHost),
		db.NewSimcardRepo(gormdb),
		sims.NewSimProvider(serviceConfig.SimTokenKey, simPool),
		simMgr,
		serviceConfig.SimManager.Name+":"+serviceConfig.SimManager.Host)

	grpcServer := ugrpc.NewGrpcServer(serviceConfig.Grpc, func(s *grpc.Server) {

		gen.RegisterUserServiceServer(s, userService)
	})

	grpcServer.StartServer()
}

func NewIccidPool(conf pkg.SimManager) (pbclient.SimPoolClient, io.Closer) {
	conn := createGrpcConn(conf)
	return pbclient.NewSimPoolClient(conn), conn
}

func newSimManagerClient(conf pkg.SimManager) (client pbclient.SimManagerServiceClient, connection io.Closer) {
	conn := createGrpcConn(conf)

	return pbclient.NewSimManagerServiceClient(conn), conn
}

func createGrpcConn(conf pkg.SimManager) *grpc.ClientConn {
	var conn *grpc.ClientConn

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	log.Infoln("Connecting to service ", conf.Host)

	conn, err := grpc.DialContext(ctx, conf.Host, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(
			grpc.ConnectParams{
				MinConnectTimeout: conf.Timeout,
			}))
	if err != nil {
		log.Fatalf("Failed to connect to service %s. Error: %v", conf.Host, err)
	}
	return conn
}
