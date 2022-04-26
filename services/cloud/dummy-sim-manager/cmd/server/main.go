package main

import (
	"github.com/sirupsen/logrus"
	"os"

	"github.com/ukama/ukamaX/cloud/dummy-sim-manager/pkg"

	"github.com/ukama/ukamaX/cloud/dummy-sim-manager/cmd/version"

	"github.com/ukama/ukamaX/cloud/hss/pb/gen/simmgr"
	ccmd "github.com/ukama/ukamaX/common/cmd"
	"github.com/ukama/ukamaX/common/config"
	ugrpc "github.com/ukama/ukamaX/common/grpc"
	"google.golang.org/grpc"
)

var serviceConfig *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	initConfig()
	runGrpcServer()
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	serviceConfig = pkg.NewConfig()
	config.LoadConfig(pkg.ServiceName, serviceConfig)
	pkg.IsDebugMode = serviceConfig.DebugMode
}

func runGrpcServer() {

	grpcServer := ugrpc.NewGrpcServer(serviceConfig.Grpc, func(s *grpc.Server) {

		var storage pkg.Storage
		if serviceConfig.EtcdEnabled {
			storage = pkg.NewEtcdStorage(serviceConfig.EtcdHost)
			logrus.Infof("Etcd storage enabled")
		} else {
			storage = pkg.NewMemStorage()
			logrus.Infof("In-memory storage enabled")
		}

		simmgr.RegisterSimManagerServiceServer(s, pkg.NewSimManagerServer(storage))
	})

	grpcServer.StartServer()
}
