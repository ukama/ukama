package main

import (
	"os"

	pb "github.com/ukama/ukama/services/cloud/net/pb/gen"
	"github.com/ukama/ukama/services/cloud/net/pkg/server"

	"github.com/ukama/ukama/services/cloud/net/pkg"

	"github.com/ukama/ukama/services/cloud/net/cmd/version"

	dnspb "github.com/coredns/coredns/pb"
	ccmd "github.com/ukama/ukama/services/common/cmd"
	"github.com/ukama/ukama/services/common/config"
	ugrpc "github.com/ukama/ukama/services/common/grpc"
	"github.com/ukama/ukama/services/common/metrics"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()

	nnsClient := pkg.NewNns(serviceConfig)
	nodeOrgMapping := pkg.NewNodeToOrgMap(serviceConfig)

	metrics.StartMetricsServer(&serviceConfig.Metrics)
	go func() {
		srv := server.NewHttpServer(serviceConfig.Http, serviceConfig.Grpc, serviceConfig.NodeMetricsPort, nnsClient, nodeOrgMapping, &serviceConfig.OrgMetricsTarget)
		srv.RunHttpServer()
	}()
	runGrpcServer(nnsClient, nodeOrgMapping)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}

func runGrpcServer(nns *pkg.Nns, nodeOrgMapping *pkg.NodeOrgMap) {

	grpcServer := ugrpc.NewGrpcServer(serviceConfig.Grpc, func(s *grpc.Server) {
		srv := server.NewNnsServer(nns, nodeOrgMapping)
		pb.RegisterNnsServer(s, srv)

		dnspb.RegisterDnsServiceServer(s, server.NewDnsServer(nns, serviceConfig.Dns))
	})

	grpcServer.StartServer()
}
