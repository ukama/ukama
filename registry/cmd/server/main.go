package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukamaX/common/config"
	"github.com/ukama/ukamaX/common/sql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"ukamaX/registry/internal"
	"ukamaX/registry/internal/db"
	"ukamaX/registry/pb/generated"
	"ukamaX/registry/pkg/server"
)

var svcConf *internal.Config

func main() {
	log.Infof("Starting the registry server")
	initConfig()
	registryDb := initDb()
	runGrpcServer(registryDb)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	svcConf = internal.NewConfig()
	config.LoadConfig("", svcConf)
}
func initDb() sql.Db {
	log.Infof("Initializing Database")
	d := sql.NewDb(svcConf.DB, svcConf.DebugMode)
	err := d.Init(&db.Org{}, &db.Network{}, &db.Site{}, &db.Node{})
	if err != nil {
		log.Fatalf("Database initialization failed. Error: %v", err)
	}
	return d
}

func runGrpcServer(gormdb sql.Db) {
	log.Infof("Starting gRpc on port " + fmt.Sprintf(":%d", svcConf.Grpc.Port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", svcConf.Grpc.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	regServer := server.NewRegistryServer(db.NewOrgRepo(gormdb), db.NewNodeRepo(gormdb))
	generated.RegisterRegistryServiceServer(s, regServer)
	generated.RegisterHealthServer(s, server.NewHealthChecker())
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
