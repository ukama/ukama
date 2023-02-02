package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	"gopkg.in/yaml.v2"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/subscriber/test-agent/cmd/version"
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg"

	generated "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"

	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg/server"
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg/storage"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"google.golang.org/grpc"
)

var svcConf = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	/* Log level */
	log.SetLevel(log.TraceLevel)
	log.Infof("Starting %s service", pkg.ServiceName)

	initConfig()

	metrics.StartMetricsServer(svcConf.Metrics)

	runGrpcServer()

	log.Infof("Exiting service %s", pkg.ServiceName)
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
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

	log.Debugf("\nService: %s Service: %+v ", pkg.ServiceName, svcConf.Service)

	pkg.IsDebugMode = svcConf.DebugMode
}

func runGrpcServer() {
	testAgentServer := server.NewTestAgentServer(storage.NewMemStorage(make(map[string]*storage.SimInfo)))

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		generated.RegisterTestAgentServiceServer(s, testAgentServer)
	})

	grpcServer.StartServer()
}
