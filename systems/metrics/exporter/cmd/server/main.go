package main

import (
	"os"

	"github.com/num30/config"
	"github.com/ukama/ukama/systems/metrics/exporter/pb/gen"
	"github.com/ukama/ukama/systems/metrics/exporter/pkg/server"
	"gopkg.in/yaml.v3"

	pkg "github.com/ukama/ukama/systems/metrics/exporter/pkg"

	"github.com/ukama/ukama/systems/metrics/exporter/cmd/version"

	log "github.com/sirupsen/logrus"
	ccmd "github.com/ukama/ukama/systems/common/cmd"
	ugrpc "github.com/ukama/ukama/systems/common/grpc"
	mb "github.com/ukama/ukama/systems/common/msgBusServiceClient"
	egen "github.com/ukama/ukama/systems/common/pb/gen/events"
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
	log.Infof("Initializing config")
	serviceConfig = pkg.NewConfig(pkg.ServiceName)
	err := config.NewConfReader(pkg.ServiceName).Read(serviceConfig)
	if err != nil {
		log.Fatal("Error reading config ", err)
	} else if pkg.IsDebugMode {
		b, err := yaml.Marshal(serviceConfig)
		if err != nil {
			log.Infof("Config:\n%s", string(b))
		}
	}
	pkg.IsDebugMode = serviceConfig.DebugMode
	log.Infof("Config: %+v", serviceConfig)
}

func runGrpcServer() {

	var mbClient mb.MsgBusServiceClient
	var instanceId string

	inst, ok := os.LookupEnv("POD_NAME")
	if ok {
		instanceId = inst
	} else {
		instanceId = pkg.InstanceId
	}

	if serviceConfig.IsMsgBus {
		mbClient = mb.NewMsgBusClient(serviceConfig.MsgClient.Timeout, pkg.SystemName,
			pkg.ServiceName, instanceId, serviceConfig.Queue.Uri,
			serviceConfig.Service.Uri, serviceConfig.MsgClient.Host, serviceConfig.MsgClient.Exchange,
			serviceConfig.MsgClient.ListenQueue, serviceConfig.MsgClient.PublishQueue,
			serviceConfig.MsgClient.RetryCount,
			serviceConfig.MsgClient.ListenerRoutes)

		log.Debugf("MessageBus Client is %+v", mbClient)
	}

	// Exporter service
	asrServer, err := server.NewExporterServer(serviceConfig.Org, mbClient)

	if err != nil {
		log.Fatalf("Exporter server initialization failed. Error: %v", err)
	}
	nSrv := server.NewExporterEventServer()

	rpcServer := ugrpc.NewGrpcServer(*serviceConfig.Grpc, func(s *grpc.Server) {
		gen.RegisterExporterServiceServer(s, asrServer)
		if serviceConfig.IsMsgBus {
			egen.RegisterEventNotificationServiceServer(s, nSrv)
		}
	})

	if serviceConfig.IsMsgBus {
		go msgBusListener(mbClient)
	}

	rpcServer.StartServer()

}

func msgBusListener(m mb.MsgBusServiceClient) {

	if err := m.Register(); err != nil {
		log.Fatalf("Failed to register to Message Client Service. Error %s", err.Error())
	}

	if err := m.Start(); err != nil {
		log.Fatalf("Failed to start to Message Client Service routine for service %s. Error %s", pkg.ServiceName, err.Error())
	}
}
