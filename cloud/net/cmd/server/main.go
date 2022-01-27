package main

import (
	pb "github.com/ukama/ukamaX/cloud/net/pb/gen"
	"github.com/ukama/ukamaX/cloud/net/pkg/server"
	"os"

	"github.com/ukama/ukamaX/cloud/net/pkg"

	"github.com/ukama/ukamaX/cloud/net/cmd/version"

	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
	ugrpc "github.com/ukama/ukamaX/common/grpc"

	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()

	nnsClient := server.NewNns(serviceConfig)

	go func() {
		srv := server.NewHttpServer(serviceConfig.Http, serviceConfig.Grpc, serviceConfig.NodeMetricsPort, nnsClient)
		srv.RunHttpServer()
	}()
	runGrpcServer(nnsClient)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}

func runGrpcServer(nns *server.Nns) {
	srv := server.NewNnsServer(nns)
	grpcServer := ugrpc.NewGrpcServer(serviceConfig.Grpc, func(s *grpc.Server) {
		pb.RegisterNnsServer(s, srv)
	})

	grpcServer.StartServer()
}
