package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	"gopkg.in/yaml.v2"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/cmd/version"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg"

	generated "github.com/ukama/ukama/systems/subscriber/sim-manager/pb/gen"

	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/subscriber/sim-manager/pkg/server"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"google.golang.org/grpc"
)

var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")

	initConfig()

	metrics.StartMetricsServer(svcConf.Metrics)

	runGrpcServer()
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
	svcConf = &pkg.Config{
		Grpc: &uconf.Grpc{
			Port: 9090,
		},
		Metrics: &uconf.Metrics{
			Port: 10250,
		},
	}

	err := config.NewConfReader(pkg.ServiceName).Read(svcConf)
	if err != nil {
		log.Fatalf("Error reading config file. Error: %v", err)
	} else if svcConf.DebugMode {
		// output config in debug mode
		b, err := yaml.Marshal(svcConf)
		if err != nil {
			logrus.Infof("Config:\n%s", string(b))
		}
	}

	pkg.IsDebugMode = svcConf.DebugMode
}

func runGrpcServer() {
	simManagerServer := server.NewSimManagerServer()

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		generated.RegisterSimManagerServiceServer(s, simManagerServer)
	})

	grpcServer.StartServer()
}
