package main

import (
	"os"

	"github.com/ukama/ukama/systems/common/metrics"
	"gopkg.in/yaml.v2"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/subscriber/test-agent/cmd/version"
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg"

	generated "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"

<<<<<<< HEAD
	uconf "github.com/ukama/ukama/systems/common/config"
=======
>>>>>>> subscriber-sys_sim-manager
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg/server"
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg/storage"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	"google.golang.org/grpc"
)

<<<<<<< HEAD
var svcConf *pkg.Config

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)
	pkg.InstanceId = os.Getenv("POD_NAME")
=======
var svcConf = pkg.NewConfig(pkg.ServiceName)

func main() {
	ccmd.ProcessVersionArgument(pkg.ServiceName, os.Args, version.Version)

	/* Log level */
	log.SetLevel(log.TraceLevel)
	log.Infof("Starting %s service", pkg.ServiceName)
>>>>>>> subscriber-sys_sim-manager

	initConfig()

	metrics.StartMetricsServer(svcConf.Metrics)

	runGrpcServer()
<<<<<<< HEAD
=======

	log.Infof("Exiting service %s", pkg.ServiceName)
>>>>>>> subscriber-sys_sim-manager
}

// initConfig reads in config file, ENV variables, and flags if set.
func initConfig() {
<<<<<<< HEAD
	svcConf = &pkg.Config{
		Grpc: &uconf.Grpc{
			Port: 9090,
		},
		Metrics: &uconf.Metrics{
			Port: 10250,
		},
	}

=======
>>>>>>> subscriber-sys_sim-manager
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

<<<<<<< HEAD
=======
	log.Debugf("\nService: %s Service: %+v ", pkg.ServiceName, svcConf.Service)

>>>>>>> subscriber-sys_sim-manager
	pkg.IsDebugMode = svcConf.DebugMode
}

func runGrpcServer() {
<<<<<<< HEAD
	testAgentServer := server.NewTestAgentServer(storage.NewMemStorage())
=======
	testAgentServer := server.NewTestAgentServer(storage.NewMemStorage(make(map[string]*storage.SimInfo)))
>>>>>>> subscriber-sys_sim-manager

	grpcServer := ugrpc.NewGrpcServer(*svcConf.Grpc, func(s *grpc.Server) {
		generated.RegisterTestAgentServiceServer(s, testAgentServer)
	})

	grpcServer.StartServer()
}
