package main

import (
	"context"
	"io"
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	"github.com/ukama/ukama/systems/registry/users/pkg/server"
	"github.com/ukama/ukama/systems/registry/users/pkg/sims"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ukama/ukama/systems/registry/users/pkg"

	"github.com/ukama/ukama/systems/registry/users/cmd/version"

	"github.com/ukama/ukama/systems/registry/users/pkg/db"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	"github.com/ukama/ukama/systems/common/config"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/registry/users/pb/gen"
	pbclient "github.com/ukama/ukama/systems/registry/users/pb/gen/simmgr"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()

	metrics.StartMetricsServer(&serviceConfig.Metrics)

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

	err := d.Init(&db.Org{}, &db.Simcard{}, &db.User{}, &db.Service{})
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
		serviceConfig.SimManager.Name+":"+serviceConfig.SimManager.Host,
	)

	grpcServer := ugrpc.NewGrpcServer(serviceConfig.Grpc, func(s *grpc.Server) {
		gen.RegisterUserServiceServer(s, userService)
	})

	grpcServer.StartServer()
}

func NewIccidPool(conf pkg.SimManager) (pbclient.SimPoolClient, io.Closer) {
	log.Info("Connecting to simPool")

	conn := createGrpcConn(conf)

	return pbclient.NewSimPoolClient(conn), conn
}

func newSimManagerClient(conf pkg.SimManager) (client pbclient.SimManagerServiceClient, connection io.Closer) {
	log.Info("Connecting to sim manager")

	conn := createGrpcConn(conf)

	return pbclient.NewSimManagerServiceClient(conn), conn
}

func createGrpcConn(conf pkg.SimManager) *grpc.ClientConn {
	var conn *grpc.ClientConn

	// connect to Grpc service
	ctx, cancel := context.WithTimeout(context.Background(), conf.Timeout)
	defer cancel()

	log.Infoln("Connecting to service ", conf.Host)

	conn, err := grpc.DialContext(ctx, conf.Host, grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: conf.Timeout}))
	if err != nil {
		cancel()
		log.Fatalf("Failed to connect to service %s. Error: %v", conf.Host, err)
	}

	return conn
}
